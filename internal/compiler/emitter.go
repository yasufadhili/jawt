package compiler

import (
	"strings"
)

type Emitter struct {
	*BaseVisitor
	output             strings.Builder
	document           *JMLDocumentNode
	indentLevel        int
	componentProcessor *ComponentProcessor
}

func NewEmitter(doc *JMLDocumentNode) *Emitter {
	return &Emitter{
		BaseVisitor:        &BaseVisitor{},
		document:           doc,
		componentProcessor: NewComponentProcessor(),
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

	e.write("<style>@keyframes rocket-launch { 0% { transform: translateY(0px) scale(1); }  25% { transform: translateY(-10px) scale(1.05); }  50% { transform: translateY(-20px) scale(1.1); }  75% { transform: translateY(-15px) scale(1.05); }100% { transform: translateY(-5px) scale(1); }}@keyframes pulse-glow {0%, 100% { opacity: 0.4; } 50% { opacity: 0.8; }} @keyframes slide-up {0% { opacity: 0; transform: translateY(30px); } 100% { opacity: 1; transform: translateY(0); }} .rocket-launch { animation: rocket-launch 2s ease-out infinite; } .pulse-glow { animation: pulse-glow 3s ease-in-out infinite; }.slide-up { animation: slide-up 0.8s ease-out forwards; } .slide-up-delay-1 { animation: slide-up 0.8s ease-out 0.2s forwards; opacity: 0; }.slide-up-delay-2 { animation: slide-up 0.8s ease-out 0.4s forwards; opacity: 0; }.slide-up-delay-3 { animation: slide-up 0.8s ease-out 0.6s forwards; opacity: 0; }</style>")

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
	// Validate component
	if errors := e.componentProcessor.ValidateComponent(component.Name, component.Properties); len(errors) > 0 {
		// TODO: Handle validation errors appropriately
		e.write("<!-- Warning: " + strings.Join(errors, "; ") + " -->")
	}

	// Get HTML tag name
	tagName := e.componentProcessor.GetHTMLTag(component.Name, component.Properties)

	// Build attributes
	attributes := e.componentProcessor.BuildAttributes(component.Name, component.Properties)
	attributeString := e.buildAttributeString(attributes)

	// Raw text content
	textContent := e.componentProcessor.GetTextContent(component.Properties)

	// Check if self-closing
	isSelfClosing := e.componentProcessor.IsSelfClosing(component.Name)

	if isSelfClosing || len(component.Children) == 0 && textContent == "" {
		// Self-closing tag
		e.write("<" + tagName + attributeString + " />")
	} else {
		// Opening tag
		e.write("<" + tagName + attributeString + ">" + textContent)
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

func (e *Emitter) buildAttributeString(attributes map[string]string) string {
	if len(attributes) == 0 {
		return ""
	}

	var attrs []string
	for name, value := range attributes {
		if value != "" {
			attrs = append(attrs, name+"=\""+value+"\"")
		}
	}

	if len(attrs) == 0 {
		return ""
	}
	return " " + strings.Join(attrs, " ")
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
	// TODO: Implement web component emission
}

// emitWasm (placeholder)
func (e *Emitter) emitWasmModule() {
	// TODO: Implement WASM module emission
}

// emitCSS (placeholder)
func (e *Emitter) emitCSS() {
	// TODO: Implement CSS emission
}

// emitJS (placeholder)
func (e *Emitter) emitJS() {
	// TODO: Implement JS emission
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
