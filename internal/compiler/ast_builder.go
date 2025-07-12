package compiler

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/yasufadhili/jawt/internal/ast"
	parser "github.com/yasufadhili/jawt/internal/compiler/parser/generated"
	"github.com/yasufadhili/jawt/internal/diagnostic"
)

// AstBuilder is a visitor that builds the AST from the ANTLR parse tree.
type AstBuilder struct {
	*parser.BaseJmlVisitor
	reporter *diagnostic.Reporter
	file     string
}

// NewAstBuilder creates a new AstBuilder.
func NewAstBuilder(reporter *diagnostic.Reporter, file string) *AstBuilder {
	return &AstBuilder{
		BaseJmlVisitor: &parser.BaseJmlVisitor{},
		reporter:       reporter,
		file:           file,
	}
}

// getPosition extracts the AST position from an ANTLR context.
func (b *AstBuilder) getPosition(ctx antlr.ParserRuleContext) ast.Position {
	return ast.Position{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
		File:   b.file,
	}
}

// VisitDocument is the entry point for building the AST for a JML document.
func (b *AstBuilder) VisitDocument(ctx *parser.DocumentContext) interface{} {
	pos := b.getPosition(ctx)

	// Doctype declaration is mandatory
	docTypeDecl := ctx.DoctypeDeclaration().Accept(b).(*struct {
		Kind ast.DocType
		Name *ast.Identifier
	})

	var body []ast.Statement
	if ctx.SourceElements() != nil {
		// SourceElements can be pageContent or componentContent
		if ctx.SourceElements().PageContent() != nil {
			// Page content has a single root component element
			compElement := ctx.SourceElements().PageContent().ComponentElement().Accept(b).(ast.Statement)
			body = append(body, compElement)
		} else if ctx.SourceElements().ComponentContent() != nil {
			// Component content can have multiple source elements
			for _, seCtx := range ctx.SourceElements().ComponentContent().AllSourceElement() {
				stmt := seCtx.Accept(b).(ast.Statement)
				body = append(body, stmt)
			}
		}
	}

	// Handle imports
	if ctx.Imports() != nil {
		for _, importDeclCtx := range ctx.Imports().AllImportDeclaration() {
			importDecl := importDeclCtx.Accept(b).(ast.Statement)
			body = append([]ast.Statement{importDecl}, body...) // Add imports at the beginning of the body
		}
	}

	return ast.NewDocument(pos, docTypeDecl.Kind, docTypeDecl.Name, body, b.file)
}

// VisitDoctypeDeclaration handles the doctype declaration.
func (b *AstBuilder) VisitDoctypeDeclaration(ctx *parser.DoctypeDeclarationContext) interface{} {
	pos := b.getPosition(ctx)
	kindStr := ctx.DoctypeKind().GetText()
	name := ast.NewIdentifier(pos, ctx.DoctypeName().GetText())

	var kind ast.DocType
	switch kindStr {
	case "page":
		kind = ast.DocTypePage
	case "component":
		kind = ast.DocTypeComponent
	default:
		b.reporter.Error(pos, fmt.Sprintf("Unknown doctype kind: %s", kindStr))
		kind = ast.DocType("") // Default to empty or handle error appropriately
	}

	return &struct {
		Kind ast.DocType
		Name *ast.Identifier
	}{Kind: kind, Name: name}
}

// VisitImportDeclaration handles import declarations.
func (b *AstBuilder) VisitImportDeclaration(ctx *parser.ImportDeclarationContext) interface{} {
	pos := b.getPosition(ctx)

	if ctx.BrowserImport() != nil {
		return ast.NewImportDeclaration(pos, nil, nil, true)
	} else if ctx.ComponentImport() != nil {
		compImportCtx := ctx.ComponentImport()
		id := ast.NewIdentifier(b.getPosition(compImportCtx.Identifier()), compImportCtx.Identifier().GetText())
		source := ast.NewLiteral(b.getPosition(compImportCtx.StringLiteral()), "string", b.stripQuotes(compImportCtx.StringLiteral().GetText()))
		return ast.NewImportDeclaration(pos, []ast.Node{ast.NewImportSpecifier(pos, id)}, source, false)
	} else if ctx.ScriptImport() != nil {
		scriptImportCtx := ctx.ScriptImport()
		id := ast.NewIdentifier(b.getPosition(scriptImportCtx.Identifier()), scriptImportCtx.Identifier().GetText())
		source := ast.NewLiteral(b.getPosition(scriptImportCtx.StringLiteral()), "string", b.stripQuotes(scriptImportCtx.StringLiteral().GetText()))
		return ast.NewImportDeclaration(pos, []ast.Node{ast.NewImportSpecifier(pos, id)}, source, false)
	} else if ctx.ModuleImport() != nil {
		moduleImportCtx := ctx.ModuleImport()
		var specifiers []ast.Node
		var source *ast.Literal

		if moduleImportCtx.StringLiteral() != nil {
			source = ast.NewLiteral(b.getPosition(moduleImportCtx.StringLiteral()), "string", b.stripQuotes(moduleImportCtx.StringLiteral().GetText()))
		}

		if moduleImportCtx.ImportClause() != nil {
			importClauseCtx := moduleImportCtx.ImportClause()
			if importClauseCtx.Identifier() != nil { // default import
				id := ast.NewIdentifier(b.getPosition(importClauseCtx.Identifier()), importClauseCtx.Identifier().GetText())
				specifiers = append(specifiers, ast.NewImportDefaultSpecifier(pos, id))
			}
			if importClauseCtx.NamedImports() != nil {
				if importClauseCtx.NamedImports().ImportsList() != nil {
					for _, specCtx := range importClauseCtx.NamedImports().ImportsList().AllImportSpecifier() {
						local := ast.NewIdentifier(b.getPosition(specCtx.AllIdentifier()[0]), specCtx.AllIdentifier()[0].GetText())
						var imported *ast.Identifier
						if len(specCtx.AllIdentifier()) > 1 {
							imported = ast.NewIdentifier(b.getPosition(specCtx.AllIdentifier()[1]), specCtx.AllIdentifier()[1].GetText())
						}
						specifiers = append(specifiers, &ast.ImportSpecifier{Position: b.getPosition(specCtx), Local: local, Imported: imported})
					}
				}
			}
		}
		return ast.NewImportDeclaration(pos, specifiers, source, false)
	}
	return nil
}

