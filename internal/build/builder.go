package build

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/yasufadhili/jawt/internal/compiler"
	"github.com/yasufadhili/jawt/internal/core"
	"github.com/yasufadhili/jawt/internal/diagnostic"
	"github.com/yasufadhili/jawt/internal/emitter"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type DocumentType int

const (
	DocumentTypePage DocumentType = iota
	DocumentTypeComponent
)

type DocumentInfo struct {
	Name         string
	RelPath      string
	AbsPath      string
	Type         DocumentType
	Dependencies []string
	DependedBy   []string
	IsCompiled   bool
	LastModified time.Time
	Hash         string
}

type ComponentInfo struct {
	DocumentInfo
	Props map[string]string
}

type PageInfo struct {
	DocumentInfo
	Route string
}

type BuildSystem struct {
	ctx        *core.JawtContext
	mu         sync.RWMutex
	docs       map[string]*DocumentInfo
	pages      map[string]*PageInfo
	comps      map[string]*ComponentInfo
	discoverer ProjectDiscoverer
	compiler   *CompilerRunner
	watcher    FileWatcher
	depGraph   DependencyGraph
}

type FileWatcher interface {
	OnEvent(handler func(fsnotify.Event))
	Start() error
	Stop() error
}

func NewBuildSystem(ctx *core.JawtContext, watcher FileWatcher) *BuildSystem {
	return &BuildSystem{
		ctx:        ctx,
		docs:       make(map[string]*DocumentInfo),
		pages:      make(map[string]*PageInfo),
		comps:      make(map[string]*ComponentInfo),
		discoverer: NewProjectDiscoverer(ctx),
		watcher:    watcher,
		compiler:   NewCompilerRunner(ctx),
		depGraph:   NewDependencyGraph(),
	}
}

// Initialise performs initial project discovery and compilation
func (bs *BuildSystem) Initialise() error {
	bs.ctx.Logger.Info("Initialising build system")

	if err := bs.generateWorkspaceConfigs(); err != nil {
		return fmt.Errorf("failed to generate workspace configs: %w", err)
	}

	if err := bs.syncWorkspaceSources(); err != nil {
		return fmt.Errorf("failed to sync workspace sources: %w", err)
	}

	if err := bs.DiscoverProject(); err != nil {
		return err
	}

	if err := bs.CompileAll(); err != nil {
		return err
	}

	bs.SetupWatcher()

	return nil
}

func (bs *BuildSystem) DiscoverProject() error {
	bs.ctx.Logger.Info("Discovering project documents")

	jmlFiles, err := bs.discoverer.DiscoverProjectFiles()
	if err != nil {
		return fmt.Errorf("failed to discover project files: %w", err)
	}

	// First pass: Add all documents to build system and dependency graph
	for _, path := range jmlFiles {
		docInfo, err := bs.discoverer.CreateDocumentInfo(path, bs.ctx.Paths.ProjectRoot)
		if err != nil {
			bs.ctx.Logger.Warn("Failed to process document",
				core.StringField("path", path),
				core.ErrorField(err))
			continue
		}

		bs.AddDocument(docInfo)

		// Add node to dependency graph
		if err := bs.depGraph.AddNode(docInfo.AbsPath, docInfo.Type); err != nil {
			bs.ctx.Logger.Warn("Failed to add document to dependency graph",
				core.StringField("path", path),
				core.ErrorField(err))
		}
	}

	// Second pass: Analyse dependencies and build graph
	if err := bs.buildDependencyGraph(); err != nil {
		return fmt.Errorf("failed to build dependency graph: %w", err)
	}

	if err := bs.depGraph.ValidateGraph(); err != nil {
		return fmt.Errorf("invalid dependency graph: %w", err)
	}

	if bs.depGraph.HasCycle() {
		cycles := bs.depGraph.GetCycles()
		bs.ctx.Logger.Warn("Circular dependencies detected",
			core.IntField("cycle_count", len(cycles)))
		// Could return error or just warn depending on requirements
	}

	bs.ctx.Logger.Info("Project discovery completed",
		core.IntField("pages", len(bs.pages)),
		core.IntField("components", len(bs.comps)))

	return nil
}

