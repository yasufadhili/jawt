# Abstract Syntax Tree (AST)

The `ast` package defines the core data structures representing the Abstract Syntax Tree (AST) for JML (JAWT Markup Language) files. It provides interfaces for various node types, concrete implementations of these nodes, factory methods for creating them, a visitor pattern for traversing the AST, and a printer for visualizing the tree structure.

## Core Concepts

*   **Node**: The fundamental interface for all elements in the AST. Every node has a `Position` (line, column, file) and can accept a `Visitor`.
*   **Statement**: Represents a complete action or instruction.
*   **Declaration**: A specific type of statement that introduces a new identifier (e.g., variable, function, component).
*   **Expression**: Represents a piece of code that produces a value.
*   **Visitor Pattern**: A way to traverse the AST and perform operations on each node without modifying the node structures themselves.

## Key Data Structures

### `Position`

Represents a location within a source file.

```go
type Position struct {
	Line   int
	Column int
	File   string
}
```

**Methods:**

*   `String() string`: Returns a formatted string of the position (e.g., `file.jml:10:5`).

### `Node` Interface

The base interface for all AST nodes.

```go
type Node interface {
	Pos() Position
	String() string
	Accept(v Visitor)
}
```

### `Program`

Represents the entire JML project, containing multiple `Document` nodes.

```go
type Program struct {
	Position
	Documents []*Document
}
```

### `Document`

Represents a single JML file (e.g., a page or a component).

```go
type Document struct {
	Position
	DocType    DocType // "page" or "component"
	Name       *Identifier
	Body       []Statement // Contains imports, exports, declarations, elements
	SourceFile string
}
```

### `DocType`

An enumeration for the type of JML document.

```go
type DocType string

const (
	DocTypePage      DocType = "page"
	DocTypeComponent DocType = "component"
)
```

### `ImportKind`

An enumeration for the type of import declaration.

```go
type ImportKind string

const (
	ImportKindComponent ImportKind = "component"
	ImportKindScript    ImportKind = "script"
	ImportKindBrowser   ImportKind = "browser"
	ImportKindModule    ImportKind = "module"
)
```

### Declarations

Nodes representing various types of declarations:

*   `ImportDeclaration`: `import { a } from 'b'`
*   `ExportDeclaration`: `export const a = 1`
*   `VariableDeclaration`: `let`, `const`, or `var` declarations.
*   `FunctionDeclaration`: Function definitions.
*   `ClassDeclaration`: Class definitions.
*   `InterfaceDeclaration`: Interface definitions.
*   `TypeAliasDeclaration`: Type aliases (e.g., `type MyType = string`).
*   `EnumDeclaration`: Enum definitions.
*   `PropertyDeclaration`: JML-specific component properties.
*   `StateDeclaration`: JML-specific component state.

### Statements

Nodes representing various types of statements:

*   `BlockStatement`: A block of statements `{ ... }`.
*   `ExpressionStatement`: An expression used as a statement.
*   `IfStatement`: `if-else` statements.
*   `ForStatement`: `for` loops.
*   `ForInStatement`: `for-in` or `for-of` loops.
*   `WhileStatement`: `while` loops.
*   `ReturnStatement`: `return` statements.
*   `BreakStatement`: `break` statements.
*   `ContinueStatement`: `continue` statements.
*   `ThrowStatement`: `throw` statements.
*   `TryStatement`: `try-catch-finally` blocks.

### Expressions

Nodes representing various types of expressions:

*   `Identifier`: A name (e.g., variable name, function name).
*   `Literal`: A literal value (string, number, boolean, etc.).
*   `ArrayLiteral`: Array creation `[1, 2, 3]`.
*   `ObjectLiteral`: Object creation `{ key: value }`.
*   `FunctionExpression`: Anonymous function expressions.
*   `ArrowFunctionExpression`: Arrow function expressions `() => {}`.
*   `UnaryExpression`: Unary operations `!a`, `typeof b`.
*   `BinaryExpression`: Binary operations `a + b`.
*   `ConditionalExpression`: Ternary operations `a ? b : c`.
*   `UpdateExpression`: Update operations `a++`, `--b`.
*   `MemberExpression`: Member access `obj.prop`, `arr[index]`.
*   `CallExpression`: Function calls `func()`.
*   `NewExpression`: Object instantiation `new Class()`.
*   `ThisExpression`: The `this` keyword.
*   `SuperExpression`: The `super` keyword.
*   `TemplateLiteral`: Template literals `` `hello ${name}` ``.

