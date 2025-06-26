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
