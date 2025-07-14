package compiler

import (
	"github.com/yasufadhili/jawt/internal/ast"
	parser "github.com/yasufadhili/jawt/internal/compiler/parser/generated"
	"github.com/yasufadhili/jawt/internal/diagnostic"
)

type AstBuilder struct {
	*parser.BaseJmlVisitor
	reporter *diagnostic.Reporter
	file     string
}

func NewAstBuilder(file string, reporter *diagnostic.Reporter) *AstBuilder {
	return &AstBuilder{
		BaseJmlVisitor: &parser.BaseJmlVisitor{},
		reporter:       reporter,
		file:           file,
	}
}

func (b *AstBuilder) VisitDocument(ctx *parser.DocumentContext) interface{} {

	return &ast.Document{}
}