func (bs *BuildSystem) buildDependencyGraph() error {
	bs.ctx.Logger.Info("Building dependency graph")

	for path, doc := range bs.docs {
		dependencies, err := bs.extractDependencies(doc)
		if err != nil {
			bs.ctx.Logger.Error("Failed to extract dependencies",
				core.StringField("path", path),
				core.ErrorField(err))
			continue
		}

		for _, dep := range dependencies {
			if err := bs.depGraph.AddDependency(path, dep); err != nil {
				bs.ctx.Logger.Error("Failed to add dependency to graph",
					core.StringField("from", path),
					core.StringField("to", dep),
					core.ErrorField(err))
			}
		}
	}

	return nil
}

func (bs *BuildSystem) extractDependencies(doc *DocumentInfo) ([]string, error) {
	content, err := os.ReadFile(doc.AbsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", doc.AbsPath, err)
	}

	return ExtractDependencies(string(content)), nil
}

func (bs *BuildSystem) SetupWatcher() {
	bs.ctx.Logger.Info("Setting up file watcher")

	bs.watcher.OnEvent(func(event fsnotify.Event) {
		bs.HandleFileEvent(event)
	})
}

func (bs *BuildSystem) CompileAll() error {
	bs.ctx.Logger.Info("Compiling all documents")

	compilationOrder, err := bs.depGraph.GetCompilationOrder()
	if err != nil {
		return fmt.Errorf("failed to determine compilation order: %w", err)
	}

	bs.ctx.Logger.Info("Compilation order determined",
		core.IntField("document_count", len(compilationOrder)))

	// Compile in dependency order
	for _, path := range compilationOrder {
		if err := bs.CompileDocument(path); err != nil {
			bs.ctx.Logger.Error("Failed to compile document",
				core.StringField("path", path),
				core.ErrorField(err))
			return err
		}
	}

	return nil
}

func (bs *BuildSystem) HandleFileEvent(event fsnotify.Event) {
	bs.ctx.Logger.Info("Handling file event",
		core.StringField("operation", event.Op.String()),
		core.StringField("file", event.Name))

	if !bs.isJMLFile(event.Name) {
		bs.ctx.Logger.Debug("Ignoring non-JML file event",
			core.StringField("file", event.Name))
		return
	}

	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		bs.HandleFileCreated(event.Name)
	case event.Op&fsnotify.Write == fsnotify.Write:
		bs.HandleFileModified(event.Name)
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		bs.HandleFileDeleted(event.Name)
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		bs.HandleFileRenamed(event.Name)
	}
}

func (bs *BuildSystem) HandleFileCreated(path string) {
	bs.ctx.Logger.Info("JML file created", core.StringField("path", path))

	// Create DocumentInfo for the new file
	docInfo, err := bs.discoverer.CreateDocumentInfo(path, bs.ctx.Paths.ProjectRoot)
	if err != nil {
		bs.ctx.Logger.Error("Failed to create document info for new file",
			core.StringField("path", path),
			core.ErrorField(err))
		return
	}

	// Add to the build system (includes adding to dependency graph)
	bs.AddDocument(docInfo)

	// Extract and add dependencies to graph
	dependencies, err := bs.extractDependencies(docInfo)
	if err != nil {
		bs.ctx.Logger.Error("Failed to extract dependencies for new file",
			core.StringField("path", path),
			core.ErrorField(err))
	} else {
		for _, dep := range dependencies {
			if err := bs.depGraph.AddDependency(path, dep); err != nil {
				bs.ctx.Logger.Error("Failed to add dependency to graph",
					core.StringField("from", path),
					core.StringField("to", dep),
					core.ErrorField(err))
			}
		}
	}

	if err := bs.CompileDocument(path); err != nil {
		bs.ctx.Logger.Error("Failed to compile new file",
			core.StringField("path", path),
			core.ErrorField(err))
	}
}

