package ast

import (
	"fmt"
)

// =================================================================================================
// == Base Interfaces
// =================================================================================================

// Position represents a location in a source file.
type Position struct {
	Line   int
	Column int
	File   string
}

func (p Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Column)
}

// Node is the base interface for all AST nodes.
type Node interface {
	Pos() Position
	String() string
	Accept(v Visitor)
}

// Statement is a node that represents a statement.
type Statement interface {
	Node
	statementNode()
}

// Declaration is a node that represents a declaration.
type Declaration interface {
	Statement
	declarationNode()
}

// Expression is a node that represents an expression.
type Expression interface {
	Node
	expressionNode()
}

// =================================================================================================
// == Visitor Interface
// =================================================================================================

// Visitor defines the interface for an AST visitor.
// Each Visit method is responsible for visiting the children of the node.
type Visitor interface {
	VisitProgram(n *Program)
	VisitDocument(n *Document)

	// Declarations
	VisitImportDeclaration(n *ImportDeclaration)
	VisitExportDeclaration(n *ExportDeclaration)
	VisitVariableDeclaration(n *VariableDeclaration)
	VisitFunctionDeclaration(n *FunctionDeclaration)
	VisitClassDeclaration(n *ClassDeclaration)
	VisitInterfaceDeclaration(n *InterfaceDeclaration)
	VisitTypeAliasDeclaration(n *TypeAliasDeclaration)
	VisitEnumDeclaration(n *EnumDeclaration)
	VisitPropertyDeclaration(n *PropertyDeclaration)
	VisitStateDeclaration(n *StateDeclaration)

	// Statements
	VisitBlockStatement(n *BlockStatement)
	VisitExpressionStatement(n *ExpressionStatement)
	VisitIfStatement(n *IfStatement)
	VisitForStatement(n *ForStatement)
	VisitForInStatement(n *ForInStatement)
	VisitWhileStatement(n *WhileStatement)
	VisitReturnStatement(n *ReturnStatement)
	VisitBreakStatement(n *BreakStatement)
	VisitContinueStatement(n *ContinueStatement)
	VisitThrowStatement(n *ThrowStatement)
	VisitTryStatement(n *TryStatement)

	// Expressions
	VisitIdentifier(n *Identifier)
	VisitLiteral(n *Literal)
	VisitArrayLiteral(n *ArrayLiteral)
	VisitObjectLiteral(n *ObjectLiteral)
	VisitFunctionExpression(n *FunctionExpression)
	VisitArrowFunctionExpression(n *ArrowFunctionExpression)
	VisitUnaryExpression(n *UnaryExpression)
	VisitBinaryExpression(n *BinaryExpression)
	VisitConditionalExpression(n *ConditionalExpression)
	VisitUpdateExpression(n *UpdateExpression)
	VisitMemberExpression(n *MemberExpression)
	VisitCallExpression(n *CallExpression)
	VisitNewExpression(n *NewExpression)
	VisitThisExpression(n *ThisExpression)
	VisitSuperExpression(n *SuperExpression)
	VisitTemplateLiteral(n *TemplateLiteral)

	// JML Specific
	VisitComponentElement(n *ComponentElement)
	VisitComponentProperty(n *ComponentProperty)
	VisitForLoop(n *ForLoop)
	VisitIfCondition(n *IfCondition)

	// Types
	VisitTypeAnnotation(n *TypeAnnotation)
	VisitTypeReference(n *TypeReference)
	VisitObjectType(n *ObjectType)
}

// =================================================================================================
// == Top-Level Nodes
// =================================================================================================

// Program represents the entire JML project, containing multiple documents.
type Program struct {
	Position
	Documents []*Document
}

func (n *Program) Pos() Position    { return n.Position }
func (n *Program) Accept(v Visitor) { v.VisitProgram(n) }
func (n *Program) String() string   { return "Program" }

