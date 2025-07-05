package ast

import (
	"fmt"
	"testing"
)

// TestUserCardComponent creates a UserCard component AST by hand
func TestUserCardComponent(t *testing.T) {
	pos := Position{Line: 1, Column: 1, File: "user-card.jml"}

	// Create the component document
	doc := NewDocument(
		pos,
		DocTypeComponent,
		NewIdentifier(pos, "UserCard"),
	)

	// Add property declaration: property userName: string
	doc.Properties = append(doc.Properties, NewPropertyDeclaration(
		pos,
		NewIdentifier(pos, "userName"),
		"string",
		nil,
		false,
	))

	// Add state declaration: state selected: boolean = false
	doc.States = append(doc.States, NewStateDeclaration(
		pos,
		NewIdentifier(pos, "selected"),
		"boolean",
		NewLiteral(pos, false, "boolean"),
	))

	// Create Container element
	containerAttrs := make(map[string]Expression)
	containerAttrs["style"] = NewLiteral(pos, "bg-white shadow-md rounded-lg p-6", "string")

	container := NewElementNode(
		pos,
		NewIdentifier(pos, "Container"),
		containerAttrs,
		nil,
	)

	// Create Text element
	textAttrs := make(map[string]Expression)
	textAttrs["content"] = NewBinding(
		pos,
		NewIdentifier(pos, "props"),
		NewIdentifier(pos, "userName"),
	)
	textAttrs["style"] = NewLiteral(pos, "text-xl font-semibold text-gray-800", "string")

	textElement := NewElementNode(
		pos,
		NewIdentifier(pos, "Text"),
		textAttrs,
		nil,
	)

	// Create Button element
	buttonAttrs := make(map[string]Expression)
	buttonAttrs["text"] = NewLiteral(pos, "Select", "string")
	buttonAttrs["onClick"] = NewLambdaExpression(
		pos,
		nil,
		"selected = !selected",
	)
	buttonAttrs["style"] = NewLiteral(pos, "mt-4 bg-blue-500 text-white px-4 py-2 rounded", "string")

	buttonElement := NewElementNode(
		pos,
		NewIdentifier(pos, "Button"),
		buttonAttrs,
		nil,
	)

	// Add children to container
	container.Children = append(container.Children, textElement, buttonElement)

	// Set root element
	doc.RootElement = container

	// Create a program with this document
	docs := []*Document{doc}
	program := NewProgram(pos, docs)

	// Test the structure
	if program == nil {
		t.Error("Program should not be nil")
	}

	if len(program.Documents) == 0 {
		t.Error("Program should have at least one document")
	}

	if len(program.Documents) != 1 {
		t.Errorf("Expected 1 document, got %d", len(program.Documents))
	}

	if program.Documents[0].Doctype != DocTypeComponent {
		t.Error("Document should be a component")
	}

	if program.Documents[0].Identifier.Name != "UserCard" {
		t.Error("Component name should be UserCard")
	}

	if len(program.Documents[0].Properties) != 1 {
		t.Error("Component should have exactly one property")
	}

	if len(program.Documents[0].States) != 1 {
		t.Error("Component should have exactly one state")
	}

	if program.Documents[0].RootElement.Tag.Name != "Container" {
		t.Error("Root element should be Container")
	}

	if len(program.Documents[0].RootElement.Children) != 2 {
		t.Error("Container should have exactly two children")
	}

	fmt.Println("✅ UserCard component AST created successfully")

	// Print the AST structure
	fmt.Println("\n--- AST Structure ---")
	visitor := NewPrintVisitor()
	Walk(program, visitor)
}

