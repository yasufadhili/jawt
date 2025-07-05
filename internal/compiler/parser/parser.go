package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/yasufadhili/jawt/internal/ast"
	"github.com/yasufadhili/jawt/internal/compiler"
	parser "github.com/yasufadhili/jawt/internal/compiler/parser/generated"
	"github.com/yasufadhili/jawt/internal/diagnostic"
)

// ParseFile parses a JML file and returns the document AST.
func ParseFile(file string, reporter *diagnostic.Reporter) (*ast.Document, error) {
	input, err := antlr.NewFileStream(file)
	if err != nil {
		return nil, err
	}

	lexer := parser.NewJmlLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewJmlParser(stream)

	p.RemoveErrorListeners()
	p.AddErrorListener(diagnostic.NewAntlrErrorListener(reporter, file))

	tree := p.Document()
	visitor := compiler.NewAstBuilder(reporter, file)
	doc := visitor.Visit(tree).(*ast.Document)

	return doc, nil
}
