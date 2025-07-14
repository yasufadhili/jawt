package build

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/yasufadhili/jawt/internal/compiler"
	"github.com/yasufadhili/jawt/internal/core"
	"github.com/yasufadhili/jawt/internal/diagnostic"
	"os"
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
	compiler   *compiler.Compiler
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
		depGraph:   NewDependencyGraph(),
	}
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

	reporter := diagnostic.NewReporter()
	_, err := bs.compiler.Compile(doc.AbsPath, reporter)

	if reporter.HasErrors() {
		printer := diagnostic.NewPrinter()
		printer.Print(reporter)
		return fmt.Errorf("failed to compile %s due to errors", doc.AbsPath)
	}

	if err != nil {
		return fmt.Errorf("failed to compile %s: %w", doc.AbsPath, err)
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
