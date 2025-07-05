package build

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/yasufadhili/jawt/internal/core"
	"sync"
	"time"
)

// DocumentType represents the type of a document (page or component)
type DocumentType int

const (
	DocumentTypePage DocumentType = iota
	DocumentTypeComponent
)

// DocumentInfo represents common information for all document types
type DocumentInfo struct {
	Name         string
	RelPath      string
	AbsPath      string
	Type         DocumentType
	Dependencies []string
	DependedBy   []string
	IsCompiled   bool
	LastModified time.Time
	Hash         string // Content hash for detecting changes
}

// ComponentInfo represents information about a component
type ComponentInfo struct {
	DocumentInfo
	Props map[string]string
}

// PageInfo represents information about a page
type PageInfo struct {
	DocumentInfo
	Route string
}

// BuildSystem manages the build process for a JAWT project
type BuildSystem struct {
	ctx      *core.JawtContext
	mu       sync.RWMutex
	docs     map[string]*DocumentInfo  // Map of a document path to DocumentInfo
	pages    map[string]*PageInfo      // Map of page path to PageInfo
	comps    map[string]*ComponentInfo // Map of a component path to ComponentInfo
	compiler Compiler
	watcher  FileWatcher
}

// Compiler is an interface for compiling documents
type Compiler interface {
	CompileDocument(doc *DocumentInfo) error
	CompilePage(page *PageInfo) error
	CompileComponent(comp *ComponentInfo) error
}

// FileWatcher is an interface for watching files for changes
type FileWatcher interface {
	OnEvent(handler func(fsnotify.Event))
	Start() error
	Stop() error
}

// NewBuildSystem creates a new BuildSystem
func NewBuildSystem(ctx *core.JawtContext, compiler Compiler, watcher FileWatcher) *BuildSystem {
	return &BuildSystem{
		ctx:      ctx,
		docs:     make(map[string]*DocumentInfo),
		pages:    make(map[string]*PageInfo),
		comps:    make(map[string]*ComponentInfo),
		compiler: compiler,
		watcher:  watcher,
	}
}

// Initialise performs initial project discovery and compilation
func (bs *BuildSystem) Initialise() error {
	bs.ctx.Logger.Info("Initialising build system")

	if err := bs.DiscoverProject(); err != nil {
		return err
	}

	if err := bs.CompileAll(); err != nil {
		return err
	}

	bs.SetupWatcher()

	return nil
}

// DiscoverProject finds all JML documents in the project
func (bs *BuildSystem) DiscoverProject() error {
	bs.ctx.Logger.Info("Discovering project documents")

	// Find all JML files in the project
	jmlFiles, err := DiscoverProjectFiles(bs.ctx)
	if err != nil {
		return fmt.Errorf("failed to discover project files: %w", err)
	}

	for _, path := range jmlFiles {
		docInfo, err := CreateDocumentInfo(path, bs.ctx.Paths.ProjectRoot)
		if err != nil {
			bs.ctx.Logger.Warn("Failed to process document",
				core.StringField("path", path),
				core.ErrorField(err))
			continue
		}

		bs.AddDocument(docInfo)
	}

	if err := AnalyseDependencies(bs.docs); err != nil {
		return fmt.Errorf("failed to analyse dependencies: %w", err)
	}

	bs.ctx.Logger.Info("Project discovery completed",
		core.IntField("pages", len(bs.pages)),
		core.IntField("components", len(bs.comps)))

	return nil
}

// CompileAll compiles all documents in the project
func (bs *BuildSystem) CompileAll() error {
	bs.ctx.Logger.Info("Compiling all documents")

	// TODO: Implement full compilation
	// 1. Determine compilation order based on dependencies
	// 2. Compile each document
	// 3. Update compilation status

	return nil
}

// SetupWatcher configures the file watcher to handle document changes
func (bs *BuildSystem) SetupWatcher() {
	bs.ctx.Logger.Info("Setting up file watcher")

	bs.watcher.OnEvent(func(event fsnotify.Event) {
		bs.HandleFileEvent(event)
	})
}

// HandleFileEvent processes file system events
func (bs *BuildSystem) HandleFileEvent(event fsnotify.Event) {
	bs.ctx.Logger.Info("Handling file event",
		core.StringField("operation", event.Op.String()),
		core.StringField("file", event.Name))

	// TODO: Implement event handling
	// 1. Determine if the file is a JML document
	// 2. Determine the type of event (create, modify, delete)
	// 3. Call the appropriate handler method

	switch event.Op {
	case fsnotify.Create:
		bs.HandleFileCreated(event.Name)
	case fsnotify.Write:
		bs.HandleFileModified(event.Name)
	case fsnotify.Remove:
		bs.HandleFileDeleted(event.Name)
	case fsnotify.Rename:
		bs.HandleFileRenamed(event.Name)
	}
}

// HandleFileCreated handles a new file being created
func (bs *BuildSystem) HandleFileCreated(path string) {
	bs.ctx.Logger.Info("File created", core.StringField("path", path))

	// TODO: Implement file creation handling
	// 1. Parse the file to determine if it's a page or component
	// 2. Create DocumentInfo, PageInfo, or ComponentInfo
	// 3. Add to the appropriate maps
	// 4. Analyze dependencies
	// 5. Compile the document
	// 6. Recompile dependent documents if necessary
}

// HandleFileModified handles a file being modified
func (bs *BuildSystem) HandleFileModified(path string) {
	bs.ctx.Logger.Info("File modified", core.StringField("path", path))

	// TODO: Implement file modification handling
	// 1. Check if the file is in our document maps
	// 2. Reparse the file to update DocumentInfo
	// 3. Recompile the document
	// 4. Recompile dependent documents if necessary
}

// HandleFileDeleted handles a file being deleted
func (bs *BuildSystem) HandleFileDeleted(path string) {
	bs.ctx.Logger.Info("File deleted", core.StringField("path", path))

	// TODO: Implement file deletion handling
	// 1. Check if the file is in our document maps
	// 2. Remove from the appropriate maps
	// 3. Update dependencies in other documents
	// 4. Recompile dependent documents if necessary
}

// HandleFileRenamed handles a file being renamed
func (bs *BuildSystem) HandleFileRenamed(path string) {
	bs.ctx.Logger.Info("File renamed", core.StringField("path", path))

	// TODO: Implement file rename handling
	// 1. This is handled as a delete followed by a create
	// 2. May need special handling for updating dependencies
}

// GetDocumentInfo retrieves document info by path
func (bs *BuildSystem) GetDocumentInfo(path string) (*DocumentInfo, bool) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	doc, exists := bs.docs[path]
	return doc, exists
}

// AddDocument adds a document to the build system
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
	}
}

// CompileDocument compiles a single document
func (bs *BuildSystem) CompileDocument(path string) error {
	bs.mu.RLock()
	doc, exists := bs.docs[path]
	bs.mu.RUnlock()

	if !exists {
		return nil // Document doesn't exist, nothing to compile
	}

	if err := bs.compiler.CompileDocument(doc); err != nil {
		return err
	}

	bs.mu.Lock()
	doc.IsCompiled = true
	bs.mu.Unlock()

	return nil
}

// RecompileDependents recompiles all documents that depend on the given document
func (bs *BuildSystem) RecompileDependents(path string) error {
	bs.mu.RLock()
	doc, exists := bs.docs[path]
	bs.mu.RUnlock()

	if !exists {
		return nil // Document doesn't exist, nothing to do
	}

	// Recompile all documents that depend on this one
	for _, depPath := range doc.DependedBy {
		if err := bs.CompileDocument(depPath); err != nil {
			return err
		}
	}

	return nil
}
