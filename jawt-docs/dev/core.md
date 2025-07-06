# Core Utilities & Configuration (`internal/core`)

The `internal/core` package provides fundamental utilities, configuration management, and the central context for the entire JAWT toolchain. It defines how JAWT projects are configured, how logging is handled, and how paths within a project are managed.

## Core Concepts

*   **Configuration**: JAWT uses two main configuration types: `JawtConfig` (global JAWT settings) and `ProjectConfig` (project-specific settings).
*   **Context**: The `JawtContext` struct acts as a central hub, passing essential configurations and shared resources (like the logger) to various subsystems.
*   **Logging**: A structured logging interface (`Logger`) and a default implementation (`DefaultLogger`) provide consistent and informative output.
*   **Paths**: The `ProjectPaths` struct centralizes the management and resolution of all file and directory paths within a JAWT project.

## Key Data Structures

### `JawtConfig`

Represents the global JAWT configuration, typically loaded from `jawt.config.json` (though currently defaults are used if no file is present). It includes paths to external tools and general JAWT settings.

```go
type JawtConfig struct {
	TypeScriptPath     string `json:"typescript_path"`
	TailwindPath       string `json:"tailwind_path"`
	NodePath           string `json:"node_path"`
	DefaultPort        int    `json:"default_port"`
	TempDir            string `json:"temp_dir"`
	CacheDir           string `json:"cache_dir"`
	EnableMinification bool `json:"enable_minification"`
	EnableSourceMaps   bool `json:"enable_source_maps"`
	EnableTreeShaking  bool `json:"enable_tree_shaking"`
}
```

### `ProjectConfig`

Represents the project-specific configuration, loaded from `jawt.project.json`. This struct is designed with nested groupings for better organization.

```go
type ProjectConfig struct {
	App struct {
		Name        string `json:"name"`
		Author      string `json:"author"`
		Version     string `json:"version"`
		Description string `json:"description"`
	} `json:"app"`
	Paths struct {
		Components string `json:"components"`
		Pages      string `json:"pages"`
		Scripts    string `json:"scripts"`
		Assets     string `json:"assets"`
	} `json:"paths"`
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	Build struct {
		OutputDir string `json:"outputDir"`
		DistDir   string `json:"distDir"`
		Minify    bool   `json:"minify"`
		ShadowDOM bool   `json:"shadowDOM"`
	} `json:"build"`
	Dev struct {
		Port       int      `json:"port"`
		EnableHMR  bool     `json:"enableHMR"`
		WatchPaths []string `json:"watchPaths"`
	} `json:"dev"`
	Tooling struct {
		TSConfigPath       string `json:"tsConfigPath"`
		TailwindConfigPath string `json:"tailwindConfigPath"`
	} `json:"tooling"`
	Scripts struct {
		PreBuild  []string `json:"preBuild"`
		PostBuild []string `json:"postBuild"`
	} `json:"scripts"`
}
```

### `BuildOptions`

Represents build-time options or detected features that influence the build process.

```go
type BuildOptions struct {
	UsesTailwindCSS bool
}
```

### `JawtContext`

The central context object passed throughout the application. It holds references to configurations, paths, the logger, and build options.

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

### `Logger` Interface & `DefaultLogger`

Defines a structured logging interface and its default implementation, supporting different log levels (Debug, Info, Warn, Error, Fatal) and structured fields.

```go
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

// DefaultLogger is a simple implementation of the Logger interface
type DefaultLogger struct {
	level  LogLevel
	logger *log.Logger
}
```

### `LogLevel`

An enumeration for different logging levels.

```go
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)
```

### `Field`

Represents a key-value pair for structured logging.

```go
type Field struct {
	Key   string
	Value interface{}
}
```

### `ProjectPaths`

Manages and resolves all relevant file and directory paths within a JAWT project, including source, build, and temporary directories.

```go
type ProjectPaths struct {
	ProjectRoot string
	WorkingDir  string

	ProjectConfig *ProjectConfig
	JawtConfig    *JawtConfig

	AppDir        string
	ComponentsDir string
	ScriptsDir    string
	AssetsDir     string

	JawtDir  string
	BuildDir string
	DistDir  string
	TempDir  string
	CacheDir string

	TypeScriptOutputDir string
	TailwindOutputDir   string
	ComponentsOutputDir string

	TSConfigPath       string
	TailwindConfigPath string
	ProjectConfigPath  string
}
```

## Functions & Methods

### Configuration (`config.go`)

