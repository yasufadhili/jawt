package compiler

type ASTVisitor interface {
	Visit(ASTNode) interface{}

	VisitDocument(*JMLDocumentNode) interface{}
	VisitDoctypeSpecifier(*DoctypeSpecifierNode) interface{}
	VisitImportStatement(*ImportStatementNode) interface{}

	VisitPageDefinition(*PageDefinitionNode) interface{}
	VisitComponentDefinition(*ComponentDefinitionNode) interface{}

	VisitComponentElement(*ComponentElementNode) interface{}
	VisitProperty(*PropertyNode) interface{}

	VisitLiteral(*LiteralNode) interface{}
}

// BaseVisitor provides a default implementation for traversing the AST by
// calling VisitChildren for composite nodes. We embed this in our
// concrete visitors and override only the methods we care about.
type BaseVisitor struct{}

func (v *BaseVisitor) Visit(n ASTNode) interface{} {
	return n.Accept(v)
}

func (v *BaseVisitor) VisitDocument(n *JMLDocumentNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitDoctypeSpecifier(n *DoctypeSpecifierNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitImportStatement(n *ImportStatementNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitDocumentContent(n *DocumentContentNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitPageDefinition(n *PageDefinitionNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitComponentDefinition(n *ComponentDefinitionNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitComponentElement(n *ComponentElementNode) interface{} {
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
	return nil
}

func (v *BaseVisitor) VisitPropertyValue(n *PropertyValueNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitLiteral(n *LiteralNode) interface{} {
	return nil
}
