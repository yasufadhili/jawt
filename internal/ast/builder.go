package ast

import (
	"github.com/antlr4-go/antlr/v4"
	parser "github.com/yasufadhili/jawt/internal/compiler/parser/generated"
)

type Builder struct {
	*parser.BaseJMLVisitor
}

func NewAstBuilder() *Builder {
	return &Builder{
		BaseJMLVisitor: &parser.BaseJMLVisitor{},
	}
}

func (b *Builder) Visit(tree antlr.ParseTree) interface{} {
	switch ctx := tree.(type) {
	case *parser.DocumentContext:
		return b.VisitDocument(ctx)
	default:
		return nil
	}
}

func (b *Builder) VisitDocument(ctx *parser.DocumentContext) interface{} {
	doc := &DocumentNode{}

	return doc
}