// TestPageWithImports creates a page with imports
func TestPageWithImports(t *testing.T) {
	pos := Position{Line: 1, Column: 1, File: "home.jml"}

	// Create the page document
	doc := NewDocument(
		pos,
		DocTypePage,
		NewIdentifier(pos, "home"),
	)

	// Add import statement: import component Layout from "components/layout"
	doc.Imports = append(doc.Imports, NewImportStatement(
		pos,
		ImportComponent,
		NewIdentifier(pos, "Layout"),
		NewLiteral(pos, "components/layout", "string"),
	))

	// Add import statement: import script analytics from "scripts/analytics"
	doc.Imports = append(doc.Imports, NewImportStatement(
		pos,
		ImportScript,
		NewIdentifier(pos, "analytics"),
		NewLiteral(pos, "scripts/analytics", "string"),
	))

	// Create a Page element
	pageAttrs := make(map[string]Expression)
	pageAttrs["title"] = NewLiteral(pos, "Welcome to My App", "string")
	pageAttrs["description"] = NewLiteral(pos, "A modern web application", "string")

	pageElement := NewElementNode(
		pos,
		NewIdentifier(pos, "Page"),
		pageAttrs,
		nil,
	)

	// Create Layout element
	layoutAttrs := make(map[string]Expression)
	layoutAttrs["showWelcome"] = NewLiteral(pos, true, "boolean")
	layoutAttrs["onClick"] = NewFunctionCall(
		pos,
		NewIdentifier(pos, "trackPageView"),
		[]Expression{NewLiteral(pos, "home", "string")},
	)

	layoutElement := NewElementNode(
		pos,
		NewIdentifier(pos, "Layout"),
		layoutAttrs,
		nil,
	)

	// Add Layout to Page
	pageElement.Children = append(pageElement.Children, layoutElement)

	// Set root element
	doc.RootElement = pageElement

	// Create a program with this document
	program := NewProgram(pos, []*Document{doc})

	// Test the structure
	if program == nil {
		t.Error("Program should not be nil")
	}

	if len(program.Documents) != 1 {
		t.Errorf("Expected 1 document, got %d", len(program.Documents))
	}

	if program.Documents[0].Doctype != DocTypePage {
		t.Error("Document should be a page")
	}

	if program.Documents[0].Identifier.Name != "home" {
		t.Error("Page name should be home")
	}

	if len(program.Documents[0].Imports) != 2 {
		t.Error("Page should have exactly two imports")
	}

	if program.Documents[0].Imports[0].Kind != ImportComponent {
		t.Error("First import should be a component")
	}

	if program.Documents[0].Imports[1].Kind != ImportScript {
		t.Error("Second import should be a script")
	}

	if program.Documents[0].RootElement.Tag.Name != "Page" {
		t.Error("Root element should be Page")
	}

	fmt.Println("✅ Home page AST created successfully")

	// Print the AST structure
	fmt.Println("\n--- AST Structure ---")
	visitor := NewPrintVisitor()
	Walk(program, visitor)
}

// TestExpressionTypes tests different expression types
func TestExpressionTypes(t *testing.T) {
	pos := Position{Line: 1, Column: 1, File: "test.jml"}

	// Test literal expressions
	stringLiteral := NewLiteral(pos, "Hello", "string")
	numberLiteral := NewLiteral(pos, 42, "number")
	booleanLiteral := NewLiteral(pos, true, "boolean")

	if stringLiteral.ExprType() != ExprLiteral {
		t.Error("String literal should have ExprLiteral type")
	}

	if numberLiteral.ExprType() != ExprLiteral {
		t.Error("Number literal should have ExprLiteral type")
	}

	if booleanLiteral.ExprType() != ExprLiteral {
		t.Error("Boolean literal should have ExprLiteral type")
	}

	// Test identifier
	identifier := NewIdentifier(pos, "myVar")
	if identifier.ExprType() != ExprIdentifier {
		t.Error("Identifier should have ExprIdentifier type")
	}

	// Test binding
	binding := NewBinding(pos, NewIdentifier(pos, "props"), NewIdentifier(pos, "name"))
	if binding.ExprType() != ExprBinding {
		t.Error("Binding should have ExprBinding type")
	}

	// Test function call
	functionCall := NewFunctionCall(pos, NewIdentifier(pos, "myFunction"), []Expression{
		NewLiteral(pos, "arg1", "string"),
		NewLiteral(pos, 123, "number"),
	})
	if functionCall.ExprType() != ExprFunctionCall {
		t.Error("Function call should have ExprFunctionCall type")
	}

	// Test lambda expression
	lambda := NewLambdaExpression(pos, []*Identifier{
		NewIdentifier(pos, "x"),
		NewIdentifier(pos, "y"),
	}, "x + y")
	if lambda.ExprType() != ExprLambda {
		t.Error("Lambda should have ExprLambda type")
	}

	fmt.Println("✅ All expression types tested successfully")
}

