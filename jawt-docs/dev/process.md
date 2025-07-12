# Process Management (`internal/process`)

This package is all about running other command-line tools. JAWT needs to run the TypeScript compiler (`tsc`) and the Tailwind CSS compiler, so I needed a way to manage these external processes.

## The Core Ideas

-   **`ManagedProcess`**: This represents a single external process that we're controlling. It knows how to start, stop, and restart the process.
-   **`ProcessManager`**: This is the main orchestrator. It keeps track of all the `ManagedProcess` instances and gives us a simple API to interact with them.
-   **`ProcessOptions`**: This is just a struct that holds all the configuration for a process, like the command to run, the arguments, the working directory, etc.

## The Key Data Structures

### `ManagedProcess`

This struct holds all the state for a single managed process.

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

The main manager for all the processes.

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

## How It's Used

When the `jawt run` command starts, it creates a `ProcessManager`. Then, it tells the process manager to start the TypeScript compiler in watch mode and the Tailwind CSS compiler in watch mode. The process manager takes care of running these commands, capturing their output, and restarting them if they crash.

This makes the main build logic much cleaner. It doesn't have to worry about the details of running external processes; it just tells the `ProcessManager` what to do.
