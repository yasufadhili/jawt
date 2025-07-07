# The Core (`internal/core`)

The `internal/core` package is the central nervous system of JAWT. It holds all the fundamental utilities, configuration, and the main context that gets passed around everywhere.

## The Core Ideas

-   **Configuration**: JAWT has two types of configuration: `JawtConfig` for global settings and `ProjectConfig` for project-specific settings.
-   **Context**: The `JawtContext` is a big struct that holds everything important: the configs, the logger, project paths, etc. It gets passed to almost every function so that all parts of the toolchain have access to the same information.
-   **Logging**: I set up a simple, structured logger to have consistent and clean output throughout the CLI.
-   **Paths**: The `ProjectPaths` struct is a helper for managing all the different file and directory paths in a JAWT project.

## The Key Data Structures

### `JawtConfig`

This holds the global configuration for JAWT, like paths to external tools and default settings.

```go
type JawtConfig struct {
	TypeScriptPath     string `json:"typescript_path"`
	TailwindPath       string `json:"tailwind_path"`
	NodePath           string `json:"node_path"`
	// ... and more
}
```

### `ProjectConfig`

This holds the configuration for a specific JAWT project, loaded from `jawt.project.json`.

```go
type ProjectConfig struct {
	App struct {
		Name        string `json:"name"`
		// ...
	} `json:"app"`
	Paths struct {
		Components string `json:"components"`
		Pages      string `json:"pages"`
		// ...
	} `json:"paths"`
	// ... and more
}
```

### `JawtContext`

The big one. This struct gets passed around everywhere.

```go
type JawtContext struct {
	ctx    context.Context
	cancel context.CancelFunc

	JawtConfig    *JawtConfig
	ProjectConfig *ProjectConfig
	Paths         *ProjectPaths
	Logger        Logger

	BuildOptions *BuildOptions

	mu       sync.RWMutex
	metadata map[string]interface{}
}
```

### `Logger` Interface

A simple interface for structured logging.

```go
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}
```

### `ProjectPaths`

A helper struct for managing all the paths in a project.

```go
type ProjectPaths struct {
	ProjectRoot string
	WorkingDir  string

	AppDir        string
	ComponentsDir string
	ScriptsDir    string
	AssetsDir     string

	BuildDir string
	DistDir  string
	// ... and more
}
```

## How It All Works

When you run a `jawt` command, the first thing that happens is a `JawtContext` gets created. It loads the global and project configs, sets up the logger, and figures out all the project paths. Then, this context is passed down to all the different parts of the toolchain, so they all have access to the same information.

This makes the code much cleaner and easier to manage. Instead of passing around a dozen different arguments to every function, I just pass the context.
