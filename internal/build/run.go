package build

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/server"
)

func (b *BuildSystem) Run() error {

	b.mu.Lock()
	b.isRunning = true
	b.mu.Unlock()

	if err := b.Build(); err != nil {
		fmt.Println("   Initial build failed, watching for changes to retry...")
	} else {
		fmt.Println("   Initial build successful!")
	}

	b.watcher = NewFileWatcher(b.project)
	b.watcher.SetErrorHandler(b.handleWatcherError)
	b.watcher.SetChangeHandler(b.handleFileChange)

	if err := b.watcher.Start(); err != nil {
		return fmt.Errorf("failed to start file watcher: %w", err)
	}

	b.server = server.NewDevServer(b.project)
	if err := b.server.Start(); err != nil {
		return fmt.Errorf("failed to start development server: %w", err)
	}

	// Keep running until stopped
	<-b.stopChan

	return nil
}
