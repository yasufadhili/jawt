package compiler

import (
	"github.com/yasufadhili/jawt/internal/ast"
	"github.com/yasufadhili/jawt/internal/compiler/parser"
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
	docTypeStr := ctx.DOCTYPE().GetText()
	docType := ast.DocumentTypePage
	if docTypeStr == "_doctype component" {
		docType = ast.DocumentTypeComponent
	}

	name := ctx.IDENTIFIER().GetText()

	// For now, just return a basic document. More complex AST building will go here.
	return &ast.Document{
		Position:   ast.Position{Line: ctx.GetStart().GetLine(), Column: ctx.GetStart().GetColumn(), File: b.file},
		DocType:    docType,
		Name:       ast.NewIdentifier(ast.Position{Line: ctx.IDENTIFIER().GetSymbol().GetLine(), Column: ctx.IDENTIFIER().GetSymbol().GetColumn(), File: b.file}, name),
		SourceFile: b.file,
	}
}