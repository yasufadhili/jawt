package ast

// BaseVisitor is a no-op implementation of the Visitor interface.
// It can be embedded in other visitors to avoid having to implement all
// methods of the interface.
type BaseVisitor struct{}

// Top-Level

func (v *BaseVisitor) VisitProgram(n *Program)   {}
func (v *BaseVisitor) VisitDocument(n *Document) {}

// Declarations

func (v *BaseVisitor) VisitImportDeclaration(n *ImportDeclaration)       {}
func (v *BaseVisitor) VisitExportDeclaration(n *ExportDeclaration)       {}
func (v *BaseVisitor) VisitVariableDeclaration(n *VariableDeclaration)   {}
func (v *BaseVisitor) VisitFunctionDeclaration(n *FunctionDeclaration)   {}
func (v *BaseVisitor) VisitClassDeclaration(n *ClassDeclaration)         {}
func (v *BaseVisitor) VisitInterfaceDeclaration(n *InterfaceDeclaration) {}
func (v *BaseVisitor) VisitTypeAliasDeclaration(n *TypeAliasDeclaration) {}
func (v *BaseVisitor) VisitEnumDeclaration(n *EnumDeclaration)           {}
func (v *BaseVisitor) VisitPropertyDeclaration(n *PropertyDeclaration)   {}
func (v *BaseVisitor) VisitStateDeclaration(n *StateDeclaration)         {}

// Statements

func (v *BaseVisitor) VisitBlockStatement(n *BlockStatement)           {}
func (v *BaseVisitor) VisitExpressionStatement(n *ExpressionStatement) {}
func (v *BaseVisitor) VisitIfStatement(n *IfStatement)                 {}
func (v *BaseVisitor) VisitForStatement(n *ForStatement)               {}
func (v *BaseVisitor) VisitForInStatement(n *ForInStatement)           {}
func (v *BaseVisitor) VisitWhileStatement(n *WhileStatement)           {}
func (v *BaseVisitor) VisitReturnStatement(n *ReturnStatement)         {}
func (v *BaseVisitor) VisitBreakStatement(n *BreakStatement)           {}
func (v *BaseVisitor) VisitContinueStatement(n *ContinueStatement)     {}
func (v *BaseVisitor) VisitThrowStatement(n *ThrowStatement)           {}
func (v *BaseVisitor) VisitTryStatement(n *TryStatement)               {}

// Expressions

func (v *BaseVisitor) VisitIdentifier(n *Identifier)                           {}
func (v *BaseVisitor) VisitLiteral(n *Literal)                                 {}
func (v *BaseVisitor) VisitArrayLiteral(n *ArrayLiteral)                       {}
func (v *BaseVisitor) VisitObjectLiteral(n *ObjectLiteral)                     {}
func (v *BaseVisitor) VisitFunctionExpression(n *FunctionExpression)           {}
func (v *BaseVisitor) VisitArrowFunctionExpression(n *ArrowFunctionExpression) {}
func (v *BaseVisitor) VisitUnaryExpression(n *UnaryExpression)                 {}
func (v *BaseVisitor) VisitBinaryExpression(n *BinaryExpression)               {}
func (v *BaseVisitor) VisitConditionalExpression(n *ConditionalExpression)     {}
func (v *BaseVisitor) VisitUpdateExpression(n *UpdateExpression)               {}
func (v *BaseVisitor) VisitMemberExpression(n *MemberExpression)               {}
func (v *BaseVisitor) VisitCallExpression(n *CallExpression)                   {}
func (v *BaseVisitor) VisitNewExpression(n *NewExpression)                     {}
func (v *BaseVisitor) VisitThisExpression(n *ThisExpression)                   {}
func (v *BaseVisitor) VisitSuperExpression(n *SuperExpression)                 {}
func (v *BaseVisitor) VisitTemplateLiteral(n *TemplateLiteral)                 {}

// JML Specific

func (v *BaseVisitor) VisitComponentElement(n *ComponentElement)   {}
func (v *BaseVisitor) VisitComponentProperty(n *ComponentProperty) {}
func (v *BaseVisitor) VisitForLoop(n *ForLoop)                     {}
func (v *BaseVisitor) VisitIfCondition(n *IfCondition)             {}

// Types

func (v *BaseVisitor) VisitTypeAnnotation(n *TypeAnnotation) {}
func (v *BaseVisitor) VisitTypeReference(n *TypeReference)   {}