// Document represents a single JML file.
type Document struct {
	Position
	DocType    DocType
	Name       *Identifier
	Body       []Statement // Can contain imports, exports, declarations, elements
	SourceFile string
}

type DocType string

const (
	DocTypePage      DocType = "page"
	DocTypeComponent DocType = "component"
)

func (n *Document) Pos() Position    { return n.Position }
func (n *Document) Accept(v Visitor) { v.VisitDocument(n) }
func (n *Document) String() string   { return fmt.Sprintf("Document(%s)", n.Name) }

// =================================================================================================
// == Declarations
// =================================================================================================

// ImportDeclaration represents an import statement. e.g. `import { a } from 'b'`
type ImportDeclaration struct {
	Position
	Specifiers []Node // Can be *ImportSpecifier, *ImportDefaultSpecifier, *ImportNamespaceSpecifier
	Source     *Literal
	IsBrowser  bool // for `import browser`
}

func (n *ImportDeclaration) Pos() Position         { return n.Position }
func (n *ImportDeclaration) Accept(v Visitor)      { v.VisitImportDeclaration(n) }
func (n *ImportDeclaration) String() string        { return "ImportDeclaration" }
func (n *ImportDeclaration) statementNode()        {}
func (n *ImportDeclaration) declarationNode()      {}
func (n *ImportDeclaration) componentBodyElement() {}

type ImportSpecifier struct {
	Position
	Local    *Identifier // The name of the imported binding
	Imported *Identifier // The name of the binding as it is exported
}

type ImportDefaultSpecifier struct {
	Position
	Local *Identifier
}

type ImportNamespaceSpecifier struct {
	Position
	Local *Identifier
}

// ExportDeclaration represents an export statement. e.g. `export const a = 1`
type ExportDeclaration struct {
	Position
	Declaration Declaration // The declaration to export
	Default     bool        // `export default`
}

func (n *ExportDeclaration) Pos() Position         { return n.Position }
func (n *ExportDeclaration) Accept(v Visitor)      { v.VisitExportDeclaration(n) }
func (n *ExportDeclaration) String() string        { return "ExportDeclaration" }
func (n *ExportDeclaration) statementNode()        {}
func (n *ExportDeclaration) declarationNode()      {}
func (n *ExportDeclaration) componentBodyElement() {}

// VariableDeclaration represents a `let`, `const`, or `var` declaration.
type VariableDeclaration struct {
	Position
	Kind          string // "var", "let", "const"
	Declarations  []*VariableDeclarator
	InjectedState bool // Special case for `state` variables
}

func (n *VariableDeclaration) Pos() Position         { return n.Position }
func (n *VariableDeclaration) Accept(v Visitor)      { v.VisitVariableDeclaration(n) }
func (n *VariableDeclaration) String() string        { return "VariableDeclaration" }
func (n *VariableDeclaration) statementNode()        {}
func (n *VariableDeclaration) declarationNode()      {}
func (n *VariableDeclaration) componentBodyElement() {}

type VariableDeclarator struct {
	Position
	ID   *Identifier
	Init Expression
}

// FunctionDeclaration represents a function declaration.
type FunctionDeclaration struct {
	Position
	ID         *Identifier
	Params     []*Identifier
	Body       *BlockStatement
	ReturnType *TypeAnnotation
	Async      bool
}

func (n *FunctionDeclaration) Pos() Position         { return n.Position }
func (n *FunctionDeclaration) Accept(v Visitor)      { v.VisitFunctionDeclaration(n) }
func (n *FunctionDeclaration) String() string        { return fmt.Sprintf("FunctionDeclaration(%s)", n.ID) }
func (n *FunctionDeclaration) statementNode()        {}
func (n *FunctionDeclaration) declarationNode()      {}
func (n *FunctionDeclaration) componentBodyElement() {}

// ClassDeclaration represents a class declaration.
type ClassDeclaration struct {
	Position
	ID         *Identifier
	SuperClass Expression
	Body       *ClassBody
}

