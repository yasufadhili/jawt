package ast

import (
	"fmt"
	"io"
	"strings"
)

// Printer is an AST visitor that prints the AST to a writer.
type Printer struct {
	BaseVisitor
	writer io.Writer
	indent int
}

// NewPrinter creates a new Printer.
func NewPrinter(writer io.Writer) *Printer {
	return &Printer{writer: writer}
}

// Print prints the given node to the writer.
func (p *Printer) Print(node Node) {
	Walk(p, node)
}

func (p *Printer) printf(format string, args ...interface{}) {
	fmt.Fprintf(p.writer, strings.Repeat("  ", p.indent)+format+"\n", args...)
}

// Visit methods

func (p *Printer) VisitProgram(n *Program) {
	p.printf("Program")
	p.indent++
	for _, doc := range n.Documents {
		p.Print(doc)
	}
	p.indent--
}

func (p *Printer) VisitDocument(n *Document) {
	p.printf("Document (%s %s)", n.DocType, n.Name.Name)
	p.indent++
	for _, stmt := range n.Body {
		p.Print(stmt)
	}
	p.indent--
}

func (p *Printer) VisitImportDeclaration(n *ImportDeclaration) {
	if n.IsBrowser {
		p.printf("ImportDeclaration (browser)")
		return
	}
	p.printf("ImportDeclaration from \"%s\"", n.Source.Value)
	p.indent++
	for _, spec := range n.Specifiers {
		p.Print(spec)
	}
	p.indent--
}

func (p *Printer) VisitImportSpecifier(n *ImportSpecifier) {
	if n.Imported != nil {
		p.printf("ImportSpecifier (%s as %s)", n.Imported.Name, n.Local.Name)
	} else {
		p.printf("ImportSpecifier (%s)", n.Local.Name)
	}
}

func (p *Printer) VisitImportDefaultSpecifier(n *ImportDefaultSpecifier) {
	p.printf("ImportDefaultSpecifier (%s)", n.Local.Name)
}

func (p *Printer) VisitImportNamespaceSpecifier(n *ImportNamespaceSpecifier) {
	p.printf("ImportNamespaceSpecifier (%s)", n.Local.Name)
}

func (p *Printer) VisitExportDeclaration(n *ExportDeclaration) {
	prefix := ""
	if n.Default {
		prefix = "default "
	}
	p.printf("ExportDeclaration (%s)", prefix)
	p.indent++
	p.Print(n.Declaration)
	p.indent--
}

func (p *Printer) VisitVariableDeclaration(n *VariableDeclaration) {
	p.printf("VariableDeclaration (%s)", n.Kind)
	p.indent++
	for _, decl := range n.Declarations {
		p.printf("Declarator (%s)", decl.ID.Name)
		if decl.Init != nil {
			p.indent++
			p.Print(decl.Init)
			p.indent--
		}
	}
	p.indent--
}

func (p *Printer) VisitFunctionDeclaration(n *FunctionDeclaration) {
	p.printf("FunctionDeclaration (%s)", n.ID.Name)
	p.indent++
	// ... (print params, return type, body)
	p.indent--
}

func (p *Printer) VisitBlockStatement(n *BlockStatement) {
	p.printf("BlockStatement")
	p.indent++
	for _, stmt := range n.List {
		p.Print(stmt)
	}
	p.indent--
}

func (p *Printer) VisitExpressionStatement(n *ExpressionStatement) {
	p.printf("ExpressionStatement")
	p.indent++
	p.Print(n.Expression)
	p.indent--
}

func (p *Printer) VisitIfStatement(n *IfStatement) {
	p.printf("IfStatement")
	p.indent++
	p.printf("Test:")
	p.indent++
	p.Print(n.Test)
	p.indent--
	p.printf("Consequent:")
	p.indent++
	p.Print(n.Consequent)
	p.indent--
	if n.Alternate != nil {
		p.printf("Alternate:")
		p.indent++
		p.Print(n.Alternate)
		p.indent--
	}
	p.indent--
}

