package diagnostic

import (
	"fmt"
)

// Error represents a single error in the toolchain.
type Error struct {
	Message  string
	Pos      Position
	Severity Severity
	Origin   string
}

// Position represents a location in a file.
type Position struct {
	Line   int
	Column int
	File   string
}

// Severity represents the severity of an error.
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

// NewError creates a new error.
func NewError(message string, pos Position, severity Severity, origin string) *Error {
	return &Error{
		Message:  message,
		Pos:      pos,
		Severity: severity,
		Origin:   origin,
	}
}

// Error returns the string representation of the error.
func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s:%d:%d: %s (%s)", e.Severity, e.Pos.File, e.Pos.Line, e.Pos.Column, e.Message, e.Origin)
}
