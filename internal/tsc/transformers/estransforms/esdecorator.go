package estransforms

import (
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/printer"
	"github.com/yasufadhili/jawt/internal/tsc/transformers"
)

type esDecoratorTransformer struct {
	transformers.Transformer
}

func (ch *esDecoratorTransformer) visit(node *ast.Node) *ast.Node {
	return node // !!!
}

func newESDecoratorTransformer(emitContext *printer.EmitContext) *transformers.Transformer {
	tx := &esDecoratorTransformer{}
	return tx.NewTransformer(tx.visit, emitContext)
}
