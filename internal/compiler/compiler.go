package compiler

import (
	"github.com/yasufadhili/jawt/internal/ast"
	"github.com/yasufadhili/jawt/internal/core"
	"github.com/yasufadhili/jawt/internal/diagnostic"
)

// Compiler compiles JML files
type Compiler struct {
	ctx      *core.JawtContext
	reporter *diagnostic.Reporter
}

// NewCompiler creates a new compiler.
func NewCompiler(ctx *core.JawtContext, reporter *diagnostic.Reporter) *Compiler {
	return &Compiler{
		ctx:      ctx,
		reporter: reporter,
	}
}

// Compile compiles a single JML file.
func (c *Compiler) Compile(file string) (*ast.Document, error) {
	// TODO: Implement the compilation logic.
	return nil, nil
}
