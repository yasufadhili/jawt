# The Abstract Syntax Tree (AST)

The `ast` package is where the magic begins. After the parser chews on a JML file, it spits out a tree-like structure that represents the code. That's the Abstract Syntax Tree, or AST. This package defines all the different types of nodes that can exist in that tree.

## The Core Idea

-	**`Node`**: This is the granddaddy of all nodes. Everything in the AST is a `Node`. It knows its position in the source file (line, column) and can be visited by a `Visitor`.
-	**`Statement`**: A statement is a complete instruction, like a variable declaration or an `if` block.
-	**`Declaration`**: A special kind of statement that creates something new, like a variable, function, or component.
-	**`Expression`**: A piece of code that results in a value, like `2 + 2` or `user.name`.
-	**The Visitor Pattern**: This is a clean way to walk through the AST and do things with the nodes without having to change the node structures themselves. It's super useful for things like type checking and code generation.

## The Key Data Structures

### `Position`

This is just a simple struct that holds the line, column, and file of a node.

```go
type Position struct {
	Line   int
	Column int
	File   string
}
```

### `Node` Interface

This is the base interface for everything in the AST.

```go
type Node interface {
	Pos() Position
	String() string
	Accept(v Visitor)
}
```

### `Program`

This represents the entire JAWT project, which is basically a collection of `Document` nodes.

```go
type Program struct {
	Position
	Documents []*Document
}
```

### `Document`

A `Document` is a single JML file, either a page or a component.

```go
type Document struct {
	Position
	DocType    DocType // "page" or "component"
	Name       *Identifier
	Body       []Statement // This holds all the good stuff: imports, declarations, elements, etc.
	SourceFile string
}
```

### Declarations

These are all the nodes that represent declarations:

-	`ImportDeclaration`: `import { a } from 'b'`
-	`ExportDeclaration`: `export const a = 1`
-	`VariableDeclaration`: `let`, `const`, or `var`.
-	And so on for functions, classes, interfaces, etc.

### Statements

These are the nodes for various statements:

-	`BlockStatement`: A block of statements in curly braces `{ ... }`.
-	`IfStatement`: An `if-else` statement.
-	`ForStatement`: A `for` loop.
-	And so on for `while`, `return`, `throw`, etc.

### Expressions

These are the nodes for expressions:

-	`Identifier`: A name, like a variable or function name.
-	`Literal`: A value, like a string, number, or boolean.
-	`BinaryExpression`: An operation with two operands, like `a + b`.
-	And so on for function calls, member access, etc.

### JML-Specific Nodes

These are nodes that are unique to JML:

-	`ComponentElement`: Represents a JML component being used, like `<MyComponent>`.
-	`ComponentProperty`: A property being assigned inside a component, like `prop: "value"`.
-	`ForLoop`: The JML `for` loop for iterating over collections.
-	`IfCondition`: The JML `if-else` block for conditional rendering.

## The Factory (`factory.go`)

`factory.go` is a set of helper functions for creating new AST nodes. It just makes the code in the `AstBuilder` cleaner and easier to read.

**Example:**

```go
// NewIdentifier creates a new Identifier node.
func NewIdentifier(pos Position, name string) *Identifier {
	return &Identifier{
		Position: pos,
		Name:     name,
	}
}
```

## The Visitor Pattern (`visitor.go`, `walk.go`)

This is how we traverse the AST.

### `Visitor` Interface (`visitor.go`)

This interface defines a `Visit` method for every single type of node in the AST.

```go
type Visitor interface {
	VisitProgram(n *Program)
	VisitDocument(n *Document)
	// ... and many more
}
```

### `BaseVisitor` (`visitor.go`)

This is a helper that implements the `Visitor` interface with empty methods. This way, when I create a new visitor, I can just embed `BaseVisitor` and only implement the `Visit` methods for the nodes I actually care about.

### `Walk` Function (`walk.go`)

The `Walk` function is the engine of the visitor pattern. It takes a `Visitor` and a `Node`, calls the right `Visit` method on the visitor, and then recursively calls `Walk` on all the children of that node.

**Example:**

To count all the identifiers in the AST:

```go
type IdentifierCounter struct {
	ast.BaseVisitor
	count int
}

func (v *IdentifierCounter) VisitIdentifier(n *ast.Identifier) {
	v.count++
}

// Later in the code:
// counter := &IdentifierCounter{}
// ast.Walk(counter, yourASTNode)
// fmt.Printf("Found %d identifiers\n", counter.count)
```

## The Printer (`printer.go`)

The `Printer` is a special visitor that prints out a human-readable representation of the AST. It's incredibly useful for debugging the parser and seeing the structure of a JML file.

**Example:**

```go
// var buf bytes.Buffer
// printer := ast.NewPrinter(&buf)
// printer.Print(myAST)
// fmt.Println(buf.String())
```