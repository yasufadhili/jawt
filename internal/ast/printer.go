package ast

import "fmt"

type Printer struct {
	indent int
}

func NewPrintVisitor() *Printer {
	return &Printer{indent: 0}
}

func (v *Printer) print(msg string) {
	for i := 0; i < v.indent; i++ {
		fmt.Print("  ")
	}
	fmt.Println(msg)
}

func (v *Printer) VisitProgram(node *Program) {
	v.print("Program")
	v.indent++
}

func (v *Printer) VisitDocument(node *Document) {
	v.print(fmt.Sprintf("Document: %s %s", node.Doctype, node.Identifier.Name))
	v.indent++
}

func (v *Printer) VisitImportStatement(node *ImportStatement) {
	v.print(fmt.Sprintf("Import: %s %s from %s", node.Kind, node.Identifier.Name, node.Path.Value))
}

func (v *Printer) VisitPropertyDeclaration(node *PropertyDeclaration) {
	v.print(fmt.Sprintf("Property: %s: %s", node.Name.Name, node.Type))
}

func (v *Printer) VisitStateDeclaration(node *StateDeclaration) {
	v.print(fmt.Sprintf("State: %s: %s", node.Name.Name, node.Type))
}

func (v *Printer) VisitElementNode(node *ElementNode) {
	v.print(fmt.Sprintf("Element: %s", node.Tag.Name))
	v.indent++
}

func (v *Printer) VisitStyleBlock(node *StyleBlock) {
	v.print("StyleBlock")
}

func (v *Printer) VisitScriptBlock(node *ScriptBlock) {
	v.print("ScriptBlock")
}

func (v *Printer) VisitLiteral(node *Literal) {
	v.print(fmt.Sprintf("Literal: %v (%s)", node.Value, node.Kind))
}

func (v *Printer) VisitIdentifier(node *Identifier) {
	v.print(fmt.Sprintf("Identifier: %s", node.Name))
}

func (v *Printer) VisitBinding(node *Binding) {
	v.print(fmt.Sprintf("Binding: %s.%s", node.Object.Name, node.Property.Name))
}

func (v *Printer) VisitFunctionCall(node *FunctionCall) {
	v.print(fmt.Sprintf("FunctionCall: %s", node.Function.Name))
	v.indent++
}

func (v *Printer) VisitLambdaExpression(node *LambdaExpression) {
	v.print("Lambda")
}
