package compiler

type ASTVisitor interface {
	Visit(ASTNode) interface{}

	VisitDocument(*JMLDocumentNode) interface{}

	VisitDoctypeSpecifier(*DoctypeSpecifierNode) interface{}
	VisitImportStatement(*ImportStatementNode) interface{}
	VisitDocumentContent(*DocumentContentNode) interface{}

	VisitPageDefinition(*PageDefinitionNode) interface{}
	VisitPageBody(*PageBodyNode) interface{}
	VisitPageProperty(*PagePropertyNode) interface{}

	VisitComponentDefinition(*ComponentDefinitionNode) interface{}
	VisitComponentElement(*ComponentElementNode) interface{}
	VisitComponentBlock(*ComponentBlockNode) interface{}
	VisitComponentBody(*ComponentBodyNode) interface{}
	VisitComponentProperty(*ComponentPropertyNode) interface{}

	VisitPropertyValue(*PropertyValueNode) interface{}
	VisitLiteralValue(*LiteralValueNode) interface{}
}

// BaseVisitor provides a default implementation for traversing the AST by
// calling VisitChildren for composite nodes. We embed this in the
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

func (v *BaseVisitor) VisitPageBody(n *PageBodyNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitPageProperty(n *PagePropertyNode) interface{} {
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

func (v *BaseVisitor) VisitPropertyValue(n *PropertyValueNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitLiteralValue(n *LiteralValueNode) interface{} {
	return nil
}
