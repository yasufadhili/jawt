package pc

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
)

type SyntaxError struct {
	Line            int
	Column          int
	Message         string
	OffendingSymbol interface{} // THe token or AST node that caused the error
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("Syntax error at line %d:%d:  %s", e.Line, e.Column, e.Message)
}

type ErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []error
}

func NewErrorListener() *ErrorListener {
	return &ErrorListener{
		DefaultErrorListener: antlr.NewDefaultErrorListener(),
		Errors:               []error{},
	}
}

// SyntaxError is called by the parser when a syntax error is found.
func (l *ErrorListener) SyntaxError(recogniser antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.Errors = append(l.Errors, &SyntaxError{
		Line:            line,
		Column:          column,
		Message:         msg,
		OffendingSymbol: offendingSymbol,
	})
}

// ReportAmbiguity, ReportAttemptingFullContext, ReportContextSensitivity
// can also be overridden if we want to capture more detailed parser diagnostics.
// For most cases, SyntaxError will be enough for reporting parse failures.
