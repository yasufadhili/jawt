package build

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yasufadhili/jawt/internal/core"
	"github.com/yasufadhili/jawt/internal/runtime"
)

// RunProject starts the JAWT project in development mode
func RunProject(ctx *core.JawtContext) error {
	ctx.Logger.Info("Starting Jawt development server",
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

	// Create a dummy compiler for now (will be replaced with real compiler later)
	compiler := &dummyCompiler{ctx: ctx}

	buildSystem := NewBuildSystem(ctx, compiler, fileWatcher)

	// Initialise the build system (discover and compile)
	if err := buildSystem.Initialise(); err != nil {
		return fmt.Errorf("failed to initialise build system: %w", err)
	}

	// Start file watcher
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

// dummyCompiler is a placeholder implementation of the Compiler interface
type dummyCompiler struct {
	ctx *core.JawtContext
}

func (c *dummyCompiler) CompileDocument(doc *DocumentInfo) error {
	c.ctx.Logger.Debug("Dummy compiler: Would compile document",
		core.StringField("name", doc.Name),
		core.StringField("path", doc.AbsPath))
	return nil
}

func (c *dummyCompiler) CompilePage(page *PageInfo) error {
	c.ctx.Logger.Debug("Dummy compiler: Would compile page",
		core.StringField("name", page.Name),
		core.StringField("path", page.AbsPath))
	return nil
}

func (c *dummyCompiler) CompileComponent(comp *ComponentInfo) error {
	c.ctx.Logger.Debug("Dummy compiler: Would compile component",
		core.StringField("name", comp.Name),
		core.StringField("path", comp.AbsPath))
	return nil
}
