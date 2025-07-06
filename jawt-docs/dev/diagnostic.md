# Diagnostic Reporting (`internal/diagnostic`)

The `internal/diagnostic` package provides a standardized way to report and manage messages related to compilation, parsing, and semantic analysis. These messages, known as diagnostics, can represent errors, warnings, or informational notes, and are crucial for providing feedback to the user during the development process.

## Core Concepts

*   **Diagnostic**: A single message indicating an issue or information, including its code, message, severity, origin, and precise location in the source code.
*   **Severity**: Defines the impact level of a diagnostic (Info, Warning, Error).
*   **Position**: Specifies the exact location (file, line, column, byte offsets) within a source file where a diagnostic occurred.
*   **Reporter**: A central collector for all diagnostics generated during a process (e.g., compilation).
*   **Printer**: Responsible for formatting and outputting diagnostics to a user-friendly interface (e.g., console).

## Key Data Structures

### `Diagnostic`

Represents a single diagnostic message.

```go
type Diagnostic struct {
	Code     DiagnosticCode
	Message  string
	Pos      Position
	Severity Severity
	Origin   string
}
```

### `Position`

Represents a precise location in a source file.

```go
type Position struct {
	Line   int
	Column int
	Start  int // 0-based byte offset of the start of the diagnostic
	End    int // 0-based byte offset of the end of the diagnostic
	File   string
}
```

### `Severity`

An enumeration for the severity level of a diagnostic.

```go
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityError
)
```

### `DiagnosticCode`

A unique code for a diagnostic message, allowing for programmatic identification and filtering of specific issues.

```go
type DiagnosticCode string
```

### `Reporter`

Collects and manages a list of diagnostics. It provides methods to add new diagnostics and retrieve them based on their severity.

```go
type Reporter struct {
	mu          sync.Mutex
	diagnostics []*Diagnostic
}
```

### `Printer`

Formats and outputs diagnostics to an `io.Writer` (e.g., `os.Stderr`). It applies color coding based on severity for better readability.

```go
type Printer struct {
	writer io.Writer
}
```

### `AntlrErrorListener`

A custom ANTLR error listener that integrates ANTLR parsing errors directly into the JAWT diagnostic system.

```go
type AntlrErrorListener struct {
	*antlr.DefaultErrorListener
	Reporter *Reporter
	File     string
}
```

## Functions & Methods

### `NewDiagnostic`

```go
func NewDiagnostic(code DiagnosticCode, message string, pos Position, severity Severity, origin string) *Diagnostic
```

Creates a new `Diagnostic` instance.

### `(*Diagnostic) Error()`

```go
func (d *Diagnostic) Error() string
```

Returns a formatted string representation of the diagnostic, making `Diagnostic` compatible with the `error` interface.

### `NewReporter`

```go
func NewReporter() *Reporter
```

Creates a new, empty `Reporter`.

### `(*Reporter) Add(d *Diagnostic)`

Adds a new diagnostic to the reporter's collection.

### `(*Reporter) All()`, `(*Reporter) Errors()`, `(*Reporter) Warnings()`, `(*Reporter) Infos()`

Methods to retrieve diagnostics based on their severity.

### `(*Reporter) HasErrors()`, `(*Reporter) HasWarnings()`

Convenience methods to quickly check if any errors or warnings have been reported.

### `(*Reporter) Reset()`

Clears all diagnostics from the reporter.

### `NewPrinter()`, `NewPrinterWithWriter(w io.Writer)`

Create new `Printer` instances. `NewPrinter()` uses `os.Stderr` by default.

### `(*Printer) Print(reporter *Reporter)`

Prints all diagnostics collected by a `Reporter`.

### `(*Printer) PrintDiagnostic(d *Diagnostic)`

Prints a single `Diagnostic` message, applying color coding based on its severity.

### `NewAntlrErrorListener(reporter *Reporter, file string)`

Creates a new `AntlrErrorListener` that will report ANTLR syntax errors to the provided `Reporter`.

### `(*AntlrErrorListener) SyntaxError(...)`

This method is automatically called by ANTLR when a syntax error is encountered during parsing. It constructs a `Diagnostic` from the ANTLR error information and adds it to the `Reporter`.

## Usage Example

```go
// Example of using the diagnostic package in a hypothetical parser:
// import (
//     "github.com/yasufadhili/jawt/internal/diagnostic"
// )

// func parseCode(code string, filename string) error {
//     reporter := diagnostic.NewReporter()

//     // Simulate a syntax error
//     pos := diagnostic.Position{Line: 5, Column: 10, File: filename}
//     diag := diagnostic.NewDiagnostic(
//         "SYNTAX_001",
//         "Unexpected token '}'",
//         pos,
//         diagnostic.SeverityError,
//         "parser",
//     )
//     reporter.Add(diag)

//     // Simulate a warning
//     warnPos := diagnostic.Position{Line: 2, Column: 1, File: filename}
//     warnDiag := diagnostic.NewDiagnostic(
//         "SEMANTIC_002",
//         "Variable 'unusedVar' declared but not used",
//         warnPos,
//         diagnostic.SeverityWarning,
//         "checker",
//     )
//     reporter.Add(warnDiag)

//     if reporter.HasErrors() {
//         printer := diagnostic.NewPrinter()
//         printer.Print(reporter)
//         return fmt.Errorf("parsing failed with errors")
//     }
//     return nil
// }

// func main() {
//     err := parseCode("some invalid code", "main.jml")
//     if err != nil {
//         fmt.Println("Process finished with errors.")
//     }
// }
```
