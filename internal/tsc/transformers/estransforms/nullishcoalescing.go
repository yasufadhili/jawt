package estransforms

import (
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/printer"
	"github.com/yasufadhili/jawt/internal/tsc/transformers"
)

type nullishCoalescingTransformer struct {
	transformers.Transformer
}

func (ch *nullishCoalescingTransformer) visit(node *ast.Node) *ast.Node {
	return node // !!!
}

func newNullishCoalescingTransformer(emitContext *printer.EmitContext) *transformers.Transformer {
	tx := &nullishCoalescingTransformer{}
	return tx.NewTransformer(tx.visit, emitContext)
}
