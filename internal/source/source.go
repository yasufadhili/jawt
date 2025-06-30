package source

import (
	"github.com/yasufadhili/jawt/internal/project"
	"time"
)

type Manager interface {
	// GetFile retrieves file content with caching
	GetFile(path string) (*File, error)

	// GetFiles retrieves multiple files efficiently
	GetFiles(paths []string) ([]*File, error)

	// InvalidateCache removes cached entries for changed files
	InvalidateCache(paths []string) error

	// WatchFiles sets up file system watching
	WatchFiles(paths []string, callback FileChangeCallback) error

	// GetFileHash returns content hash for cache invalidation
	GetFileHash(path string) (string, error)

	// CreateVirtualFile creates an in-memory file for testing
	CreateVirtualFile(path string, content string) *File
}

type File struct {
	RelPath      string    `json:"rel_path"`
	AbsPath      string    `json:"abs_path"`
	Hash         string    `json:"hash"`
	ModTime      time.Time `json:"mod_time"`
	Dependencies []string  `json:"dependencies"`
	Language     Language  `json:"language"`
}

type Language int

const (
	LangJml Language = iota
	LangTypeScript
	LangCSS
	LangJSON
)

type FileChangeCallback func(path string, event project.ChangeType) error
