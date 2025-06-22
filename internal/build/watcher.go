package build

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/project"
	"time"
)

// FileWatcher monitors file changes and triggers recompilation
type FileWatcher struct {
	project  *project.Structure
	compiler *CompilerManager
	stopChan chan bool
}

// NewFileWatcher creates a new file watcher
func NewFileWatcher(p *project.Structure, compiler *CompilerManager) *FileWatcher {
	return &FileWatcher{
		project:  p,
		compiler: compiler,
		stopChan: make(chan bool),
	}
}

// Start begins watching for file changes
func (fw *FileWatcher) Start() error {

	// TODO: use fsnotify library
	go fw.watchLoop()

	return nil
}

// Stop stops the file watcher
func (fw *FileWatcher) Stop() {
	fmt.Println("🛑 Stopping file watcher...")
	fw.stopChan <- true
}

// watchLoop is the main watching loop (placeholder)
func (fw *FileWatcher) watchLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check for changes and recompile if necessary
			if err := fw.compiler.CompileChanged(); err != nil {
				fmt.Printf("❌ Recompilation failed: %v\n", err)
			}
		case <-fw.stopChan:
			return
		}
	}
}
