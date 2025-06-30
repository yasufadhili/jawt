package ast

type Node interface {
	Accept(Visitor) interface{}
}

type DocumentNode struct {
	Doctype *DoctypeSpecifierNode
	Imports []*ImportStatementNode
	Content DocumentContentNode // Interface for page/component
}

func (n *DocumentNode) Accept(v Visitor) interface{} {
	return v.VisitDocument(n)
}

type DoctypeSpecifierNode struct {
	Doctype string // "page", "component"
	Name    string
}

func (n *DoctypeSpecifierNode) Accept(v Visitor) interface{} {
	return v.VisitDoctypeSpecifier(n)
}

type ImportStatementNode struct {
	Type       string // "component", "styles"
	Identifier string
	From       string
}

func (n *ImportStatementNode) Accept(v Visitor) interface{} {
	return v.VisitImportStatement(n)
}

// DocumentContentNode - implemented by page, component definitions
type DocumentContentNode interface {
	Node
	IsDocumentContent()
}

type PageDefinitionNode struct {
	Properties []*PropertyNode
	Child      *ComponentElementNode // Pages can only have one child
}

func (n *PageDefinitionNode) Accept(v Visitor) interface{} {
	return v.VisitPageDefinition(n)
}

func (n *PageDefinitionNode) IsDocumentContent() {}

type ComponentDefinitionNode struct {
	Element *ComponentElementNode
}

func (n *ComponentDefinitionNode) Accept(v Visitor) interface{} {
	return v.VisitComponentDefinition(n)
}

func (n *ComponentDefinitionNode) IsDocumentContent() {}

type ComponentElementNode struct {
	Name       string
	Properties []*PropertyNode
	Children   []*ComponentElementNode
}

func (n *ComponentElementNode) Accept(v Visitor) interface{} {
	return v.VisitComponentElement(n)
}

func (n *ComponentElementNode) IsDocumentContent() {}

type ComponentBlockNode struct {
}

func (c *ComponentBlockNode) Accept(v Visitor) interface{} {
	return v.Visit(c)
}

type ComponentBodyNode struct {
}

func (c *ComponentBodyNode) Accept(v Visitor) interface{} {
	return v.Visit(c)
}

type ComponentPropertyNode struct {
}

func (c *ComponentPropertyNode) Accept(v Visitor) interface{} {
	return v.Visit(c)
}

type PropertyNode struct {
	Name  string
	Value PropertyValueNode
}

func (p *PropertyNode) Accept(v Visitor) interface{} {
	return v.VisitProperty(p)
}

type PropertyValueNode interface {
	Node
	IsPropertyValue()
}

type LiteralNode struct {
	Type  string // "integer", "float", "string", "boolean", "null"
	Value interface{}
}

func (l *LiteralNode) Accept(v Visitor) interface{} {
	return v.VisitLiteral(l)
}

func (l *LiteralNode) IsPropertyValue() {}
