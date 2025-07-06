package compiler

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/yasufadhili/jawt/internal/ast"
	parser "github.com/yasufadhili/jawt/internal/compiler/parser/generated"
	"github.com/yasufadhili/jawt/internal/core"
	"github.com/yasufadhili/jawt/internal/diagnostic"
)

// Compiler compiles JML files
type Compiler struct {
	ctx *core.JawtContext
}

// NewCompiler creates a new compiler.
func NewCompiler(ctx *core.JawtContext) *Compiler {
	return &Compiler{
		ctx: ctx,
	}
}

// Compile compiles a single JML file.
func (c *Compiler) Compile(file string, reporter *diagnostic.Reporter) (*ast.Document, error) {
	tree, err := parseFile(file, reporter)
	if err != nil {
		return nil, err
	}

	// Build the AST from the parse tree
	builder := NewAstBuilder(reporter, file)
	doc := builder.Visit(tree).(*ast.Document)

	return doc, nil
}

// parseFile parses a JML file and returns the ANTLR parse tree.
func parseFile(file string, reporter *diagnostic.Reporter) (antlr.Tree, error) {
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

	return tree, nil
}

// AstBuilder is a visitor that builds the AST from the ANTLR parse tree.
type AstBuilder struct {
	*parser.BaseJmlVisitor
	reporter *diagnostic.Reporter
	file     string
}

// NewAstBuilder creates a new AstBuilder.
func NewAstBuilder(reporter *diagnostic.Reporter, file string) *AstBuilder {
	return &AstBuilder{
		BaseJmlVisitor: &parser.BaseJmlVisitor{},
		reporter:       reporter,
		file:           file,
	}
}

// VisitDocument is called when visiting the document rule.
func (b *AstBuilder) VisitDocument(ctx *parser.DocumentContext) interface{} {
	// Placeholder for actual AST construction
	return &ast.Document{}
}

// Visit is the main entry point for the visitor.
func (b *AstBuilder) Visit(tree antlr.Tree) interface{} {
	// return tree.Accept(b)
	return nil
}