// TestVisitorPattern tests the visitor pattern implementation
func TestVisitorPattern(t *testing.T) {
	pos := Position{Line: 1, Column: 1, File: "test.jml"}

	// Create a simple AST
	doc := NewDocument(pos, DocTypeComponent, NewIdentifier(pos, "TestComponent"))

	// Add a property
	doc.Properties = append(doc.Properties, NewPropertyDeclaration(
		pos,
		NewIdentifier(pos, "value"),
		"string",
		NewLiteral(pos, "default", "string"),
		false,
	))

	// Add root element
	attrs := make(map[string]Expression)
	attrs["content"] = NewBinding(pos, NewIdentifier(pos, "props"), NewIdentifier(pos, "value"))

	doc.RootElement = NewElementNode(pos, NewIdentifier(pos, "Text"), attrs, nil)

	program := NewProgram(pos, []*Document{doc})

	// Test visitor pattern
	visitor := NewPrintVisitor()
	Walk(program, visitor)

	fmt.Println("✅ Visitor pattern tested successfully")
}

// TestFactoryMethods tests all factory methods
func TestFactoryMethods(t *testing.T) {
	pos := Position{Line: 1, Column: 1, File: "test.jml"}

	// Test all factory methods
	program := NewProgram(pos, nil)
	if program == nil {
		t.Error("NewProgram should not return nil")
	}

	doc := NewDocument(pos, DocTypeComponent, NewIdentifier(pos, "Test"))
	if doc == nil {
		t.Error("NewDocument should not return nil")
	}

	import_ := NewImportStatement(pos, ImportScript, NewIdentifier(pos, "script"), NewLiteral(pos, "path", "string"))
	if import_ == nil {
		t.Error("NewImportStatement should not return nil")
	}

	prop := NewPropertyDeclaration(pos, NewIdentifier(pos, "prop"), "string", nil, false)
	if prop == nil {
		t.Error("NewPropertyDeclaration should not return nil")
	}

	state := NewStateDeclaration(pos, NewIdentifier(pos, "state"), "boolean", NewLiteral(pos, false, "boolean"))
	if state == nil {
		t.Error("NewStateDeclaration should not return nil")
	}

	element := NewElementNode(pos, NewIdentifier(pos, "Element"), nil, nil)
	if element == nil {
		t.Error("NewElementNode should not return nil")
	}

	literal := NewLiteral(pos, "test", "string")
	if literal == nil {
		t.Error("NewLiteral should not return nil")
	}

	identifier := NewIdentifier(pos, "id")
	if identifier == nil {
		t.Error("NewIdentifier should not return nil")
	}

	binding := NewBinding(pos, NewIdentifier(pos, "obj"), NewIdentifier(pos, "prop"))
	if binding == nil {
		t.Error("NewBinding should not return nil")
	}

	functionCall := NewFunctionCall(pos, NewIdentifier(pos, "func"), nil)
	if functionCall == nil {
		t.Error("NewFunctionCall should not return nil")
	}

	lambda := NewLambdaExpression(pos, nil, "body")
	if lambda == nil {
		t.Error("NewLambdaExpression should not return nil")
	}

	fmt.Println("✅ All factory methods tested successfully")
}
