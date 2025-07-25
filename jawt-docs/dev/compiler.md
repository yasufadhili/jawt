# The Compiler (`internal/compiler`)

The `compiler` package is where the raw JML code gets turned into something the rest of the toolchain can understand: the Abstract Syntax Tree (AST). It uses ANTLR4 to do the heavy lifting of parsing the code.

## The Core Ideas

-   **JML Grammar**: The `Jml.g4` file (in the top-level `grammar/` directory) is the rulebook for the JML language. It defines all the valid syntax.
-   **ANTLR4**: This is a tool that takes the `Jml.g4` grammar and generates a Go parser for it. This saves me from having to write the parser by hand, which is a huge pain.
-   **Lexer**: This is the first stage of the parser. It scans the JML code and breaks it down into a stream of tokens (like keywords, identifiers, operators, etc.).
-   **Parser**: This takes the stream of tokens from the lexer and builds a parse tree. A parse tree is a very detailed, concrete representation of the code.
-   **`AstBuilder`**: This is a custom visitor that walks the parse tree and builds our own, more abstract and useful AST (the one defined in the `internal/ast` package).

## The Key Data Structures

### `Compiler`

The main struct for the compiler. It's pretty simple right now.

```go
type Compiler struct {
	ctx *core.JawtContext
}
```

### `AstBuilder`

This is where the real work happens. It's a visitor that walks the ANTLR parse tree and builds our AST.

```go
type AstBuilder struct {
	*parser.BaseJmlVisitor // Embeds the base visitor generated by ANTLR
	reporter *diagnostic.Reporter
	file     string
}
```

## The Process

### `Compile`

This is the main entry point for the compiler. It takes a file path and does two things:

1.  Calls `parseFile` to get the ANTLR parse tree.
2.  Uses `AstBuilder` to walk the parse tree and build our AST.

It returns the `ast.Document` node, which is the root of our AST for that file.

### `parseFile`

This is an internal helper that does the actual parsing with ANTLR. It sets up the lexer and parser and attaches a custom error listener so we can report syntax errors in a nice, user-friendly way.

## The Parser Generation Script (`parser/generate.sh`)

This script is a lifesaver. Whenever I make a change to the `Jml.g4` grammar file, I just run this script, and it automatically regenerates the Go parser code. This makes it super easy to evolve the JML language without having to do a bunch of manual, error-prone work.

**How to use it:**

```bash
./internal/compiler/parser/generate.sh
```

This script uses the `antlr-4.13.2-complete.jar` to generate the Go files in the `internal/compiler/parser/generated/` directory. These generated files are then used by the `compiler.go` to do the parsing.
