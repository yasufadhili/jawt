package compiler

import (
	"strings"
)

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

	if e.document.Content != nil {
		switch content := e.document.Content.(type) {
		case *PageDefinitionNode:
			e.emitHTMLPage()
		case *ComponentDefinitionNode:
			e.emitWebComponent()
		default:
			// TODO: Handle other content types
			_ = content
		}
	}

	return e.output.String()
}

func (e *Emitter) emitHTMLPage() {

	e.write("<!DOCTYPE html>")
	e.write("<html lang=\"en\">")
	e.indent()
	e.write("<head>")
	e.indent()
	e.write("<meta charset=\"UTF-8\">")
	e.write("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">")
	e.write("<script src=\"https://cdn.tailwindcss.com\"></script>")

	// Extract title from page properties if available
	if pageContent, ok := e.document.Content.(*PageDefinitionNode); ok {
		title := e.getPageTitle(pageContent)
		if title != "" {
			e.write("<title>" + title + "</title>")
		}

	}

	e.dedent()
	e.write("</head>")

	e.write("<body>")
	e.indent()

	if pageContent, ok := e.document.Content.(*PageDefinitionNode); ok {
		e.emitPageBody(pageContent)
	}

	e.dedent()
	e.write("</body>")
	e.dedent()
	e.write("</html>")

}

func (e *Emitter) emitPageBody(page *PageDefinitionNode) {
	if page.Child != nil {
		e.emitComponentElement(page.Child)
	}
}

func (e *Emitter) emitComponentElement(component *ComponentElementNode) {
	// Start the HTML element
	tagName := e.getHTMLTagName(component.Name)
	attributes := e.buildAttributes(component.Properties)

	if len(component.Children) == 0 {
		// Self-closing tag
		e.write("<" + tagName + attributes + " />")
	} else {
		// Opening tag
		e.write("<" + tagName + attributes + ">")
		e.indent()

		// Emit children
		for _, child := range component.Children {
			e.emitComponentElement(child)
		}

		e.dedent()
		// Closing tag
		e.write("</" + tagName + ">")
	}
}

func (e *Emitter) getPageTitle(page *PageDefinitionNode) string {
	for _, prop := range page.Properties {
		if prop.Name == "title" {
			if literal, ok := prop.Value.(*LiteralNode); ok {
				if str, ok := literal.Value.(string); ok {
					return str
				}
			}
		}
	}
	return ""
}

func (e *Emitter) getHTMLTagName(componentName string) string {
	// TODO: make this more sophisticated
	switch componentName {
	case "Button":
		return "button"
	case "Input":
		return "input"
	case "Text":
		return "p"
	case "Div":
		return "div"
	case "Header":
		return "header"
	case "Footer":
		return "footer"
	case "Nav":
		return "nav"
	case "Section":
		return "section"
	case "Article":
		return "article"
	case "Aside":
		return "aside"
	case "Main":
		return "main"
	default:
		// Default to div for unknown components
		return "div"
	}
}

func (e *Emitter) buildAttributes(properties []*PropertyNode) string {
	if len(properties) == 0 {
		return ""
	}

	var attrs []string
	for _, prop := range properties {
		if literal, ok := prop.Value.(*LiteralNode); ok {
			attrName := e.mapPropertyToAttribute(prop.Name)
			attrValue := e.formatAttributeValue(literal)
			if attrValue != "" {
				attrs = append(attrs, attrName+"=\""+attrValue+"\"")
			}
		}
	}

	if len(attrs) == 0 {
		return ""
	}
	return " " + strings.Join(attrs, " ")
}

func (e *Emitter) mapPropertyToAttribute(propName string) string {
	// TODO: Improve map property names to HTML attributes
	switch propName {
	case "style":
		return "class"
	case "onClick":
		return "onclick"
	case "onChange":
		return "onchange"
	case "onSubmit":
		return "onsubmit"
	default:
		return propName
	}
}

func (e *Emitter) formatAttributeValue(literal *LiteralNode) string {
	switch v := literal.Value.(type) {
	case string:
		return v
	case int:
		return strings.Repeat(" ", 0) + string(rune(v+'0')) // Simple int to string
	case float64:
		return strings.Repeat(" ", 0) + "0" // Placeholder for float formatting
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func (e *Emitter) emitWebComponent() {

}

// emitWasm (placeholder)
func (e *Emitter) emitWasmModule() {

}

// emitCSS (placeholder)
func (e *Emitter) emitCSS() {

}

// emitJS (placeholder)
func (e *Emitter) emitJS() {

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
