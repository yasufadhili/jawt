package page_compiler

import (
	"github.com/antlr4-go/antlr/v4"
	parser "github.com/yasufadhili/jawt/internal/pc/parser/generated"
)

type AstBuilder struct {
	*parser.BaseJMLPageVisitor
}

func NewAstBuilder() *AstBuilder {
	return &AstBuilder{
		BaseJMLPageVisitor: &parser.BaseJMLPageVisitor{},
	}
}

func (ab *AstBuilder) Visit(tree antlr.ParseTree) interface{} {

	switch ctx := tree.(type) {
	case *parser.PageContext:
		return ab.VisitPage(ctx)
	case *parser.DoctypeSpecifierContext:
		return ab.VisitDoctypeSpecifier(ctx)
	case *parser.ImportStatementContext:
		return ab.VisitImportStatement(ctx)
	case *parser.PageDefinitionContext:
		return ab.VisitPageDefinition(ctx)
	case *parser.PagePropertyContext:
		return ab.VisitPageProperty(ctx)
	case *parser.LiteralContext:
		return ab.VisitLiteral(ctx)
	case *parser.PropertyValueContext:
		return ab.VisitPropertyValue(ctx)
	}

	return nil // For any unhandled types
}
