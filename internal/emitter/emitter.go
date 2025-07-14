package emitter

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/ast"
	"github.com/yasufadhili/jawt/internal/core"
	os
	"path/filepath"
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
		core.StringField("name", doc.Name.Name),
		core.StringField("path", doc.SourceFile))

	// Determine the output path in the user source directory
	relPath, err := filepath.Rel(e.ctx.Paths.ProjectRoot, doc.SourceFile)
	if err != nil {
		return fmt.Errorf("failed to get relative path for %s: %w", doc.SourceFile, err)
	}

	// Change .jml extension to .ts
	outputPath := filepath.Join(e.ctx.Paths.UserSrcDir, strings.TrimSuffix(relPath, ".jml")+".ts")

	// Ensure the output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory for %s: %w", outputPath, err)
	}

	// For now, just write a placeholder TypeScript file
	content := []byte(fmt.Sprintf("// Generated from JML file: %s\n\nconsole.log(\"Hello from %s!\");\n", doc.SourceFile, doc.Name.Name))

	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write emitted TypeScript to %s: %w", outputPath, err)
	}

	return nil
}
