package page_compiler

type Node interface {
	Accept(Visitor) interface{}
}

type Page struct {
	Doctype        *DoctypeSpecifier
	Imports        []*ImportStatement
	PageDefinition *PageDefinition
}

func (p *Page) Accept(v Visitor) interface{} {
	return v.VisitPage(p)
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

type PageDefinition struct {
	Name       string
	Properties []*PageProperty
}

func (d *PageDefinition) Accept(v Visitor) interface{} {
	return v.VisitPageDefinition(d)
}

type PageProperty struct {
	Key   string
	Value interface{} // Can be string, int, bool or another AST node
}

func (p *PageProperty) Accept(v Visitor) interface{} {
	return v.VisitPageProperty(p)
}

type PageBody struct {
	Properties []*PageProperty
	Child      *ComponentElement
}

type ComponentElement struct {
	Name  string
	Block *ComponentBlock
}

type ComponentBlock struct {
	Properties []*ComponentProperty
	Children   []*ComponentElement
	Functions  []*ScriptFunction
}

type ComponentProperty struct {
	Key   string
	Value interface{}
}

type ScriptFunction struct {
	Name       string
	Parameters []*Parameter
	ReturnType string
	Body       []Statement
}

type Parameter struct {
	Name string
	Type string
}

type Statement interface {
	StatementNode()
}
