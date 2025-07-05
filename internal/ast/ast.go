package ast

import (
	"fmt"
	"strings"
)

type Position struct {
	Line   int
	Column int
	File   string
}

func (p Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Column)
}

type Node interface {
	Pos() Position
	Accept(Visitor)
	String() string
}

type Visitor interface {
	DocumentVisitor
	DeclarationVisitor
	BlockVisitor
	ExpressionVisitor
}

type ExpressionVisitor interface {
	VisitLiteral(*Literal)
	VisitIdentifier(*Identifier)
	VisitBinding(*Binding)
	VisitFunctionCall(*FunctionCall)
	VisitLambdaExpression(node *LambdaExpression)
}

type DeclarationVisitor interface {
	VisitPropertyDeclaration(*PropertyDeclaration)
	VisitStateDeclaration(*StateDeclaration)
}

type BlockVisitor interface {
	VisitStyleBlock(node *StyleBlock)
	VisitScriptBlock(node *ScriptBlock)
	VisitElementNode(node *ElementNode)
}

type DocumentVisitor interface {
	VisitProgram(*Program)
	VisitDocument(*Document)
	VisitImportStatement(*ImportStatement)
}

// Document Types

type DocType int

const (
	DoctypePage DocType = iota
	DoctypeComponent
)

func (d DocType) String() string {
	switch d {
	case DoctypePage:
		return "page"
	case DoctypeComponent:
		return "component"
	default:
		return "unknown"
	}
}

// Import Types

type ImportKind int

const (
	ImportComponent ImportKind = iota
	ImportScript
)

func (i ImportKind) String() string {
	switch i {
	case ImportComponent:
		return "component"
	case ImportScript:
		return "script"
	default:
		return "unknown"
	}
}

// Expression Types

type ExprType int

const (
	ExprLiteral ExprType = iota
	ExprIdentifier
	ExprBinding
	ExprFunctionCall
	ExprLambda
)

type Expression interface {
	Node
	ExprType() ExprType
}

// ================= ROOT NODE ======================

// Program represents he entire JML project
type Program struct {
	Position
	Documents []*Document
}

func (p *Program) Pos() Position {
	return p.Position
}

func (p *Program) Accept(v Visitor) {
	v.VisitProgram(p)
}

func (p *Program) String() string {
	var docs []string
	for _, doc := range p.Documents {
		docs = append(docs, doc.String())
	}
	return fmt.Sprintf("Program{\n%s\n}", strings.Join(docs, "\n"))
}

// ======== DOCUMENT NODE ========

// Document represents a single JML file
type Document struct {
	Doctype     DocType
	Position    Position
	Identifier  *Identifier
	Imports     []*ImportStatement
	Properties  []*PropertyDeclaration
	States      []*StateDeclaration
	RootElement *ElementNode
	ScriptBlock *ScriptBlock // FUTURE
	StyleBlock  *StyleBlock  // FUTURE
}

func (d *Document) Pos() Position          { return d.Position }
func (d *Document) Accept(visitor Visitor) { visitor.VisitDocument(d) }
func (d *Document) String() string {
	return fmt.Sprintf("Document{%s %s}", d.Doctype, d.Identifier.Name)
}

// ===== IMPORT STATEMENT ======

// ImportStatement represents an import declaration
type ImportStatement struct {
	Position
	Kind       ImportKind
	Identifier *Identifier
	Path       *Literal
}

func (i *ImportStatement) Pos() Position          { return i.Position }
func (i *ImportStatement) Accept(visitor Visitor) { visitor.VisitImportStatement(i) }
func (i *ImportStatement) String() string {
	return fmt.Sprintf("Import{%s %s from %s}", i.Kind, i.Identifier.Name, i.Path.Value)
}

// ========== DECLARATIONS ==========

type PropertyDeclaration struct {
	Position
	Name         *Identifier
	Type         string
	DefaultValue Expression
	Optional     bool
}

func (p *PropertyDeclaration) Pos() Position          { return p.Position }
func (p *PropertyDeclaration) Accept(visitor Visitor) { visitor.VisitPropertyDeclaration(p) }
func (p *PropertyDeclaration) String() string {
	optional := ""
	if p.Optional {
		optional = "?"
	}
	return fmt.Sprintf("Property{%s%s: %s}", p.Name.Name, optional, p.Type)
}

type StateDeclaration struct {
	Position
	Name         *Identifier
	Type         string
	DefaultValue Expression
}

func (s *StateDeclaration) Pos() Position          { return s.Position }
func (s *StateDeclaration) Accept(visitor Visitor) { visitor.VisitStateDeclaration(s) }
func (s *StateDeclaration) String() string {
	return fmt.Sprintf("State{%s: %s}", s.Name.Name, s.Type)
}

