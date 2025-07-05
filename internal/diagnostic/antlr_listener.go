package diagnostic

import (
	"github.com/antlr4-go/antlr/v4"
)

// AntlrErrorListener is a custom error listener for Antlr4.
type AntlrErrorListener struct {
	*antlr.DefaultErrorListener
	Reporter *Reporter
	File     string
}

// NewAntlrErrorListener creates a new AntlrErrorListener.
func NewAntlrErrorListener(reporter *Reporter, file string) *AntlrErrorListener {
	return &AntlrErrorListener{
		Reporter: reporter,
		File:     file,
	}
}

// SyntaxError is called by Antlr4 when a syntax error is encountered.
func (l *AntlrErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	pos := Position{
		Line:   line,
		Column: column,
		File:   l.File,
	}

	err := NewError(msg, pos, SeverityError, "parser")
	l.Reporter.Add(err)
}
