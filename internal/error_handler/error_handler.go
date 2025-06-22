package error_handler

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
)

type SyntaxErrorListener struct {
	antlr.DefaultErrorListener
	Errors []SyntaxError
}

type SyntaxError struct {
	Line    int
	Column  int
	Message string
	Symbol  string
	Context string
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("syntax error at line %d:%d - %s", e.Line, e.Column, e.Message)
}

// SyntaxError is called when a syntax error occurs during parsing
func (l *SyntaxErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// Extract the offending symbol as string
	var symbolText string
	if token, ok := offendingSymbol.(antlr.Token); ok {
		symbolText = token.GetText()
	}

	// Get context from the input stream if available
	var context string
	if lexer, ok := recognizer.(antlr.Lexer); ok {
		input := lexer.GetInputStream()
		if input != nil {
			// Get surrounding context (10 characters before and after)
			start := max(0, column-10)
			end := min(input.Size(), column+10)
			context = input.GetText(start, end)
		}
	}

	syntaxErr := SyntaxError{
		Line:    line,
		Column:  column,
		Message: msg,
		Symbol:  symbolText,
		Context: context,
	}

	l.Errors = append(l.Errors, syntaxErr)
}

// HasErrors returns true if any syntax errors were encountered
func (l *SyntaxErrorListener) HasErrors() bool {
	return len(l.Errors) > 0
}

// GetErrors returns all collected syntax errors
func (l *SyntaxErrorListener) GetErrors() []SyntaxError {
	return l.Errors
}

func (l *SyntaxErrorListener) Reset() {
	l.Errors = []SyntaxError{}
}

// ErrorStrategy implements antlr.ErrorStrategy for recovery strategies
type ErrorStrategy struct {
	antlr.DefaultErrorStrategy
	MaxRecoveryAttempts int
	RecoveryAttempts    int
}

// NewErrorStrategy creates a new error strategy
func NewErrorStrategy(maxAttempts int) *ErrorStrategy {
	return &ErrorStrategy{
		DefaultErrorStrategy: antlr.DefaultErrorStrategy{},
		MaxRecoveryAttempts:  maxAttempts,
		RecoveryAttempts:     0,
	}
}

// Recover attempts to recover from a parsing error
func (s *ErrorStrategy) Recover(recognizer antlr.Parser, e antlr.RecognitionException) {
	s.RecoveryAttempts++

	if s.RecoveryAttempts > s.MaxRecoveryAttempts {
		panic(fmt.Sprintf("Too many recovery attempts (%d), giving up", s.RecoveryAttempts))
	}

	// Call the default recovery mechanism
	s.DefaultErrorStrategy.Recover(recognizer, e)
}

// RecoverInline attempts inline recovery for missing tokens
func (s *ErrorStrategy) RecoverInline(recognizer antlr.Parser) antlr.Token {
	// Try default inline recovery first
	return s.DefaultErrorStrategy.RecoverInline(recognizer)
}

func (s *ErrorStrategy) Reset() {
	s.RecoveryAttempts = 0
}

func NewSyntaxErrorListener() *SyntaxErrorListener {
	return &SyntaxErrorListener{
		Errors:               []SyntaxError{},
		DefaultErrorListener: antlr.DefaultErrorListener{},
	}

}
