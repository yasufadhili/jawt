package page_compiler

import (
	"fmt"
	"strings"
)

type HTMLEmitter struct {
	*BaseVisitor
	output      strings.Builder
	indentLevel int
}

func NewHTMLEmitter() *HTMLEmitter {
	return &HTMLEmitter{
		BaseVisitor: &BaseVisitor{},
	}
}

func (e *HTMLEmitter) EmitHTML(page *Page) string {

	e.output.Reset()
	e.indentLevel = 0

	e.writeHTML("<!DOCTYPE html>")
	e.writeHTML("<html lang=\"en\">")
	e.indent()

	e.emitHead(page)
	e.emitBody(page)

	e.dedent()
	e.writeHTML("</html>")

	return e.output.String()
}

// emitHead generates the HTML head section
func (e *HTMLEmitter) emitHead(page *Page) {
	e.writeHTML("<head>")
	e.indent()

	e.writeHTML("<meta charset=\"UTF-8\">")
	e.writeHTML("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">")

	if page.PageDefinition != nil && page.PageDefinition.Properties != nil {
		for _, prop := range page.PageDefinition.Properties {
			if prop.Key == "title" {
				e.writeHTML(fmt.Sprintf("<title>%s</title>", prop.Value))
			}
		}
	}

	e.writeHTML("<script src=\"https://cdn.tailwindcss.com\"></script>")

	e.dedent()
	e.writeHTML("</head>")
}

// emitBody generates the HTML body section
func (e *HTMLEmitter) emitBody(page *Page) {
	e.writeHTML("<body class=\"bg-gray-50 min-h-screen\">")
	e.indent()

	// Visit the page definition to emit content
	if page.PageDefinition != nil {
		page.PageDefinition.Accept(e)
	}

	e.dedent()
	e.writeHTML("</body>")
}

// VisitPageDefinition emits a container div for the page
func (e *HTMLEmitter) VisitPageDefinition(node *PageDefinition) interface{} {
	e.writeHTML("<div class=\"container mx-auto px-4 py-8\">")
	e.indent()

	for _, prop := range node.Properties {
		prop.Accept(e)
	}

	e.dedent()
	e.writeHTML("</div>")

	return nil
}

// VisitPageProperty emits HTML based on the property type
func (e *HTMLEmitter) VisitPageProperty(node *PageProperty) interface{} {
	switch node.Key {
	case "title":
		e.emitTitle(node.Value)
	case "subtitle":
		e.emitSubtitle(node.Value)
	case "content":
		e.emitContent(node.Value)
	case "layout":
		e.emitLayout(node.Value)
	default:
		// For unknown properties, emit as a data attribute div
		e.emitGenericProperty(node.Key, node.Value)
	}

	return nil
}

// emitTitle creates a h1 element
func (e *HTMLEmitter) emitTitle(value interface{}) {
	if str, ok := value.(string); ok {
		e.writeHTML(fmt.Sprintf("<h1 class=\"text-4xl font-bold text-gray-900 mb-6\">%s</h1>", e.escapeHTML(str)))
	}
}

// emitSubtitle creates an h2 element
func (e *HTMLEmitter) emitSubtitle(value interface{}) {
	if str, ok := value.(string); ok {
		e.writeHTML(fmt.Sprintf("<h2 class=\"text-2xl font-semibold text-gray-700 mb-4\">%s</h2>", e.escapeHTML(str)))
	}
}

// writeHTML writes a line of HTML with proper indentation
func (e *HTMLEmitter) writeHTML(html string) {
	indent := strings.Repeat("  ", e.indentLevel)
	e.output.WriteString(indent + html + "\n")
}

// emitContent creates a paragraph or div
func (e *HTMLEmitter) emitContent(value interface{}) {
	if str, ok := value.(string); ok {
		e.writeHTML(fmt.Sprintf("<p class=\"text-lg text-gray-600 leading-relaxed mb-4\">%s</p>", e.escapeHTML(str)))
	}
}

// emitLayout applies layout-specific classes
func (e *HTMLEmitter) emitLayout(value interface{}) {
	if str, ok := value.(string); ok {
		// TODO: modify the container classes based on layout type
		e.writeHTML(fmt.Sprintf("<!-- Layout: %s -->", e.escapeHTML(str)))
	}
}

// emitGenericProperty creates a generic div with data attribute
func (e *HTMLEmitter) emitGenericProperty(key string, value interface{}) {
	switch v := value.(type) {
	case string:
		e.writeHTML(fmt.Sprintf("<div class=\"mb-2\" data-%s=\"%s\">", key, e.escapeHTML(v)))
		e.indent()
		e.writeHTML(fmt.Sprintf("<span class=\"font-medium text-gray-800\">%s:</span> ", strings.Title(key)))
		e.writeHTML(fmt.Sprintf("<span class=\"text-gray-600\">%s</span>", e.escapeHTML(v)))
		e.dedent()
		e.writeHTML("</div>")
	case int:
		e.writeHTML(fmt.Sprintf("<div class=\"mb-2\" data-%s=\"%d\">", key, v))
		e.indent()
		e.writeHTML(fmt.Sprintf("<span class=\"font-medium text-gray-800\">%s:</span> ", strings.Title(key)))
		e.writeHTML(fmt.Sprintf("<span class=\"text-gray-600\">%d</span>", v))
		e.dedent()
		e.writeHTML("</div>")
	}
}

// indent increases the indentation level
func (e *HTMLEmitter) indent() {
	e.indentLevel++
}

// dedent decreases the indentation level
func (e *HTMLEmitter) dedent() {
	if e.indentLevel > 0 {
		e.indentLevel--
	}
}

// escapeHTML escapes HTML special characters
func (e *HTMLEmitter) escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// EmitterConfig allows customisation of the HTML output
type EmitterConfig struct {
	TailwindCDN  string
	CustomCSS    []string
	CustomJS     []string
	BodyClasses  []string
	WrapperClass string
}

// ConfigurableHTMLEmitter extends HTMLEmitter with configuration options
type ConfigurableHTMLEmitter struct {
	*HTMLEmitter
	config *EmitterConfig
}

// NewConfigurableHTMLEmitter creates a new configurable HTML emitter
func NewConfigurableHTMLEmitter(config *EmitterConfig) *ConfigurableHTMLEmitter {
	if config == nil {
		config = &EmitterConfig{
			TailwindCDN:  "https://cdn.tailwindcss.com",
			BodyClasses:  []string{"bg-gray-50", "min-h-screen"},
			WrapperClass: "container mx-auto px-4 py-8",
		}
	}

	return &ConfigurableHTMLEmitter{
		HTMLEmitter: NewHTMLEmitter(),
		config:      config,
	}
}

// EmitHTML generates customised HTML from the AST
func (e *ConfigurableHTMLEmitter) EmitHTML(page *Page) string {
	e.output.Reset()
	e.indentLevel = 0

	e.writeHTML("<!DOCTYPE html>")
	e.writeHTML("<html lang=\"en\">")
	e.indent()

	e.emitConfigurableHead(page)
	e.emitConfigurableBody(page)

	e.dedent()
	e.writeHTML("</html>")

	return e.output.String()
}

// emitConfigurableHead generates head with custom configuration
func (e *ConfigurableHTMLEmitter) emitConfigurableHead(page *Page) {
	e.writeHTML("<head>")
	e.indent()

	e.writeHTML("<meta charset=\"UTF-8\">")
	e.writeHTML("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">")

	title := "Page"
	if page.Doctype != nil && page.Doctype.Name != "" {
		title = page.Doctype.Name
	}
	e.writeHTML(fmt.Sprintf("<title>%s</title>", title))

	// TailwindCSS
	if e.config.TailwindCDN != "" {
		e.writeHTML(fmt.Sprintf("<script src=\"%s\"></script>", e.config.TailwindCDN))
	}

	// Custom CSS
	for _, css := range e.config.CustomCSS {
		e.writeHTML(fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\">", css))
	}

	// Custom JS
	for _, js := range e.config.CustomJS {
		e.writeHTML(fmt.Sprintf("<script src=\"%s\"></script>", js))
	}

	e.dedent()
	e.writeHTML("</head>")
}

// emitConfigurableBody generates body with custom classes
func (e *ConfigurableHTMLEmitter) emitConfigurableBody(page *Page) {
	bodyClass := strings.Join(e.config.BodyClasses, " ")
	e.writeHTML(fmt.Sprintf("<body class=\"%s\">", bodyClass))
	e.indent()

	// Wrapper div
	e.writeHTML(fmt.Sprintf("<div class=\"%s\">", e.config.WrapperClass))
	e.indent()

	if page.PageDefinition != nil {
		// Process page properties directly without additional container
		for _, prop := range page.PageDefinition.Properties {
			prop.Accept(e)
		}
	}

	e.dedent()
	e.writeHTML("</div>")

	e.dedent()
	e.writeHTML("</body>")
}
