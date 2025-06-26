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

func (v *BaseVisitor) Visit(node ASTNode) interface{} {
	return node.Accept(v)
}

func (v *BaseVisitor) VisitDocument(node *JMLDocumentNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitDoctypeSpecifier(node *DoctypeSpecifierNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitImportStatement(node *ImportStatementNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitDocumentContent(node *DocumentContentNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitPageDefinition(node *PageDefinitionNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitPageBody(node *PageBodyNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitPageProperty(node *PagePropertyNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitComponentDefinition(node *ComponentDefinitionNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitComponentElement(node *ComponentElementNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitComponentBlock(node *ComponentBlockNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitComponentBody(node *ComponentBodyNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitComponentProperty(node *ComponentPropertyNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitPropertyValue(node *PropertyValueNode) interface{} {
	return nil
}

func (v *BaseVisitor) VisitLiteralValue(node *LiteralValueNode) interface{} {
	return nil
}
