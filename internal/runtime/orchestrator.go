package runtime

import (
	"context"
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/yasufadhili/jawt/internal/core"
	"github.com/yasufadhili/jawt/internal/process"
)

// Orchestrator manages the entire runtime environment
type Orchestrator struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger core.Logger

	processManager *process.ProcessManager
	fileWatcher    *FileWatcher
	devServer      *DevServer
	jawtContext    *core.JawtContext

	wg sync.WaitGroup
}

// NewOrchestrator creates a new Orchestrator
func NewOrchestrator(ctx context.Context, logger core.Logger, jawtCtx *core.JawtContext) (*Orchestrator, error) {
	orchCtx, cancel := context.WithCancel(ctx)

	pm := process.NewProcessManager(orchCtx, logger, jawtCtx)

	fw, err := NewFileWatcher(orchCtx, logger)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	return &Orchestrator{
		ctx:            orchCtx,
		cancel:         cancel,
		logger:         logger,
		processManager: pm,
		fileWatcher:    fw,
		devServer:      NewDevServer(orchCtx, logger),
		jawtContext:    jawtCtx,
	}, nil
}

// StartAll starts all managed processes and watchers
func (o *Orchestrator) StartAll() error {
	o.logger.Info("Starting orchestrator")

	// Start the process manager
	if err := o.processManager.Start(); err != nil {
		return fmt.Errorf("failed to start process manager: %w", err)
	}

	// Start the file watcher
	if err := o.fileWatcher.Start(); err != nil {
		return fmt.Errorf("failed to start file watcher: %w", err)
	}

	// Start TypeScript watch
	if err := o.processManager.StartTypeScriptWatch(o.jawtContext); err != nil {
		o.logger.Error("Failed to start TypeScript watch", core.ErrorField(err))
	}

	// Start TSServer
	if err := o.startTSServer(); err != nil {
		o.logger.Error("Failed to start TSServer", core.ErrorField(err))
	}

	// Configure and start JML watcher
	o.fileWatcher.SetWatchPatterns([]string{".jml"})
	o.fileWatcher.OnEvent(o.handleJmlFileEvent)
	if err := o.fileWatcher.AddPathsRecursive([]string{o.jawtContext.Paths.ProjectRoot}); err != nil {
		return fmt.Errorf("failed to add paths to JML watcher: %w", err)
	}

	// Start dev server
	go func() {
		if err := o.devServer.Start(o.jawtContext.ProjectConfig.GetDevServerAddress()); err != nil {
			o.logger.Error("Failed to start dev server", core.ErrorField(err))
		}
	}()

	o.logger.Info("Orchestrator started")
	return nil
}

// StopAll stops all managed processes and watchers
func (o *Orchestrator) StopAll() error {
	o.logger.Info("Stopping orchestrator")
	o.cancel()

	if err := o.processManager.Stop(); err != nil {
		o.logger.Error("Failed to stop process manager", core.ErrorField(err))
	}

	if err := o.fileWatcher.Stop(); err != nil {
		o.logger.Error("Failed to stop file watcher", core.ErrorField(err))
	}

	o.wg.Wait()
	o.logger.Info("Orchestrator stopped")
	return nil
}

// RestartProcess restarts a managed process
func (o *Orchestrator) RestartProcess(name string) error {
	return o.processManager.RestartProcess(name)
}

// startTSServer starts the TypeScript server
func (o *Orchestrator) startTSServer() error {
	return o.processManager.StartNodeProcess(
		"tsserver",
		[]string{"node_modules/typescript/bin/tsserver"},
		o.jawtContext.Paths.ProjectRoot,
		func(output string) {
			// Handle tsserver output (JSON protocol)
			o.logger.Debug("TSServer output", core.StringField("output", output))
		},
		func(err error) {
			o.logger.Error("TSServer error", core.ErrorField(err))
		},
	)
}

// handleJmlFileEvent handles file events for JML files
func (o *Orchestrator) handleJmlFileEvent(event fsnotify.Event) {
	o.logger.Info("JML file changed, triggering compiler", core.StringField("file", event.Name))
	// Here we shall trigger the JML compiler pipeline
	// For now, we'll just log the event and broadcast a reload message
	o.devServer.Broadcast([]byte("reload"))
}