func (bs *BuildSystem) HandleFileModified(path string) {
	bs.ctx.Logger.Info("JML file modified", core.StringField("path", path))

	// Get old dependencies before updating
	oldDeps := bs.depGraph.GetDependencies(path)

	// Re-parse and update document info
	docInfo, err := bs.discoverer.CreateDocumentInfo(path, bs.ctx.Paths.ProjectRoot)
	if err != nil {
		bs.ctx.Logger.Error("Failed to update document info for modified file",
			core.StringField("path", path),
			core.ErrorField(err))
		return
	}

	newDeps, err := bs.extractDependencies(docInfo)
	if err != nil {
		bs.ctx.Logger.Error("Failed to extract new dependencies",
			core.StringField("path", path),
			core.ErrorField(err))
		newDeps = []string{}
	}

	bs.updateDependenciesInGraph(path, oldDeps, newDeps)

	bs.AddDocument(docInfo)

	if err := bs.CompileDocument(path); err != nil {
		bs.ctx.Logger.Error("Failed to recompile modified file",
			core.StringField("path", path),
			core.ErrorField(err))
	}

	if err := bs.RecompileDependents(path); err != nil {
		bs.ctx.Logger.Error("Failed to recompile dependents",
			core.StringField("path", path),
			core.ErrorField(err))
	}
}

func (bs *BuildSystem) HandleFileDeleted(path string) {
	bs.ctx.Logger.Info("JML file deleted", core.StringField("path", path))

	// Check if we know about this file
	if _, exists := bs.GetDocumentInfo(path); !exists {
		bs.ctx.Logger.Debug("Deleted file not in build system, ignoring",
			core.StringField("path", path))
		return
	}

	// Remove from the build system
	bs.RemoveDocument(path)

	// TODO: Update dependencies in other documents that might reference this file
	// TODO: Recompile dependent documents if necessary

	bs.ctx.Logger.Info("Successfully removed deleted file from build system",
		core.StringField("path", path))
}

// HandleFileRenamed handles a file being renamed
func (bs *BuildSystem) HandleFileRenamed(path string) {
	bs.ctx.Logger.Info("JML file renamed", core.StringField("path", path))

	// For fsnotify, rename events are a bit tricky
	// handle this as a potential delete followed by a create
	// The actual implementation depends on whether the file still exists

	if _, err := os.Stat(path); err == nil {
		// File exists, treat as created/modified
		bs.HandleFileModified(path)
	} else {
		// File doesn't exist, treat as deleted
		bs.HandleFileDeleted(path)
	}
}

func (bs *BuildSystem) GetDocumentInfo(path string) (*DocumentInfo, bool) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	doc, exists := bs.docs[path]
	return doc, exists
}

func (bs *BuildSystem) AddDocument(doc *DocumentInfo) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bs.docs[doc.AbsPath] = doc

	// Add to type-specific maps
	switch doc.Type {
	case DocumentTypePage:
		bs.pages[doc.AbsPath] = &PageInfo{DocumentInfo: *doc}
	case DocumentTypeComponent:
		bs.comps[doc.AbsPath] = &ComponentInfo{DocumentInfo: *doc}
	}

	if err := bs.depGraph.AddNode(doc.AbsPath, doc.Type); err != nil {
		bs.ctx.Logger.Error("Failed to add document to dependency graph",
			core.StringField("path", doc.AbsPath),
			core.ErrorField(err))
	}
}

// RemoveDocument removes a document from the build system
func (bs *BuildSystem) RemoveDocument(path string) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if doc, exists := bs.docs[path]; exists {
		// Remove from type-specific maps
		switch doc.Type {
		case DocumentTypePage:
			delete(bs.pages, path)
		case DocumentTypeComponent:
			delete(bs.comps, path)
		}

		// Remove from main document map
		delete(bs.docs, path)

		if err := bs.depGraph.RemoveNode(path); err != nil {
			bs.ctx.Logger.Error("Failed to remove document from dependency graph",
				core.StringField("path", path),
				core.ErrorField(err))
		}
	}
}

