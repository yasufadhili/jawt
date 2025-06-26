package compiler

import "strings"

type Emitter struct {
	*BaseVisitor
	output      strings.Builder
	document    *JMLDocumentNode
	indentLevel int
}

func NewEmitter(doc *JMLDocumentNode) *Emitter {
	return &Emitter{
		BaseVisitor: &BaseVisitor{},
		document:    doc,
	}
}

func (e *Emitter) Emit() string {

	e.output.Reset()
	e.indentLevel = 0

	return e.output.String()
}

func (e *Emitter) emitHTMLPage() string {

	return ""
}

func (e *Emitter) emitWebComponent() string {

	return ""
}

// emitWasm (placeholder)
func (e *Emitter) emitWasmModule() string {
	return ""
}

// emitCSS (placeholder)
func (e *Emitter) emitCSS() string {
	return ""
}

// emitJS (placeholder)
func (e *Emitter) emitJS() string {
	return ""
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