*   `DefaultJawtConfig()`: Returns a `JawtConfig` with default values.
*   `DefaultProjectConfig()`: Returns a `ProjectConfig` with default values, reflecting the new structured format.
*   `LoadJawtConfig(configPath string)`: Loads global JAWT configuration from a specified path or returns defaults.
*   `LoadProjectConfig(projectDir string)`: Loads project-specific configuration from `jawt.project.json` within the given directory or returns defaults.
*   `(*JawtConfig) Save(configPath string)`: Saves the global JAWT configuration to a file.
*   `(*ProjectConfig) Save(projectDir string)`: Saves the project configuration to `jawt.project.json`.
*   `(*JawtConfig) Validate()`: Validates the global JAWT configuration.
*   `(*ProjectConfig) Validate()`: Validates the project configuration, ensuring all required paths and settings are valid.
*   **Getter Methods on `ProjectConfig`**: Numerous methods like `GetComponentsPath`, `GetPagesPath`, `GetBuildOutputDir`, `GetServerAddress`, `IsMinificationEnabled`, `IsHMRenabled`, etc., provide convenient access to specific configuration values, often resolving them to absolute paths.
*   **Setter Methods on `ProjectConfig`**: Methods like `SetProjectName`, `SetServerPort`, `SetMinification`, etc., allow programmatic modification of project settings.

### Context (`context.go`)

*   `NewJawtContext(jawtConfig *JawtConfig, projectConfig *ProjectConfig, paths *ProjectPaths, logger Logger, buildOptions *BuildOptions)`: Creates a new `JawtContext` instance.
*   `(*JawtContext) Context()`: Returns the underlying `context.Context` for cancellation and deadlines.
*   `(*JawtContext) Cancel()`: Cancels the context, signaling all dependent goroutines to shut down.
*   `(*JawtContext) SetMetadata(key string, value interface{})`, `(*JawtContext) GetMetadata(key string)`: Allow storing and retrieving arbitrary metadata within the context.

### Logging (`logger.go`)

*   `NewDefaultLogger(level LogLevel)`: Creates a new `DefaultLogger` with a specified minimum log level.
*   `(*DefaultLogger) Debug`, `Info`, `Warn`, `Error`, `Fatal`: Methods for logging messages at different severity levels. `Fatal` also exits the application.
*   `(*DefaultLogger) SetLevel(level LogLevel)`, `(*DefaultLogger) GetLevel()`: Control the verbosity of the logger.
*   `StringField`, `IntField`, `BoolField`, `ErrorField`, `DurationField`: Helper functions to create structured log `Field`s.

### Paths (`paths.go`)

*   `NewProjectPaths(projectRoot string, projectConfig *ProjectConfig, jawtConfig *JawtConfig)`: Creates a new `ProjectPaths` instance, initializing all relevant paths based on the provided configurations.
*   `ResolveExecutablePath(cmd string)`: Resolves the absolute path to an executable, searching in various locations (absolute path, system PATH, relative to JAWT executable).
*   `(*ProjectPaths) EnsureDirectories()`: Creates all necessary build, temporary, and cache directories for the project.
*   `(*ProjectPaths) GetRelativePath(path string)`: Returns a path relative to the project root.
*   `(*ProjectPaths) GetAbsolutePath(relativePath string)`: Returns an absolute path from a relative path.
*   `(*ProjectPaths) GetJMLFiles()`, `(*ProjectPaths) GetTypeScriptFiles()`: Discover and return lists of JML and TypeScript files within the project's configured source directories.
*   `(*ProjectPaths) GetWatchPaths()`: Returns a list of all directories and files that should be watched for changes during development, based on `ProjectConfig.Dev.WatchPaths` and other config files.
*   `(*ProjectPaths) GetTempFile(filename string)`, `(*ProjectPaths) GetCacheFile(filename string)`: Generate paths for temporary and cache files within the `.jawt` directory.
*   `(*ProjectPaths) Clean()`: Removes all generated build, distribution, and internal JAWT directories.

## Usage Examples

### Initializing a JawtContext

```go
// In cmd/run.go or similar entry point
// ...
// cfg, err := core.LoadJawtConfig("") // Load global config
// projectConfig, err := core.LoadProjectConfig(projectDir) // Load project config
// paths, err := core.NewProjectPaths(projectDir, projectConfig, cfg)
// logger := core.NewDefaultLogger(core.InfoLevel)
// buildOptions := core.NewBuildOptions()
// ctx := core.NewJawtContext(cfg, projectConfig, paths, logger, buildOptions)
// ...
```

### Logging a message

```go
// In any package with access to JawtContext.Logger
// ctx.Logger.Info("Application started", core.StringField("version", "1.0.0"))
// ctx.Logger.Error("Failed to load file", core.StringField("file", "config.json"), core.ErrorField(err))
```

### Accessing Project Paths

```go
// In any package with access to JawtContext.Paths
// jmlFiles, err := ctx.Paths.GetJMLFiles()
// if err == nil {
//     fmt.Printf("Found %d JML files.\n", len(jmlFiles))
// }
// buildDir := ctx.Paths.BuildDir
// fmt.Printf("Build output will go to: %s\n", buildDir)
```