// VisitSourceElement handles a generic source element.
func (b *AstBuilder) VisitSourceElement(ctx *parser.SourceElementContext) interface{} {
	if ctx.Statement() != nil {
		return ctx.Statement().Accept(b).(ast.Statement)
	}
	if ctx.Declaration() != nil {
		return ctx.Declaration().Accept(b).(ast.Statement)
	}
	if ctx.ExportDeclaration() != nil {
		return ctx.ExportDeclaration().Accept(b).(ast.Statement)
	}
	if ctx.ComponentElement() != nil {
		return ctx.ComponentElement().Accept(b).(ast.Statement)
	}
	return nil
}

// VisitComponentElement handles JML component instantiation.
func (b *AstBuilder) VisitComponentElement(ctx *parser.ComponentElementContext) interface{} {
	pos := b.getPosition(ctx)
	tag := ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
	var body []ast.ComponentBodyElement

	if ctx.ComponentBody() != nil {
		for _, bodyElemCtx := range ctx.ComponentBody().AllComponentBodyElement() {
			elem := bodyElemCtx.Accept(b).(ast.ComponentBodyElement)
			body = append(body, elem)
		}
	}
	return ast.NewComponentElement(pos, tag, body)
}

// VisitComponentBodyElement handles elements within a component body.
func (b *AstBuilder) VisitComponentBodyElement(ctx *parser.ComponentBodyElementContext) interface{} {
	if ctx.ComponentProperty() != nil {
		return ctx.ComponentProperty().Accept(b).(ast.ComponentBodyElement)
	}
	if ctx.ComponentElement() != nil {
		return ctx.ComponentElement().Accept(b).(ast.ComponentBodyElement)
	}
	if ctx.Statement() != nil {
		return ctx.Statement().Accept(b).(ast.ComponentBodyElement)
	}
	if ctx.ForLoop() != nil {
		return ctx.ForLoop().Accept(b).(ast.ComponentBodyElement)
	}
	if ctx.IfCondition() != nil {
		return ctx.IfCondition().Accept(b).(ast.ComponentBodyElement)
	}
	return nil
}

// VisitComponentProperty handles component properties.
func (b *AstBuilder) VisitComponentProperty(ctx *parser.ComponentPropertyContext) interface{} {
	pos := b.getPosition(ctx)
	key := ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
	value := ctx.Expression().Accept(b).(ast.Expression)
	return ast.NewComponentProperty(pos, key, value)
}

// VisitForLoop handles for loops within component bodies.
func (b *AstBuilder) VisitForLoop(ctx *parser.ForLoopContext) interface{} {
	pos := b.getPosition(ctx)
	variable := ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
	source := ctx.Expression().Accept(b).(ast.Expression)
	var body []ast.ComponentBodyElement
	if ctx.ComponentBody() != nil {
		for _, bodyElemCtx := range ctx.ComponentBody().AllComponentBodyElement() {
			elem := bodyElemCtx.Accept(b).(ast.ComponentBodyElement)
			body = append(body, elem)
		}
	}
	return ast.NewForLoop(pos, variable, source, body)
}

// VisitIfCondition handles if conditions within component bodies.
func (b *AstBuilder) VisitIfCondition(ctx *parser.IfConditionContext) interface{} {
	pos := b.getPosition(ctx)
	test := ctx.Expression().Accept(b).(ast.Expression)
	var consequent []ast.ComponentBodyElement
	if ctx.ComponentBody(0) != nil {
		for _, bodyElemCtx := range ctx.ComponentBody(0).AllComponentBodyElement() {
			elem := bodyElemCtx.Accept(b).(ast.ComponentBodyElement)
			consequent = append(consequent, elem)
		}
	}
	var alternate []ast.ComponentBodyElement
	if ctx.ELSE() != nil && ctx.ComponentBody(1) != nil {
		for _, bodyElemCtx := range ctx.ComponentBody(1).AllComponentBodyElement() {
			elem := bodyElemCtx.Accept(b).(ast.ComponentBodyElement)
			alternate = append(alternate, elem)
		}
	}
	return ast.NewIfCondition(pos, test, consequent, alternate)
}

// VisitExportDeclaration handles export declarations.
func (b *AstBuilder) VisitExportDeclaration(ctx *parser.ExportDeclarationContext) interface{} {
	pos := b.getPosition(ctx)
	var decl ast.Declaration
	isDefault := false

	if ctx.DEFAULT() != nil {
		isDefault = true
		if ctx.Declaration() != nil {
			decl = ctx.Declaration().Accept(b).(ast.Declaration)
		} else if ctx.Expression() != nil {
			// Export default expression, need to wrap it in a dummy declaration or handle differently
			// For now, assume it's always a declaration for simplicity based on AST structure
			// This might need refinement based on how `export default expression` is used.
			b.reporter.Error(pos, "Export default expression not fully supported yet, expecting a declaration.")
			return nil
		}
	} else if ctx.Declaration() != nil {
		decl = ctx.Declaration().Accept(b).(ast.Declaration)
	} else if ctx.ExportsList() != nil {
		// Handle named exports: export { a, b as c }
		// This requires creating a VariableDeclaration or similar to hold the exported identifiers.
		// For simplicity, just create a dummy VariableDeclaration for now.
		b.reporter.Error(pos, "Named exports not fully supported yet.")
		return nil
	}
	return ast.NewExportDeclaration(pos, decl, isDefault)
}

