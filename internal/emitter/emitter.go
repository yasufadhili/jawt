package emitter

import (
	"github.com/yasufadhili/jawt/internal/ast"
	"github.com/yasufadhili/jawt/internal/common"
	"github.com/yasufadhili/jawt/internal/emitter/registry"
	"strings"
)

type Emitter struct {
	*ast.BaseVisitor
	document           *ast.DocumentNode
	indentLevel        int
	target             common.BuildTarget
	output             strings.Builder
	componentProcessor *registry.ComponentProcessor
}

func NewEmitter(doc *ast.DocumentNode, target common.BuildTarget) *Emitter {
	return &Emitter{
		BaseVisitor:        &ast.BaseVisitor{},
		document:           doc,
		target:             target,
		componentProcessor: registry.NewComponentProcessor(),
	}
}

func (e *Emitter) Emit() string {
	switch e.target {
	case common.TargetPage:
		e.emitPage()
	case common.TargetComponent:
		e.emitComponent()
	default:
		return ""
	}
	return e.output.String()
}

func (e *Emitter) write(s string) {
	indent := strings.Repeat("  ", e.indentLevel)
	e.output.WriteString(indent + s + "\n")
}

// indent increases the indentation level
func (e *Emitter) indent() {
	e.indentLevel++
}

// dedent decreases the indentation level
func (e *Emitter) dedent() {
	if e.indentLevel > 0 {
		e.indentLevel--
	}
}
