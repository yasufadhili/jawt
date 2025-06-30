package emitter

import "github.com/yasufadhili/jawt/internal/ast"

func (e *Emitter) emitPage() {
	p := e.document.Content.(*ast.PageDefinitionNode)
	stringProps := e.getPageStringProperties(p)
	e.write("<!DOCTYPE html>")
	e.write("<html lang=\"en\">")
	e.indent()

	e.emitPageHead(stringProps)

	e.write("<body>")
	e.indent()

	if pageContent, ok := e.document.Content.(*ast.PageDefinitionNode); ok {
		e.emitPageBody(pageContent)
	}

	e.dedent()
	e.write("</body>")
	e.dedent()
	e.write("</html>")
}

func (e *Emitter) emitPageHead(stringProps map[string]string) {
	e.write("<head>")
	e.indent()
	e.write("<meta charset=\"utf-8\"/>")
	for k, v := range stringProps {
		switch k {
		case "title":
			e.write("<title>" + v + "</title>")
		}
	}
	e.dedent()
	e.write("</head>")
}

func (e *Emitter) emitPageBody(p *ast.PageDefinitionNode) {

}

func (e *Emitter) getPageTitle(p *ast.PageDefinitionNode) string {
	return ""
}

func (e *Emitter) getPageStringProperties(n *ast.PageDefinitionNode) map[string]string {
	res := make(map[string]string)
	for _, prop := range n.Properties {
		if literal, ok := prop.Value.(*ast.LiteralNode); ok {
			if str, ok := literal.Value.(string); ok {
				res[prop.Name] = str
			}
		}
	}
	return res
}