// VisitDeclaration handles generic declarations.
func (b *AstBuilder) VisitDeclaration(ctx *parser.DeclarationContext) interface{} {
	if ctx.VariableDeclaration() != nil {
		return ctx.VariableDeclaration().Accept(b).(ast.Declaration)
	}
	if ctx.FunctionDeclaration() != nil {
		return ctx.FunctionDeclaration().Accept(b).(ast.Declaration)
	}
	if ctx.ClassDeclaration() != nil {
		return ctx.ClassDeclaration().Accept(b).(ast.Declaration)
	}
	if ctx.InterfaceDeclaration() != nil {
		return ctx.InterfaceDeclaration().Accept(b).(ast.Declaration)
	}
	if ctx.TypeAliasDeclaration() != nil {
		return ctx.TypeAliasDeclaration().Accept(b).(ast.Declaration)
	}
	if ctx.EnumDeclaration() != nil {
		return ctx.EnumDeclaration().Accept(b).(ast.Declaration)
	}
	if ctx.PropertyDeclaration() != nil {
		return ctx.PropertyDeclaration().Accept(b).(ast.Declaration)
	}
	if ctx.StateDeclaration() != nil {
		return ctx.StateDeclaration().Accept(b).(ast.Declaration)
	}
	return nil
}

// VisitPropertyDeclaration handles property declarations.
func (b *AstBuilder) VisitPropertyDeclaration(ctx *parser.PropertyDeclarationContext) interface{} {
	pos := b.getPosition(ctx)
	id := ast.NewIdentifier(b.getPosition(ctx.Identifier(0)), ctx.Identifier(0).GetText())
	var typeAnn *ast.TypeAnnotation
	if ctx.TypeAnnotation() != nil {
		typeAnn = ctx.TypeAnnotation().Accept(b).(*ast.TypeAnnotation)
	}
	var defaultValue ast.Expression
	if ctx.Expression() != nil {
		defaultValue = ctx.Expression().Accept(b).(ast.Expression)
	}
	return ast.NewPropertyDeclaration(pos, id, typeAnn, defaultValue, false) // Optional not yet supported in grammar
}

// VisitStateDeclaration handles state declarations.
func (b *AstBuilder) VisitStateDeclaration(ctx *parser.StateDeclarationContext) interface{} {
	pos := b.getPosition(ctx)
	id := ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
	// State declarations in JML grammar don't have direct type/default value in the declaration itself,
	// but rather in the stateBody. For now, create a basic StateDeclaration.
	// The stateBody elements will be handled separately when needed for a more detailed AST.
	// For simplicity, assume no type/default value at this top level for now.
	return ast.NewStateDeclaration(pos, id, nil, nil)
}

// VisitVariableDeclaration handles variable declarations.
func (b *AstBuilder) VisitVariableDeclaration(ctx *parser.VariableDeclarationContext) interface{} {
	pos := b.getPosition(ctx)
	kind := ctx.VariableDeclarationList().VariableDeclarationKind().GetText()
	var declarators []*ast.VariableDeclarator
	for _, declCtx := range ctx.VariableDeclarationList().AllVariableDeclarator() {
		id := ast.NewIdentifier(b.getPosition(declCtx.Identifier()), declCtx.Identifier().GetText())
		var init ast.Expression
		if declCtx.Expression() != nil {
			init = declCtx.Expression().Accept(b).(ast.Expression)
		}
		declarators = append(declarators, ast.NewVariableDeclarator(b.getPosition(declCtx), id, init))
	}
	return ast.NewVariableDeclaration(pos, kind, declarators)
}

// VisitFunctionDeclaration handles function declarations.
func (b *AstBuilder) VisitFunctionDeclaration(ctx *parser.FunctionDeclarationContext) interface{} {
	pos := b.getPosition(ctx)
	id := ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
	var params []*ast.Identifier
	if ctx.ParameterList() != nil {
		for _, paramCtx := range ctx.ParameterList().AllParameter() {
			// Assuming simple identifier parameters for now
			params = append(params, ast.NewIdentifier(b.getPosition(paramCtx.Identifier()), paramCtx.Identifier().GetText()))
		}
	}
	body := ctx.FunctionBody().Accept(b).(*ast.BlockStatement)
	var returnType *ast.TypeAnnotation
	if ctx.TypeAnnotation() != nil {
		returnType = ctx.TypeAnnotation().Accept(b).(*ast.TypeAnnotation)
	}
	return ast.NewFunctionDeclaration(pos, id, params, body, returnType, false) // Async not in grammar yet
}

// VisitClassDeclaration handles class declarations.
func (b *AstBuilder) VisitClassDeclaration(ctx *parser.ClassDeclarationContext) interface{} {
	pos := b.getPosition(ctx)
	id := ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
	var superClass ast.Expression
	if ctx.ClassHeritage() != nil && ctx.ClassHeritage().EXTENDS() != nil {
		// Assuming typeReference is an Identifier for simplicity for now
		superClass = ast.NewIdentifier(b.getPosition(ctx.ClassHeritage().TypeReference().Identifier()), ctx.ClassHeritage().TypeReference().Identifier().GetText())
	}
	// ClassBody not yet implemented in AST, so passing nil for now.
	return ast.NewClassDeclaration(pos, id, superClass, nil)
}

// VisitInterfaceDeclaration handles interface declarations.
func (b *AstBuilder) VisitInterfaceDeclaration(ctx *parser.InterfaceDeclarationContext) interface{} {
	pos := b.getPosition(ctx)
	id := ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
	var extends []*ast.TypeReference
	if ctx.InterfaceHeritage() != nil {
		for _, trCtx := range ctx.InterfaceHeritage().AllTypeReference() {
			extends = append(extends, trCtx.Accept(b).(*ast.TypeReference))
		}
	}
	// InterfaceBody not yet implemented in AST, so passing nil for now.
	return ast.NewInterfaceDeclaration(pos, id, nil, extends)
}

// VisitTypeAliasDeclaration handles type alias declarations.
func (b *AstBuilder) VisitTypeAliasDeclaration(ctx *parser.TypeAliasDeclarationContext) interface{} {
	pos := b.getPosition(ctx)
	id := ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
	typeNode := ctx.Type().Accept(b).(ast.Node)
	return ast.NewTypeAliasDeclaration(pos, id, typeNode)
}

