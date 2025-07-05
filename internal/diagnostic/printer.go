package diagnostic

import (
	"fmt"
	"io"
	"os"
)

// Printer prints diagnostics to a writer.
type Printer struct {
	writer io.Writer
}

// NewPrinter creates a new printer.
func NewPrinter() *Printer {
	return &Printer{
		writer: os.Stderr,
	}
}

// Print prints all diagnostics from a reporter.
func (p *Printer) Print(reporter *Reporter) {
	for _, err := range reporter.Errors() {
		p.printError(err)
	}

	for _, warn := range reporter.Warnings() {
		p.printWarning(warn)
	}

	for _, info := range reporter.Infos() {
		p.printInfo(info)
	}
}

// printError prints a single error.
func (p *Printer) printError(err *Error) {
	fmt.Fprintf(p.writer, "\033[31mError\033[0m: %s\n", err.Message)
	fmt.Fprintf(p.writer, "   at %s:%d:%d\n", err.Pos.File, err.Pos.Line, err.Pos.Column)
}

// printWarning prints a single warning.
func (p *Printer) printWarning(warn *Error) {
	fmt.Fprintf(p.writer, "\033[33mWarning\033[0m: %s\n", warn.Message)
	fmt.Fprintf(p.writer, "   at %s:%d:%d\n", warn.Pos.File, warn.Pos.Line, warn.Pos.Column)
}

// printInfo prints a single informational message.
func (p *Printer) printInfo(info *Error) {
	fmt.Fprintf(p.writer, "\033[34mInfo\033[0m: %s\n", info.Message)
	fmt.Fprintf(p.writer, "   at %s:%d:%d\n", info.Pos.File, info.Pos.Line, info.Pos.Column)
}
