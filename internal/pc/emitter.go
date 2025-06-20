package pc

import "fmt"

type Emitter struct {
	BaseVisitor
	Emitted string
}

func NewEmitter() *Emitter {
	return &Emitter{}
}

func (e *Emitter) VisitProgram(node *Program) interface{} {
	e.Emitted += "// Generated JML Page Code\n\n"
	e.Emitted += fmt.Sprintf("PAGE_NAME = \"%s\"\n", node.Doctype.Name)
	e.BaseVisitor.VisitProgram(node) // Continue visiting
	e.Emitted += "\n// End of Generated Code\n"

	return nil
}

func (e *Emitter) VisitImportStatement(node *ImportStatement) interface{} {
	e.Emitted += fmt.Sprintf("INCLUDE_%s(\"%s\") // From %s\n", node.Doctype, node.Identifier, node.From)
	return nil
}

func (e *Emitter) VisitPageProperty(node *PageProperty) interface{} {
	e.Emitted += fmt.Sprintf("  SET_PROPERTY(\"%s\", %v)\n", node.Key, node.Value)
	return nil
}