func (bs *BuildSystem) CompileDocument(path string) error {
	bs.mu.RLock()
	doc, exists := bs.docs[path]
	bs.mu.RUnlock()

	if !exists {
		return nil // Document doesn't exist, nothing to compile
	}

	// 1. Compile JML to TypeScript
	reporter := diagnostic.NewReporter()
	jmlCompiler := compiler.NewCompiler(bs.ctx)
	ast, err := jmlCompiler.Compile(doc.AbsPath, reporter)
	if err != nil {
		return fmt.Errorf("failed to compile JML file %s: %w", doc.AbsPath, err)
	}
	if reporter.HasErrors() {
		printer := diagnostic.NewPrinter()
		printer.Print(reporter)
		return fmt.Errorf("compilation of %s failed with errors", doc.AbsPath)
	}

	// 2. Emit TypeScript from the AST to the .jawt/src/user directory
	emitter := emitter.NewEmitter(bs.ctx)
	if err := emitter.Emit(ast); err != nil {
		return fmt.Errorf("failed to emit TypeScript for %s: %w", doc.AbsPath, err)
	}

	// 3. Run external compilers
	if err := bs.compiler.RunTSC(); err != nil {
		return fmt.Errorf("failed to run tsc: %w", err)
	}

	if err := bs.compiler.RunTailwind(); err != nil {
		return fmt.Errorf("failed to run tailwind: %w", err)
	}

	bs.mu.Lock()
	doc.IsCompiled = true
	bs.mu.Unlock()

	return nil
}

// RecompileDependents recompiles all documents that depend on the given document
func (bs *BuildSystem) RecompileDependents(path string) error {
	dependents := bs.depGraph.GetDependents(path)

	bs.ctx.Logger.Info("Recompiling dependents",
		core.StringField("changed_file", path),
		core.IntField("dependent_count", len(dependents)))

	// Get compilation order for just the dependents
	compilationOrder, err := bs.getCompilationOrderForNodes(dependents)
	if err != nil {
		return fmt.Errorf("failed to get compilation order for dependents: %w", err)
	}

	// Recompile in dependency order
	for _, depPath := range compilationOrder {
		if err := bs.CompileDocument(depPath); err != nil {
			return fmt.Errorf("failed to recompile dependent %s: %w", depPath, err)
		}
	}

	return nil
}

func (bs *BuildSystem) getCompilationOrderForNodes(nodes []string) ([]string, error) {
	// This would create a subgraph with just the specified nodes
	// and their dependencies, then return topological order

	// Placeholder implementation
	return nodes, nil
}

func (bs *BuildSystem) updateDependenciesInGraph(path string, oldDeps, newDeps []string) {
	// Remove old dependencies that are no longer present
	for _, oldDep := range oldDeps {
		found := false
		for _, newDep := range newDeps {
			if oldDep == newDep {
				found = true
				break
			}
		}
		if !found {
			bs.depGraph.RemoveDependency(path, oldDep)
		}
	}

	// Add new dependencies
	for _, newDep := range newDeps {
		found := false
		for _, oldDep := range oldDeps {
			if newDep == oldDep {
				found = true
				break
			}
		}
		if !found {
			bs.depGraph.AddDependency(path, newDep)
		}
	}
}

func (bs *BuildSystem) isJMLFile(filePath string) bool {
	if !strings.HasSuffix(strings.ToLower(filePath), ".jml") {
		return false
	}

	if info, err := os.Stat(filePath); err != nil || info.IsDir() {
		return false
	}

	return true
}

// getDocumentTypeString returns a string representation of the document type
func (bs *BuildSystem) getDocumentTypeString(docType DocumentType) string {
	switch docType {
	case DocumentTypePage:
		return "page"
	case DocumentTypeComponent:
		return "component"
	default:
		return "unknown"
	}
}

