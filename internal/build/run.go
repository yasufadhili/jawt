package build

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yasufadhili/jawt/internal/core"
	"github.com/yasufadhili/jawt/internal/runtime"
)

// RunProject starts the JAWT project in development mode
func RunProject(ctx *core.JawtContext) error {
	ctx.Logger.Info("Starting JAWT development server",
		core.StringField("host", ctx.ProjectConfig.Server.Host),
		core.IntField("port", ctx.ProjectConfig.Server.Port))

	fileWatcher, err := runtime.NewFileWatcher(ctx.Context(), ctx.Logger)
	if err != nil {
		return fmt.Errorf("failed to initialize file watcher: %w", err)
	}

	watchPaths := ctx.Paths.GetWatchPaths()
	if err := fileWatcher.AddPathsRecursive(watchPaths); err != nil {
		return fmt.Errorf("failed to add watch paths: %w", err)
	}

	// Override defaults
	fileWatcher.SetWatchPatterns([]string{".jml", ".ts", ".tsx", ".js", ".jsx", ".json", ".css"})
	fileWatcher.SetIgnorePatterns([]string{
		".git/", "node_modules/", ".jawt/", "dist/", "build/",
		".DS_Store", "*.tmp", "*.swp", "*.swo",
	})

	// Register event handler
	fileWatcher.OnEvent(func(event fsnotify.Event) {
		ctx.Logger.Info("Detected file event",
			core.StringField("operation", event.Op.String()),
			core.StringField("file", event.Name))
		// TODO: Implement build/reload logic based on event type
		// - For .ts/.tsx files: Recompile TypeScript to ctx.Paths.TypeScriptOutputDir
		// - For .jml files: Reprocess JML to ctx.Paths.ComponentsOutputDir or ctx.Paths.DistDir
		// - For .css files: Reprocess Tailwind CSS to ctx.Paths.TailwindOutputDir
		// - Trigger browser reload if ctx.ProjectConfig.EnableHMR is true
		ctx.Logger.Debug("Placeholder: Would trigger build/reload for event",
			core.StringField("operation", event.Op.String()),
			core.StringField("file", event.Name))
	})

	if err := fileWatcher.Start(); err != nil {
		return fmt.Errorf("failed to start file watcher: %w", err)
	}

	// Set up HTTP server
	serverAddr := ctx.ProjectConfig.GetServerAddress()
	server := &http.Server{
		Addr:    serverAddr,
		Handler: http.FileServer(http.Dir(ctx.Paths.DistDir)),
	}

	// Start HTTP server in a goroutine
	go func() {
		ctx.Logger.Info("Starting HTTP server", core.StringField("address", serverAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ctx.Logger.Error("HTTP server error", core.ErrorField(err))
		}
	}()

	// Wait for context cancellation (e.g., Ctrl+C)
	<-ctx.Context().Done()

	// Shutdown sequence
	ctx.Logger.Info("Initiating shutdown")

	// Stop file watcher
	if err := fileWatcher.Stop(); err != nil {
		ctx.Logger.Error("Failed to stop file watcher", core.ErrorField(err))
	}

	// Shutdown HTTP server gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		ctx.Logger.Error("Failed to shutdown HTTP server", core.ErrorField(err))
	}

	ctx.Logger.Info("JAWT development server stopped")
	return nil
}
