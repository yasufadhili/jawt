package estransforms

import (
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/printer"
	"github.com/yasufadhili/jawt/internal/tsc/transformers"
)

func convertClassDeclarationToClassExpression(emitContext *printer.EmitContext, node *ast.ClassDeclaration) *ast.Expression {
	updated := emitContext.Factory.NewClassExpression(
		transformers.ExtractModifiers(emitContext, node.Modifiers(), ^ast.ModifierFlagsExportDefault),
		node.Name(),
		node.TypeParameters,
		node.HeritageClauses,
		node.Members,
	)
	emitContext.SetOriginal(updated, node.AsNode())
	updated.Loc = node.Loc
	return updated
}
