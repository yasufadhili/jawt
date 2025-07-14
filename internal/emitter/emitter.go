package emitter

import (
	"github.com/yasufadhili/jawt/internal/ast"
	"github.com/yasufadhili/jawt/internal/core"
)

type Emitter struct {
	ctx *core.JawtContext
}

func NewEmitter(ctx *core.JawtContext) *Emitter {
	return &Emitter{ctx: ctx}
}

// Emit takes an AST document and emits TypeScript code to the workspace.
func (e *Emitter) Emit(doc *ast.Document) error {
	e.ctx.Logger.Info("Emitting TypeScript for JML document",
		core.StringField("name", ""),
		core.StringField("path", ""))

	return nil
}
