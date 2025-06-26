package compiler

import "fmt"

type PrettyPrintVisitor struct {
	BaseVisitor
	indent int
}

func NewPrettyPrintVisitor() *PrettyPrintVisitor {
	return &PrettyPrintVisitor{indent: 0}
}

func (v *PrettyPrintVisitor) VisitDocument(n *JMLDocumentNode) interface{} {
	fmt.Printf("Document:\n")
	v.indent++
	v.BaseVisitor.VisitDocument(n)
	v.indent--
	return nil
}

func (v *PrettyPrintVisitor) VisitComponentElement(n *ComponentElementNode) interface{} {
	fmt.Printf("%s%s {\n", v.getIndent(), n.Name)
	v.indent++
	v.BaseVisitor.VisitComponentElement(n)
	v.indent--
	fmt.Printf("%s}\n", v.getIndent())
	return nil
}

func (v *PrettyPrintVisitor) VisitProperty(n *PropertyNode) interface{} {
	fmt.Printf("%s%s: ", v.getIndent(), n.Name)
	v.BaseVisitor.VisitProperty(n)
	fmt.Println()
	return nil
}

func (v *PrettyPrintVisitor) VisitLiteral(n *LiteralNode) interface{} {
	fmt.Printf("%v", n.Value)
	return nil
}

func (v *PrettyPrintVisitor) getIndent() string {
	result := ""
	for i := 0; i < v.indent; i++ {
		result += "  "
	}
	return result
}
