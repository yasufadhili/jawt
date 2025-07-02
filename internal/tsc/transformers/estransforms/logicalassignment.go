package estransforms

import (
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/printer"
	"github.com/yasufadhili/jawt/internal/tsc/transformers"
)

type logicalAssignmentTransformer struct {
	transformers.Transformer
}

func (ch *logicalAssignmentTransformer) visit(node *ast.Node) *ast.Node {
	return node // !!!
}

func newLogicalAssignmentTransformer(emitContext *printer.EmitContext) *transformers.Transformer {
	tx := &logicalAssignmentTransformer{}
	return tx.NewTransformer(tx.visit, emitContext)
}
