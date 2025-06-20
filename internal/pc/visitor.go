package pc

type Visitor interface {
	Visit(Node) interface{}
	VisitProgram(*Program) interface{}
	VisitDoctypeSpecifier(*DoctypeSpecifier) interface{}
	VisitImportStatement(*ImportStatement) interface{}
	VisitPage(*Page) interface{}
	VisitPageProperty(*PageProperty) interface{}
}

// BaseVisitor provides a default implementation for traversing the AST by
// calling VisitChildren for composite nodes. We embed this in our
// concrete visitors and override only the methods we care about.
type BaseVisitor struct{}

func (v *BaseVisitor) Visit(node Node) interface{} {
	return node.Accept(v)
}

func (v *BaseVisitor) VisitProgram(node *Program) interface{} {
	if node.Doctype != nil {
		node.Doctype.Accept(v)
	}
	for _, imp := range node.Imports {
		imp.Accept(v)
	}
	if node.Page != nil {
		node.Page.Accept(v)
	}
	return nil
}

func (v *BaseVisitor) VisitDoctypeSpecifier(node *DoctypeSpecifier) interface{} {
	return nil // Leaf node, nothing to visit further
}

func (v *BaseVisitor) VisitImportStatement(node *ImportStatement) interface{} {
	return nil // Leaf node
}

func (v *BaseVisitor) VisitPage(node *Page) interface{} {
	for _, prop := range node.Properties {
		prop.Accept(v)
	}
	return nil
}

func (v *BaseVisitor) VisitPageProperty(node *PageProperty) interface{} {
	// For now the value is a primitive type (string, int)
	// FUTURE: Have expressions, complex types as properties
	return nil
}