// ========== ELEMENT NODE ==========

// ElementNode represents a Jml tag structure
type ElementNode struct {
	Position   Position
	Tag        *Identifier
	Attributes map[string]Expression
	Children   []*ElementNode
}

func (e *ElementNode) Pos() Position          { return e.Position }
func (e *ElementNode) Accept(visitor Visitor) { visitor.VisitElementNode(e) }
func (e *ElementNode) String() string {
	return fmt.Sprintf("Element{%s}", e.Tag.Name)
}

// ========== STYLE AND SCRIPT BLOCKS ==========

// StyleBlock represents a style block (FUTURE)
type StyleBlock struct {
	Position
	Content string
}

func (s *StyleBlock) Pos() Position          { return s.Position }
func (s *StyleBlock) Accept(visitor Visitor) { visitor.VisitStyleBlock(s) }
func (s *StyleBlock) String() string {
	return "StyleBlock{...}"
}

// ScriptBlock represents a script block (FUTURE)
type ScriptBlock struct {
	Position
	Content string
}

func (s *ScriptBlock) Pos() Position          { return s.Position }
func (s *ScriptBlock) Accept(visitor Visitor) { visitor.VisitScriptBlock(s) }
func (s *ScriptBlock) String() string {
	return "ScriptBlock{...}"
}

// ========== EXPRESSIONS ==========

// Literal represents a literal value
type Literal struct {
	Position
	Value interface{}
	Kind  string // "string", "number", "boolean"
}

func (l *Literal) Pos() Position          { return l.Position }
func (l *Literal) Accept(visitor Visitor) { visitor.VisitLiteral(l) }
func (l *Literal) ExprType() ExprType     { return ExprLiteral }
func (l *Literal) String() string {
	return fmt.Sprintf("Literal{%v}", l.Value)
}

// Identifier represents an identifier
type Identifier struct {
	Position Position
	Name     string
}

func (i *Identifier) Pos() Position          { return i.Position }
func (i *Identifier) Accept(visitor Visitor) { visitor.VisitIdentifier(i) }
func (i *Identifier) ExprType() ExprType     { return ExprIdentifier }
func (i *Identifier) String() string {
	return fmt.Sprintf("Identifier{%s}", i.Name)
}

// Binding represents a binding like props.name or state.count
type Binding struct {
	Position Position
	Object   *Identifier
	Property *Identifier
}

func (b *Binding) Pos() Position          { return b.Position }
func (b *Binding) Accept(visitor Visitor) { visitor.VisitBinding(b) }
func (b *Binding) ExprType() ExprType     { return ExprBinding }
func (b *Binding) String() string {
	return fmt.Sprintf("Binding{%s.%s}", b.Object.Name, b.Property.Name)
}

type FunctionCall struct {
	Position  Position
	Function  *Identifier
	Arguments []Expression
}

func (f *FunctionCall) Pos() Position          { return f.Position }
func (f *FunctionCall) Accept(visitor Visitor) { visitor.VisitFunctionCall(f) }
func (f *FunctionCall) ExprType() ExprType     { return ExprFunctionCall }
func (f *FunctionCall) String() string {
	return fmt.Sprintf("FunctionCall{%s()}", f.Function.Name)
}

// LambdaExpression represents a lambda/arrow function
type LambdaExpression struct {
	Position   Position
	Parameters []*Identifier
	Body       string // For now, store as string - could be expanded to an expression tree
}

func (l *LambdaExpression) Pos() Position          { return l.Position }
func (l *LambdaExpression) Accept(visitor Visitor) { visitor.VisitLambdaExpression(l) }
func (l *LambdaExpression) ExprType() ExprType     { return ExprLambda }
func (l *LambdaExpression) String() string {
	return "Lambda{...}"
}

// ========== FACTORY METHODS ==========

// NewProgram creates a new Program node
func NewProgram(pos Position, documents []*Document) *Program {
	return &Program{
		Position:  pos,
		Documents: documents,
	}
}

// NewDocument creates a new Document node
func NewDocument(pos Position, doctype DocType, identifier *Identifier) *Document {
	return &Document{
		Position:    pos,
		Doctype:     doctype,
		Identifier:  identifier,
		Imports:     make([]*ImportStatement, 0),
		Properties:  make([]*PropertyDeclaration, 0),
		States:      make([]*StateDeclaration, 0),
		RootElement: nil,
		ScriptBlock: nil,
		StyleBlock:  nil,
	}
}

// NewImportStatement creates a new ImportStatement node
func NewImportStatement(pos Position, kind ImportKind, identifier *Identifier, path *Literal) *ImportStatement {
	return &ImportStatement{
		Position:   pos,
		Kind:       kind,
		Identifier: identifier,
		Path:       path,
	}
}

