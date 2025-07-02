package emitter

import (
	"github.com/yasufadhili/jawt/internal/ast"
	"strings"
)

func (e *Emitter) emitComponent() {

}

func (e *Emitter) emitComponentElement(c *ast.ComponentElementNode) {
	// Validate component
	if errors := e.componentProcessor.ValidateComponent(c.Name, c.Properties); len(errors) > 0 {
		// TODO: Handle validation errors appropriately
		e.write("<!-- Warning: " + strings.Join(errors, "; ") + " -->")
	}

	// Get HTML tag name
	tagName := e.componentProcessor.GetHTMLTag(c.Name, c.Properties)

	// Build attributes
	attributes := e.componentProcessor.BuildAttributes(c.Name, c.Properties)
	attributeString := e.buildAttributeString(attributes)

	// Raw text content
	textContent := e.componentProcessor.GetTextContent(c.Properties)

	isSelfClosing := e.componentProcessor.IsSelfClosing(c.Name)

	if isSelfClosing || len(c.Children) == 0 && textContent == "" {
		e.write("<" + tagName + attributeString + " />")
	} else {
		e.write("<" + tagName + attributeString + ">" + textContent)
		e.indent()

		for _, child := range c.Children {
			e.emitComponentElement(child)
		}

		e.dedent()
		e.write("</" + tagName + ">")
	}

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

func (e *Emitter) formatAttributeValue(literal *ast.LiteralNode) string {
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
