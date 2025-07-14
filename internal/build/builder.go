package build

import (
	"github.com/fsnotify/fsnotify"
	"github.com/yasufadhili/jawt/internal/core"
	"sync"
	"time"
)

// DocumentType represents the type of document (page or component)
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
	Hash         string
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

type BuildSystem struct {
	ctx   *core.JawtContext
	mu    sync.RWMutex
	docs  map[string]*DocumentInfo
	pages map[string]*PageInfo
	comps map[string]*ComponentInfo
	// compiler *compiler.Compiler
	watcher  FileWatcher
	depGraph DependencyGraph
}

// FileWatcher is an interface for watching files for changes
type FileWatcher interface {
	OnEvent(handler func(fsnotify.Event))
	Start() error
	Stop() error
}

// NewBuildSystem creates a new BuildSystem
func NewBuildSystem(ctx *core.JawtContext, watcher FileWatcher) *BuildSystem {
	return &BuildSystem{
		ctx:      ctx,
		docs:     make(map[string]*DocumentInfo),
		pages:    make(map[string]*PageInfo),
		comps:    make(map[string]*ComponentInfo),
		watcher:  watcher,
		depGraph: NewDependencyGraph(),
	}
}
