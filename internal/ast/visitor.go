package ast

// BaseVisitor provides default no-op implementations for all visitor methods
// Embed this in your concrete visitors to avoid implementing unused methods
type BaseVisitor struct{}

// DocumentVisitor methods

func (b *BaseVisitor) VisitProgram(*Program)                 {}
func (b *BaseVisitor) VisitDocument(*Document)               {}
func (b *BaseVisitor) VisitImportStatement(*ImportStatement) {}

// DeclarationVisitor methods

func (b *BaseVisitor) VisitPropertyDeclaration(*PropertyDeclaration) {}
func (b *BaseVisitor) VisitStateDeclaration(*StateDeclaration)       {}

// BlockVisitor methods

func (b *BaseVisitor) VisitStyleBlock(*StyleBlock)   {}
func (b *BaseVisitor) VisitScriptBlock(*ScriptBlock) {}
func (b *BaseVisitor) VisitElementNode(*ElementNode) {}

// ExpressionVisitor methods

func (b *BaseVisitor) VisitLiteral(*Literal)                   {}
func (b *BaseVisitor) VisitIdentifier(*Identifier)             {}
func (b *BaseVisitor) VisitBinding(*Binding)                   {}
func (b *BaseVisitor) VisitFunctionCall(*FunctionCall)         {}
func (b *BaseVisitor) VisitLambdaExpression(*LambdaExpression) {}
