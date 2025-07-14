package process

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yasufadhili/jawt/internal/core"
)

// Manager manages all external processes
type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger core.Logger

	mu        sync.RWMutex
	processes map[string]*ManagedProcess

	wg          sync.WaitGroup
	jawtContext *core.JawtContext
}

func NewProcessManager(ctx context.Context, logger core.Logger, jawtContext *core.JawtContext) *Manager {
	managerCtx, cancel := context.WithCancel(ctx)

	return &Manager{
		ctx:         managerCtx,
		cancel:      cancel,
		logger:      logger,
		processes:   make(map[string]*ManagedProcess),
		jawtContext: jawtContext,
	}
}

// Start starts the process manager
func (pm *Manager) Start() error {
	pm.logger.Info("Starting process manager")
	return nil
}

// Stop stops all managed processes
func (pm *Manager) Stop() error {
	pm.logger.Info("Stopping process manager")
	pm.cancel()

	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Stop all processes
	for _, process := range pm.processes {
		if err := process.Stop(); err != nil {
			pm.logger.Error("Failed to stop process",
				core.StringField("process", process.Name()),
				core.ErrorField(err))
		}
	}

	pm.wg.Wait()
	return nil
}

// StartProcess starts a new managed process
func (pm *Manager) StartProcess(name string, options Options) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Check if a process already exists
	if _, exists := pm.processes[name]; exists {
		return fmt.Errorf("process %s already exists", name)
	}

	// Create a new managed process
	process := NewManagedProcess(name, options, pm.ctx, pm.logger)
	pm.processes[name] = process

	// Start the process
	pm.wg.Add(1)
	go func() {
		defer pm.wg.Done()
		if err := process.Start(); err != nil {
			pm.logger.Error("Failed to start process",
				core.StringField("process", name),
				core.ErrorField(err))
		}
	}()

	pm.logger.Info("Started process", core.StringField("process", name))
	return nil
}

// StopProcess stops a managed process
func (pm *Manager) StopProcess(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	process, exists := pm.processes[name]
	if !exists {
		return fmt.Errorf("process %s not found", name)
	}

	if err := process.Stop(); err != nil {
		return fmt.Errorf("failed to stop process %s: %w", name, err)
	}

	delete(pm.processes, name)
	pm.logger.Info("Stopped process", core.StringField("process", name))
	return nil
}

// RestartProcess restarts a managed process
func (pm *Manager) RestartProcess(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	process, exists := pm.processes[name]
	if !exists {
		return fmt.Errorf("process %s not found", name)
	}

	if err := process.Restart(); err != nil {
		return fmt.Errorf("failed to restart process %s: %w", name, err)
	}

	pm.logger.Info("Restarted process", core.StringField("process", name))
	return nil
}

// GetProcess returns a managed process by name
func (pm *Manager) GetProcess(name string) (*ManagedProcess, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	process, exists := pm.processes[name]
	return process, exists
}

// GetProcesses returns all managed processes
func (pm *Manager) GetProcesses() map[string]*ManagedProcess {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make(map[string]*ManagedProcess)
	for name, process := range pm.processes {
		result[name] = process
	}
	return result
}

// IsRunning checks if a process is running
func (pm *Manager) IsRunning(name string) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	process, exists := pm.processes[name]
	if !exists {
		return false
	}

	return process.IsRunning()
}

// GetProcessStats returns statistics for all processes
func (pm *Manager) GetProcessStats() map[string]Stats {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	stats := make(map[string]Stats)
	for name, process := range pm.processes {
		stats[name] = process.GetStats()
	}
	return stats
}

// StartTypeScriptWatch starts TypeScript compiler in watch mode
func (pm *Manager) StartTypeScriptWatch(ctx *core.JawtContext) error {
	tsPath, err := core.ResolveExecutablePath(ctx.JawtConfig.TypeScriptPath, pm.jawtContext)
	if err != nil {
		return fmt.Errorf("failed to resolve TypeScript executable: %w", err)
	}

	options := Options{
		Command:          tsPath,
		Args:             []string{"--watch", "--project", ctx.Paths.TSConfigPath},
		WorkingDir:       ctx.Paths.ProjectRoot,
		Env:              nil,
		RestartOnFailure: true,
		RestartDelay:     2 * time.Second,
		MaxRestarts:      3,
		OutputHandler: func(output string) {
			// pm.eventBus.Publish(events.CreateProcessOutputEvent("process_manager", "typescript", output))
		},
		ErrorHandler: func(err error) {
			// pm.eventBus.Publish(events.CreateProcessErrorEvent("process_manager", "typescript", err))
		},
	}

	return pm.StartProcess("typescript", options)
}

// StartTailwindWatch starts Tailwind CSS compiler in watch mode
func (pm *Manager) StartTailwindWatch(ctx *core.JawtContext) error {
	tailwindPath, err := core.ResolveExecutablePath(ctx.JawtConfig.TailwindPath, pm.jawtContext)
	if err != nil {
		return fmt.Errorf("failed to resolve Tailwind CSS executable: %w", err)
	}

	options := Options{
		Command:          tailwindPath,
		Args:             []string{"--watch", "--config", ctx.Paths.TailwindConfigPath},
		WorkingDir:       ctx.Paths.ProjectRoot,
		Env:              nil,
		RestartOnFailure: true,
		RestartDelay:     2 * time.Second,
		MaxRestarts:      10,
		OutputHandler: func(output string) {
			// pm.eventBus.Publish(events.CreateProcessOutputEvent("process_manager", "tailwind", output))
		},
		ErrorHandler: func(err error) {
			// pm.eventBus.Publish(events.CreateProcessErrorEvent("process_manager", "tailwind", err))
		},
	}

	return pm.StartProcess("tailwind", options)
}

// StartNodeProcess starts a Node.js process
func (pm *Manager) StartNodeProcess(name string, args []string, workingDir string, outputHandler func(string), errorHandler func(error)) error {
	nodePath, err := core.ResolveExecutablePath(pm.jawtContext.JawtConfig.NodePath, pm.jawtContext)
	if err != nil {
		return fmt.Errorf("failed to resolve Node.js executable: %w", err)
	}

	options := Options{
		Command:          nodePath,
		Args:             args,
		WorkingDir:       workingDir,
		Env:              nil,
		RestartOnFailure: true,
		RestartDelay:     5 * time.Second,
		MaxRestarts:      5,
		OutputHandler:    outputHandler,
		ErrorHandler:     errorHandler,
	}

	return pm.StartProcess(name, options)
}

// StartDevServer starts the development server
func (pm *Manager) StartDevServer(ctx *core.JawtContext, serverBinary string) error {
	options := Options{
		Command:          serverBinary,
		Args:             []string{"--port", fmt.Sprintf("%d", ctx.ProjectConfig.Server.Port)},
		WorkingDir:       ctx.Paths.ProjectRoot,
		Env:              nil,
		RestartOnFailure: true,
		RestartDelay:     5 * time.Second,
		MaxRestarts:      5,
		OutputHandler: func(output string) {
			// pm.eventBus.Publish(events.CreateProcessOutputEvent("process_manager", "devserver", output))
		},
		ErrorHandler: func(err error) {
			// pm.eventBus.Publish(events.CreateProcessErrorEvent("process_manager", "devserver", err))
		},
	}

	return pm.StartProcess("devserver", options)
}

// Health checks the health of all processes
func (pm *Manager) Health() map[string]bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	health := make(map[string]bool)
	for name, process := range pm.processes {
		health[name] = process.IsRunning()
	}
	return health
}
