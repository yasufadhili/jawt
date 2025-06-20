package pc

type Node interface {
	Accept(Visitor) interface{}
}

type Program struct {
	Doctype *DoctypeSpecifier
	Imports []*ImportStatement
	Page    *Page
}

func (p *Program) Accept(v Visitor) interface{} {
	return v.VisitProgram(p)
}

type DoctypeSpecifier struct {
	Doctype string
	Name    string
}

func (d *DoctypeSpecifier) Accept(v Visitor) interface{} {
	return v.VisitDoctypeSpecifier(d)
}

type ImportStatement struct {
	Doctype    string
	Identifier string
	From       string
}

func (i *ImportStatement) Accept(v Visitor) interface{} {
	return v.VisitImportStatement(i)
}

type Page struct {
	Name       string
	RelPath    string
	AbsPath    string
	Properties []PageProperty
}

func (p *Page) Accept(v Visitor) interface{} {
	return v.VisitPage(p)
}

type PageProperty struct {
	Key   string
	Value interface{} // Can be string, int, bool or another AST node
}

func (pp *PageProperty) Accept(v Visitor) interface{} {
	return v.VisitPageProperty(pp)
}
