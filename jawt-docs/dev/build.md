# Build System (`internal/build`)

The `internal/build` package is responsible for managing the entire build process of a JAWT project. This includes discovering JML files, managing their dependencies, compiling them, and reacting to file system changes to trigger recompilations.

## Core Concepts

*   **DocumentInfo**: Represents metadata about a JML file (page or component), including its path, type, and compilation status.
*   **DependencyGraph**: Manages the relationships between JML documents, allowing for topological sorting and cycle detection.
*   **BuildSystem**: The central orchestrator that ties together document discovery, compilation, and file watching.

## Key Data Structures

### `DocumentType`

An enumeration representing the type of a JML document.

```go
type DocumentType int

const (
	DocumentTypePage DocumentType = iota
	DocumentTypeComponent
)
```

### `DocumentInfo`

Common information for all JML document types.

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
	Hash         string // Content hash for detecting changes
}
```

### `ComponentInfo`

Specific information for JML components.

```go
type ComponentInfo struct {
	DocumentInfo
	Props map[string]string
}
```

### `PageInfo`

Specific information for JML pages.

```go
type PageInfo struct {
	DocumentInfo
	Route string
}
```

### `BuildSystem`

The main struct that manages the build process.

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

### `FileWatcher` Interface

An interface defining the contract for file watching, allowing for different implementations.

```go
type FileWatcher interface {
	OnEvent(handler func(fsnotify.Event))
	Start() error
	Stop() error
}
```

### `DependencyGraph` Interface

An interface for managing the dependency graph between documents.

```go
type DependencyGraph interface {
	AddNode(path string, docType DocumentType) error
	RemoveNode(path string) error
	AddDependency(from, to string) error
	RemoveDependency(from, to string) error
	GetDependencies(path string) []string
	GetDependents(path string) []string
	GetAllNodes() []string
	HasCycle() bool
	GetCycles() [][]string
	GetTopologicalOrder() ([]string, error)
	GetCompilationOrder() ([]string, error)
	ValidateGraph() error
	IsConnected(from, to string) bool
	GetShortestPath(from, to string) []string
	GetNodesByType(docType DocumentType) []string
	GetTransitiveDependencies(path string) []string
	GetTransitiveDependents(path string) []string
}
```

### `GraphNode`

Represents a node within the `DependencyGraph`.

```go
type GraphNode struct {
	Path     string
	DocType  DocumentType
	Metadata map[string]interface{}
}
```

## Functions & Methods

### `NewBuildSystem`

```go
func NewBuildSystem(ctx *core.JawtContext, compiler *compiler.Compiler, watcher FileWatcher) *BuildSystem
```

Creates a new `BuildSystem` instance.

### `Initialise`

```go
func (bs *BuildSystem) Initialise() error
```

Performs initial project discovery and compilation. This is typically called once at the start of the development server.

### `DiscoverProject`

```go
func (bs *BuildSystem) DiscoverProject() error
```

Finds all JML documents in the project, creates `DocumentInfo` for each, and builds the internal dependency graph.

### `CompileAll`

```go
func (bs *BuildSystem) CompileAll() error
```

Compiles all documents in the project based on their dependency order.

### `SetupWatcher`

```go
func (bs *BuildSystem) SetupWatcher()
```

Configures the file watcher to listen for changes in JML files and triggers appropriate build actions.

### `HandleFileEvent`

```go
func (bs *BuildSystem) HandleFileEvent(event fsnotify.Event)
```

Processes file system events (create, write, remove, rename) for JML files, delegating to more specific handlers.

### `HandleFileCreated`, `HandleFileModified`, `HandleFileDeleted`, `HandleFileRenamed`

Specific handlers for different file system events, updating the build system and triggering recompilations as needed.

### `GetDocumentInfo`

```go
func (bs *BuildSystem) GetDocumentInfo(path string) (*DocumentInfo, bool)
```

Retrieves `DocumentInfo` for a given file path.

### `AddDocument`, `RemoveDocument`

Methods to add or remove documents from the build system's internal tracking.

### `CompileDocument`

```go
func (bs *BuildSystem) CompileDocument(path string) error
```

Compiles a single JML document. Reports and prints any compilation errors.

### `RecompileDependents`

```go
func (bs *BuildSystem) RecompileDependents(path string) error
```

Recompiles all documents that directly or indirectly depend on the given document. This is crucial for hot-reloading and ensuring consistency after a change.

### `ExtractDependencies`

```go
func ExtractDependencies(content string) []string
```

Extracts component and script import paths from the content of a JML file using regular expressions.

### `NewDependencyGraph`

```go
func NewDependencyGraph() DependencyGraph
```

Creates a new, empty `DependencyGraph`.

### `AddNode`, `RemoveNode`, `AddDependency`, `RemoveDependency`

Core methods for manipulating the dependency graph by adding/removing nodes and edges.

### `GetDependencies`, `GetDependents`, `GetAllNodes`

Query methods for retrieving information from the dependency graph.

### `HasCycle`, `GetCycles`

Methods for detecting and retrieving circular dependencies within the graph.

### `GetTopologicalOrder`, `GetCompilationOrder`

Methods for performing topological sorting on the graph, which is essential for determining the correct compilation order of documents.

### `ValidateGraph`

```go
func (dg *dependencyGraph) ValidateGraph() error
```

Checks the internal consistency of the dependency graph.

### `IsConnected`, `GetShortestPath`

Utility methods for graph traversal and pathfinding.

### `GetNodesByType`, `GetTransitiveDependencies`, `GetTransitiveDependents`

Additional utility methods for querying nodes based on type and finding all direct and indirect dependencies/dependents.

### `RunProject` (`run.go`)

```go
func RunProject(ctx *core.JawtContext) error
```

This function is the entry point for running the JAWT project in development mode. It initializes the build system, sets up file watchers, and starts the HTTP development server.

**Flow:**
1.  Initializes `FileWatcher` and adds project paths for watching.
2.  Creates a `compiler.Compiler` instance.
3.  Creates and initializes a `BuildSystem`.
4.  Starts the `FileWatcher`.
5.  Sets up and starts an HTTP server to serve the built files.
6.  Waits for context cancellation (e.g., Ctrl+C) to gracefully shut down.

### `InitProject` (`init.go`)

```go
func InitProject(ctx *core.JawtContext, projectName string, targetDir string) error
```

Initializes a new JAWT project by creating the directory structure, generating configuration files (`jawt.project.json`, `tsconfig.json`, `tailwind.config.js`), and creating initial JML and TypeScript files from embedded templates.

**Key Steps:**
1.  Validates and sanitizes the `projectName`.
2.  Determines the `targetDir` for the project.
3.  Handles existing directories (creates if not exists, checks for conflicts if not empty).
4.  Creates the basic project directory structure (`app`, `components`, `scripts`, `assets`).
5.  Generates `jawt.project.json` with default settings.
6.  Creates initial `index.jml`, `layout.jml`, and `main.ts` files from embedded templates.
7.  Prints a success message with next steps for the user.

### `DiscoverProjectFiles` (`discovery.go`)

```go
func DiscoverProjectFiles(ctx *core.JawtContext) ([]string, error)
```

Discovers all `.jml` files within the configured pages and components directories of the project. It uses `findJMLFiles` internally to recursively search directories.

### `CreateDocumentInfo` (`discovery.go`)

```go
func CreateDocumentInfo(path string, projectRoot string) (*DocumentInfo, error)
```

Creates a `DocumentInfo` struct for a given JML file path. It determines the `DocumentType` (page or component) based on the file's location within the project structure and extracts the document's name and relative path.

## Usage Examples

### Initializing a new project

```go
// In cmd/init.go
// ...
// err = build.InitProject(ctx, projectName, projectDir)
// if err != nil {
//     logger.Error("Failed to initialise project", core.ErrorField(err))
//     os.Exit(1)
// }
// ...
```

### Running the development server

```go
// In cmd/run.go
// ...
// err = build.RunProject(ctx)
// if err != nil {
//     logger.Error("Failed to run project", core.ErrorField(err))
//     os.Exit(1)
// }
// ...
```

### Accessing Build System Information

```go
// Example of getting a document's info
// docInfo, exists := buildSystem.GetDocumentInfo("/path/to/my/page.jml")
// if exists {
//     fmt.Printf("Document Name: %s, Type: %s\n", docInfo.Name, docInfo.Type)
// }

// Example of getting compilation order
// order, err := buildSystem.depGraph.GetCompilationOrder()
// if err == nil {
//     fmt.Printf("Compilation Order: %v\n", order)
// }
```