# The Build System (`internal/build`)

The `internal/build` package is the brains of the operation. It's responsible for taking a JAWT project, figuring out what's what, and turning it into a real website. This means finding all the JML files, understanding how they depend on each other, compiling them, and watching for changes.

## The Core Ideas

-   **`DocumentInfo`**: This is like a passport for a JML file. It holds all the important metadata, like its path, whether it's a page or a component, and if it's been compiled yet.
-   **`DependencyGraph`**: This is a map of how all the JML files are connected. It's crucial for figuring out the right order to compile things and for detecting any circular dependencies (which are bad news).
-   **`BuildSystem`**: This is the main orchestrator. It brings everything together: file discovery, compilation, and watching for changes.

## The Key Data Structures

### `DocumentInfo`

This struct holds all the common info for a JML document.

```go
type DocumentInfo struct {
	Name         string
	RelPath      string
	AbsPath      string
	Type         DocumentType
	Dependencies []string
	DependedBy   []string
	IsCompiled   bool
	LastModified time.Time
	Hash         string // A hash of the file content to detect changes
}
```

### `BuildSystem`

The main struct that manages the whole build process.

```go
type BuildSystem struct {
	ctx      *core.JawtContext
	mu       sync.RWMutex
	docs     map[string]*DocumentInfo
	pages    map[string]*PageInfo
	comps    map[string]*ComponentInfo
	compiler *compiler.Compiler
	watcher  FileWatcher
	depGraph DependencyGraph
}
```

### `DependencyGraph` Interface

This interface defines all the things the dependency graph needs to be able to do.

```go
type DependencyGraph interface {
	AddNode(path string, docType DocumentType) error
	AddDependency(from, to string) error
	GetDependencies(path string) []string
	GetDependents(path string) []string
	HasCycle() bool
	GetCompilationOrder() ([]string, error)
    // ... and more
}
```

## The Process

### `Initialise`

This is where it all starts. When you run `jawt run`, this method gets called to do the first scan of the project and compile everything.

### `DiscoverProject`

This method walks through the project directory, finds all the `.jml` files, creates a `DocumentInfo` for each one, and builds the dependency graph.

### `CompileAll`

This method compiles all the documents in the project. It uses the dependency graph to make sure everything is compiled in the right order.

### `SetupWatcher`

This sets up the file watcher to keep an eye on all the JML files. When a file changes, it triggers the right build actions.

### `HandleFileEvent`

This is the callback for the file watcher. When a file is created, modified, or deleted, this method figures out what to do next.

### `CompileDocument`

This method compiles a single JML file. It's what gets called when a file is changed or when it's part of the initial build.

### `RecompileDependents`

This is a really important one. When a component changes, we need to recompile not just that component, but also every page or component that uses it. This method figures out all the dependents and recompiles them.

### `RunProject` (`run.go`)

This is the entry point for the `jawt run` command. It sets up the build system, starts the file watcher, and kicks off the dev server.

### `InitProject` (`init.go`)

This is the entry point for the `jawt init` command. It creates the basic project structure and all the starting files from templates.