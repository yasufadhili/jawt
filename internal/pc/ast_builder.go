package pc

import (
	"github.com/antlr4-go/antlr/v4"
	parser "github.com/yasufadhili/jawt/internal/pc/parser/generated"
	"strconv"
	"strings"
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
	case *parser.ProgramContext:
		return ab.VisitProgram(ctx)
	case *parser.DoctypeSpecifierContext:
		return ab.VisitDoctypeSpecifier(ctx)
	case *parser.ImportStatementContext:
		return ab.VisitImportStatement(ctx)
	case *parser.PageContext:
		return ab.VisitPage(ctx)
	case *parser.PagePropertyContext:
		return ab.VisitPageProperty(ctx)
	case *parser.LiteralContext:
		return ab.VisitLiteral(ctx)
	case *parser.PropertyValueContext:
		return ab.VisitPropertyValue(ctx)
	default:
		// For any unhandled types, return nil
		return nil
	}
}

func (ab *AstBuilder) VisitProgram(ctx *parser.ProgramContext) interface{} {
	p := &Program{}

	// Visit doctypeSpecifier
	if ctx.DoctypeSpecifier() != nil {
		p.Doctype = ctx.DoctypeSpecifier().Accept(ab).(*DoctypeSpecifier)
	}

	// Visit imports
	if ctx.Imports() != nil {
		importsCtx := ctx.Imports().(*parser.ImportsContext)
		for _, impCtx := range importsCtx.AllImportStatement() {
			p.Imports = append(p.Imports, impCtx.Accept(ab).(*ImportStatement))
		}
	}

	// Visit page
	if ctx.Page() != nil {
		p.Page = ctx.Page().Accept(ab).(*Page)
	}

	return p
}

func (ab *AstBuilder) VisitDoctypeSpecifier(ctx *parser.DoctypeSpecifierContext) interface{} {
	return &DoctypeSpecifier{
		Doctype: ctx.Doctype().GetText(),
		Name:    ctx.IDENTIFIER().GetText(),
	}
}

func (ab *AstBuilder) VisitImportStatement(ctx *parser.ImportStatementContext) interface{} {
	// Remove quotes from the string literal
	fromPath := strings.Trim(ctx.STRING().GetText(), `"`)

	return &ImportStatement{
		Doctype:    ctx.Doctype().GetText(),
		Identifier: ctx.IDENTIFIER().GetText(),
		From:       fromPath,
	}
}

func (ab *AstBuilder) VisitPage(ctx *parser.PageContext) interface{} {
	p := &Page{}

	if ctx.PageBody() != nil {
		bodyCtx := ctx.PageBody().(*parser.PageBodyContext)
		for _, propCtx := range bodyCtx.AllPageProperty() {
			p.Properties = append(p.Properties, ab.Visit(propCtx).(*PageProperty))
		}
	}

	return p
}

func (ab *AstBuilder) VisitPageProperty(ctx *parser.PagePropertyContext) interface{} {
	val := ctx.PropertyValue().Accept(ab)
	return &PageProperty{
		Key:   ctx.IDENTIFIER().GetText(),
		Value: val,
	}
}

func (ab *AstBuilder) VisitLiteral(ctx *parser.LiteralContext) interface{} {
	if ctx.INTEGER() != nil {
		val, err := strconv.Atoi(ctx.INTEGER().GetText())
		if err != nil {
			// Handle parsing error appropriately (log, return error, etc.)
			return 0
		}
		return val
	} else if ctx.STRING() != nil {
		// Remove quotes from the string literal
		str := ctx.STRING().GetText()
		if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
			str = str[1 : len(str)-1]
		}
		return str
	} else if ctx.IDENTIFIER() != nil {
		return ctx.IDENTIFIER().GetText() // Treat as raw identifier for now
	}
	return nil
}

// VisitPropertyValue just delegates to VisitLiteral
func (ab *AstBuilder) VisitPropertyValue(ctx *parser.PropertyValueContext) interface{} {
	return ab.Visit(ctx.Literal())
}
