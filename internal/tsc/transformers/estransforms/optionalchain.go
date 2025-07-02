package estransforms

import (
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/printer"
	"github.com/yasufadhili/jawt/internal/tsc/transformers"
)

type optionalChainTransformer struct {
	transformers.Transformer
}

func (ch *optionalChainTransformer) visit(node *ast.Node) *ast.Node {
	return node // !!!
}

func newOptionalChainTransformer(emitContext *printer.EmitContext) *transformers.Transformer {
	tx := &optionalChainTransformer{}
	return tx.NewTransformer(tx.visit, emitContext)
}
