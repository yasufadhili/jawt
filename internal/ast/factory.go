package ast

// =================================================================================================
// == Factory Methods
// =================================================================================================

// NewProgram creates a new Program node.
func NewProgram(pos Position, documents []*Document) *Program {
	return &Program{
		Position:  pos,
		Documents: documents,
	}
}

// NewDocument creates a new Document node.
func NewDocument(pos Position, docType DocType, name *Identifier, body []Statement, sourceFile string) *Document {
	return &Document{
		Position:   pos,
		DocType:    docType,
		Name:       name,
		Body:       body,
		SourceFile: sourceFile,
	}
}

// NewImportDeclaration creates a new ImportDeclaration node.
func NewImportDeclaration(pos Position, specifiers []Node, source *Literal, isBrowser bool) *ImportDeclaration {
	return &ImportDeclaration{
		Position:   pos,
		Specifiers: specifiers,
		Source:     source,
		IsBrowser:  isBrowser,
	}
}

// NewExportDeclaration creates a new ExportDeclaration node.
func NewExportDeclaration(pos Position, declaration Declaration, isDefault bool) *ExportDeclaration {
	return &ExportDeclaration{
		Position:    pos,
		Declaration: declaration,
		Default:     isDefault,
	}
}

func NewImportDefaultSpecifier(pos Position, local *Identifier) *ImportDefaultSpecifier {
	return &ImportDefaultSpecifier{
		Position: pos,
		Local:    local,
	}
}

func NewImportNamespaceSpecifier(pos Position, local *Identifier) *ImportNamespaceSpecifier {
	return &ImportNamespaceSpecifier{
		Position: pos,
	}
}

func NewImportSpecifier(pos Position, local *Identifier) *ImportSpecifier {
	return &ImportSpecifier{
		Position: pos,
		Local:    local,
	}
}

// NewVariableDeclaration creates a new VariableDeclaration node.
func NewVariableDeclaration(pos Position, kind string, declarations []*VariableDeclarator) *VariableDeclaration {
	return &VariableDeclaration{
		Position:     pos,
		Kind:         kind,
		Declarations: declarations,
	}
}

// NewVariableDeclarator creates a new VariableDeclarator node.
func NewVariableDeclarator(pos Position, id *Identifier, init Expression) *VariableDeclarator {
	return &VariableDeclarator{
		Position: pos,
		ID:       id,
		Init:     init,
	}
}

// NewFunctionDeclaration creates a new FunctionDeclaration node.
func NewFunctionDeclaration(pos Position, id *Identifier, params []*Identifier, body *BlockStatement, returnType *TypeAnnotation, async bool) *FunctionDeclaration {
	return &FunctionDeclaration{
		Position:   pos,
		ID:         id,
		Params:     params,
		Body:       body,
		ReturnType: returnType,
		Async:      async,
	}
}

// NewClassDeclaration creates a new ClassDeclaration node.
func NewClassDeclaration(pos Position, id *Identifier, superClass Expression, body *ClassBody) *ClassDeclaration {
	return &ClassDeclaration{
		Position:   pos,
		ID:         id,
		SuperClass: superClass,
		Body:       body,
	}
}

// NewInterfaceDeclaration creates a new InterfaceDeclaration node.
func NewInterfaceDeclaration(pos Position, id *Identifier, body *ObjectType, extends []*TypeReference) *InterfaceDeclaration {
	return &InterfaceDeclaration{
		Position: pos,
		ID:       id,
		Body:     body,
		Extends:  extends,
	}
}

// NewTypeAliasDeclaration creates a new TypeAliasDeclaration node.
func NewTypeAliasDeclaration(pos Position, id *Identifier, typeNode Node) *TypeAliasDeclaration {
	return &TypeAliasDeclaration{
		Position: pos,
		ID:       id,
		Type:     typeNode,
	}
}

// NewEnumDeclaration creates a new EnumDeclaration node.
func NewEnumDeclaration(pos Position, id *Identifier, members []*EnumMember) *EnumDeclaration {
	return &EnumDeclaration{
		Position: pos,
		ID:       id,
		Members:  members,
	}
}

// NewPropertyDeclaration creates a new PropertyDeclaration node.
func NewPropertyDeclaration(pos Position, id *Identifier, typeAnn *TypeAnnotation, defaultValue Expression, optional bool) *PropertyDeclaration {
	return &PropertyDeclaration{
		Position:     pos,
		ID:           id,
		Type:         typeAnn,
		DefaultValue: defaultValue,
		Optional:     optional,
	}
}

