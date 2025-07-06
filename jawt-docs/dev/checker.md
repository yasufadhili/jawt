# Checker (`internal/checker`)

The `internal/checker` package is responsible for semantic analysis of JML code. It ensures that the code is not only syntactically correct (parsed by the compiler) but also semantically valid. This involves building and managing a symbol table to track declarations and their scopes, and then using this information to identify issues like undefined variables, type mismatches, and incorrect usage of components.

## Core Concepts

*   **Symbol Table**: A data structure that stores information about identifiers (symbols) in the program, such as their name, type, kind (e.g., variable, function, component), and scope.
*   **Scope**: A region of the program where a declared name is valid. The checker manages different types of scopes (global, document, function, block, element).
*   **Symbol**: Represents a declared entity in the code. The checker defines various types of symbols (e.g., `PropertySymbol`, `StateSymbol`, `FunctionSymbol`).
*   **Semantic Analysis**: The process of checking the meaning and consistency of the code, going beyond just its syntax.

## Key Data Structures

### `SymbolKind`

An enumeration representing the different types of symbols that can be stored in the symbol table.

```go
type SymbolKind int

const (
	SymbolProperty SymbolKind = iota
	SymbolState
	SymbolComponent
	SymbolFunction
	SymbolVariable
	SymbolParameter
	SymbolImport
	SymbolBuiltIn
)
```

### `ScopeKind`

An enumeration representing the different types of scopes in the program.

```go
type ScopeKind int

const (
	ScopeGlobal ScopeKind = iota
	ScopeDocument
	ScopeFunction
	ScopeBlock
	ScopeElement
)
```

### `Symbol` Interface

The interface that all symbol types must implement.

```go
type Symbol interface {
	Name() string
	Kind() SymbolKind
	Type() string
	Position() ast.Position
	String() string
}
```

### `BaseSymbol`

Provides common fields and methods for all concrete symbol implementations.

```go
type BaseSymbol struct {
	name     string
	kind     SymbolKind
	typeStr  string
	position ast.Position
}
```

### Concrete Symbol Types

The `checker` package defines several concrete implementations of the `Symbol` interface:

*   `PropertySymbol`: Represents a component property.
*   `StateSymbol`: Represents a component state variable.
*   `ComponentSymbol`: Represents a JML component (document).
*   `FunctionSymbol`: Represents a function declaration.
*   `ParameterSymbol`: Represents a function parameter.
*   `VariableSymbol`: Represents a general variable declaration (`let`, `const`).
*   `ImportSymbol`: Represents an import statement.
*   `BuiltInSymbol`: Represents built-in types and functions (e.g., `string`, `console.log`).

### `Scope`

Represents a lexical scope, holding a collection of symbols and a reference to its parent scope.

```go
type Scope struct {
	kind     ScopeKind
	parent   *Scope
	children []*Scope
	symbols  map[string]Symbol
	name     string // Optional name for debugging
}
```

### `SymbolTable`

Manages the entire symbol table, including the global scope and the current active scope. It provides methods for entering and exiting scopes, defining new symbols, and looking up existing symbols.

```go
type SymbolTable struct {
	global     *Scope
	current    *Scope
	components map[string]*ComponentSymbol
	builtIns   map[string]*BuiltInSymbol
}
```

### `Checker`

The main struct for performing semantic checks. It embeds `ast.BaseVisitor` to traverse the AST and uses a `SymbolTable` to manage symbols and scopes, and a `diagnostic.Reporter` to report issues.

```go
type Checker struct {
	*ast.BaseVisitor
	reporter *diagnostic.Reporter
	table    *SymbolTable
}
```

## Functions & Methods

### `NewChecker`

```go
func NewChecker(reporter *diagnostic.Reporter) *Checker
```

Creates a new `Checker` instance, initializing its internal `SymbolTable`.

### `NewSymbolTable`

```go
func NewSymbolTable() *SymbolTable
```

Creates a new `SymbolTable`, pre-populating it with built-in types and functions.

### `EnterScope`, `ExitScope`

```go
func (st *SymbolTable) EnterScope(kind ScopeKind, name string) *Scope
func (st *SymbolTable) ExitScope() error
```

Methods of `SymbolTable` to manage the scope stack. `EnterScope` creates a new child scope and makes it the current scope. `ExitScope` reverts to the parent scope.

### `Define`

```go
func (st *SymbolTable) Define(symbol Symbol) error
func (s *Scope) Define(symbol Symbol) error
```

Adds a new symbol to the current scope. Returns an error if a symbol with the same name is already defined in the current scope.

### `Lookup`, `LookupRecursive`

```go
func (st *SymbolTable) Lookup(name string) (Symbol, bool)
func (s *Scope) Lookup(name string) (Symbol, bool)
func (s *Scope) LookupRecursive(name string) (Symbol, bool)
```

`Lookup` searches for a symbol only in the current scope. `LookupRecursive` searches in the current scope and then recursively in parent scopes until the symbol is found or the global scope is reached.

### `DefineComponent`, `GetComponent`, `GetAllComponents`

Methods of `SymbolTable` specifically for managing `ComponentSymbol`s, which are stored globally.

### `IsBuiltIn`, `GetBuiltIn`

Methods of `SymbolTable` to check for and retrieve built-in symbols.

### `Debug`

```go
func (st *SymbolTable) Debug() string
```

Returns a string representation of the entire symbol table structure, useful for debugging.

## Usage Example

The `Checker` would typically be used after the AST has been built by the compiler. It traverses the AST, populating the symbol table and reporting any semantic errors.

```go
// Example (conceptual usage within a larger compilation pipeline):
// import (
//     "github.com/yasufadhili/jawt/internal/ast"
//     "github.com/yasufadhili/jawt/internal/checker"
//     "github.com/yasufadhili/jawt/internal/diagnostic"
// )

// func performSemanticAnalysis(documentAST *ast.Document) {
//     reporter := diagnostic.NewReporter()
//     checker := checker.NewChecker(reporter)

//     // Traverse the AST to perform checks and populate the symbol table
//     ast.Walk(checker, documentAST)

//     if reporter.HasErrors() {
//         fmt.Println("Semantic errors found:")
//         printer := diagnostic.NewPrinter()
//         printer.Print(reporter)
//     } else {
//         fmt.Println("Semantic analysis completed with no errors.")
//     }

//     // You can also inspect the symbol table for debugging or further analysis
//     // fmt.Println(checker.table.Debug())
// }
```