// VisitEnumDeclaration handles enum declarations.
func (b *AstBuilder) VisitEnumDeclaration(ctx *parser.EnumDeclarationContext) interface{} {
	pos := b.getPosition(ctx)
	id := ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
	var members []*ast.EnumMember
	if ctx.EnumBody() != nil {
		for _, memberCtx := range ctx.EnumBody().AllEnumMember() {
			memberID := ast.NewIdentifier(b.getPosition(memberCtx.Identifier()), memberCtx.Identifier().GetText())
			var init ast.Expression
			if memberCtx.Expression() != nil {
				init = memberCtx.Expression().Accept(b).(ast.Expression)
			}
			members = append(members, &ast.EnumMember{Position: b.getPosition(memberCtx), ID: memberID, Init: init})
		}
	}
	return ast.NewEnumDeclaration(pos, id, members)
}

// VisitBlockStatement handles block statements.
func (b *AstBuilder) VisitBlockStatement(ctx *parser.BlockStatementContext) interface{} {
	pos := b.getPosition(ctx)
	var statements []ast.Statement
	if ctx.StatementList() != nil {
		for _, stmtCtx := range ctx.StatementList().AllStatement() {
			statements = append(statements, stmtCtx.Accept(b).(ast.Statement))
		}
	}
	return ast.NewBlockStatement(pos, statements)
}

// VisitExpressionStatement handles expression statements.
func (b *AstBuilder) VisitExpressionStatement(ctx *parser.ExpressionStatementContext) interface{} {
	pos := b.getPosition(ctx)
	expr := ctx.Expression().Accept(b).(ast.Expression)
	return ast.NewExpressionStatement(pos, expr)
}

// VisitIfStatement handles if statements.
func (b *AstBuilder) VisitIfStatement(ctx *parser.IfStatementContext) interface{} {
	pos := b.getPosition(ctx)
	test := ctx.Expression().Accept(b).(ast.Expression)
	consequent := ctx.Statement(0).Accept(b).(ast.Statement)
	var alternate ast.Statement
	if ctx.ELSE() != nil && ctx.Statement(1) != nil {
		alternate = ctx.Statement(1).Accept(b).(ast.Statement)
	}
	return ast.NewIfStatement(pos, test, consequent, alternate)
}

// VisitIterationStatement handles iteration statements (for, while, do-while).
func (b *AstBuilder) VisitIterationStatement(ctx *parser.IterationStatementContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.FOR() != nil {
		if ctx.IN() != nil || ctx.OF() != nil { // for-in or for-of
			var leftNode ast.Node
			if ctx.VariableDeclarationList() != nil {
				leftNode = ctx.VariableDeclarationList().Accept(b).(ast.Node)
			} else if ctx.LeftHandSideExpression() != nil {
				leftNode = ctx.LeftHandSideExpression().Accept(b).(ast.Expression)
			}
			right := ctx.AllExpression()[len(ctx.AllExpression())-1].Accept(b).(ast.Expression)
			body := ctx.Statement().Accept(b).(ast.Statement)
			isOf := ctx.OF() != nil
			return ast.NewForInStatement(pos, leftNode, right, body, isOf)
		} else { // traditional for loop
			var init ast.Node
			if ctx.VariableDeclarationList() != nil {
				init = ctx.VariableDeclarationList().Accept(b).(ast.Node)
			} else if len(ctx.AllExpression()) > 0 && ctx.AllExpression()[0] != nil {
				init = ctx.AllExpression()[0].Accept(b).(ast.Expression)
			}
			var test ast.Expression
			if len(ctx.AllExpression()) > 1 && ctx.AllExpression()[1] != nil {
				test = ctx.AllExpression()[1].Accept(b).(ast.Expression)
			}
			var update ast.Expression
			if len(ctx.AllExpression()) > 2 && ctx.AllExpression()[2] != nil {
				update = ctx.AllExpression()[2].Accept(b).(ast.Expression)
			}
			body := ctx.Statement().Accept(b).(ast.Statement)
			return ast.NewForStatement(pos, init, test, update, body)
		}
	} else if ctx.WHILE() != nil {
		test := ctx.Expression(0).Accept(b).(ast.Expression)
		body := ctx.Statement().Accept(b).(ast.Statement)
		return ast.NewWhileStatement(pos, test, body)
	} else if ctx.DO() != nil {
		body := ctx.Statement().Accept(b).(ast.Statement)
		test := ctx.Expression(0).Accept(b).(ast.Expression)
		return ast.NewWhileStatement(pos, test, body) // Do-while is represented as a while loop with body executed once
	}
	return nil
}

// VisitReturnStatement handles return statements.
func (b *AstBuilder) VisitReturnStatement(ctx *parser.ReturnStatementContext) interface{} {
	pos := b.getPosition(ctx)
	var arg ast.Expression
	if ctx.Expression() != nil {
		arg = ctx.Expression().Accept(b).(ast.Expression)
	}
	return ast.NewReturnStatement(pos, arg)
}

// VisitBreakStatement handles break statements.
func (b *AstBuilder) VisitBreakStatement(ctx *parser.BreakStatementContext) interface{} {
	pos := b.getPosition(ctx)
	var label *ast.Identifier
	// Grammar doesn't yet support labels for break/continue directly, but AST has it.
	return ast.NewBreakStatement(pos, label)
}

// VisitContinueStatement handles continue statements.
func (b *AstBuilder) VisitContinueStatement(ctx *parser.ContinueStatementContext) interface{} {
	pos := b.getPosition(ctx)
	var label *ast.Identifier
	return ast.NewContinueStatement(pos, label)
}

// VisitThrowStatement handles throw statements.
func (b *AstBuilder) VisitThrowStatement(ctx *parser.ThrowStatementContext) interface{} {
	pos := b.getPosition(ctx)
	arg := ctx.Expression().Accept(b).(ast.Expression)
	return ast.NewThrowStatement(pos, arg)
}

// VisitTryStatement handles try-catch-finally statements.
func (b *AstBuilder) VisitTryStatement(ctx *parser.TryStatementContext) interface{} {
	pos := b.getPosition(ctx)
	block := ctx.Block().Accept(b).(*ast.BlockStatement)
	var handler *ast.CatchClause
	if ctx.CatchClause() != nil {
		catchCtx := ctx.CatchClause()
		param := ast.NewIdentifier(b.getPosition(catchCtx.Identifier()), catchCtx.Identifier().GetText())
		body := catchCtx.Block().Accept(b).(*ast.BlockStatement)
		handler = &ast.CatchClause{Position: b.getPosition(catchCtx), Param: param, Body: body}
	}
	var finalizer *ast.BlockStatement
	if ctx.FinallyClause() != nil {
		finalizer = ctx.FinallyClause().Block().Accept(b).(*ast.BlockStatement)
	}
	return ast.NewTryStatement(pos, block, handler, finalizer)
}

