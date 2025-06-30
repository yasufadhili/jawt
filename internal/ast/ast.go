package ast

type Node interface {
	Accept(Visitor) interface{}
}

type DocumentNode struct {
}

func (n *DocumentNode) Accept(v Visitor) interface{} {
	return v.VisitDocument(n)
}
