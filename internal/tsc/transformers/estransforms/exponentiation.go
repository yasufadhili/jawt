package estransforms

import (
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/printer"
	"github.com/yasufadhili/jawt/internal/tsc/transformers"
)

type exponentiationTransformer struct {
	transformers.Transformer
}

func (ch *exponentiationTransformer) visit(node *ast.Node) *ast.Node {
	return node // !!!
}

func newExponentiationTransformer(emitContext *printer.EmitContext) *transformers.Transformer {
	tx := &exponentiationTransformer{}
	return tx.NewTransformer(tx.visit, emitContext)
}