// VisitIdentifier handles identifiers.
func (b *AstBuilder) VisitIdentifier(ctx *parser.IdentifierContext) interface{} {
	pos := b.getPosition(ctx)
	return ast.NewIdentifier(pos, ctx.GetText())
}

// VisitLiteral handles literal expressions.
func (b *AstBuilder) VisitLiteral(ctx *parser.LiteralContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.StringLiteral() != nil {
		return ast.NewLiteral(pos, "string", b.stripQuotes(ctx.StringLiteral().GetText()))
	}
	if ctx.NumericLiteral() != nil {
		return ast.NewLiteral(pos, "number", ctx.NumericLiteral().GetText())
	}
	if ctx.TRUE() != nil {
		return ast.NewLiteral(pos, "boolean", "true")
	}
	if ctx.FALSE() != nil {
		return ast.NewLiteral(pos, "boolean", "false")
	}
	if ctx.NULL() != nil {
		return ast.NewLiteral(pos, "null", "null")
	}
	if ctx.UNDEFINED() != nil {
		return ast.NewLiteral(pos, "undefined", "undefined")
	}
	if ctx.RegexLiteral() != nil {
		return ast.NewLiteral(pos, "regexp", ctx.RegexLiteral().GetText())
	}
	return nil
}

// VisitArrayLiteral handles array literals.
func (b *AstBuilder) VisitArrayLiteral(ctx *parser.ArrayLiteralContext) interface{} {
	pos := b.getPosition(ctx)
	var elements []ast.Expression
	if ctx.ElementList() != nil {
		for _, exprCtx := range ctx.ElementList().AllExpression() {
			elements = append(elements, exprCtx.Accept(b).(ast.Expression))
		}
	}
	return ast.NewArrayLiteral(pos, elements)
}

// VisitObjectLiteral handles object literals.
func (b *AstBuilder) VisitObjectLiteral(ctx *parser.ObjectLiteralContext) interface{} {
	pos := b.getPosition(ctx)
	var properties []*ast.Property
	if ctx.PropertyNameAndValueList() != nil {
		for _, propAssignCtx := range ctx.PropertyNameAndValueList().AllPropertyAssignment() {
			propPos := b.getPosition(propAssignCtx)
			if propAssignCtx.PropertyName() != nil {
				key := propAssignCtx.PropertyName().Accept(b).(ast.Expression)
				value := propAssignCtx.Expression().Accept(b).(ast.Expression)
				properties = append(properties, &ast.Property{Position: propPos, Key: key, Value: value, Kind: "init"})
			} else if propAssignCtx.Identifier() != nil { // Shorthand property
				id := ast.NewIdentifier(b.getPosition(propAssignCtx.Identifier()), propAssignCtx.Identifier().GetText())
				properties = append(properties, &ast.Property{Position: propPos, Key: id, Value: id, Kind: "init"})
			} else if propAssignCtx.ELLIPSIS() != nil { // Spread property
				value := propAssignCtx.Expression().Accept(b).(ast.Expression)
				properties = append(properties, &ast.Property{Position: propPos, Key: nil, Value: value, Kind: "spread"})
			}
		}
	}
	return ast.NewObjectLiteral(pos, properties)
}

// VisitPropertyName handles property names in object literals.
func (b *AstBuilder) VisitPropertyName(ctx *parser.PropertyNameContext) interface{} {
	if ctx.Identifier() != nil {
		return ctx.Identifier().Accept(b).(ast.Expression)
	}
	if ctx.StringLiteral() != nil {
		return ast.NewLiteral(b.getPosition(ctx.StringLiteral()), "string", b.stripQuotes(ctx.StringLiteral().GetText()))
	}
	if ctx.NumericLiteral() != nil {
		return ast.NewLiteral(b.getPosition(ctx.NumericLiteral()), "number", ctx.NumericLiteral().GetText())
	}
	if ctx.Expression() != nil { // Computed property name
		return ctx.Expression().Accept(b).(ast.Expression)
	}
	return nil
}

// VisitFunctionExpression handles function expressions.
func (b *AstBuilder) VisitFunctionExpression(ctx *parser.FunctionExpressionContext) interface{} {
	pos := b.getPosition(ctx)
	var id *ast.Identifier
	if ctx.Identifier() != nil {
		id = ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
	}
	var params []*ast.Identifier
	if ctx.ParameterList() != nil {
		for _, paramCtx := range ctx.ParameterList().AllParameter() {
			params = append(params, ast.NewIdentifier(b.getPosition(paramCtx.Identifier()), paramCtx.Identifier().GetText()))
		}
	}
	body := ctx.FunctionBody().Accept(b).(*ast.BlockStatement)
	return ast.NewFunctionExpression(pos, id, params, body, false) // Async not in grammar
}

// VisitArrowFunction handles arrow function expressions.
func (b *AstBuilder) VisitArrowFunction(ctx *parser.ArrowFunctionContext) interface{} {
	pos := b.getPosition(ctx)
	var params []*ast.Identifier
	if ctx.ArrowFunctionParameters().Identifier() != nil {
		params = append(params, ast.NewIdentifier(b.getPosition(ctx.ArrowFunctionParameters().Identifier()), ctx.ArrowFunctionParameters().Identifier().GetText()))
	} else if ctx.ArrowFunctionParameters().ParameterList() != nil {
		if ctx.ArrowFunctionParameters().ParameterList() != nil {
			for _, paramCtx := range ctx.ArrowFunctionParameters().ParameterList().AllParameter() {
				params = append(params, ast.NewIdentifier(b.getPosition(paramCtx.Identifier()), paramCtx.Identifier().GetText()))
			}
		}
	}
	var body ast.Node
	if ctx.FunctionBody() != nil {
		body = ctx.FunctionBody().Accept(b).(*ast.BlockStatement)
	} else if ctx.Expression() != nil {
		body = ctx.Expression().Accept(b).(ast.Expression)
	}
	return ast.NewArrowFunctionExpression(pos, params, body, false) // Async not in grammar
}