// syncWorkspaceSources copies user code and extracts internal code into the .jawt/src directory
func (bs *BuildSystem) syncWorkspaceSources() error {
	bs.ctx.Logger.Info("Synchronizing workspace sources")

	// Clean the user and internal source directories
	if err := os.RemoveAll(bs.ctx.Paths.UserSrcDir); err != nil {
		return fmt.Errorf("failed to clean user source directory: %w", err)
	}
	if err := os.RemoveAll(bs.ctx.Paths.InternalSrcDir); err != nil {
		return fmt.Errorf("failed to clean internal source directory: %w", err)
	}
	if err := bs.ctx.Paths.EnsureDirectories(); err != nil {
		return fmt.Errorf("failed to ensure workspace directories: %w", err)
	}

	// Copy user scripts
	if err := bs.copyUserScripts(); err != nil {
		if !os.IsNotExist(err) { // It's okay if the scripts directory doesn't exist
			return fmt.Errorf("failed to copy user scripts: %w", err)
		}
	}

	// Extract internal Jawt components/scripts
	if err := bs.extractInternalScripts(); err != nil {
		return fmt.Errorf("failed to extract internal scripts: %w", err)
	}

	return nil
}

func (bs *BuildSystem) copyUserScripts() error {
	return filepath.Walk(bs.ctx.Paths.ScriptsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(bs.ctx.Paths.ScriptsDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(bs.ctx.Paths.UserSrcDir, relPath)

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, content, 0644)
	})
}

func (bs *BuildSystem) extractInternalScripts() error {
	internalScripts := map[string]string{
		"browser.ts": `// Placeholder for Jawt's internal browser API
export function setTitle(title: string) { console.log("Setting title: " + title); }
export function scrollToTop() { console.log('Scrolling to top'); }
`,
		"store.ts":   `// Placeholder for Jawt's internal store API
export function get(key: string) { console.log("Getting key: " + key); return null; }
export function set(key: string, value: any) { console.log("Setting key: " + key + " with value: " + value); }
`,
	}

	for filename, content := range internalScripts {
		destPath := filepath.Join(bs.ctx.Paths.InternalSrcDir, filename)
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write internal script %s: %w", filename, err)
		}
	}
	return nil
}

// generateWorkspaceConfigs creates the necessary config files in the .jawt directory
func (bs *BuildSystem) generateWorkspaceConfigs() error {
	bs.ctx.Logger.Info("Generating workspace configurations")

	// Create tsconfig.json
	tsconfigContent := `{
	  "compilerOptions": {
	    "target": "ESNext",
	    "module": "ESNext",
	    "moduleResolution": "node",
	    "strict": true,
	    "jsx": "preserve",
	    "importHelpers": true,
	    "experimentalDecorators": true,
	    "esModuleInterop": true,
	    "allowSyntheticDefaultImports": true,
	    "sourceMap": true,
	    "baseUrl": ".",
	    "paths": {
	      "@/*": ["src/user/*"],
	      "@jawt/*": ["src/internal/*"]
	    },
	    "lib": ["ESNext", "DOM"],
	    "outDir": "../build",
	    "rootDir": "src"
	  },
	  "include": ["src/**/*.ts", "src/**/*.tsx"],
	  "exclude": ["node_modules"]
	}`
	if err := os.WriteFile(bs.ctx.Paths.TSConfigPath, []byte(tsconfigContent), 0644); err != nil {
		return fmt.Errorf("failed to write tsconfig.json: %w", err)
	}

	// Create tailwind.config.js
	tailwindConfigContent := `/** @type {import('tailwindcss').Config} */
	module.exports = {
	  content: ["./src/user/**/*.ts", "./src/internal/**/*.ts"],
	  theme: {
	    extend: {},
	  },
	  plugins: [],
	}`
	if err := os.WriteFile(bs.ctx.Paths.TailwindConfigPath, []byte(tailwindConfigContent), 0644); err != nil {
		return fmt.Errorf("failed to write tailwind.config.js: %w", err)
	}

	return nil
}
