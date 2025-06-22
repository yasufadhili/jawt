package page_compiler

type Visitor interface {
	Visit(Node) interface{}
	visitPage(*Page) interface{}
	visitDoctypeSpecifier(*DoctypeSpecifier) interface{}
	visitImportStatement(*ImportStatement) interface{}
	visitPageDefinition(*PageDefinition) interface{}
	visitPageProperty(*PageProperty) interface{}
}

// BaseVisitor provides a default implementation for traversing the AST by
// calling VisitChildren for composite nodes. We embed this in our
// concrete visitors and override only the methods we care about.
type BaseVisitor struct{}

func (v *BaseVisitor) Visit(node Node) interface{} {
	return node.Accept(v)
}

func (v *BaseVisitor) visitPage(node *Page) interface{} {

	return nil
}

func (v *BaseVisitor) visitDoctypeSpecifier(node *DoctypeSpecifier) interface{} {
	return nil
}

func (v *BaseVisitor) visitImportStatement(node *ImportStatement) interface{} {
	return nil
}

func (v *BaseVisitor) visitPageDefinition(node *PageDefinition) interface{} {
	return nil
}

func (v *BaseVisitor) visitPageProperty(node *PageProperty) interface{} {
	return nil
}
