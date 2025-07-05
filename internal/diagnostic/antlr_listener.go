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
	var start, end int
	if offendingSymbol != nil {
		if token, ok := offendingSymbol.(antlr.Token); ok {
			start = token.GetStart()
			end = token.GetStop() + 1 // Antlr's Stop is inclusive, so add 1 for exclusive end
		}
	}

	pos := Position{
		Line:   line,
		Column: column,
		Start:  start,
		End:    end,
		File:   l.File,
	}

	diag := NewDiagnostic("SYNTAX_ERROR", msg, pos, SeverityError, "parser")
	l.Reporter.Add(diag)
}