// VisitTemplateLiteral handles template literals.
func (b *AstBuilder) VisitTemplateLiteral(ctx *parser.TemplateLiteralContext) interface{} {
	pos := b.getPosition(ctx)
	// Simplified handling. The full implementation will parse the TemplateStringLiteral
	// to extract quasis and expressions. For now, take the full text.
	// The grammar has `TemplateStringCharacter*` which includes `${` ... `}`
	// This needs more sophisticated parsing to separate quasis and expressions.
	// For now,  create a single quasi and put all expressions in it.
	// This is a placeholder and needs proper implementation.
	b.reporter.Warn(pos, "Simplified handling of TemplateLiteral. Full parsing of quasis and expressions is not yet implemented.")
	return ast.NewTemplateLiteral(pos, nil, nil)
}

// VisitUnaryExpression handles unary expressions.
func (b *AstBuilder) VisitUnaryExpression(ctx *parser.UnaryExpressionContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.PostfixExpression() != nil {
		return ctx.PostfixExpression().Accept(b).(ast.Expression)
	}
	op := ctx.GetChild(0).(antlr.Token).GetText()
	arg := ctx.UnaryExpression().Accept(b).(ast.Expression)
	return ast.NewUnaryExpression(pos, op, arg, true) // Prefix is true for these operators
}

// VisitPostfixExpression handles postfix expressions.
func (b *AstBuilder) VisitPostfixExpression(ctx *parser.PostfixExpressionContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.LeftHandSideExpression() != nil {
		expr := ctx.LeftHandSideExpression().Accept(b).(ast.Expression)
		if ctx.GetChildCount() > 1 { // Has ++ or --
			op := ctx.GetChild(1).(antlr.Token).GetText()
			return ast.NewUpdateExpression(pos, op, expr, false) // Prefix is false for postfix
		}
		return expr
	}
	return nil
}

// VisitMultiplicativeExpression handles multiplicative expressions.
func (b *AstBuilder) VisitMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.GetChildCount() == 1 {
		return ctx.UnaryExpression(0).Accept(b).(ast.Expression)
	}
	left := ctx.MultiplicativeExpression().Accept(b).(ast.Expression)
	op := ctx.GetChild(1).(antlr.Token).GetText()
	right := ctx.UnaryExpression(0).Accept(b).(ast.Expression)
	return ast.NewBinaryExpression(pos, op, left, right)
}

// VisitAdditiveExpression handles additive expressions.
func (b *AstBuilder) VisitAdditiveExpression(ctx *parser.AdditiveExpressionContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.GetChildCount() == 1 {
		return ctx.MultiplicativeExpression(0).Accept(b).(ast.Expression)
	}
	left := ctx.AdditiveExpression().Accept(b).(ast.Expression)
	op := ctx.GetChild(1).(antlr.Token).GetText()
	right := ctx.MultiplicativeExpression(0).Accept(b).(ast.Expression)
	return ast.NewBinaryExpression(pos, op, left, right)
}

// VisitRelationalExpression handles relational expressions.
func (b *AstBuilder) VisitRelationalExpression(ctx *parser.RelationalExpressionContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.GetChildCount() == 1 {
		return ctx.AdditiveExpression(0).Accept(b).(ast.Expression)
	}
	left := ctx.RelationalExpression().Accept(b).(ast.Expression)
	op := ctx.GetChild(1).(antlr.Token).GetText()
	right := ctx.AdditiveExpression(0).Accept(b).(ast.Expression)
	return ast.NewBinaryExpression(pos, op, left, right)
}

// VisitEqualityExpression handles equality expressions.
func (b *AstBuilder) VisitEqualityExpression(ctx *parser.EqualityExpressionContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.GetChildCount() == 1 {
		return ctx.RelationalExpression(0).Accept(b).(ast.Expression)
	}
	left := ctx.EqualityExpression().Accept(b).(ast.Expression)
	op := ctx.GetChild(1).(antlr.Token).GetText()
	right := ctx.RelationalExpression(0).Accept(b).(ast.Expression)
	return ast.NewBinaryExpression(pos, op, left, right)
}

// VisitLogicalAndExpression handles logical AND expressions.
func (b *AstBuilder) VisitLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.GetChildCount() == 1 {
		return ctx.EqualityExpression(0).Accept(b).(ast.Expression)
	}
	left := ctx.LogicalAndExpression().Accept(b).(ast.Expression)
	op := ctx.GetChild(1).(antlr.Token).GetText()
	right := ctx.EqualityExpression(0).Accept(b).(ast.Expression)
	return ast.NewBinaryExpression(pos, op, left, right)
}

// VisitLogicalOrExpression handles logical OR expressions.
func (b *AstBuilder) VisitLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.GetChildCount() == 1 {
		return ctx.LogicalAndExpression(0).Accept(b).(ast.Expression)
	}
	left := ctx.LogicalOrExpression().Accept(b).(ast.Expression)
	op := ctx.GetChild(1).(antlr.Token).GetText()
	right := ctx.LogicalAndExpression(0).Accept(b).(ast.Expression)
	return ast.NewBinaryExpression(pos, op, left, right)
}

// VisitConditionalExpression handles conditional (ternary) expressions.
func (b *AstBuilder) VisitConditionalExpression(ctx *parser.ConditionalExpressionContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.GetChildCount() == 1 {
		return ctx.LogicalOrExpression().Accept(b).(ast.Expression)
	}
	test := ctx.LogicalOrExpression().Accept(b).(ast.Expression)
	consequent := ctx.Expression(0).Accept(b).(ast.Expression)
	alternate := ctx.Expression(1).Accept(b).(ast.Expression)
	return ast.NewConditionalExpression(pos, test, consequent, alternate)
}

