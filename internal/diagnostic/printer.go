package diagnostic

import (
	"fmt"
	"io"
	"os"
)

// Printer prints diagnostics to a writer.
type Printer struct {
	writer io.Writer
	// TODO: Add more configuration options here, e.g., colourisation, verbosity
}

// NewPrinter creates a new printer with os.Stderr as the default writer.
func NewPrinter() *Printer {
	return &Printer{
		writer: os.Stderr,
	}
}

// NewPrinterWithWriter creates a new printer with a custom writer.
func NewPrinterWithWriter(w io.Writer) *Printer {
	return &Printer{
		writer: w,
	}
}

func (p *Printer) Print(reporter *Reporter) {
	for _, d := range reporter.All() {
		p.PrintDiagnostic(d)
	}
}

func (p *Printer) PrintDiagnostic(d *Diagnostic) {
	switch d.Severity {
	case SeverityError:
		p.printError(d)
	case SeverityWarning:
		p.printWarning(d)
	case SeverityInfo:
		p.printInfo(d)
	}
}

func (p *Printer) printError(d *Diagnostic) {
	fmt.Fprintf(p.writer, "\033[31mError\033[0m [%s]: %s\n", d.Code, d.Message)
	fmt.Fprintf(p.writer, "   at %s:%d:%d\n", d.Pos.File, d.Pos.Line, d.Pos.Column)
}

func (p *Printer) printWarning(d *Diagnostic) {
	fmt.Fprintf(p.writer, "\033[33mWarning\033[0m [%s]: %s\n", d.Code, d.Message)
	fmt.Fprintf(p.writer, "   at %s:%d:%d\n", d.Pos.File, d.Pos.Line, d.Pos.Column)
}

func (p *Printer) printInfo(d *Diagnostic) {
	fmt.Fprintf(p.writer, "\033[34mInfo\033[0m [%s]: %s\n", d.Code, d.Message)
	fmt.Fprintf(p.writer, "   at %s:%d:%d\n", d.Pos.File, d.Pos.Line, d.Pos.Column)
}
