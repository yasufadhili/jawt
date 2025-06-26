package compiler

type ASTNode interface {
	Accept(ASTVisitor) interface{}
}

type JMLDocumentNode struct {
	Doctype *DoctypeSpecifierNode
	Imports []*ImportStatementNode
	Content DocumentContentNode // Interface for page/component/module
}

func (n *JMLDocumentNode) Accept(v ASTVisitor) interface{} {
	return v.VisitDocument(n)
}

type DoctypeSpecifierNode struct {
	Doctype string // "page", "component", "module"
	Name    string
}

func (n *DoctypeSpecifierNode) Accept(v ASTVisitor) interface{} {
	return v.VisitDoctypeSpecifier(n)
}

type ImportStatementNode struct {
	Type       string // "page", "component", "module", "browser"
	Identifier string // empty for browser imports
	From       string // file path, empty for browser imports
	IsBrowser  bool   // true for browser API imports
}

func (n *ImportStatementNode) Accept(v ASTVisitor) interface{} {
	return v.VisitImportStatement(n)
}

// DocumentContentNode - implemented by page, component, module definitions
type DocumentContentNode interface {
	ASTNode
	IsDocumentContent()
}

type PageDefinitionNode struct {
	Properties []*PropertyNode
	Child      *ComponentElementNode // Pages can only have one child
}

func (p *PageDefinitionNode) Accept(v ASTVisitor) interface{} {
	return v.VisitPageDefinition(p)
}

func (p *PageDefinitionNode) IsDocumentContent() {}

type ComponentDefinitionNode struct {
	Element *ComponentElementNode
}

func (c *ComponentDefinitionNode) Accept(v ASTVisitor) interface{} {
	return v.VisitComponentDefinition(c)
}

func (c *ComponentDefinitionNode) IsDocumentContent() {}

type ComponentElementNode struct {
	Name       string
	Properties []*PropertyNode
	Children   []*ComponentElementNode
}

func (c *ComponentElementNode) Accept(v ASTVisitor) interface{} {
	return v.VisitComponentElement(c)
}

type ComponentBlockNode struct {
}

func (c *ComponentBlockNode) Accept(v ASTVisitor) interface{} {
	return v.Visit(c)
}

type ComponentBodyNode struct {
}

func (c *ComponentBodyNode) Accept(v ASTVisitor) interface{} {
	return v.Visit(c)
}

type ComponentPropertyNode struct {
}

func (c *ComponentPropertyNode) Accept(v ASTVisitor) interface{} {
	return v.Visit(c)
}

type PropertyNode struct {
	Name  string
	Value PropertyValueNode
}

func (p *PropertyNode) Accept(v ASTVisitor) interface{} {
	return v.VisitProperty(p)
}

type PropertyValueNode interface {
	ASTNode
	IsPropertyValue()
}

type LiteralNode struct {
	Type  string // "integer", "float", "string", "boolean", "null"
	Value interface{}
}

func (l *LiteralNode) Accept(v ASTVisitor) interface{} {
	return v.VisitLiteral(l)
}

func (l *LiteralNode) IsPropertyValue() {}