// VisitLeftHandSideExpression handles left-hand side expressions (calls, member access, new).
func (b *AstBuilder) VisitLeftHandSideExpression(ctx *parser.LeftHandSideExpressionContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.MemberExpression() != nil {
		expr := ctx.MemberExpression().Accept(b).(ast.Expression)
		if ctx.Arguments() != nil { // Call expression
			args := b.extractArguments(ctx.Arguments())
			return ast.NewCallExpression(pos, expr, args)
		}
		return expr
	}
	if ctx.NEW() != nil {
		callee := ctx.LeftHandSideExpression(0).Accept(b).(ast.Expression)
		var args []ast.Expression
		if ctx.Arguments() != nil {
			args = b.extractArguments(ctx.Arguments())
		}
		return ast.NewNewExpression(pos, callee, args)
	}
	// Chained calls, member access, property access
	if ctx.LeftHandSideExpression() != nil {
		base := ctx.LeftHandSideExpression().Accept(b).(ast.Expression)
		if ctx.Arguments() != nil { // Chained call
			args := b.extractArguments(ctx.Arguments())
			return ast.NewCallExpression(pos, base, args)
		} else if ctx.Expression() != nil { // Member access: base[expression]
			prop := ctx.Expression().Accept(b).(ast.Expression)
			return ast.NewMemberExpression(pos, base, prop, true)
		} else if ctx.Identifier() != nil { // Property access: base.identifier
			prop := ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
			return ast.NewMemberExpression(pos, base, prop, false)
		}
	}
	return nil
}

// VisitMemberExpression handles member expressions (dot and bracket notation).
func (b *AstBuilder) VisitMemberExpression(ctx *parser.MemberExpressionContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.PrimaryExpression() != nil {
		return ctx.PrimaryExpression().Accept(b).(ast.Expression)
	}
	base := ctx.MemberExpression().Accept(b).(ast.Expression)
	if ctx.Expression() != nil { // base[expression]
		prop := ctx.Expression().Accept(b).(ast.Expression)
		return ast.NewMemberExpression(pos, base, prop, true)
	} else if ctx.Identifier() != nil { // base.identifier
		prop := ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
		return ast.NewMemberExpression(pos, base, prop, false)
	}
	return nil
}

// VisitPrimaryExpression handles primary expressions.
func (b *AstBuilder) VisitPrimaryExpression(ctx *parser.PrimaryExpressionContext) interface{} {
	if ctx.THIS() != nil {
		return ast.NewThisExpression(b.getPosition(ctx))
	}
	if ctx.SUPER() != nil {
		return ast.NewSuperExpression(b.getPosition(ctx))
	}
	if ctx.Identifier() != nil {
		return ctx.Identifier().Accept(b).(ast.Expression)
	}
	if ctx.Literal() != nil {
		return ctx.Literal().Accept(b).(ast.Expression)
	}
	if ctx.ArrayLiteral() != nil {
		return ctx.ArrayLiteral().Accept(b).(ast.Expression)
	}
	if ctx.ObjectLiteral() != nil {
		return ctx.ObjectLiteral().Accept(b).(ast.Expression)
	}
	if ctx.FunctionExpression() != nil {
		return ctx.FunctionExpression().Accept(b).(ast.Expression)
	}
	if ctx.ArrowFunction() != nil {
		return ctx.ArrowFunction().Accept(b).(ast.Expression)
	}
	if ctx.Expression() != nil { // Parenthesized expression
		return ctx.Expression().Accept(b).(ast.Expression)
	}
	if ctx.TemplateLiteral() != nil {
		return ctx.TemplateLiteral().Accept(b).(ast.Expression)
	}
	return nil
}

// extractArguments extracts arguments from an ArgumentsContext.
func (b *AstBuilder) extractArguments(ctx *parser.ArgumentsContext) []ast.Expression {
	var args []ast.Expression
	if ctx.ArgumentList() != nil {
		for _, exprCtx := range ctx.ArgumentList().AllExpression() {
			args = append(args, exprCtx.Accept(b).(ast.Expression))
		}
	}
	return args
}

// VisitTypeAnnotation handles type annotations.
func (b *AstBuilder) VisitTypeAnnotation(ctx *parser.TypeAnnotationContext) interface{} {
	pos := b.getPosition(ctx)
	typeNode := ctx.Type().Accept(b).(ast.Node)
	return ast.NewTypeAnnotation(pos, typeNode)
}

// VisitType handles generic types.
func (b *AstBuilder) VisitType(ctx *parser.TypeContext) interface{} {
	return ctx.UnionType().Accept(b).(ast.Node)
}

// VisitUnionType handles union types.
func (b *AstBuilder) VisitUnionType(ctx *parser.UnionTypeContext) interface{} {
	// For simplicity, if it's a single type, return that type.
	// If it's a union, might need a dedicated UnionType AST node,
	// but for now, just return the first intersection type.
	// This needs to be refined when union types are to be fully represented.
	if len(ctx.AllIntersectionType()) == 1 {
		return ctx.IntersectionType(0).Accept(b).(ast.Node)
	}
	b.reporter.Warn(b.getPosition(ctx), "Union types are simplified to the first type in AST.")
	return ctx.IntersectionType(0).Accept(b).(ast.Node)
}

// VisitIntersectionType handles intersection types.
func (b *AstBuilder) VisitIntersectionType(ctx *parser.IntersectionTypeContext) interface{} {
	// Similar to union types, simplifying for now.
	if len(ctx.AllPrimaryType()) == 1 {
		return ctx.PrimaryType(0).Accept(b).(ast.Node)
	}
	b.reporter.Warn(b.getPosition(ctx), "Intersection types are simplified to the first type in AST.")
	return ctx.PrimaryType(0).Accept(b).(ast.Node)
}

// VisitPrimaryType handles primary types.
func (b *AstBuilder) VisitPrimaryType(ctx *parser.PrimaryTypeContext) interface{} {
	// Array types (e.g., string[]) are handled here.
	// For now, return the base type. Array dimensions are not explicitly in AST yet.
	if ctx.BaseType() != nil {
		return ctx.BaseType().Accept(b).(ast.Node)
	}
	return nil
}

