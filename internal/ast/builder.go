package ast

import (
	"github.com/antlr4-go/antlr/v4"
	parser "github.com/yasufadhili/jawt/internal/compiler/parser/generated"
	"strconv"
	"strings"
)

type Builder struct {
	*parser.BaseJMLVisitor
}

func NewBuilder() *Builder {
	return &Builder{
		BaseJMLVisitor: &parser.BaseJMLVisitor{},
	}
}

func (b *Builder) Visit(tree antlr.ParseTree) interface{} {

	switch ctx := tree.(type) {
	case *parser.DocumentContext:
		return b.VisitDocument(ctx)
	case *parser.DocumentContentContext:
		return b.VisitDocumentContent(ctx)
	case *parser.DoctypeSpecifierContext:
		return b.VisitDoctypeSpecifier(ctx)
	case *parser.ImportStatementContext:
		return b.VisitImportStatement(ctx)
	case *parser.PageDefinitionContext:
		return b.VisitPageDefinition(ctx)
	case *parser.ComponentDefinitionContext:
		return b.VisitComponentDefinition(ctx)
	case *parser.ComponentElementContext:
		return b.VisitComponentElement(ctx)
	case *parser.ComponentPropertyContext:
		return b.VisitComponentProperty(ctx)
	case *parser.PagePropertyContext:
		return b.VisitPageProperty(ctx)
	case *parser.PropertyValueContext:
		return b.VisitPropertyValue(ctx)
	case *parser.LiteralContext:
		return b.VisitLiteral(ctx)
	default:
		return nil
	}
}

func (b *Builder) VisitDocument(ctx *parser.DocumentContext) interface{} {
	doc := &DocumentNode{}

	if ctx.DoctypeSpecifier() != nil {
		if doctype := b.Visit(ctx.DoctypeSpecifier()); doctype != nil {
			doc.Doctype = doctype.(*DoctypeSpecifierNode)
		}
	}

	if ctx.Imports() != nil {
		importsCtx := ctx.Imports().(*parser.ImportsContext)
		for _, importStmt := range importsCtx.AllImportStatement() {
			if imp := b.Visit(importStmt); imp != nil {
				doc.Imports = append(doc.Imports, imp.(*ImportStatementNode))
			}
		}
	}

	if ctx.DocumentContent() != nil {
		if content := b.Visit(ctx.DocumentContent()); content != nil {
			doc.Content = content.(DocumentContentNode)
		}
	}

	return doc
}

func (b *Builder) VisitDoctypeSpecifier(ctx *parser.DoctypeSpecifierContext) interface{} {
	doctype := &DoctypeSpecifierNode{}

	// Get doctype ("page", "component", "module")
	if ctx.Doctype() != nil {
		doctype.Doctype = ctx.Doctype().GetText()
	}

	if ctx.IDENTIFIER() != nil {
		doctype.Name = ctx.IDENTIFIER().GetText()
	}

	return doctype
}

func (b *Builder) VisitImportStatement(ctx *parser.ImportStatementContext) interface{} {
	imp := &ImportStatementNode{}

	//if ctx.GetText() == "importbrowser" { // This matches the grammar pattern
	//	imp.IsBrowser = true
	//	imp.Type = "browser"
	//	return imp
	//}

	// Regular import: import doctype IDENTIFIER from STRING
	if ctx.Doctype() != nil {
		imp.Type = ctx.Doctype().GetText()
	}

	if ctx.IDENTIFIER() != nil {
		imp.Identifier = ctx.IDENTIFIER().GetText()
	}

	if ctx.STRING() != nil {
		// Remove quotes from string literal
		str := ctx.STRING().GetText()
		imp.From = strings.Trim(str, `"'`)
	}

	return imp
}

func (b *Builder) VisitDocumentContent(ctx *parser.DocumentContentContext) interface{} {
	if ctx.PageDefinition() != nil {
		return b.Visit(ctx.PageDefinition())
	}
	if ctx.ComponentDefinition() != nil {
		return b.Visit(ctx.ComponentDefinition())
	}
	// TODO: Module Definition
	return nil
}

