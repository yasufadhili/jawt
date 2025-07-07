# The Checker (`internal/checker`)

The `checker` is where we make sure the JML code actually makes sense. The parser ensures the syntax is right, but the checker handles the semantics. It answers questions like, "Is this variable declared?" or "Are you trying to add a string to a number?"

## The Core Ideas

-   **Symbol Table**: This is the heart of the checker. It's a data structure that keeps track of every identifier (variable, function, component, etc.) in the code. It knows what they are, what their type is, and where they live.
-   **Scope**: This is all about context. A variable declared inside a function only exists within that function. The checker manages a stack of scopes (global, document, function, block) to keep track of what's visible where.
-   **Semantic Analysis**: This is the fancy term for what the checker does. It walks the AST, uses the symbol table to understand the code, and reports any issues it finds.

## The Key Data Structures

### `SymbolTable`

This is the main manager for all the symbols and scopes. It has a global scope, and it keeps track of the current scope as it walks the AST.

```go
type SymbolTable struct {
	global     *Scope
	current    *Scope
	components map[string]*ComponentSymbol
	builtIns   map[string]*BuiltInSymbol
}
```

### `Scope`

A `Scope` holds all the symbols declared within it and has a pointer to its parent scope. This is how we get lexical scoping.

```go
type Scope struct {
	kind     ScopeKind
	parent   *Scope
	children []*Scope
	symbols  map[string]Symbol
	name     string // Just for debugging
}
```

### `Symbol` Interface

This is the base interface for all symbols. It just defines the common methods that every symbol must have.

```go
type Symbol interface {
	Name() string
	Kind() SymbolKind
	Type() string
	Position() ast.Position
}
```

### `Checker`

The main struct for the checker. It's an `ast.Visitor`, so it can walk the AST. It has a `SymbolTable` to manage symbols and a `diagnostic.Reporter` to report any errors it finds.

```go
type Checker struct {
	*ast.BaseVisitor
	reporter *diagnostic.Reporter
	table    *SymbolTable
}
```

## The Process

### `EnterScope` and `ExitScope`

As the checker walks the AST, it calls `EnterScope` whenever it enters a new scope (like a function or a block). When it leaves that scope, it calls `ExitScope`. This keeps the `current` scope in the `SymbolTable` up to date.

### `Define`

When the checker sees a declaration (like `let x = 10`), it creates a new symbol and calls `Define` to add it to the current scope.

### `Lookup` and `LookupRecursive`

When the checker sees an identifier being used (like `y = x + 5`), it calls `LookupRecursive` to find the declaration of that identifier. It starts in the current scope and then walks up the parent scopes until it finds the symbol or reaches the global scope.

## How It's Used

The checker runs after the parser has built the AST. It walks the AST, fills up the symbol table, and reports any semantic errors it finds. If the checker passes, we can be pretty confident that the code is valid and ready for the next stage: code generation.