func (n *ClassDeclaration) Pos() Position         { return n.Position }
func (n *ClassDeclaration) Accept(v Visitor)      { v.VisitClassDeclaration(n) }
func (n *ClassDeclaration) String() string        { return fmt.Sprintf("ClassDeclaration(%s)", n.ID) }
func (n *ClassDeclaration) statementNode()        {}
func (n *ClassDeclaration) declarationNode()      {}
func (n *ClassDeclaration) componentBodyElement() {}

type ClassBody struct {
	Position
	Body []Node // MethodDefinition, PropertyDefinition
}

// InterfaceDeclaration represents an interface declaration.
type InterfaceDeclaration struct {
	Position
	ID      *Identifier
	Body    *ObjectType
	Extends []*TypeReference
}

func (n *InterfaceDeclaration) Pos() Position    { return n.Position }
func (n *InterfaceDeclaration) Accept(v Visitor) { v.VisitInterfaceDeclaration(n) }
func (n *InterfaceDeclaration) String() string   { return fmt.Sprintf("InterfaceDeclaration(%s)", n.ID) }
func (n *InterfaceDeclaration) statementNode()   {}
func (n *InterfaceDeclaration) declarationNode() {}

// TypeAliasDeclaration represents a type alias. e.g. `type A = string`
type TypeAliasDeclaration struct {
	Position
	ID   *Identifier
	Type Node // A Type node
}

func (n *TypeAliasDeclaration) Pos() Position    { return n.Position }
func (n *TypeAliasDeclaration) Accept(v Visitor) { v.VisitTypeAliasDeclaration(n) }
func (n *TypeAliasDeclaration) String() string   { return "TypeAliasDeclaration" }
func (n *TypeAliasDeclaration) statementNode()   {}
func (n *TypeAliasDeclaration) declarationNode() {}

// EnumDeclaration represents an enum declaration.
type EnumDeclaration struct {
	Position
	ID      *Identifier
	Members []*EnumMember
}

func (n *EnumDeclaration) Pos() Position    { return n.Position }
func (n *EnumDeclaration) Accept(v Visitor) { v.VisitEnumDeclaration(n) }
func (n *EnumDeclaration) String() string   { return "EnumDeclaration" }
func (n *EnumDeclaration) statementNode()   {}
func (n *EnumDeclaration) declarationNode() {}

type EnumMember struct {
	Position
	ID   *Identifier
	Init Expression
}

// PropertyDeclaration is a JML-specific declaration for component properties.
type PropertyDeclaration struct {
	Position
	ID           *Identifier
	Type         *TypeAnnotation
	DefaultValue Expression
	Optional     bool
}

func (n *PropertyDeclaration) Pos() Position    { return n.Position }
func (n *PropertyDeclaration) Accept(v Visitor) { v.VisitPropertyDeclaration(n) }
func (n *PropertyDeclaration) String() string   { return "PropertyDeclaration" }
func (n *PropertyDeclaration) statementNode()   {}
func (n *PropertyDeclaration) declarationNode() {}

// StateDeclaration is a JML-specific declaration for component state.
type StateDeclaration struct {
	Position
	ID           *Identifier
	Type         *TypeAnnotation
	DefaultValue Expression
}

func (n *StateDeclaration) Pos() Position    { return n.Position }
func (n *StateDeclaration) Accept(v Visitor) { v.VisitStateDeclaration(n) }
func (n *StateDeclaration) String() string   { return "StateDeclaration" }
func (n *StateDeclaration) statementNode()   {}
func (n *StateDeclaration) declarationNode() {}

// =================================================================================================
// == Statements
// =================================================================================================

// BlockStatement represents a block of statements. e.g. `{ ... }`
type BlockStatement struct {
	Position
	List []Statement
}

