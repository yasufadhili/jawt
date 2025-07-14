package compiler

import (
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/yasufadhili/jawt/internal/ast"
	"github.com/yasufadhili/jawt/internal/compiler/parser"
	"github.com/yasufadhili/jawt/internal/core"
	"github.com/yasufadhili/jawt/internal/diagnostic"
)

type Compiler struct {
	ctx *core.JawtContext
}

func NewCompiler(ctx *core.JawtContext) *Compiler {
	return &Compiler{
		ctx: ctx,
	}
}

// Compile compiles a single JML file.
func (c *Compiler) Compile(file string, reporter *diagnostic.Reporter) (*ast.Document, error) {
	input, err := antlr.NewFileStream(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read JML file %s: %w", file, err)
	}

	lexer := parser.NewJmlLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := parser.NewJmlParser(stream)

	// Remove default error listeners and add our custom one
	parser.RemoveErrorListeners()
	parser.AddErrorListener(diagnostic.NewAntlrErrorListener(file, reporter))

	// Parse the input
	tree := parser.Document()

	// Build the AST
	builder := NewAstBuilder(file, reporter)
	astDoc := builder.Visit(tree).(*ast.Document)

	return astDoc, nil
}
