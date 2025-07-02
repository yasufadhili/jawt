package estransforms

import (
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/printer"
	"github.com/yasufadhili/jawt/internal/tsc/transformers"
)

type forawaitTransformer struct {
	transformers.Transformer
}

func (ch *forawaitTransformer) visit(node *ast.Node) *ast.Node {
	return node // !!!
}

func newforawaitTransformer(emitContext *printer.EmitContext) *transformers.Transformer {
	tx := &forawaitTransformer{}
	return tx.NewTransformer(tx.visit, emitContext)
}
