package page_compiler

type Visitor interface {
	Visit(Node) interface{}
	VisitPage(*Page) interface{}
	VisitDoctypeSpecifier(*DoctypeSpecifier) interface{}
	VisitImportStatement(*ImportStatement) interface{}
	VisitPageDefinition(*PageDefinition) interface{}
	VisitPageProperty(*PageProperty) interface{}
}

// BaseVisitor provides a default implementation for traversing the AST by
// calling VisitChildren for composite nodes. We embed this in our
// concrete visitors and override only the methods we care about.
type BaseVisitor struct{}

func (v *BaseVisitor) Visit(node Node) interface{} {
	return node.Accept(v)
}

func (v *BaseVisitor) VisitPage(node *Page) interface{} {
	if node.Doctype != nil {
		node.Doctype.Accept(v)
	}
	for _, imp := range node.Imports {
		imp.Accept(v)
	}

	if node.PageDefinition != nil {
		node.PageDefinition.Accept(v)
	}

	return nil
}

func (v *BaseVisitor) VisitDoctypeSpecifier(node *DoctypeSpecifier) interface{} {
	return nil
}

func (v *BaseVisitor) VisitImportStatement(node *ImportStatement) interface{} {
	return nil
}

func (v *BaseVisitor) VisitPageDefinition(node *PageDefinition) interface{} {
	for _, prop := range node.Properties {
		prop.Accept(v)
	}
	return nil
}

func (v *BaseVisitor) VisitPageProperty(node *PageProperty) interface{} {
	return nil
}