func (n *BlockStatement) Pos() Position    { return n.Position }
func (n *BlockStatement) Accept(v Visitor) { v.VisitBlockStatement(n) }
func (n *BlockStatement) String() string   { return "BlockStatement" }
func (n *BlockStatement) statementNode()   {}

// ExpressionStatement represents an expression that is used as a statement.
type ExpressionStatement struct {
	Position
	Expression Expression
}

func (n *ExpressionStatement) Pos() Position    { return n.Position }
func (n *ExpressionStatement) Accept(v Visitor) { v.VisitExpressionStatement(n) }
func (n *ExpressionStatement) String() string   { return "ExpressionStatement" }
func (n *ExpressionStatement) statementNode()   {}

// IfStatement represents an if-else statement.
type IfStatement struct {
	Position
	Test       Expression
	Consequent Statement
	Alternate  Statement
}

func (n *IfStatement) Pos() Position    { return n.Position }
func (n *IfStatement) Accept(v Visitor) { v.VisitIfStatement(n) }
func (n *IfStatement) String() string   { return "IfStatement" }
func (n *IfStatement) statementNode()   {}

// ForStatement represents a for loop. e.g. `for (let i = 0; i < 10; i++)`
type ForStatement struct {
	Position
	Init   Node // VariableDeclaration or Expression
	Test   Expression
	Update Expression
	Body   Statement
}

func (n *ForStatement) Pos() Position    { return n.Position }
func (n *ForStatement) Accept(v Visitor) { v.VisitForStatement(n) }
func (n *ForStatement) String() string   { return "ForStatement" }
func (n *ForStatement) statementNode()   {}

// ForInStatement represents a for-in or for-of loop.
type ForInStatement struct {
	Position
	Left  Node // VariableDeclaration or Expression
	Right Expression
	Body  Statement
	Of    bool // true for "for-of", false for "for-in"
}

func (n *ForInStatement) Pos() Position    { return n.Position }
func (n *ForInStatement) Accept(v Visitor) { v.VisitForInStatement(n) }
func (n *ForInStatement) String() string   { return "ForInStatement" }
func (n *ForInStatement) statementNode()   {}

// WhileStatement represents a while loop.
type WhileStatement struct {
	Position
	Test Expression
	Body Statement
}

func (n *WhileStatement) Pos() Position    { return n.Position }
func (n *WhileStatement) Accept(v Visitor) { v.VisitWhileStatement(n) }
func (n *WhileStatement) String() string   { return "WhileStatement" }
func (n *WhileStatement) statementNode()   {}

// ReturnStatement represents a return statement.
type ReturnStatement struct {
	Position
	Argument Expression
}

func (n *ReturnStatement) Pos() Position    { return n.Position }
func (n *ReturnStatement) Accept(v Visitor) { v.VisitReturnStatement(n) }
func (n *ReturnStatement) String() string   { return "ReturnStatement" }
func (n *ReturnStatement) statementNode()   {}

// BreakStatement represents a break statement.
type BreakStatement struct {
	Position
	Label *Identifier
}

func (n *BreakStatement) Pos() Position    { return n.Position }
func (n *BreakStatement) Accept(v Visitor) { v.VisitBreakStatement(n) }
func (n *BreakStatement) String() string   { return "BreakStatement" }
func (n *BreakStatement) statementNode()   {}

// ContinueStatement represents a continue statement.
type ContinueStatement struct {
	Position
	Label *Identifier
}

func (n *ContinueStatement) Pos() Position    { return n.Position }
func (n *ContinueStatement) Accept(v Visitor) { v.VisitContinueStatement(n) }
func (n *ContinueStatement) String() string   { return "ContinueStatement" }
func (n *ContinueStatement) statementNode()   {}

// ThrowStatement represents a throw statement.
type ThrowStatement struct {
	Position
	Argument Expression
}

