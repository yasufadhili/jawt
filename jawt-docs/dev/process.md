# Process Management (`internal/process`)

The `internal/process` package is responsible for managing external processes that are part of the JAWT development and build workflow. This includes tools like the TypeScript compiler (`tsc`), Tailwind CSS, and potentially other Node.js-based services. It provides robust mechanisms for starting, stopping, restarting, and monitoring these processes.

## Core Concepts

*   **ManagedProcess**: Represents a single external process that is controlled by the `ProcessManager`. It handles the lifecycle of the process, including starting, stopping, and automatic restarts.
*   **ProcessManager**: The central orchestrator for all `ManagedProcess` instances. It provides an API to interact with and monitor multiple external tools.
*   **ProcessOptions**: Defines the configuration for a `ManagedProcess`, such as the command to run, arguments, working directory, environment variables, and restart policies.
*   **ProcessStatus**: An enumeration representing the current state of a `ManagedProcess` (e.g., stopped, running, failed).
*   **ProcessStats**: Provides runtime statistics for a `ManagedProcess`, including PID, uptime, and restart count.

## Key Data Structures

### `ProcessStatus`

An enumeration for the different states a managed process can be in.

```go
type ProcessStatus int

const (
	StatusStopped ProcessStatus = iota
	StatusStarting
	StatusRunning
	StatusStopping
	StatusRestarting
	StatusFailed
)
```

### `ProcessStats`

Contains runtime statistics and information about a managed process.

```go
type ProcessStats struct {
	Name         string        `json:"name"`
	Status       ProcessStatus `json:"status"`
	PID          int           `json:"pid"`
	StartTime    time.Time     `json:"start_time"`
	Uptime       time.Duration `json:"uptime"`
	RestartCount int           `json:"restart_count"`
	LastError    string        `json:"last_error,omitempty"`
}
```

### `ProcessOptions`

Configures how a `ManagedProcess` should behave.

```go
type ProcessOptions struct {
	Command string
	Args    []string

	WorkingDir string
	Env        []string

	RestartOnFailure bool
	RestartDelay     time.Duration
	MaxRestarts      int

	OutputHandler func(string)
	ErrorHandler  func(error)
}
```

### `ManagedProcess`

Represents an individual external process managed by the `ProcessManager`.

```go
type ManagedProcess struct {
	name    string
	opts    ProcessOptions
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
```

### `ProcessManager`

The central component for managing multiple external processes.

```go
type ProcessManager struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger core.Logger

	mu        sync.RWMutex
	processes map[string]*ManagedProcess

	wg          sync.WaitGroup
	jawtContext *core.JawtContext
}
```

## Functions & Methods

### `NewProcessManager`

```go
func NewProcessManager(ctx context.Context, logger core.Logger, jawtContext *core.JawtContext) *ProcessManager
```

Creates a new `ProcessManager` instance.

### `(*ProcessManager) Start()`

Starts the process manager. This typically involves initializing internal structures.

### `(*ProcessManager) Stop()`

Stops all currently managed processes and gracefully shuts down the manager.

### `(*ProcessManager) StartProcess(name string, options ProcessOptions) error`

Starts a new external process with the given name and options. The process will be managed by the `ProcessManager`.

### `(*ProcessManager) StopProcess(name string) error`

Stops a specific managed process by its name.

### `(*ProcessManager) RestartProcess(name string) error`

Restarts a specific managed process by its name.

### `(*ProcessManager) GetProcess(name string) (*ManagedProcess, bool)`

Retrieves a `ManagedProcess` instance by its name.

### `(*ProcessManager) GetProcesses() map[string]*ManagedProcess`

Returns a map of all currently managed processes.

### `(*ProcessManager) IsRunning(name string) bool`

Checks if a specific process is currently running.

### `(*ProcessManager) GetProcessStats() map[string]ProcessStats`

Returns a map of `ProcessStats` for all managed processes.

### `(*ProcessManager) StartTypeScriptWatch(ctx *core.JawtContext) error`

Convenience method to start the TypeScript compiler in watch mode, configured with project-specific paths.

### `(*ProcessManager) StartTailwindWatch(ctx *core.JawtContext) error`

Convenience method to start the Tailwind CSS compiler in watch mode, configured with project-specific paths.

### `(*ProcessManager) StartNodeProcess(name string, args []string, workingDir string, outputHandler func(string), errorHandler func(error)) error`

Starts a generic Node.js process with specified arguments, working directory, and output/error handlers.

### `(*ProcessManager) StartDevServer(ctx *core.JawtContext, serverBinary string) error`

Convenience method to start the development server process.

### `(*ProcessManager) Health() map[string]bool`

Returns a map indicating the running status (health) of all managed processes.

### `NewManagedProcess`

```go
func NewManagedProcess(name string, options ProcessOptions, ctx context.Context, logger core.Logger) *ManagedProcess
```

Creates a new `ManagedProcess` instance.

### `(*ManagedProcess) Start()`

Starts the underlying external command. This method is typically called by the `ProcessManager`.

### `(*ManagedProcess) Stop()`

Stops the underlying external command.

### `(*ManagedProcess) Restart()`

Restarts the underlying external command.

### `(*ManagedProcess) IsRunning()`

Checks if the process is currently running.

### `(*ManagedProcess) Name()`

Returns the name of the managed process.

### `(*ManagedProcess) GetStats()`

Returns the `ProcessStats` for the managed process.

### `DefaultProcessOptions()`

```go
func DefaultProcessOptions() ProcessOptions
```

Returns a `ProcessOptions` struct with default values.

### `WithCommand`, `WithWorkingDir`, `WithEnv`, `WithRestart`, `WithOutputHandler`, `WithErrorHandler`

Chainable methods on `ProcessOptions` to easily configure process settings.

### `TypeScriptWatchOptions`, `TailwindWatchOptions`, `DevServerOptions`, `NodeScriptOptions`

Convenience functions to create pre-configured `ProcessOptions` for common development tools.

## Usage Examples

### Starting a TypeScript Watch Process

```go
// In internal/runtime/orchestrator.go
// ...
// if err := o.processManager.StartTypeScriptWatch(o.jawtContext); err != nil {
//     o.logger.Error("Failed to start TypeScript watch", core.ErrorField(err))
// }
// ...
```

### Starting a Node.js Server

```go
// In internal/runtime/orchestrator.go
// ...
// return o.processManager.StartNodeProcess(
//     "tsserver",
//     []string{"node_modules/typescript/bin/tsserver"},
//     o.jawtContext.Paths.ProjectRoot,
//     func(output string) {
//         o.logger.Debug("TSServer output", core.StringField("output", output))
//     },
//     func(err error) {
//         o.logger.Error("TSServer error", core.ErrorField(err))
//     },
// )
// ...
```
