package ast

type Visitor interface {
	Visit(Node) interface{}
	VisitDocument(*DocumentNode) interface{}
}

// BaseVisitor provides a default implementation for traversing the AST by
// calling Accept on child nodes for composite nodes. Concrete visitors can
// embed this and override only the methods they need.
type BaseVisitor struct{}

func (v *BaseVisitor) Visit(n Node) interface{} {
	return n.Accept(v)
}

func (v *BaseVisitor) VisitDocument(n *DocumentNode) interface{} {
	return nil
}
