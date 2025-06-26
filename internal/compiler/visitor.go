package compiler

type ASTVisitor interface {
	Visit(ASTNode) interface{}
	VisitDocument(*JMLDocumentNode) interface{}
	VisitDoctypeSpecifier(*DoctypeSpecifierNode) interface{}
	VisitImportStatement(*ImportStatementNode) interface{}
	VisitPageDefinition(*PageDefinitionNode) interface{}
	VisitComponentDefinition(*ComponentDefinitionNode) interface{}
	VisitComponentElement(*ComponentElementNode) interface{}
	VisitComponentBlock(*ComponentBlockNode) interface{}
	VisitComponentBody(*ComponentBodyNode) interface{}
	VisitComponentProperty(*ComponentPropertyNode) interface{}
	VisitProperty(*PropertyNode) interface{}
	VisitLiteral(*LiteralNode) interface{}
}

// BaseVisitor provides a default implementation for traversing the AST by
// calling Accept on child nodes for composite nodes. Concrete visitors can
// embed this and override only the methods they need.
type BaseVisitor struct{}

func (v *BaseVisitor) Visit(n ASTNode) interface{} {
	return n.Accept(v)
}

func (v *BaseVisitor) VisitDocument(n *JMLDocumentNode) interface{} {
	if n.Doctype != nil {
		n.Doctype.Accept(v)
	}
	for _, imp := range n.Imports {
		imp.Accept(v)
	}
	if n.Content != nil {
		n.Content.Accept(v)
	}
	return nil
}

func (v *BaseVisitor) VisitDoctypeSpecifier(n *DoctypeSpecifierNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitImportStatement(n *ImportStatementNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitPageDefinition(n *PageDefinitionNode) interface{} {
	for _, prop := range n.Properties {
		prop.Accept(v)
	}
	if n.Child != nil {
		n.Child.Accept(v)
	}
	return nil
}

func (v *BaseVisitor) VisitComponentDefinition(n *ComponentDefinitionNode) interface{} {
	if n.Element != nil {
		n.Element.Accept(v)
	}
	return nil
}

func (v *BaseVisitor) VisitComponentElement(n *ComponentElementNode) interface{} {
	for _, prop := range n.Properties {
		prop.Accept(v)
	}
	for _, child := range n.Children {
		child.Accept(v)
	}
	return nil
}

func (v *BaseVisitor) VisitComponentBlock(n *ComponentBlockNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitComponentBody(n *ComponentBodyNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitComponentProperty(n *ComponentPropertyNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitProperty(n *PropertyNode) interface{} {
	if n.Value != nil {
		n.Value.Accept(v)
	}
	return nil
}

func (v *BaseVisitor) VisitLiteral(n *LiteralNode) interface{} {
	return nil
}
