package process

import (
	"time"
)

// ProcessOptions contains configuration for a managed process
type ProcessOptions struct {
	// Command and arguments
	Command string   `json:"command"`
	Args    []string `json:"args"`

	// Working directory and environment
	WorkingDir string   `json:"working_dir"`
	Env        []string `json:"env"`

	// Restart configuration
	RestartOnFailure bool          `json:"restart_on_failure"`
	RestartDelay     time.Duration `json:"restart_delay"`
	MaxRestarts      int           `json:"max_restarts"`

	// Handlers
	OutputHandler func(string) `json:"-"`
	ErrorHandler  func(error)  `json:"-"`
}

// DefaultProcessOptions returns default process options
func DefaultProcessOptions() ProcessOptions {
	return ProcessOptions{
		RestartOnFailure: false,
		RestartDelay:     5 * time.Second,
		MaxRestarts:      3,
	}
}

// WithCommand sets the command and arguments
func (po ProcessOptions) WithCommand(command string, args ...string) ProcessOptions {
	po.Command = command
	po.Args = args
	return po
}

// WithWorkingDir sets the working directory
func (po ProcessOptions) WithWorkingDir(dir string) ProcessOptions {
	po.WorkingDir = dir
	return po
}

// WithEnv sets environment variables
func (po ProcessOptions) WithEnv(env []string) ProcessOptions {
	po.Env = env
	return po
}

// WithRestart enables restart on failure
func (po ProcessOptions) WithRestart(enabled bool, delay time.Duration, maxRestarts int) ProcessOptions {
	po.RestartOnFailure = enabled
	po.RestartDelay = delay
	po.MaxRestarts = maxRestarts
	return po
}

// WithOutputHandler sets the output handler
func (po ProcessOptions) WithOutputHandler(handler func(string)) ProcessOptions {
	po.OutputHandler = handler
	return po
}

// WithErrorHandler sets the error handler
func (po ProcessOptions) WithErrorHandler(handler func(error)) ProcessOptions {
	po.ErrorHandler = handler
	return po
}

// TypeScriptWatchOptions returns options for TypeScript watch mode
func TypeScriptWatchOptions(tscPath, configPath, workingDir string) ProcessOptions {
	return DefaultProcessOptions().
		WithCommand(tscPath, "--watch", "--project", configPath).
		WithWorkingDir(workingDir).
		WithRestart(true, 2*time.Second, 10)
}

// TailwindWatchOptions returns options for Tailwind watch mode
func TailwindWatchOptions(tailwindPath, configPath, workingDir string) ProcessOptions {
	return DefaultProcessOptions().
		WithCommand(tailwindPath, "--watch", "--config", configPath).
		WithWorkingDir(workingDir).
		WithRestart(true, 2*time.Second, 10)
}

// DevServerOptions returns options for the development server
func DevServerOptions(serverBinary, workingDir string, port int) ProcessOptions {
	return DefaultProcessOptions().
		WithCommand(serverBinary, "--port", string(rune(port))).
		WithWorkingDir(workingDir).
		WithRestart(true, 5*time.Second, 5)
}

// NodeScriptOptions returns options for running Node.js scripts
func NodeScriptOptions(nodePath, scriptPath, workingDir string, args ...string) ProcessOptions {
	allArgs := append([]string{scriptPath}, args...)
	return DefaultProcessOptions().
		WithCommand(nodePath, allArgs...).
		WithWorkingDir(workingDir).
		WithRestart(false, 0, 0)
}
