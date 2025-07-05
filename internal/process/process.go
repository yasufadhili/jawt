package process

import (
	"bufio"
	"context"
	"fmt"
	"github.com/yasufadhili/jawt/internal/core"
	"io"
	"os/exec"
	"sync"
	"time"
)

// ProcessStatus represents the status of a managed process
type ProcessStatus int

const (
	StatusStopped ProcessStatus = iota
	StatusStarting
	StatusRunning
	StatusStopping
	StatusRestarting
	StatusFailed
)

// String returns the string representation of the process status
func (s ProcessStatus) String() string {
	switch s {
	case StatusStopped:
		return "stopped"
	case StatusStarting:
		return "starting"
	case StatusRunning:
		return "running"
	case StatusStopping:
		return "stopping"
	case StatusRestarting:
		return "restarting"
	case StatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// ProcessStats contains statistics about a managed process
type ProcessStats struct {
	Name         string        `json:"name"`
	Status       ProcessStatus `json:"status"`
	PID          int           `json:"pid"`
	StartTime    time.Time     `json:"start_time"`
	Uptime       time.Duration `json:"uptime"`
	RestartCount int           `json:"restart_count"`
	LastError    string        `json:"last_error,omitempty"`
}

// ManagedProcess represents a managed external process
type ManagedProcess struct {
	name    string
	options ProcessOptions
	ctx     context.Context
	cancel  context.CancelFunc
	logger  core.Logger

	mu           sync.RWMutex
	cmd          *exec.Cmd
	status       ProcessStatus
	pid          int
	startTime    time.Time
	restartCount int
	lastError    error

	stopChan chan struct{}
	doneChan chan struct{}
}

// NewManagedProcess creates a new managed process
func NewManagedProcess(name string, options ProcessOptions, ctx context.Context, logger core.Logger) *ManagedProcess {
	processCtx, cancel := context.WithCancel(ctx)

	return &ManagedProcess{
		name:     name,
		options:  options,
		ctx:      processCtx,
		cancel:   cancel,
		logger:   logger,
		status:   StatusStopped,
		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}
}

// Start starts the managed process
func (mp *ManagedProcess) Start() error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.status != StatusStopped {
		return fmt.Errorf("process %s is already running", mp.name)
	}

	mp.status = StatusStarting
	mp.startTime = time.Now()

	go mp.run()

	return nil
}

// Stop stops the managed process
func (mp *ManagedProcess) Stop() error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.status == StatusStopped {
		return nil
	}

	mp.status = StatusStopping
	mp.cancel()

	// Signal stop and wait for completion
	close(mp.stopChan)
	<-mp.doneChan

	mp.status = StatusStopped
	return nil
}

// Restart restarts the managed process
func (mp *ManagedProcess) Restart() error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.status == StatusStopped {
		return mp.Start()
	}

	mp.status = StatusRestarting
	mp.restartCount++

	// Kill current process
	if mp.cmd != nil && mp.cmd.Process != nil {
		if err := mp.cmd.Process.Kill(); err != nil {
			mp.logger.Warn("Failed to kill process during restart",
				core.StringField("process", mp.name),
				core.ErrorField(err))
		}
	}

	// Start new process
	go mp.run()

	return nil
}

// IsRunning returns true if the process is running
func (mp *ManagedProcess) IsRunning() bool {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return mp.status == StatusRunning
}

// Name returns the process name
func (mp *ManagedProcess) Name() string {
	return mp.name
}

// GetStats returns process statistics
func (mp *ManagedProcess) GetStats() ProcessStats {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	stats := ProcessStats{
		Name:         mp.name,
		Status:       mp.status,
		PID:          mp.pid,
		StartTime:    mp.startTime,
		RestartCount: mp.restartCount,
	}

	if mp.status == StatusRunning {
		stats.Uptime = time.Since(mp.startTime)
	}

	if mp.lastError != nil {
		stats.LastError = mp.lastError.Error()
	}

	return stats
}

// run is the main process loop
func (mp *ManagedProcess) run() {
	defer close(mp.doneChan)

	for {
		select {
		case <-mp.ctx.Done():
			return
		case <-mp.stopChan:
			return
		default:
			if err := mp.runOnce(); err != nil {
				mp.mu.Lock()
				mp.lastError = err
				mp.status = StatusFailed
				mp.mu.Unlock()

				mp.logger.Error("Process failed",
					core.StringField("process", mp.name),
					core.ErrorField(err))

				// mp.eventBus.Publish(events.ProcessErrorEvent("managed_process", mp.name, err))

				// Check if we should restart
				if mp.options.RestartOnFailure && mp.restartCount < mp.options.MaxRestarts {
					mp.logger.Info("Restarting process",
						core.StringField("process", mp.name),
						core.IntField("restart_count", mp.restartCount+1))

					time.Sleep(mp.options.RestartDelay)
					mp.mu.Lock()
					mp.restartCount++
					mp.status = StatusRestarting
					mp.mu.Unlock()
					continue
				}

				return
			}
		}
	}
}

// runOnce runs the process once
func (mp *ManagedProcess) runOnce() error {
	mp.mu.Lock()
	mp.cmd = exec.CommandContext(mp.ctx, mp.options.Command, mp.options.Args...)
	mp.cmd.Dir = mp.options.WorkingDir
	mp.cmd.Env = mp.options.Env
	mp.mu.Unlock()

	// Set up pipes for stdout and stderr
	stdout, err := mp.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := mp.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the process
	if err := mp.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}

	mp.mu.Lock()
	mp.status = StatusRunning
	mp.pid = mp.cmd.Process.Pid
	mp.mu.Unlock()

	mp.logger.Info("Process started",
		core.StringField("process", mp.name),
		core.IntField("pid", mp.pid))

	// mp.eventBus.Publish(events.CreateProcessStopEvent("managed_process", mp.name, mp.pid))

	// Handle output
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		mp.handleOutput(stdout, false)
	}()

	go func() {
		defer wg.Done()
		mp.handleOutput(stderr, true)
	}()

	// Wait for process to complete
	err = mp.cmd.Wait()

	wg.Wait()

	mp.mu.Lock()
	mp.status = StatusStopped
	mp.mu.Unlock()

	mp.logger.Info("Process stopped",
		core.StringField("process", mp.name),
		core.IntField("pid", mp.pid))

	// mp.eventBus.Publish(events.CreateProcessStopEvent("managed_process", mp.name, mp.pid))

	return err
}

// handleOutput handles process output
func (mp *ManagedProcess) handleOutput(reader io.Reader, isError bool) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		// Log the output
		if isError {
			mp.logger.Error("Process error output",
				core.StringField("process", mp.name),
				core.StringField("output", line))
		} else {
			mp.logger.Debug("Process output",
				core.StringField("process", mp.name),
				core.StringField("output", line))
		}

		// Call output handler if provided
		if mp.options.OutputHandler != nil {
			mp.options.OutputHandler(line)
		}
	}

	if err := scanner.Err(); err != nil {
		mp.logger.Error("Error reading process output",
			core.StringField("process", mp.name),
			core.ErrorField(err))
	}
}
