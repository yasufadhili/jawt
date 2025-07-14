package build

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/core"
	"os/exec"
)

// CompilerRunner manages the execution of external compilers like tsc and tailwind.
type CompilerRunner struct {
	ctx *core.JawtContext
}

// NewCompilerRunner creates a new CompilerRunner.
func NewCompilerRunner(ctx *core.JawtContext) *CompilerRunner {
	return &CompilerRunner{ctx: ctx}
}

// RunTSC runs the TypeScript compiler.
func (cr *CompilerRunner) RunTSC() error {
	cr.ctx.Logger.Info("Running TypeScript compiler")

	// TODO: Implement logic to find and run the tsc executable

	return nil
}

// RunTailwind runs the Tailwind CSS compiler.
func (cr *CompilerRunner) RunTailwind() error {
	cr.ctx.Logger.Info("Running Tailwind CSS compiler")

	// TODO: Implement logic to find and run the tailwindcss executable

	return nil
}
