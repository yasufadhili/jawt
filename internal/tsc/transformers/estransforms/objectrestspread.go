package estransforms

import (
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/printer"
	"github.com/yasufadhili/jawt/internal/tsc/transformers"
)

type objectRestSpreadTransformer struct {
	transformers.Transformer
}

func (ch *objectRestSpreadTransformer) visit(node *ast.Node) *ast.Node {
	return node // !!!
}

func newObjectRestSpreadTransformer(emitContext *printer.EmitContext) *transformers.Transformer {
	tx := &objectRestSpreadTransformer{}
	return tx.NewTransformer(tx.visit, emitContext)
}