// NewStateDeclaration creates a new StateDeclaration node.
func NewStateDeclaration(pos Position, id *Identifier, typeAnn *TypeAnnotation, defaultValue Expression) *StateDeclaration {
	return &StateDeclaration{
		Position:     pos,
		ID:           id,
		Type:         typeAnn,
		DefaultValue: defaultValue,
	}
}

// NewBlockStatement creates a new BlockStatement node.
func NewBlockStatement(pos Position, list []Statement) *BlockStatement {
	return &BlockStatement{
		Position: pos,
		List:     list,
	}
}

// NewExpressionStatement creates a new ExpressionStatement node.
func NewExpressionStatement(pos Position, expr Expression) *ExpressionStatement {
	return &ExpressionStatement{
		Position:   pos,
		Expression: expr,
	}
}

// NewIfStatement creates a new IfStatement node.
func NewIfStatement(pos Position, test Expression, consequent Statement, alternate Statement) *IfStatement {
	return &IfStatement{
		Position:   pos,
		Test:       test,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

// NewForStatement creates a new ForStatement node.
func NewForStatement(pos Position, init Node, test Expression, update Expression, body Statement) *ForStatement {
	return &ForStatement{
		Position: pos,
		Init:     init,
		Test:     test,
		Update:   update,
		Body:     body,
	}
}

// NewForInStatement creates a new ForInStatement node.
func NewForInStatement(pos Position, left Node, right Expression, body Statement, of bool) *ForInStatement {
	return &ForInStatement{
		Position: pos,
		Left:     left,
		Right:    right,
		Body:     body,
		Of:       of,
	}
}

// NewWhileStatement creates a new WhileStatement node.
func NewWhileStatement(pos Position, test Expression, body Statement) *WhileStatement {
	return &WhileStatement{
		Position: pos,
		Test:     test,
		Body:     body,
	}
}

// NewReturnStatement creates a new ReturnStatement node.
func NewReturnStatement(pos Position, arg Expression) *ReturnStatement {
	return &ReturnStatement{
		Position: pos,
		Argument: arg,
	}
}

// NewBreakStatement creates a new BreakStatement node.
func NewBreakStatement(pos Position, label *Identifier) *BreakStatement {
	return &BreakStatement{
		Position: pos,
		Label:    label,
	}
}

// NewContinueStatement creates a new ContinueStatement node.
func NewContinueStatement(pos Position, label *Identifier) *ContinueStatement {
	return &ContinueStatement{
		Position: pos,
		Label:    label,
	}
}

// NewThrowStatement creates a new ThrowStatement node.
func NewThrowStatement(pos Position, arg Expression) *ThrowStatement {
	return &ThrowStatement{
		Position: pos,
		Argument: arg,
	}
}

// NewTryStatement creates a new TryStatement node.
func NewTryStatement(pos Position, block *BlockStatement, handler *CatchClause, finalizer *BlockStatement) *TryStatement {
	return &TryStatement{
		Position:  pos,
		Block:     block,
		Handler:   handler,
		Finalizer: finalizer,
	}
}

// NewIdentifier creates a new Identifier node.
func NewIdentifier(pos Position, name string) *Identifier {
	return &Identifier{
		Position: pos,
		Name:     name,
	}
}

// NewLiteral creates a new Literal node.
func NewLiteral(pos Position, kind string, value string) *Literal {
	return &Literal{
		Position: pos,
		Kind:     kind,
		Value:    value,
	}
}

// NewArrayLiteral creates a new ArrayLiteral node.
func NewArrayLiteral(pos Position, elements []Expression) *ArrayLiteral {
	return &ArrayLiteral{
		Position: pos,
		Elements: elements,
	}
}

// NewObjectLiteral creates a new ObjectLiteral node.
func NewObjectLiteral(pos Position, properties []*Property) *ObjectLiteral {
	return &ObjectLiteral{
		Position:   pos,
		Properties: properties,
	}
}

// NewFunctionExpression creates a new FunctionExpression node.
func NewFunctionExpression(pos Position, id *Identifier, params []*Identifier, body *BlockStatement, async bool) *FunctionExpression {
	return &FunctionExpression{
		Position: pos,
		ID:       id,
		Params:   params,
		Body:     body,
		Async:    async,
	}
}

// NewArrowFunctionExpression creates a new ArrowFunctionExpression node.
func NewArrowFunctionExpression(pos Position, params []*Identifier, body Node, async bool) *ArrowFunctionExpression {
	return &ArrowFunctionExpression{
		Position: pos,
		Params:   params,
		Body:     body,
		Async:    async,
	}
}

// NewUnaryExpression creates a new UnaryExpression node.
func NewUnaryExpression(pos Position, op string, arg Expression, prefix bool) *UnaryExpression {
	return &UnaryExpression{
		Position: pos,
		Operator: op,
		Argument: arg,
		Prefix:   prefix,
	}
}

// NewBinaryExpression creates a new BinaryExpression node.
func NewBinaryExpression(pos Position, op string, left Expression, right Expression) *BinaryExpression {
	return &BinaryExpression{
		Position: pos,
		Operator: op,
		Left:     left,
		Right:    right,
	}
}

// NewConditionalExpression creates a new ConditionalExpression node.
func NewConditionalExpression(pos Position, test Expression, consequent Expression, alternate Expression) *ConditionalExpression {
	return &ConditionalExpression{
		Position:   pos,
		Test:       test,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

// NewUpdateExpression creates a new UpdateExpression node.
func NewUpdateExpression(pos Position, op string, arg Expression, prefix bool) *UpdateExpression {
	return &UpdateExpression{
		Position: pos,
		Operator: op,
		Argument: arg,
		Prefix:   prefix,
	}
}

// NewMemberExpression creates a new MemberExpression node.
func NewMemberExpression(pos Position, obj Expression, prop Expression, computed bool) *MemberExpression {
	return &MemberExpression{
		Position: pos,
		Object:   obj,
		Property: prop,
		Computed: computed,
	}
}

// NewCallExpression creates a new CallExpression node.
func NewCallExpression(pos Position, callee Expression, args []Expression) *CallExpression {
	return &CallExpression{
		Position:  pos,
		Callee:    callee,
		Arguments: args,
	}
}

// NewNewExpression creates a new NewExpression node.
func NewNewExpression(pos Position, callee Expression, args []Expression) *NewExpression {
	return &NewExpression{
		Position:  pos,
		Callee:    callee,
		Arguments: args,
	}
}

// NewThisExpression creates a new ThisExpression node.
func NewThisExpression(pos Position) *ThisExpression {
	return &ThisExpression{Position: pos}
}

// NewSuperExpression creates a new SuperExpression node.
func NewSuperExpression(pos Position) *SuperExpression {
	return &SuperExpression{Position: pos}
}

// NewTemplateLiteral creates a new TemplateLiteral node.
func NewTemplateLiteral(pos Position, quasis []*TemplateElement, expressions []Expression) *TemplateLiteral {
	return &TemplateLiteral{
		Position:    pos,
		Quasis:      quasis,
		Expressions: expressions,
	}
}

// NewComponentElement creates a new ComponentElement node.
func NewComponentElement(pos Position, tag *Identifier, body []ComponentBodyElement) *ComponentElement {
	return &ComponentElement{
		Position: pos,
		Tag:      tag,
		Body:     body,
	}
}

// NewComponentProperty creates a new ComponentProperty node.
func NewComponentProperty(pos Position, key *Identifier, value Expression) *ComponentProperty {
	return &ComponentProperty{
		Position: pos,
		Key:      key,
		Value:    value,
	}
}

// NewForLoop creates a new ForLoop node.
func NewForLoop(pos Position, variable *Identifier, source Expression, body []ComponentBodyElement) *ForLoop {
	return &ForLoop{
		Position: pos,
		Variable: variable,
		Source:   source,
		Body:     body,
	}
}

// NewIfCondition creates a new IfCondition node.
func NewIfCondition(pos Position, test Expression, consequent []ComponentBodyElement, alternate []ComponentBodyElement) *IfCondition {
	return &IfCondition{
		Position:   pos,
		Test:       test,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

// NewTypeAnnotation creates a new TypeAnnotation node.
func NewTypeAnnotation(pos Position, typeNode Node) *TypeAnnotation {
	return &TypeAnnotation{
		Position: pos,
		Type:     typeNode,
	}
}

// NewTypeReference creates a new TypeReference node.
func NewTypeReference(pos Position, name *Identifier, typeParams []*TypeAnnotation) *TypeReference {
	return &TypeReference{
		Position:   pos,
		Name:       name,
		TypeParams: typeParams,
	}
}

// NewObjectType creates a new ObjectType node.
func NewObjectType(pos Position, members []Node) *ObjectType {
	return &ObjectType{
		Position: pos,
		Members:  members,
	}
}
