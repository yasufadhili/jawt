package page_compiler

type Node interface {
	Accept(Visitor) interface{}
}

type Page struct {
	Doctype        *DoctypeSpecifier
	Imports        []*ImportStatement
	PageDefinition *PageDefinition
}

type DoctypeSpecifier struct {
	Doctype string
	Name    string
}

type ImportStatement struct {
	Doctype    string
	Identifier string
	From       string
}

type PageDefinition struct {
	Name       string
	Properties []*PageProperty
}

type PageProperty struct {
	Key   string
	Value interface{} // Can be string, int, bool or another AST node
}
