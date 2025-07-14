package diagnostic

import (
	"fmt"
)

// Diagnostic represents a single diagnostic message (error, warning, or info).
type Diagnostic struct {
	Code     DiagnosticCode
	Message  string
	Pos      Position
	Severity Severity
	Origin   string
}

// Position represents a location in a file, with start and end offsets.
type Position struct {
	Line   int
	Column int
	Start  int // 0-based byte offset of the start of the diagnostic
	End    int // 0-based byte offset of the end of the diagnostic
	File   string
}

// Severity represents the severity of a diagnostic.
type Severity int

const (
	// SeverityInfo is for informational messages.
	SeverityInfo Severity = iota
	// SeverityWarning is for warnings.
	SeverityWarning
	// SeverityError is for errors.
	SeverityError
)

// String returns the string representation of the severity.
func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// DiagnosticCode represents a unique code for a diagnostic message.
type DiagnosticCode string

// NewDiagnostic creates a new diagnostic.
func NewDiagnostic(code DiagnosticCode, message string, pos Position, severity Severity, origin string) *Diagnostic {
	return &Diagnostic{
		Code:     code,
		Message:  message,
		Pos:      pos,
		Severity: severity,
		Origin:   origin,
	}
}

// Error returns the string representation of the diagnostic.
func (d *Diagnostic) Error() string {
	return fmt.Sprintf("[%s] %s:%d:%d: %s (%s) [%s]", d.Severity, d.Pos.File, d.Pos.Line, d.Pos.Column, d.Message, d.Origin, d.Code)
}
