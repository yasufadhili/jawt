package page_compiler

import (
	"github.com/antlr4-go/antlr/v4"
	parser "github.com/yasufadhili/jawt/internal/page_compiler/parser/generated"
	"strconv"
	"strings"
)

type AstBuilder struct {
	*parser.BaseJMLVisitor
}

func NewAstBuilder() *AstBuilder {
	return &AstBuilder{
		BaseJMLVisitor: &parser.BaseJMLVisitor{},
	}
}

func (ab *AstBuilder) Visit(tree antlr.ParseTree) interface{} {
	switch ctx := tree.(type) {
	case *parser.JmlDocumentContext:
		return ab.VisitJmlDocument(ctx)
	case *parser.DocumentContentContext:
		return ab.VisitDocumentContent(ctx)
	case *parser.PageContentContext:
		return ab.VisitPageContent(ctx)
	case *parser.DoctypeSpecifierContext:
		return ab.VisitDoctypeSpecifier(ctx)
	case *parser.ImportsContext:
		return ab.VisitImports(ctx)
	case *parser.ImportStatementContext:
		return ab.VisitImportStatement(ctx)
	case *parser.PageBodyContext:
		return ab.VisitPageBody(ctx)
	case *parser.PagePropertiesContext:
		return ab.VisitPageProperties(ctx)
	case *parser.PagePropertyContext:
		return ab.VisitPageProperty(ctx)
	case *parser.SingleComponentChildContext:
		return ab.VisitSingleComponentChild(ctx)
	case *parser.ComponentElementContext:
		return ab.VisitComponentElement(ctx)
	case *parser.LiteralContext:
		return ab.VisitLiteral(ctx)
	case *parser.PropertyValueContext:
		return ab.VisitPropertyValue(ctx)
	}

	return nil // For any unhandled types
}

func (ab *AstBuilder) VisitJmlDocument(ctx *parser.JmlDocumentContext) interface{} {
	p := &Page{}

	// Handle doctype specifier
	if ctx.DoctypeSpecifier() != nil {
		p.Doctype = ctx.DoctypeSpecifier().Accept(ab).(*DoctypeSpecifier)
	}

	// Handle imports
	if ctx.Imports() != nil {
		importsResult := ctx.Imports().Accept(ab)
		if imports, ok := importsResult.([]*ImportStatement); ok {
			p.Imports = imports
		}
	}

	// Handle document content - extract page content
	if ctx.DocumentContent() != nil {
		contentResult := ctx.DocumentContent().Accept(ab)
		if pageDefinition, ok := contentResult.(*PageDefinition); ok {
			p.PageDefinition = pageDefinition
		}
	}

	return p
}

func (ab *AstBuilder) VisitDocumentContent(ctx *parser.DocumentContentContext) interface{} {
	// For page compiler, we only care about pageContent
	if ctx.PageContent() != nil {
		return ctx.PageContent().Accept(ab)
	}
	return nil
}

func (ab *AstBuilder) VisitPageContent(ctx *parser.PageContentContext) interface{} {
	pd := &PageDefinition{}

	if ctx.PageBody() != nil {
		bodyResult := ctx.PageBody().Accept(ab)
		if pageBody, ok := bodyResult.(*PageBody); ok {
			pd.Properties = pageBody.Properties
			pd.Child = pageBody.Child
		}
	}

	return pd
}

func (ab *AstBuilder) VisitDoctypeSpecifier(ctx *parser.DoctypeSpecifierContext) interface{} {
	return &DoctypeSpecifier{
		Doctype: ctx.Doctype().GetText(),
		Name:    ctx.IDENTIFIER().GetText(),
	}
}

func (ab *AstBuilder) VisitImports(ctx *parser.ImportsContext) interface{} {
	var imports []*ImportStatement
	for _, impCtx := range ctx.AllImportStatement() {
		import_ := impCtx.Accept(ab).(*ImportStatement)
		imports = append(imports, import_)
	}
	return imports
}

