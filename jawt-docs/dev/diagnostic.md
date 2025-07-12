# Diagnostic Reporting (`internal/diagnostic`)

This package is all about handling errors and warnings in a clean, consistent way. When something goes wrong during compilation, parsing, or checking, we need to be able to report it to the user in a way that's easy to understand.

## The Core Ideas

-   **`Diagnostic`**: This is a single error or warning message. It knows what the message is, where it happened (file, line, column), how severe it is, and where it came from.
-   **`Severity`**: This is just an enum for the level of the diagnostic: Info, Warning, or Error.
-   **`Reporter`**: This is a central place to collect all the diagnostics that are generated during a process. For example, the compiler will have a reporter that it passes to the parser and the checker.
-   **`Printer`**: This is what takes the diagnostics from the reporter and prints them to the console in a nice, colorful format.

## The Key Data Structures

### `Diagnostic`

This struct holds all the info for a single diagnostic message.

```go
type Diagnostic struct {
	Code     DiagnosticCode
	Message  string
	Pos      Position
	Severity Severity
	Origin   string
}
```

### `Reporter`

This collects all the diagnostics.

```go
type Reporter struct {
	mu          sync.Mutex
	diagnostics []*Diagnostic
}
```

### `Printer`

This prints the diagnostics to the console.

```go
type Printer struct {
	writer io.Writer
}
```

### `AntlrErrorListener`

This is a custom error listener for ANTLR. When the ANTLR parser finds a syntax error, it calls a method on this listener. I've implemented that method to create a `Diagnostic` and add it to our `Reporter`. This way, we can handle syntax errors from the parser in the same way we handle all other errors.

## How It's Used

When a process like compilation starts, it creates a new `Reporter`. This reporter is then passed down to all the different parts of the process. If any part of the process finds an error, it creates a `Diagnostic` and adds it to the reporter.

At the end of the process, we check if the reporter has any errors. If it does, we create a `Printer` and use it to print all the diagnostics to the console. This gives the user a nice, clean list of all the issues that were found.