func (n *ThrowStatement) Pos() Position    { return n.Position }
func (n *ThrowStatement) Accept(v Visitor) { v.VisitThrowStatement(n) }
func (n *ThrowStatement) String() string   { return "ThrowStatement" }
func (n *ThrowStatement) statementNode()   {}

// TryStatement represents a try-catch-finally statement.
type TryStatement struct {
	Position
	Block     *BlockStatement
	Handler   *CatchClause
	Finalizer *BlockStatement
}

func (n *TryStatement) Pos() Position    { return n.Position }
func (n *TryStatement) Accept(v Visitor) { v.VisitTryStatement(n) }
func (n *TryStatement) String() string   { return "TryStatement" }
func (n *TryStatement) statementNode()   {}

type CatchClause struct {
	Position
	Param *Identifier
	Body  *BlockStatement
}

// =================================================================================================
// == Expressions
// =================================================================================================

// Identifier represents an identifier.
type Identifier struct {
	Position
	Name string
}

func (n *Identifier) Pos() Position    { return n.Position }
func (n *Identifier) Accept(v Visitor) { v.VisitIdentifier(n) }
func (n *Identifier) String() string   { return n.Name }
func (n *Identifier) expressionNode()  {}

// Literal represents a literal value.
type Literal struct {
	Position
	Kind  string // "string", "number", "boolean", "null", "undefined", "regexp"
	Value string
}

func (n *Literal) Pos() Position    { return n.Position }
func (n *Literal) Accept(v Visitor) { v.VisitLiteral(n) }
func (n *Literal) String() string   { return n.Value }
func (n *Literal) expressionNode()  {}

// ArrayLiteral represents an array literal.
type ArrayLiteral struct {
	Position
	Elements []Expression
}

func (n *ArrayLiteral) Pos() Position    { return n.Position }
func (n *ArrayLiteral) Accept(v Visitor) { v.VisitArrayLiteral(n) }
func (n *ArrayLiteral) String() string   { return "ArrayLiteral" }
func (n *ArrayLiteral) expressionNode()  {}

// ObjectLiteral represents an object literal.
type ObjectLiteral struct {
	Position
	Properties []*Property
}

func (n *ObjectLiteral) Pos() Position    { return n.Position }
func (n *ObjectLiteral) Accept(v Visitor) { v.VisitObjectLiteral(n) }
func (n *ObjectLiteral) String() string   { return "ObjectLiteral" }
func (n *ObjectLiteral) expressionNode()  {}

type Property struct {
	Position
	Key   Expression
	Value Expression
	Kind  string // "init", "get", "set"
}

// FunctionExpression represents a function expression.
type FunctionExpression struct {
	Position
	ID     *Identifier
	Params []*Identifier
	Body   *BlockStatement
	Async  bool
}

func (n *FunctionExpression) Pos() Position    { return n.Position }
func (n *FunctionExpression) Accept(v Visitor) { v.VisitFunctionExpression(n) }
func (n *FunctionExpression) String() string   { return "FunctionExpression" }
func (n *FunctionExpression) expressionNode()  {}

// ArrowFunctionExpression represents an arrow function. e.g. `(a, b) => a + b`
type ArrowFunctionExpression struct {
	Position
	Params []*Identifier
	Body   Node // BlockStatement or Expression
	Async  bool
}

func (n *ArrowFunctionExpression) Pos() Position    { return n.Position }
func (n *ArrowFunctionExpression) Accept(v Visitor) { v.VisitArrowFunctionExpression(n) }
func (n *ArrowFunctionExpression) String() string   { return "ArrowFunctionExpression" }
func (n *ArrowFunctionExpression) expressionNode()  {}

// UnaryExpression represents a unary operation. e.g. `!a`, `typeof b`
type UnaryExpression struct {
	Position
	Operator string
	Argument Expression
	Prefix   bool
}

