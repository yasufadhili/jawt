package compiler

import (
	"github.com/yasufadhili/jawt/internal/ast"
	"github.com/yasufadhili/jawt/internal/core"
	"github.com/yasufadhili/jawt/internal/diagnostic"
)

type Compiler struct {
	ctx *core.JawtContext
}

func NewCompiler(ctx *core.JawtContext) *Compiler {
	return &Compiler{
		ctx: ctx,
	}
}

// Compile compiles a single JML file.
func (c *Compiler) Compile(file string, reporter *diagnostic.Reporter) (*ast.Document, error) {
	return nil, nil
}
