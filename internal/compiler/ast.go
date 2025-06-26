package compiler

type ASTNode interface {
	Accept(ASTVisitor) interface{}
}

type JMLDocumentNode struct {
	Doctype *DoctypeSpecifierNode
	Imports []*ImportStatementNode
	Content *DocumentContentNode
}

func (n *JMLDocumentNode) Accept(v ASTVisitor) interface{} {
	return v.Visit(n)
}

type DoctypeSpecifierNode struct {
}

type ImportStatementNode struct {
}

type DocumentContentNode struct {
}

type PageDefinitionNode struct {
}

type PageBodyNode struct {
}

type PagePropertyNode struct {
}

type ComponentDefinitionNode struct {
}

type ComponentElementNode struct {
}

type ComponentBlockNode struct {
}

type ComponentBodyNode struct {
}

type ComponentPropertyNode struct {
}

type PropertyValueNode struct {
}

type LiteralValueNode struct {
}
