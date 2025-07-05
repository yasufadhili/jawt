package build

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/yasufadhili/jawt/internal/core"
	"os"
	"strings"
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

	// Check if this is a JML file we care about
	if !bs.isJMLFile(event.Name) {
		bs.ctx.Logger.Debug("Ignoring non-JML file event",
			core.StringField("file", event.Name))
		return
	}

	// Handle different event types
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

// HandleFileCreated handles a new file being created
func (bs *BuildSystem) HandleFileCreated(path string) {
	bs.ctx.Logger.Info("JML file created", core.StringField("path", path))

	// Create DocumentInfo for the new file
	docInfo, err := CreateDocumentInfo(path, bs.ctx.Paths.ProjectRoot)
	if err != nil {
		bs.ctx.Logger.Error("Failed to create document info for new file",
			core.StringField("path", path),
			core.ErrorField(err))
		return
	}

	bs.AddDocument(docInfo)

	// Analyse dependencies for the new document
	docs := map[string]*DocumentInfo{path: docInfo}
	if err := AnalyseDependencies(docs); err != nil {
		bs.ctx.Logger.Error("Failed to analyse dependencies for new file",
			core.StringField("path", path),
			core.ErrorField(err))
	}

	if err := bs.CompileDocument(path); err != nil {
		bs.ctx.Logger.Error("Failed to compile new file",
			core.StringField("path", path),
			core.ErrorField(err))
	}

	bs.ctx.Logger.Info("Successfully processed new JML file",
		core.StringField("path", path),
		core.StringField("type", bs.getDocumentTypeString(docInfo.Type)))
}

// HandleFileModified handles a file being modified
func (bs *BuildSystem) HandleFileModified(path string) {
	bs.ctx.Logger.Info("JML file modified", core.StringField("path", path))

	// Check if we know about this file
	if _, exists := bs.GetDocumentInfo(path); !exists {
		bs.ctx.Logger.Info("Modified file not in build system, treating as new file",
			core.StringField("path", path))
		bs.HandleFileCreated(path)
		return
	}

	// Re-parse and update document info
	docInfo, err := CreateDocumentInfo(path, bs.ctx.Paths.ProjectRoot)
	if err != nil {
		bs.ctx.Logger.Error("Failed to update document info for modified file",
			core.StringField("path", path),
			core.ErrorField(err))
		return
	}

	// Update in build system
	bs.AddDocument(docInfo) // This will overwrite the existing entry

	// Recompile the document
	if err := bs.CompileDocument(path); err != nil {
		bs.ctx.Logger.Error("Failed to recompile modified file",
			core.StringField("path", path),
			core.ErrorField(err))
	}

	// Recompile dependent documents
	if err := bs.RecompileDependents(path); err != nil {
		bs.ctx.Logger.Error("Failed to recompile dependents",
			core.StringField("path", path),
			core.ErrorField(err))
	}
}

// HandleFileDeleted handles a file being deleted
func (bs *BuildSystem) HandleFileDeleted(path string) {
	bs.ctx.Logger.Info("JML file deleted", core.StringField("path", path))

	// Check if we know about this file
	if _, exists := bs.GetDocumentInfo(path); !exists {
		bs.ctx.Logger.Debug("Deleted file not in build system, ignoring",
			core.StringField("path", path))
		return
	}

	// Remove from build system
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
	// We'll handle this as a potential delete followed by a create
	// The actual implementation depends on whether the file still exists

	if _, err := os.Stat(path); err == nil {
		// File exists, treat as created/modified
		bs.HandleFileModified(path)
	} else {
		// File doesn't exist, treat as deleted
		bs.HandleFileDeleted(path)
	}
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

// isJMLFile checks if a file is a JML file we should process
func (bs *BuildSystem) isJMLFile(filePath string) bool {
	// Check if file has .jml extension
	if !strings.HasSuffix(strings.ToLower(filePath), ".jml") {
		return false
	}

	// Check if file exists and is not a directory
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
