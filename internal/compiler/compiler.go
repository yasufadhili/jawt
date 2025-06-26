package compiler

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	parser "github.com/yasufadhili/jawt/internal/compiler/parser/generated"
	"github.com/yasufadhili/jawt/internal/error_handler"
)

type Compiler struct {
	parser   *Parser
	FileType string
}

func NewCompiler(fileType string) (*Compiler, error) {
	if fileType != "Component" && fileType != "Page" {
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}
	return &Compiler{
		FileType: fileType,
	}, nil
}

func (c *Compiler) Compile() (string, error) {

	return "", nil
}

// parseFile handles the actual parsing with custom error handling
func (c *Compiler) parseFile(input antlr.CharStream) ParseResult {
	// Reset parser state for a new file
	c.parser.errorListener.Reset()
	c.parser.errorStrategy.Reset()

	lexer := parser.NewJMLLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(c.parser.errorListener)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewJMLParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(c.parser.errorListener)
	p.SetErrorHandler(c.parser.errorStrategy)

	tree := p.JmlDocument()

	return ParseResult{
		Success: !c.parser.errorListener.HasErrors(),
		Errors:  c.parser.errorListener.GetErrors(),
		Tree:    tree,
	}
}

// ParseResult holds the result of parsing with error information
type ParseResult struct {
	Success bool
	Errors  []error_handler.SyntaxError
	Tree    antlr.ParseTree
}

type Parser struct {
	errorListener *error_handler.SyntaxErrorListener
	errorStrategy *error_handler.ErrorStrategy
}

// newParser creates a new parser with error handling
func newParser() *Parser {
	return &Parser{
		errorListener: error_handler.NewSyntaxErrorListener(),
		errorStrategy: error_handler.NewErrorStrategy(3), // Allow up to 3 recovery attempts
	}
}