func (ab *AstBuilder) VisitImportStatement(ctx *parser.ImportStatementContext) interface{} {
	import_ := &ImportStatement{}

	// Handle a regular import statement with doctype, identifier, and from clause
	if ctx.Doctype() != nil && ctx.IDENTIFIER() != nil && ctx.STRING() != nil {
		import_.Doctype = ctx.Doctype().GetText()
		import_.Identifier = ctx.IDENTIFIER().GetText()
		// Remove quotes from the string literal
		import_.From = strings.Trim(ctx.STRING().GetText(), `"'`)
	} else if ctx.GetText() == "import browser" {
		// Handle browser import
		import_.Doctype = "browser"
		import_.Identifier = "browser"
		import_.From = ""
	}

	return import_
}

func (ab *AstBuilder) VisitPageBody(ctx *parser.PageBodyContext) interface{} {
	pageBody := &PageBody{}

	// Handle page properties
	if ctx.PageProperties() != nil {
		propertiesResult := ctx.PageProperties().Accept(ab)
		if properties, ok := propertiesResult.([]*PageProperty); ok {
			pageBody.Properties = properties
		}
	}

	// Handle single component child
	if ctx.SingleComponentChild() != nil {
		childResult := ctx.SingleComponentChild().Accept(ab)
		if child, ok := childResult.(*ComponentElement); ok {
			pageBody.Child = child
		}
	}

	return pageBody
}

func (ab *AstBuilder) VisitPageProperties(ctx *parser.PagePropertiesContext) interface{} {
	var properties []*PageProperty
	for _, propCtx := range ctx.AllPageProperty() {
		prop := propCtx.Accept(ab).(*PageProperty)
		properties = append(properties, prop)
	}
	return properties
}

func (ab *AstBuilder) VisitPageProperty(ctx *parser.PagePropertyContext) interface{} {
	val := ctx.PropertyValue().Accept(ab)
	return &PageProperty{
		Key:   ctx.IDENTIFIER().GetText(),
		Value: val,
	}
}

func (ab *AstBuilder) VisitSingleComponentChild(ctx *parser.SingleComponentChildContext) interface{} {
	return ctx.ComponentElement().Accept(ab)
}

func (ab *AstBuilder) VisitComponentElement(ctx *parser.ComponentElementContext) interface{} {
	element := &ComponentElement{
		Name: ctx.IDENTIFIER().GetText(),
	}

	// Handle component block if present
	if ctx.ComponentBlock() != nil {
		blockResult := ctx.ComponentBlock().Accept(ab)
		if block, ok := blockResult.(*ComponentBlock); ok {
			element.Block = block
		}
	}

	return element
}

func (ab *AstBuilder) VisitLiteral(ctx *parser.LiteralContext) interface{} {
	if ctx.INTEGER() != nil {
		val, err := strconv.Atoi(ctx.INTEGER().GetText())
		if err != nil {
			// Handle parsing error appropriately (log, return error, etc.)
			return 0
		}
		return val
	} else if ctx.FLOAT() != nil {
		val, err := strconv.ParseFloat(ctx.FLOAT().GetText(), 64)
		if err != nil {
			return 0.0
		}
		return val
	} else if ctx.STRING() != nil {
		// Remove quotes from the string literal
		str := ctx.STRING().GetText()
		if len(str) >= 2 && ((str[0] == '"' && str[len(str)-1] == '"') || (str[0] == '\'' && str[len(str)-1] == '\'')) {
			str = str[1 : len(str)-1]
		}
		return str
	} else if ctx.BOOLEAN() != nil {
		return ctx.BOOLEAN().GetText() == "true"
	} else if ctx.NULL() != nil {
		return nil
	}
	return nil
}

// VisitPropertyValue handles different types of property values
func (ab *AstBuilder) VisitPropertyValue(ctx *parser.PropertyValueContext) interface{} {
	if ctx.Literal() != nil {
		return ab.Visit(ctx.Literal())
	} else if ctx.Expression() != nil {
		return ab.Visit(ctx.Expression())
	} else if ctx.ComponentElement() != nil {
		return ab.Visit(ctx.ComponentElement())
	} else if ctx.ArrayLiteral() != nil {
		return ab.Visit(ctx.ArrayLiteral())
	} else if ctx.ObjectLiteral() != nil {
		return ab.Visit(ctx.ObjectLiteral())
	}
	return nil
}
