package compiler

import (
	"fmt"
	"strings"
)

// JMLPrinter implements ASTVisitor to print JML source code from AST
type JMLPrinter struct {
	*BaseVisitor
	output      strings.Builder
	indentLevel int
}

// NewJMLPrinter creates a new printer instance
func NewJMLPrinter() *JMLPrinter {
	return &JMLPrinter{
		BaseVisitor: &BaseVisitor{},
		indentLevel: 0,
	}
}

// Print takes an AST node and returns the formatted JML source code
func (p *JMLPrinter) Print(node ASTNode) string {
	p.output.Reset()
	p.indentLevel = 0
	node.Accept(p)
	return p.output.String()
}

// Helper methods for formatting
func (p *JMLPrinter) writeIndent() {
	for i := 0; i < p.indentLevel; i++ {
		p.output.WriteString("  ")
	}
}

func (p *JMLPrinter) writeLine(text string) {
	p.writeIndent()
	p.output.WriteString(text)
	p.output.WriteString("\n")
}

func (p *JMLPrinter) write(text string) {
	p.output.WriteString(text)
}

func (p *JMLPrinter) increaseIndent() {
	p.indentLevel++
}

func (p *JMLPrinter) decreaseIndent() {
	if p.indentLevel > 0 {
		p.indentLevel--
	}
}

// Visitor implementations
func (p *JMLPrinter) VisitDocument(n *JMLDocumentNode) interface{} {
	// Print doctype if present
	if n.Doctype != nil {
		n.Doctype.Accept(p)
		p.output.WriteString("\n")
	}

	// Print imports if present
	if len(n.Imports) > 0 {
		for _, imp := range n.Imports {
			imp.Accept(p)
		}
		p.output.WriteString("\n")
	}

	// Print document content
	if n.Content != nil {
		n.Content.Accept(p)
	}

	return nil
}

func (p *JMLPrinter) VisitDoctypeSpecifier(n *DoctypeSpecifierNode) interface{} {
	p.write(fmt.Sprintf("doctype %s", n.Doctype))
	if n.Name != "" {
		p.write(fmt.Sprintf(" %s", n.Name))
	}
	return nil
}

func (p *JMLPrinter) VisitImportStatement(n *ImportStatementNode) interface{} {
	if n.IsBrowser {
		p.writeLine("import browser")
	} else {
		line := fmt.Sprintf("import %s %s from \"%s\"", n.Type, n.Identifier, n.From)
		p.writeLine(line)
	}
	return nil
}

func (p *JMLPrinter) VisitPageDefinition(n *PageDefinitionNode) interface{} {
	p.writeLine("page {")
	p.increaseIndent()

	if len(n.Properties) > 0 {
		for _, prop := range n.Properties {
			prop.Accept(p)
		}
	}

	if n.Child != nil {
		if len(n.Properties) > 0 {
			p.output.WriteString("\n")
		}
		n.Child.Accept(p)
	}

	p.decreaseIndent()
	p.writeLine("}")
	return nil
}

func (p *JMLPrinter) VisitComponentDefinition(n *ComponentDefinitionNode) interface{} {
	p.writeLine("component {")
	p.increaseIndent()

	if n.Element != nil {
		n.Element.Accept(p)
	}

	p.decreaseIndent()
	p.writeLine("}")
	return nil
}

func (p *JMLPrinter) VisitComponentElement(n *ComponentElementNode) interface{} {
	// Start the component element
	p.writeIndent()
	p.write(n.Name)

	// If we have properties or children, add a block
	if len(n.Properties) > 0 || len(n.Children) > 0 {
		p.write(" {\n")
		p.increaseIndent()

		// Print properties first
		for _, prop := range n.Properties {
			prop.Accept(p)
		}

		// Add spacing between properties and children if both exist
		if len(n.Properties) > 0 && len(n.Children) > 0 {
			p.output.WriteString("\n")
		}

		for _, child := range n.Children {
			child.Accept(p)
		}

		p.decreaseIndent()
		p.writeLine("}")
	} else {
		// Self-closing component
		p.write("\n")
	}

	return nil
}

func (p *JMLPrinter) VisitProperty(n *PropertyNode) interface{} {
	p.writeIndent()
	p.write(fmt.Sprintf("%s: ", n.Name))

	if n.Value != nil {
		n.Value.Accept(p)
	}

	p.write("\n")
	return nil
}

func (p *JMLPrinter) VisitLiteral(n *LiteralNode) interface{} {
	switch n.Type {
	case "string":
		p.write(fmt.Sprintf("\"%s\"", n.Value))
	case "integer":
		p.write(fmt.Sprintf("%d", n.Value))
	case "float":
		p.write(fmt.Sprintf("%g", n.Value))
	case "boolean":
		p.write(fmt.Sprintf("%t", n.Value))
	case "null":
		p.write("null")
	default:
		p.write(fmt.Sprintf("%v", n.Value))
	}
	return nil
}

func ExamplePrinterUsage() {

	ast := &JMLDocumentNode{
		Doctype: &DoctypeSpecifierNode{
			Doctype: "page",
			Name:    "HomePage",
		},
		Imports: []*ImportStatementNode{
			{
				Type:       "component",
				Identifier: "Header",
				From:       "./components/Header.jml",
				IsBrowser:  false,
			},
			{
				Type:      "browser",
				IsBrowser: true,
			},
		},
		Content: &PageDefinitionNode{
			Properties: []*PropertyNode{
				{
					Name: "title",
					Value: &LiteralNode{
						Type:  "string",
						Value: "Welcome to My Site",
					},
				},
				{
					Name: "showNav",
					Value: &LiteralNode{
						Type:  "boolean",
						Value: true,
					},
				},
			},
			Child: &ComponentElementNode{
				Name: "div",
				Properties: []*PropertyNode{
					{
						Name: "className",
						Value: &LiteralNode{
							Type:  "string",
							Value: "container",
						},
					},
				},
				Children: []*ComponentElementNode{
					{
						Name: "Header",
						Properties: []*PropertyNode{
							{
								Name: "title",
								Value: &LiteralNode{
									Type:  "string",
									Value: "My Website",
								},
							},
						},
					},
					{
						Name: "main",
						Children: []*ComponentElementNode{
							{
								Name: "h1",
								Properties: []*PropertyNode{
									{
										Name: "text",
										Value: &LiteralNode{
											Type:  "string",
											Value: "Hello World",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	printer := NewJMLPrinter()
	output := printer.Print(ast)

	fmt.Println("Generated JML:")
	fmt.Println(output)
}
