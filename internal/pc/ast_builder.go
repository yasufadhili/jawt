package pc

import parser "github.com/yasufadhili/jawt/internal/pc/parser/generated"

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
