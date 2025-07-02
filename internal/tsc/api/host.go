package api

import "github.com/yasufadhili/jawt/internal/tsc/vfs"

type APIHost interface {
	FS() vfs.FS
	DefaultLibraryPath() string
	GetCurrentDirectory() string
	NewLine() string
}