### JML Specific Nodes

*   `ComponentBodyElement`: Marker interface for nodes that can appear within a component's body.
*   `ComponentElement`: Represents a JML component instantiation `<MyComponent>`.
*   `ComponentProperty`: Represents a property assignment within a component body `prop: "value"`.
*   `ForLoop`: JML-specific `for` loop for iterating over collections.
*   `IfCondition`: JML-specific `if-else` block for conditional rendering.

### Type Nodes

*   `TypeAnnotation`: Represents a type annotation `: string`.
*   `TypeReference`: A reference to a type (e.g., `string`, `MyType`).
*   `ObjectType`: An object type literal `{ a: string, b: number }`.

## Factory Methods (`factory.go`)

The `factory.go` file provides convenience functions (prefixed with `New`) for creating instances of all AST node types. These methods simplify the construction of the AST programmatically.

**Example:**

```go
// NewIdentifier creates a new Identifier node.
func NewIdentifier(pos Position, name string) *Identifier {
	return &Identifier{
		Position: pos,
		Name:     name,
	}
}

// NewVariableDeclaration creates a new VariableDeclaration node.
func NewVariableDeclaration(pos Position, kind string, declarations []*VariableDeclarator) *VariableDeclaration {
	return &VariableDeclaration{
		Position:     pos,
		Kind:         kind,
		Declarations: declarations,
	}
}
```

## Visitor Pattern (`visitor.go`, `walk.go`)

The `ast` package implements the visitor pattern for traversing the AST.

### `Visitor` Interface (`visitor.go`)

Defines methods for visiting each type of AST node.

```go
type Visitor interface {
	VisitProgram(n *Program)
	VisitDocument(n *Document)
	// ... (many more Visit methods for each node type)
}
```

### `BaseVisitor` (`visitor.go`)

A no-op implementation of the `Visitor` interface. This can be embedded in custom visitors to avoid implementing all methods, allowing developers to focus only on the node types they care about.

```go
type BaseVisitor struct{}

// Example:
func (v *BaseVisitor) VisitProgram(n *Program)   {}
func (v *BaseVisitor) VisitDocument(n *Document) {}
// ...
```

### `Walk` Function (`walk.go`)

The `Walk` function is responsible for traversing the AST. It takes a `Visitor` and a `Node` as input, calls the appropriate `Visit` method on the visitor for the current node, and then recursively calls `Walk` for all child nodes.

```go
func Walk(v Visitor, node Node) {
	if node == nil {
		return
	}

	node.Accept(v) // Calls the specific Visit method on the visitor

	switch n := node.(type) {
	case *Program:
		for _, doc := range n.Documents {
			Walk(v, doc) // Recursively walk children
		}
	// ... (logic for other node types)
	}
}
```

**Usage Example:**

To create a custom visitor that counts the number of `Identifier` nodes:

```go
type IdentifierCounter struct {
	ast.BaseVisitor
	count int
}

func (v *IdentifierCounter) VisitIdentifier(n *ast.Identifier) {
	v.count++
}

// In your code:
// counter := &IdentifierCounter{}
// ast.Walk(counter, yourASTNode)
// fmt.Printf("Total identifiers: %d\n", counter.count)
```

## Printer (`printer.go`)

The `Printer` is an implementation of the `Visitor` interface that prints a human-readable representation of the AST to an `io.Writer`. It's useful for debugging and understanding the structure of a parsed JML file.

```go
type Printer struct {
	ast.BaseVisitor
	writer io.Writer
	indent int
}

// NewPrinter creates a new Printer.
func NewPrinter(writer io.Writer) *Printer {
	return &Printer{writer: writer}
}

// Print prints the given node to the writer.
func (p *Printer) Print(node ast.Node) {
	ast.Walk(p, node)
}
```

**Usage Example:**

```go
// Assuming 'myAST' is an *ast.Program or *ast.Document
// var buf bytes.Buffer
// printer := ast.NewPrinter(&buf)
// printer.Print(myAST)
// fmt.Println(buf.String())
```