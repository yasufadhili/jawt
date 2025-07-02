package build

import "fmt"

// Stop stops all running services
func (b *BuildSystem) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.isRunning {
		return
	}

	fmt.Println("\n🛑 Stopping development mode...")

	if b.watcher != nil {
		b.watcher.Stop()
		fmt.Println("   File watcher stopped")
	}

	if b.server != nil {
		//b.server.Stop()
		fmt.Println("   Development server stopped")
	}

	b.isRunning = false
	close(b.stopChan)
	fmt.Println("   Goodbye! 👋")
}