func (p *Printer) VisitReturnStatement(n *ReturnStatement) {
	p.printf("ReturnStatement")
	if n.Argument != nil {
		p.indent++
		p.Print(n.Argument)
		p.indent--
	}
}

func (p *Printer) VisitIdentifier(n *Identifier) {
	p.printf("Identifier (%s)", n.Name)
}

func (p *Printer) VisitLiteral(n *Literal) {
	p.printf("Literal (%s: %s)", n.Kind, n.Value)
}

func (p *Printer) VisitBinaryExpression(n *BinaryExpression) {
	p.printf("BinaryExpression (%s)", n.Operator)
	p.indent++
	p.Print(n.Left)
	p.Print(n.Right)
	p.indent--
}

func (p *Printer) VisitCallExpression(n *CallExpression) {
	p.printf("CallExpression")
	p.indent++
	p.printf("Callee:")
	p.indent++
	p.Print(n.Callee)
	p.indent--
	if len(n.Arguments) > 0 {
		p.printf("Arguments:")
		p.indent++
		for _, arg := range n.Arguments {
			p.Print(arg)
		}
		p.indent--
	}
	p.indent--
}

func (p *Printer) VisitComponentElement(n *ComponentElement) {
	p.printf("ComponentElement <%s>", n.Tag.Name)
	p.indent++
	for _, child := range n.Body {
		p.Print(child)
	}
	p.indent--
}

func (p *Printer) VisitComponentProperty(n *ComponentProperty) {
	p.printf("ComponentProperty (%s)", n.Key.Name)
	p.indent++
	p.Print(n.Value)
	p.indent--
}

func (p *Printer) VisitObjectType(n *ObjectType) {
	p.printf("ObjectType")
	p.indent++
	for _, member := range n.Members {
		p.Print(member)
	}
	p.indent--
}

func (p *Printer) VisitInterfaceDeclaration(n *InterfaceDeclaration) {
	p.printf("InterfaceDeclaration (%s)", n.ID.Name)
	p.indent++
	if len(n.Extends) > 0 {
		p.printf("Extends:")
		p.indent++
		for _, ext := range n.Extends {
			p.Print(ext)
		}
		p.indent--
	}
	p.printf("Body:")
	p.indent++
	p.Print(n.Body)
	p.indent--
	p.indent--
}

func (p *Printer) VisitForLoop(n *ForLoop) {
	p.printf("ForLoop (%s in ...)", n.Variable.Name)
	p.indent++
	p.printf("Source:")
	p.indent++
	p.Print(n.Source)
	p.indent--
	p.printf("Body:")
	p.indent++
	for _, child := range n.Body {
		p.Print(child)
	}
	p.indent--
	p.indent--
}

func (p *Printer) VisitIfCondition(n *IfCondition) {
	p.printf("IfCondition")
	p.indent++
	p.printf("Test:")
	p.indent++
	p.Print(n.Test)
	p.indent--
	p.printf("Consequent:")
	p.indent++
	for _, child := range n.Consequent {
		p.Print(child)
	}
	p.indent--
	if n.Alternate != nil {
		p.printf("Alternate:")
		p.indent++
		for _, child := range n.Alternate {
			p.Print(child)
		}
		p.indent--
	}
	p.indent--
}

func (p *Printer) VisitTypeAnnotation(n *TypeAnnotation) {
	p.printf("TypeAnnotation")
	p.indent++
	p.Print(n.Type)
	p.indent--
}

func (p *Printer) VisitTypeReference(n *TypeReference) {
	p.printf("TypeReference (%s)", n.Name.Name)
	if len(n.TypeParams) > 0 {
		p.indent++
		p.printf("Params:")
		p.indent++
		for _, param := range n.TypeParams {
			p.Print(param)
		}
		p.indent--
		p.indent--
	}
}