func (n *UnaryExpression) Pos() Position    { return n.Position }
func (n *UnaryExpression) Accept(v Visitor) { v.VisitUnaryExpression(n) }
func (n *UnaryExpression) String() string   { return fmt.Sprintf("Unary(%s)", n.Operator) }
func (n *UnaryExpression) expressionNode()  {}

// BinaryExpression represents a binary operation. e.g. `a + b`
type BinaryExpression struct {
	Position
	Operator string
	Left     Expression
	Right    Expression
}

func (n *BinaryExpression) Pos() Position    { return n.Position }
func (n *BinaryExpression) Accept(v Visitor) { v.VisitBinaryExpression(n) }
func (n *BinaryExpression) String() string   { return fmt.Sprintf("Binary(%s)", n.Operator) }
func (n *BinaryExpression) expressionNode()  {}

// ConditionalExpression represents a ternary operation. e.g. `a ? b : c`
type ConditionalExpression struct {
	Position
	Test       Expression
	Consequent Expression
	Alternate  Expression
}

func (n *ConditionalExpression) Pos() Position    { return n.Position }
func (n *ConditionalExpression) Accept(v Visitor) { v.VisitConditionalExpression(n) }
func (n *ConditionalExpression) String() string   { return "ConditionalExpression" }
func (n *ConditionalExpression) expressionNode()  {}

// UpdateExpression represents an update operation. e.g. `a++`, `--b`
type UpdateExpression struct {
	Position
	Operator string
	Argument Expression
	Prefix   bool
}

func (n *UpdateExpression) Pos() Position    { return n.Position }
func (n *UpdateExpression) Accept(v Visitor) { v.VisitUpdateExpression(n) }
func (n *UpdateExpression) String() string   { return "UpdateExpression" }
func (n *UpdateExpression) expressionNode()  {}

// MemberExpression represents a member access. e.g. `a.b`, `a[b]`
type MemberExpression struct {
	Position
	Object   Expression
	Property Expression
	Computed bool // true for `a[b]`, false for `a.b`
}

func (n *MemberExpression) Pos() Position    { return n.Position }
func (n *MemberExpression) Accept(v Visitor) { v.VisitMemberExpression(n) }
func (n *MemberExpression) String() string   { return "MemberExpression" }
func (n *MemberExpression) expressionNode()  {}

// CallExpression represents a function call. e.g. `a(b, c)`
type CallExpression struct {
	Position
	Callee    Expression
	Arguments []Expression
}

func (n *CallExpression) Pos() Position    { return n.Position }
func (n *CallExpression) Accept(v Visitor) { v.VisitCallExpression(n) }
func (n *CallExpression) String() string   { return "CallExpression" }
func (n *CallExpression) expressionNode()  {}

// NewExpression represents a new expression. e.g. `new A()`
type NewExpression struct {
	Position
	Callee    Expression
	Arguments []Expression
}

func (n *NewExpression) Pos() Position    { return n.Position }
func (n *NewExpression) Accept(v Visitor) { v.VisitNewExpression(n) }
func (n *NewExpression) String() string   { return "NewExpression" }
func (n *NewExpression) expressionNode()  {}

// ThisExpression represents the `this` keyword.
type ThisExpression struct {
	Position
}

func (n *ThisExpression) Pos() Position    { return n.Position }
func (n *ThisExpression) Accept(v Visitor) { v.VisitThisExpression(n) }
func (n *ThisExpression) String() string   { return "this" }
func (n *ThisExpression) expressionNode()  {}

// SuperExpression represents the `super` keyword.
type SuperExpression struct {
	Position
}

func (n *SuperExpression) Pos() Position    { return n.Position }
func (n *SuperExpression) Accept(v Visitor) { v.VisitSuperExpression(n) }
func (n *SuperExpression) String() string   { return "super" }
func (n *SuperExpression) expressionNode()  {}

// TemplateLiteral represents a template literal. e.g. `hello ${world}`
type TemplateLiteral struct {
	Position
	Quasis      []*TemplateElement
	Expressions []Expression
}

