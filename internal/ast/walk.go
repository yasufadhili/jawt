package ast

// Walk traverses an AST node and its children, calling the appropriate
// visitor methods for each node.
func Walk(v Visitor, node Node) {
	if node == nil {
		return
	}

	node.Accept(v)

	switch n := node.(type) {
	// Top-Level
	case *Program:
		for _, doc := range n.Documents {
			Walk(v, doc)
		}
	case *Document:
		Walk(v, n.Name)
		for _, stmt := range n.Body {
			Walk(v, stmt)
		}

	// Declarations
	case *ImportDeclaration:
		for _, spec := range n.Specifiers {
			Walk(v, spec)
		}
		Walk(v, n.Source)
	case *ImportSpecifier:
		Walk(v, n.Local)
		if n.Imported != nil {
			Walk(v, n.Imported)
		}
	case *ImportDefaultSpecifier:
		Walk(v, n.Local)
	case *ImportNamespaceSpecifier:
		Walk(v, n.Local)
	case *ExportDeclaration:
		Walk(v, n.Declaration)
	case *VariableDeclaration:
		for _, decl := range n.Declarations {
			Walk(v, decl.ID)
			if decl.Init != nil {
				Walk(v, decl.Init)
			}
		}
	case *FunctionDeclaration:
		Walk(v, n.ID)
		for _, param := range n.Params {
			Walk(v, param)
		}
		if n.ReturnType != nil {
			Walk(v, n.ReturnType)
		}
		Walk(v, n.Body)
	case *ClassDeclaration:
		Walk(v, n.ID)
		if n.SuperClass != nil {
			Walk(v, n.SuperClass)
		}
		// Walk class body elements
	case *InterfaceDeclaration:
		Walk(v, n.ID)
		for _, ext := range n.Extends {
			Walk(v, ext)
		}
		Walk(v, n.Body)
	case *TypeAliasDeclaration:
		Walk(v, n.ID)
		Walk(v, n.Type)
	case *EnumDeclaration:
		Walk(v, n.ID)
		for _, member := range n.Members {
			Walk(v, member.ID)
			if member.Init != nil {
				Walk(v, member.Init)
			}
		}
	case *PropertyDeclaration:
		Walk(v, n.ID)
		if n.Type != nil {
			Walk(v, n.Type)
		}
		if n.DefaultValue != nil {
			Walk(v, n.DefaultValue)
		}
	case *StateDeclaration:
		Walk(v, n.ID)
		if n.Type != nil {
			Walk(v, n.Type)
		}
		if n.DefaultValue != nil {
			Walk(v, n.DefaultValue)
		}

	// Statements
	case *BlockStatement:
		for _, stmt := range n.List {
			Walk(v, stmt)
		}
	case *ExpressionStatement:
		Walk(v, n.Expression)
	case *IfStatement:
		Walk(v, n.Test)
		Walk(v, n.Consequent)
		if n.Alternate != nil {
			Walk(v, n.Alternate)
		}
	case *ForStatement:
		if n.Init != nil {
			Walk(v, n.Init)
		}
		if n.Test != nil {
			Walk(v, n.Test)
		}
		if n.Update != nil {
			Walk(v, n.Update)
		}
		Walk(v, n.Body)
	case *ForInStatement:
		Walk(v, n.Left)
		Walk(v, n.Right)
		Walk(v, n.Body)
	case *WhileStatement:
		Walk(v, n.Test)
		Walk(v, n.Body)
	case *ReturnStatement:
		if n.Argument != nil {
			Walk(v, n.Argument)
		}
	//case *BreakStatement:
	//	if n.Label != nil {
	//		Walk(v, n.Label)
	//	}
	//case *ContinueStatement:
	//	if n.Label != nil {
	//		Walk(v, n.Label)
	//	}
	case *ThrowStatement:
		Walk(v, n.Argument)
	case *TryStatement:
		Walk(v, n.Block)
		if n.Handler != nil {
			Walk(v, n.Handler.Param)
			Walk(v, n.Handler.Body)
		}
		if n.Finalizer != nil {
			Walk(v, n.Finalizer)
		}

	// Expressions
	case *ArrayLiteral:
		for _, elem := range n.Elements {
			Walk(v, elem)
		}
	case *ObjectLiteral:
		for _, prop := range n.Properties {
			Walk(v, prop.Key)
			Walk(v, prop.Value)
		}
	case *FunctionExpression:
		if n.ID != nil {
			Walk(v, n.ID)
		}
		for _, param := range n.Params {
			Walk(v, param)
		}
		Walk(v, n.Body)
	case *ArrowFunctionExpression:
		for _, param := range n.Params {
			Walk(v, param)
		}
		Walk(v, n.Body)
	case *UnaryExpression:
		Walk(v, n.Argument)
	case *UpdateExpression:
		Walk(v, n.Argument)
	case *BinaryExpression:
		Walk(v, n.Left)
		Walk(v, n.Right)
	case *ConditionalExpression:
		Walk(v, n.Test)
		Walk(v, n.Consequent)
		Walk(v, n.Alternate)
	case *MemberExpression:
		Walk(v, n.Object)
		Walk(v, n.Property)
	case *CallExpression:
		Walk(v, n.Callee)
		for _, arg := range n.Arguments {
			Walk(v, arg)
		}
	case *NewExpression:
		Walk(v, n.Callee)
		for _, arg := range n.Arguments {
			Walk(v, arg)
		}
	case *TemplateLiteral:
		for _, expr := range n.Expressions {
			Walk(v, expr)
		}

	// JML Specific
	case *ComponentElement:
		Walk(v, n.Tag)
		for _, child := range n.Body {
			Walk(v, child)
		}
	case *ComponentProperty:
		Walk(v, n.Key)
		Walk(v, n.Value)
	case *ForLoop:
		Walk(v, n.Variable)
		Walk(v, n.Source)
		for _, child := range n.Body {
			Walk(v, child)
		}
	case *IfCondition:
		Walk(v, n.Test)
		for _, child := range n.Consequent {
			Walk(v, child)
		}
		if n.Alternate != nil {
			for _, child := range n.Alternate {
				Walk(v, child)
			}
		}

	// Types
	case *TypeAnnotation:
		Walk(v, n.Type)
	case *TypeReference:
		Walk(v, n.Name)
		for _, param := range n.TypeParams {
			Walk(v, param)
		}
	case *ObjectType:
		for _, member := range n.Members {
			Walk(v, member)
		}

	// Leaf nodes (no children to walk)
	case *Identifier, *Literal, *ThisExpression, *SuperExpression, *BreakStatement, *ContinueStatement:
		// No-op
	}
}
