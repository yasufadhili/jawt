package build

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/core"
	"github.com/yasufadhili/jawt/internal/process"
	"os/exec"
	"sync"
)

// CompilerRunner manages the execution of external compilers like tsc and tailwind.
type CompilerRunner struct {
	ctx *core.JawtContext
}

func NewCompilerRunner(ctx *core.JawtContext) *CompilerRunner {
	return &CompilerRunner{ctx: ctx}
}

func (cr *CompilerRunner) RunTSC() error {
	cr.ctx.Logger.Info("Running TypeScript compiler")

	tscPath, err := core.ResolveExecutablePath("tsc")
	if err != nil {
		return fmt.Errorf("tsc not found: %w. Please ensure TypeScript is installed", err)
	}

	cmd := exec.Command(tscPath, "--project", cr.ctx.Paths.TSConfigPath)
	cmd.Dir = cr.ctx.Paths.JawtDir // Run from the .jawt directory

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		process.ProcessLogger(stdout, cr.ctx.Logger, "tsc")
	}()

	go func() {
		defer wg.Done()
		process.ProcessLogger(stderr, cr.ctx.Logger, "tsc-err")
	}()

	wg.Wait()

	return cmd.Wait()
}

func (cr *CompilerRunner) RunTailwind() error {
	cr.ctx.Logger.Info("Running Tailwind CSS compiler")

	tailwindPath, err := core.ResolveExecutablePath("tailwindcss")
	if err != nil {
		return fmt.Errorf("tailwindcss not found: %w. Please ensure tailwindcss is installed", err)
	}

	cmd := exec.Command(tailwindPath,
		"-i", cr.ctx.Paths.TailwindConfigPath,
		"-o", cr.ctx.Paths.TailwindCSSPath,
		"--config", cr.ctx.Paths.TailwindConfigPath,
	)
	cmd.Dir = cr.ctx.Paths.JawtDir // Run from the .jawt directory

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		process.ProcessLogger(stdout, cr.ctx.Logger, "tailwind")
	}()

	go func() {
		defer wg.Done()
		process.ProcessLogger(stderr, cr.ctx.Logger, "tailwind-err")
	}()

	wg.Wait()

	return cmd.Wait()
}