func (n *TemplateLiteral) Pos() Position    { return n.Position }
func (n *TemplateLiteral) Accept(v Visitor) { v.VisitTemplateLiteral(n) }
func (n *TemplateLiteral) String() string   { return "TemplateLiteral" }
func (n *TemplateLiteral) expressionNode()  {}

type TemplateElement struct {
	Position
	Value string
	Tail  bool
}

// =================================================================================================
// == JML Specific Nodes
// =================================================================================================

// ComponentBodyElement is a marker interface for nodes that can appear in a component body.
type ComponentBodyElement interface {
	Node
	componentBodyElement()
}

// ComponentElement represents a JML component instantiation. e.g. `MyComponent { ... }`
type ComponentElement struct {
	Position
	Tag  *Identifier
	Body []ComponentBodyElement
}

func (n *ComponentElement) Pos() Position         { return n.Position }
func (n *ComponentElement) Accept(v Visitor)      { v.VisitComponentElement(n) }
func (n *ComponentElement) String() string        { return fmt.Sprintf("ComponentElement(%s)", n.Tag) }
func (n *ComponentElement) statementNode()        {}
func (n *ComponentElement) componentBodyElement() {}

// ComponentProperty represents a property assignment in a component body. e.g. `prop: "value"`
type ComponentProperty struct {
	Position
	Key   *Identifier
	Value Expression
}

func (n *ComponentProperty) Pos() Position         { return n.Position }
func (n *ComponentProperty) Accept(v Visitor)      { v.VisitComponentProperty(n) }
func (n *ComponentProperty) String() string        { return "ComponentProperty" }
func (n *ComponentProperty) componentBodyElement() {}

// ForLoop represents a for loop in a component body. e.g. `for (item in items) { ... }`
type ForLoop struct {
	Position
	Variable *Identifier
	Source   Expression
	Body     []ComponentBodyElement
}

func (n *ForLoop) Pos() Position         { return n.Position }
func (n *ForLoop) Accept(v Visitor)      { v.VisitForLoop(n) }
func (n *ForLoop) String() string        { return "ForLoop" }
func (n *ForLoop) componentBodyElement() {}

// IfCondition represents an if-else block in a component body.
type IfCondition struct {
	Position
	Test       Expression
	Consequent []ComponentBodyElement
	Alternate  []ComponentBodyElement // Can be nil
}

func (n *IfCondition) Pos() Position         { return n.Position }
func (n *IfCondition) Accept(v Visitor)      { v.VisitIfCondition(n) }
func (n *IfCondition) String() string        { return "IfCondition" }
func (n *IfCondition) componentBodyElement() {}

// =================================================================================================
// == Type Nodes
// =================================================================================================

// TypeAnnotation represents a type annotation. e.g. `: string`
type TypeAnnotation struct {
	Position
	Type Node // A Type node
}

func (n *TypeAnnotation) Pos() Position    { return n.Position }
func (n *TypeAnnotation) Accept(v Visitor) { v.VisitTypeAnnotation(n) }
func (n *TypeAnnotation) String() string   { return "TypeAnnotation" }

// TypeReference represents a reference to a type. e.g. `string`, `MyType`
type TypeReference struct {
	Position
	Name       *Identifier
	TypeParams []*TypeAnnotation
}

func (n *TypeReference) Pos() Position    { return n.Position }
func (n *TypeReference) Accept(v Visitor) { v.VisitTypeReference(n) }
func (n *TypeReference) String() string   { return "TypeReference" }

// ObjectType represents an object type literal. e.g. `{ a: string, b: number }`
type ObjectType struct {
	Position
	Members []Node // PropertySignature, MethodSignature, etc.
}

func (n *ObjectType) Pos() Position { return n.Position }

func (n *ObjectType) Accept(v Visitor) { v.VisitObjectType(n) }
