package pc

import (
	parser "github.com/yasufadhili/jawt/internal/pc/parser/generated"
	"strconv"
)

type AstBuilder struct {
	*parser.BaseJMLPageVisitor
}

func NewAstBuilder() *AstBuilder {
	return &AstBuilder{
		BaseJMLPageVisitor: &parser.BaseJMLPageVisitor{},
	}
}

func (ab *AstBuilder) VisitProgram(ctx *parser.ProgramContext) interface{} {
	p := &Program{}

	// Visit doctypeSpecifier
	if ctx.DoctypeSpecifier() != nil {
		p.Doctype = ab.Visit(ctx.DoctypeSpecifier()).(*DoctypeSpecifier)
	}

	// Visit imports
	for ctx.Imports() != nil {
		importsCtx := ctx.Imports().(*parser.ImportsContext)
		for _, impCtx := range importsCtx.AllImportStatement() {
			p.Imports = append(p.Imports, ab.Visit(impCtx).(*ImportStatement))
		}
	}

	// Visit page
	if ctx.Page() != nil {
		p.Page = ab.Visit(ctx.Page()).(*Page)
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
	fromPath := ctx.STRING().GetText()
	if len(fromPath) >= 2 && fromPath[0] == '"' && fromPath[len(fromPath)-1] == '"' {
		fromPath = fromPath[1 : len(ctx.STRING().GetText())-1]
	}

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
	value := ab.Visit(ctx.PropertyValue()).(interface{}) // Get the raw value from the literal
	return &PageProperty{
		Key:   ctx.IDENTIFIER().GetText(),
		Value: value,
	}
}

func (ab *AstBuilder) VisitLiteral(ctx *parser.LiteralContext) interface{} {
	if ctx.INTEGER() != nil {
		val, _ := strconv.Atoi(ctx.INTEGER().GetText())
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