// VisitBaseType handles base types.
func (b *AstBuilder) VisitBaseType(ctx *parser.BaseTypeContext) interface{} {
	if ctx.TypeReference() != nil {
		return ctx.TypeReference().Accept(b).(ast.Node)
	}
	if ctx.ObjectType() != nil {
		return ctx.ObjectType().Accept(b).(ast.Node)
	}
	if ctx.TupleType() != nil {
		// TupleType not yet in AST
		b.reporter.Warn(b.getPosition(ctx), "Tuple types are not yet fully supported in AST.")
		return nil
	}
	if ctx.PrimitiveType() != nil {
		return ast.NewLiteral(b.getPosition(ctx.PrimitiveType()), "type", ctx.PrimitiveType().GetText())
	}
	if ctx.LiteralType() != nil {
		return ctx.LiteralType().Accept(b).(ast.Node)
	}
	if ctx.Type() != nil { // Parenthesized type
		return ctx.Type().Accept(b).(ast.Node)
	}
	return nil
}

// VisitPrimitiveType handles primitive types.
func (b *AstBuilder) VisitPrimitiveType(ctx *parser.PrimitiveTypeContext) interface{} {
	pos := b.getPosition(ctx)
	return ast.NewLiteral(pos, "type", ctx.GetText())
}

// VisitLiteralType handles literal types.
func (b *AstBuilder) VisitLiteralType(ctx *parser.LiteralTypeContext) interface{} {
	pos := b.getPosition(ctx)
	if ctx.StringLiteral() != nil {
		return ast.NewLiteral(pos, "string", b.stripQuotes(ctx.StringLiteral().GetText()))
	}
	if ctx.NumericLiteral() != nil {
		return ast.NewLiteral(pos, "number", ctx.NumericLiteral().GetText())
	}
	if ctx.TRUE() != nil {
		return ast.NewLiteral(pos, "boolean", "true")
	}
	if ctx.FALSE() != nil {
		return ast.NewLiteral(pos, "boolean", "false")
	}
	return nil
}

// VisitTypeReference handles type references.
func (b *AstBuilder) VisitTypeReference(ctx *parser.TypeReferenceContext) interface{} {
	pos := b.getPosition(ctx)
	name := ast.NewIdentifier(b.getPosition(ctx.Identifier()), ctx.Identifier().GetText())
	var typeParams []*ast.TypeAnnotation
	if ctx.TypeArguments() != nil {
		for _, typeCtx := range ctx.TypeArguments().AllType() {
			typeParams = append(typeParams, ast.NewTypeAnnotation(b.getPosition(typeCtx), typeCtx.Accept(b).(ast.Node)))
		}
	}
	return ast.NewTypeReference(pos, name, typeParams)
}

// VisitObjectType handles object types.
func (b *AstBuilder) VisitObjectType(ctx *parser.ObjectTypeContext) interface{} {
	pos := b.getPosition(ctx)
	var members []ast.Node
	if ctx.ObjectTypeBody() != nil {
		for _, memberCtx := range ctx.ObjectTypeBody().AllObjectTypeMember() {
			if memberCtx.PropertySignature() != nil {
				propSigCtx := memberCtx.PropertySignature()
				id := ast.NewIdentifier(b.getPosition(propSigCtx.Identifier()), propSigCtx.Identifier().GetText())
				var typeAnn *ast.TypeAnnotation
				if propSigCtx.TypeAnnotation() != nil {
					typeAnn = propSigCtx.TypeAnnotation().Accept(b).(*ast.TypeAnnotation)
				}
				// PropertySignature is not a direct AST node, but its components can form a PropertyDeclaration-like structure
				// For now, represent it as a PropertyDeclaration without a default value.
				members = append(members, ast.NewPropertyDeclaration(b.getPosition(propSigCtx), id, typeAnn, nil, propSigCtx.QUESTION() != nil))
			} else if memberCtx.MethodSignature() != nil {
				// MethodSignature not yet in AST
				b.reporter.Warn(b.getPosition(memberCtx), "Method signatures in object types are not yet fully supported in AST.")
			} else if memberCtx.IndexSignature() != nil {
				// IndexSignature not yet in AST
				b.reporter.Warn(b.getPosition(memberCtx), "Index signatures in object types are not yet fully supported in AST.")
			}
		}
	}
	return ast.NewObjectType(pos, members)
}

// stripQuotes removes quotes from string literals.
func (b *AstBuilder) stripQuotes(s string) string {
	if len(s) >= 2 && (s[0] == '"' && s[len(s)-1] == '"' || s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}

// VisitStatement handles generic statements.
func (b *AstBuilder) VisitStatement(ctx *parser.StatementContext) interface{} {
	if ctx.Block() != nil {
		return ctx.Block().Accept(b).(ast.Statement)
	}
	if ctx.ExpressionStatement() != nil {
		return ctx.ExpressionStatement().Accept(b).(ast.Statement)
	}
	if ctx.IfStatement() != nil {
		return ctx.IfStatement().Accept(b).(ast.Statement)
	}
	if ctx.IterationStatement() != nil {
		return ctx.IterationStatement().Accept(b).(ast.Statement)
	}
	if ctx.ReturnStatement() != nil {
		return ctx.ReturnStatement().Accept(b).(ast.Statement)
	}
	if ctx.BreakStatement() != nil {
		return ctx.BreakStatement().Accept(b).(ast.Statement)
	}
	if ctx.ContinueStatement() != nil {
		return ctx.ContinueStatement().Accept(b).(ast.Statement)
	}
	if ctx.ThrowStatement() != nil {
		return ctx.ThrowStatement().Accept(b).(ast.Statement)
	}
	if ctx.TryStatement() != nil {
		return ctx.TryStatement().Accept(b).(ast.Statement)
	}
	// Handle empty statement (semicolon)
	if ctx.SEMI() != nil {
		return ast.NewExpressionStatement(b.getPosition(ctx), nil) // Represent as an empty expression statement
	}
	return nil
}

// VisitExpression handles generic expressions.
func (b *AstBuilder) VisitExpression(ctx *parser.ExpressionContext) interface{} {
	return ctx.ConditionalExpression().Accept(b).(ast.Expression)
}