// NewPropertyDeclaration creates a new PropertyDeclaration node
func NewPropertyDeclaration(pos Position, name *Identifier, typeStr string, defaultValue Expression, optional bool) *PropertyDeclaration {
	return &PropertyDeclaration{
		Position:     pos,
		Name:         name,
		Type:         typeStr,
		DefaultValue: defaultValue,
		Optional:     optional,
	}
}

// NewStateDeclaration creates a new StateDeclaration node
func NewStateDeclaration(pos Position, name *Identifier, typeStr string, defaultValue Expression) *StateDeclaration {
	return &StateDeclaration{
		Position:     pos,
		Name:         name,
		Type:         typeStr,
		DefaultValue: defaultValue,
	}
}

// NewElementNode creates a new ElementNode
func NewElementNode(pos Position, tag *Identifier, attributes map[string]Expression, children []*ElementNode) *ElementNode {
	if attributes == nil {
		attributes = make(map[string]Expression)
	}
	if children == nil {
		children = make([]*ElementNode, 0)
	}
	return &ElementNode{
		Position:   pos,
		Tag:        tag,
		Attributes: attributes,
		Children:   children,
	}
}

// NewLiteral creates a new Literal node
func NewLiteral(pos Position, value interface{}, kind string) *Literal {
	return &Literal{
		Position: pos,
		Value:    value,
		Kind:     kind,
	}
}

// NewIdentifier creates a new Identifier node
func NewIdentifier(pos Position, name string) *Identifier {
	return &Identifier{
		Position: pos,
		Name:     name,
	}
}

// NewBinding creates a new Binding node
func NewBinding(pos Position, object *Identifier, property *Identifier) *Binding {
	return &Binding{
		Position: pos,
		Object:   object,
		Property: property,
	}
}

// NewFunctionCall creates a new FunctionCall node
func NewFunctionCall(pos Position, function *Identifier, arguments []Expression) *FunctionCall {
	if arguments == nil {
		arguments = make([]Expression, 0)
	}
	return &FunctionCall{
		Position:  pos,
		Function:  function,
		Arguments: arguments,
	}
}

// NewLambdaExpression creates a new LambdaExpression node
func NewLambdaExpression(pos Position, parameters []*Identifier, body string) *LambdaExpression {
	if parameters == nil {
		parameters = make([]*Identifier, 0)
	}
	return &LambdaExpression{
		Position:   pos,
		Parameters: parameters,
		Body:       body,
	}
}

// ========== WALKER FUNCTION ==========

// Walk traverses the AST and calls the visitor for each node
func Walk(node Node, visitor Visitor) {
	if node == nil {
		return
	}

	node.Accept(visitor)

	// Handle specific node types for deep traversal
	switch n := node.(type) {
	case *Program:
		for _, doc := range n.Documents {
			Walk(doc, visitor)
		}
	case *Document:
		for _, imp := range n.Imports {
			Walk(imp, visitor)
		}
		for _, prop := range n.Properties {
			Walk(prop, visitor)
		}
		for _, state := range n.States {
			Walk(state, visitor)
		}
		if n.RootElement != nil {
			Walk(n.RootElement, visitor)
		}
		if n.ScriptBlock != nil {
			Walk(n.ScriptBlock, visitor)
		}
		if n.StyleBlock != nil {
			Walk(n.StyleBlock, visitor)
		}
	case *ElementNode:
		for _, attr := range n.Attributes {
			Walk(attr, visitor)
		}
		for _, child := range n.Children {
			Walk(child, visitor)
		}
	case *PropertyDeclaration:
		if n.Name != nil {
			Walk(n.Name, visitor)
		}
		if n.DefaultValue != nil {
			Walk(n.DefaultValue, visitor)
		}
	case *StateDeclaration:
		if n.Name != nil {
			Walk(n.Name, visitor)
		}
		if n.DefaultValue != nil {
			Walk(n.DefaultValue, visitor)
		}
	case *ImportStatement:
		if n.Identifier != nil {
			Walk(n.Identifier, visitor)
		}
		if n.Path != nil {
			Walk(n.Path, visitor)
		}
	case *FunctionCall:
		if n.Function != nil {
			Walk(n.Function, visitor)
		}
		for _, arg := range n.Arguments {
			Walk(arg, visitor)
		}
	case *Binding:
		if n.Object != nil {
			Walk(n.Object, visitor)
		}
		if n.Property != nil {
			Walk(n.Property, visitor)
		}
	case *LambdaExpression:
		for _, param := range n.Parameters {
			Walk(param, visitor)
		}
	}
}
