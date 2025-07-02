package project

import (
	"context"

	"github.com/yasufadhili/jawt/internal/tsc/lsp/lsproto"
	"github.com/yasufadhili/jawt/internal/tsc/vfs"
)

type WatcherHandle string

type Client interface {
	WatchFiles(ctx context.Context, watchers []*lsproto.FileSystemWatcher) (WatcherHandle, error)
	UnwatchFiles(ctx context.Context, handle WatcherHandle) error
	RefreshDiagnostics(ctx context.Context) error
}

type ServiceHost interface {
	FS() vfs.FS
	DefaultLibraryPath() string
	TypingsLocation() string
	GetCurrentDirectory() string
	NewLine() string

	Client() Client
}