func (b *Builder) VisitPageDefinition(ctx *parser.PageDefinitionContext) interface{} {
	page := &PageDefinitionNode{}

	if ctx.PageBody() != nil {
		pageBodyCtx := ctx.PageBody().(*parser.PageBodyContext)

		if pageBodyCtx.PageProperties() != nil {
			propsCtx := pageBodyCtx.PageProperties().(*parser.PagePropertiesContext)
			for _, propCtx := range propsCtx.AllPageProperty() {
				if prop := b.Visit(propCtx); prop != nil {
					page.Properties = append(page.Properties, prop.(*PropertyNode))
				}
			}
		}

		if pageBodyCtx.SingleComponentChild() != nil {
			childCtx := pageBodyCtx.SingleComponentChild().(*parser.SingleComponentChildContext)
			if childCtx.ComponentElement() != nil {
				if child := b.Visit(childCtx.ComponentElement()); child != nil {
					page.Child = child.(*ComponentElementNode)
				}
			}
		}
	}

	return page
}

func (b *Builder) VisitComponentDefinition(ctx *parser.ComponentDefinitionContext) interface{} {
	comp := &ComponentDefinitionNode{}

	if ctx.ComponentElement() != nil {
		if element := b.Visit(ctx.ComponentElement()); element != nil {
			comp.Element = element.(*ComponentElementNode)
		}
	}

	return comp
}

func (b *Builder) VisitComponentElement(ctx *parser.ComponentElementContext) interface{} {
	element := &ComponentElementNode{}

	if ctx.IDENTIFIER() != nil {
		element.Name = ctx.IDENTIFIER().GetText()
	}

	if ctx.ComponentBlock() != nil {
		blockCtx := ctx.ComponentBlock().(*parser.ComponentBlockContext)
		if blockCtx.ComponentBody() != nil {
			bodyCtx := blockCtx.ComponentBody().(*parser.ComponentBodyContext)

			for _, propCtx := range bodyCtx.AllComponentProperty() {
				if prop := b.Visit(propCtx); prop != nil {
					element.Properties = append(element.Properties, prop.(*PropertyNode))
				}
			}

			for _, childCtx := range bodyCtx.AllComponentElement() {
				if child := b.Visit(childCtx); child != nil {
					element.Children = append(element.Children, child.(*ComponentElementNode))
				}
			}
		}
	}

	return element
}

func (b *Builder) VisitPageProperty(ctx *parser.PagePropertyContext) interface{} {
	prop := &PropertyNode{}

	if ctx.IDENTIFIER() != nil {
		prop.Name = ctx.IDENTIFIER().GetText()
	}

	if ctx.PropertyValue() != nil {
		if value := b.Visit(ctx.PropertyValue()); value != nil {
			prop.Value = value.(PropertyValueNode)
		}
	}

	return prop
}

func (b *Builder) VisitComponentProperty(ctx *parser.ComponentPropertyContext) interface{} {
	prop := &PropertyNode{}

	if ctx.IDENTIFIER() != nil {
		prop.Name = ctx.IDENTIFIER().GetText()
	}

	if ctx.PropertyValue() != nil {
		if value := b.Visit(ctx.PropertyValue()); value != nil {
			prop.Value = value.(PropertyValueNode)
		}
	}

	return prop
}

func (b *Builder) VisitPropertyValue(ctx *parser.PropertyValueContext) interface{} {
	if ctx.Literal() != nil {
		if literal := b.Visit(ctx.Literal()); literal != nil {
			return literal.(*LiteralNode)
		}
	}
	if ctx.ComponentElement() != nil {
		if element := b.Visit(ctx.ComponentElement()); element != nil {
			return element.(*ComponentElementNode)
		}
	}
	// Expression, array literal, and object literal will be handled here
	// when those node types are implemented
	return nil
}

func (b *Builder) VisitLiteral(ctx *parser.LiteralContext) interface{} {
	literal := &LiteralNode{}

	switch {
	case ctx.INTEGER() != nil:
		literal.Type = "integer"
		if val, err := strconv.Atoi(ctx.INTEGER().GetText()); err == nil {
			literal.Value = val
		}
	case ctx.FLOAT() != nil:
		literal.Type = "float"
		if val, err := strconv.ParseFloat(ctx.FLOAT().GetText(), 64); err == nil {
			literal.Value = val
		}
	case ctx.STRING() != nil:
		literal.Type = "string"
		// Remove quotes from string literal
		str := ctx.STRING().GetText()
		literal.Value = strings.Trim(str, `"'`)
	case ctx.BOOLEAN() != nil:
		literal.Type = "boolean"
		literal.Value = ctx.BOOLEAN().GetText() == "true"
	case ctx.NULL() != nil:
		literal.Type = "null"
		literal.Value = nil
	}

	return literal
}
