package compiler

import (
	"encoding/json"
)

// JSONPrinter is a visitor that constructs a JSON representation of the AST.
type JSONPrinter struct {
	BaseVisitor
}

// NewJSONPrinter creates a new JSONPrinter instance.
func NewJSONPrinter() *JSONPrinter {
	return &JSONPrinter{}
}

// Print converts the AST rooted at the given node to a JSON string.
func (p *JSONPrinter) Print(node ASTNode) (string, error) {
	result := node.Accept(p)
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func (p *JSONPrinter) VisitDocument(n *JMLDocumentNode) interface{} {
	node := make(map[string]interface{})
	node["type"] = "JMLDocument"
	if n.Doctype != nil {
		node["doctype"] = n.Doctype.Accept(p)
	}
	if len(n.Imports) > 0 {
		imports := make([]interface{}, len(n.Imports))
		for i, imp := range n.Imports {
			imports[i] = imp.Accept(p)
		}
		node["imports"] = imports
	}
	if n.Content != nil {
		node["content"] = n.Content.Accept(p)
	}
	return node
}

func (p *JSONPrinter) VisitDoctypeSpecifier(n *DoctypeSpecifierNode) interface{} {
	node := make(map[string]interface{})
	node["type"] = "DoctypeSpecifier"
	node["doctype"] = n.Doctype
	if n.Name != "" {
		node["name"] = n.Name
	}
	return node
}

func (p *JSONPrinter) VisitImportStatement(n *ImportStatementNode) interface{} {
	node := make(map[string]interface{})
	node["type"] = "ImportStatement"
	node["importType"] = n.Type
	if n.Identifier != "" {
		node["identifier"] = n.Identifier
	}
	if n.From != "" {
		node["from"] = n.From
	}
	node["isBrowser"] = n.IsBrowser
	return node
}

func (p *JSONPrinter) VisitPageDefinition(n *PageDefinitionNode) interface{} {
	node := make(map[string]interface{})
	node["type"] = "PageDefinition"
	if len(n.Properties) > 0 {
		props := make([]interface{}, len(n.Properties))
		for i, prop := range n.Properties {
			props[i] = prop.Accept(p)
		}
		node["properties"] = props
	}
	if n.Child != nil {
		node["child"] = n.Child.Accept(p)
	}
	return node
}

func (p *JSONPrinter) VisitComponentDefinition(n *ComponentDefinitionNode) interface{} {
	node := make(map[string]interface{})
	node["type"] = "ComponentDefinition"
	if n.Element != nil {
		node["element"] = n.Element.Accept(p)
	}
	return node
}

func (p *JSONPrinter) VisitComponentElement(n *ComponentElementNode) interface{} {
	node := make(map[string]interface{})
	node["type"] = "ComponentElement"
	node["name"] = n.Name
	if len(n.Properties) > 0 {
		props := make([]interface{}, len(n.Properties))
		for i, prop := range n.Properties {
			props[i] = prop.Accept(p)
		}
		node["properties"] = props
	}
	if len(n.Children) > 0 {
		children := make([]interface{}, len(n.Children))
		for i, child := range n.Children {
			children[i] = child.Accept(p)
		}
		node["children"] = children
	}
	return node
}

func (p *JSONPrinter) VisitComponentBlock(n *ComponentBlockNode) interface{} {
	node := make(map[string]interface{})
	node["type"] = "ComponentBlock"
	return node
}

func (p *JSONPrinter) VisitComponentBody(n *ComponentBodyNode) interface{} {
	node := make(map[string]interface{})
	node["type"] = "ComponentBody"
	return node
}

func (p *JSONPrinter) VisitComponentProperty(n *ComponentPropertyNode) interface{} {
	node := make(map[string]interface{})
	node["type"] = "ComponentProperty"
	return node
}

func (p *JSONPrinter) VisitProperty(n *PropertyNode) interface{} {
	node := make(map[string]interface{})
	node["type"] = "Property"
	node["name"] = n.Name
	if n.Value != nil {
		node["value"] = n.Value.Accept(p)
	}
	return node
}

func (p *JSONPrinter) VisitLiteral(n *LiteralNode) interface{} {
	node := make(map[string]interface{})
	node["type"] = "Literal"
	node["literalType"] = n.Type
	node["value"] = n.Value
	return node
}
