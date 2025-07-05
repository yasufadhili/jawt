package checker

import (
	"github.com/yasufadhili/jawt/internal/ast"
	"strings"
	"testing"
)

// Test helper functions
func createTestPosition(line, col int) ast.Position {
	return ast.Position{Line: line, Column: col, File: "test.jml"}
}

func TestSymbolKind_String(t *testing.T) {
	tests := []struct {
		kind     SymbolKind
		expected string
	}{
		{SymbolProperty, "property"},
		{SymbolState, "state"},
		{SymbolComponent, "component"},
		{SymbolFunction, "function"},
		{SymbolVariable, "variable"},
		{SymbolParameter, "parameter"},
		{SymbolImport, "import"},
		{SymbolBuiltIn, "built-in"},
	}

	for _, test := range tests {
		if got := test.kind.String(); got != test.expected {
			t.Errorf("SymbolKind.String() = %q, want %q", got, test.expected)
		}
	}
}

func TestScopeKind_String(t *testing.T) {
	tests := []struct {
		kind     ScopeKind
		expected string
	}{
		{ScopeGlobal, "global"},
		{ScopeDocument, "document"},
		{ScopeFunction, "function"},
		{ScopeBlock, "block"},
		{ScopeElement, "element"},
	}

	for _, test := range tests {
		if got := test.kind.String(); got != test.expected {
			t.Errorf("ScopeKind.String() = %q, want %q", got, test.expected)
		}
	}
}

func TestPropertySymbol(t *testing.T) {
	pos := createTestPosition(1, 1)

	// Test required property
	prop := NewPropertySymbol("title", "string", pos, false, nil)

	if prop.Name() != "title" {
		t.Errorf("PropertySymbol.Name() = %q, want %q", prop.Name(), "title")
	}

	if prop.Kind() != SymbolProperty {
		t.Errorf("PropertySymbol.Kind() = %v, want %v", prop.Kind(), SymbolProperty)
	}

	if prop.Type() != "string" {
		t.Errorf("PropertySymbol.Type() = %q, want %q", prop.Type(), "string")
	}

	if prop.IsOptional {
		t.Error("PropertySymbol.IsOptional should be false")
	}

	expected := "property title: string"
	if got := prop.String(); got != expected {
		t.Errorf("PropertySymbol.String() = %q, want %q", got, expected)
	}

	// Test optional property
	optProp := NewPropertySymbol("subtitle", "string", pos, true, nil)
	expected = "property subtitle?: string"
	if got := optProp.String(); got != expected {
		t.Errorf("Optional PropertySymbol.String() = %q, want %q", got, expected)
	}
}

func TestStateSymbol(t *testing.T) {
	pos := createTestPosition(1, 1)
	state := NewStateSymbol("count", "number", pos, nil)

	if state.Name() != "count" {
		t.Errorf("StateSymbol.Name() = %q, want %q", state.Name(), "count")
	}

	if state.Kind() != SymbolState {
		t.Errorf("StateSymbol.Kind() = %v, want %v", state.Kind(), SymbolState)
	}

	if state.Type() != "number" {
		t.Errorf("StateSymbol.Type() = %q, want %q", state.Type(), "number")
	}

	expected := "state count: number"
	if got := state.String(); got != expected {
		t.Errorf("StateSymbol.String() = %q, want %q", got, expected)
	}
}

func TestComponentSymbol(t *testing.T) {
	pos := createTestPosition(1, 1)
	comp := NewComponentSymbol("Button", pos, ast.DocTypeComponent)

	if comp.Name() != "Button" {
		t.Errorf("ComponentSymbol.Name() = %q, want %q", comp.Name(), "Button")
	}

	if comp.Kind() != SymbolComponent {
		t.Errorf("ComponentSymbol.Kind() = %v, want %v", comp.Kind(), SymbolComponent)
	}

	if comp.DocType != ast.DocTypeComponent {
		t.Errorf("ComponentSymbol.DocType = %v, want %v", comp.DocType, ast.DocTypeComponent)
	}

	expected := "component Button (component)"
	if got := comp.String(); got != expected {
		t.Errorf("ComponentSymbol.String() = %q, want %q", got, expected)
	}
}

func TestFunctionSymbol(t *testing.T) {
	pos := createTestPosition(1, 1)
	params := []*ParameterSymbol{
		NewParameterSymbol("x", "number", pos, false),
		NewParameterSymbol("y", "number", pos, true),
	}

	fn := NewFunctionSymbol("add", "number", pos, params)

	if fn.Name() != "add" {
		t.Errorf("FunctionSymbol.Name() = %q, want %q", fn.Name(), "add")
	}

	if fn.Kind() != SymbolFunction {
		t.Errorf("FunctionSymbol.Kind() = %v, want %v", fn.Kind(), SymbolFunction)
	}

	if len(fn.Parameters) != 2 {
		t.Errorf("FunctionSymbol.Parameters length = %d, want 2", len(fn.Parameters))
	}

	expected := "function add(x: number, y?: number): number"
	if got := fn.String(); got != expected {
		t.Errorf("FunctionSymbol.String() = %q, want %q", got, expected)
	}
}

func TestParameterSymbol(t *testing.T) {
	pos := createTestPosition(1, 1)

	// Required parameter
	param := NewParameterSymbol("x", "number", pos, false)
	expected := "x: number"
	if got := param.String(); got != expected {
		t.Errorf("ParameterSymbol.String() = %q, want %q", got, expected)
	}

	// Optional parameter
	optParam := NewParameterSymbol("y", "string", pos, true)
	expected = "y?: string"
	if got := optParam.String(); got != expected {
		t.Errorf("Optional ParameterSymbol.String() = %q, want %q", got, expected)
	}
}

func TestVariableSymbol(t *testing.T) {
	pos := createTestPosition(1, 1)

	// Variable
	variable := NewVariableSymbol("data", "string", pos, false)
	expected := "let data: string"
	if got := variable.String(); got != expected {
		t.Errorf("VariableSymbol.String() = %q, want %q", got, expected)
	}

	// Constant
	constant := NewVariableSymbol("API_URL", "string", pos, true)
	expected = "const API_URL: string"
	if got := constant.String(); got != expected {
		t.Errorf("Constant VariableSymbol.String() = %q, want %q", got, expected)
	}
}

func TestImportSymbol(t *testing.T) {
	pos := createTestPosition(1, 1)
	imp := NewImportSymbol("Button", "./Button.jml", pos, ast.ImportComponent)

	if imp.Name() != "Button" {
		t.Errorf("ImportSymbol.Name() = %q, want %q", imp.Name(), "Button")
	}

	if imp.Path != "./Button.jml" {
		t.Errorf("ImportSymbol.Path = %q, want %q", imp.Path, "./Button.jml")
	}

	if imp.ImportKind != ast.ImportComponent {
		t.Errorf("ImportSymbol.ImportKind = %v, want %v", imp.ImportKind, ast.ImportComponent)
	}

	expected := "import component Button from ./Button.jml"
	if got := imp.String(); got != expected {
		t.Errorf("ImportSymbol.String() = %q, want %q", got, expected)
	}
}

func TestBuiltInSymbol(t *testing.T) {
	builtin := NewBuiltInSymbol("string", "type", "String type")

	if builtin.Name() != "string" {
		t.Errorf("BuiltInSymbol.Name() = %q, want %q", builtin.Name(), "string")
	}

	if builtin.Kind() != SymbolBuiltIn {
		t.Errorf("BuiltInSymbol.Kind() = %v, want %v", builtin.Kind(), SymbolBuiltIn)
	}

	expected := "built-in string: type"
	if got := builtin.String(); got != expected {
		t.Errorf("BuiltInSymbol.String() = %q, want %q", got, expected)
	}
}

func TestScope_DefineAndLookup(t *testing.T) {
	scope := NewScope(ScopeGlobal, nil, "test")
	pos := createTestPosition(1, 1)

	// Test define
	symbol := NewVariableSymbol("test", "string", pos, false)
	err := scope.Define(symbol)
	if err != nil {
		t.Errorf("Scope.Define() error = %v", err)
	}

	// Test lookup
	found, exists := scope.Lookup("test")
	if !exists {
		t.Error("Scope.Lookup() should find symbol")
	}

	if found.Name() != "test" {
		t.Errorf("Scope.Lookup() returned wrong symbol name: %q", found.Name())
	}

	// Test lookup non-existent
	_, exists = scope.Lookup("nonexistent")
	if exists {
		t.Error("Scope.Lookup() should not find non-existent symbol")
	}

	// Test duplicate define
	duplicate := NewVariableSymbol("test", "number", pos, false)
	err = scope.Define(duplicate)
	if err == nil {
		t.Error("Scope.Define() should error on duplicate symbol")
	}
}

func TestScope_LookupRecursive(t *testing.T) {
	parent := NewScope(ScopeGlobal, nil, "parent")
	child := NewScope(ScopeFunction, parent, "child")
	pos := createTestPosition(1, 1)

	// Define symbol in parent
	parentSymbol := NewVariableSymbol("parent_var", "string", pos, false)
	parent.Define(parentSymbol)

	// Define symbol in child
	childSymbol := NewVariableSymbol("child_var", "number", pos, false)
	child.Define(childSymbol)

	// Test child can find its own symbol
	found, exists := child.LookupRecursive("child_var")
	if !exists {
		t.Error("Child scope should find its own symbol")
	}
	if found.Name() != "child_var" {
		t.Errorf("Found wrong symbol: %q", found.Name())
	}

	// Test child can find parent symbol
	found, exists = child.LookupRecursive("parent_var")
	if !exists {
		t.Error("Child scope should find parent symbol")
	}
	if found.Name() != "parent_var" {
		t.Errorf("Found wrong symbol: %q", found.Name())
	}

	// Test parent cannot find child symbol
	_, exists = parent.LookupRecursive("child_var")
	if exists {
		t.Error("Parent scope should not find child symbol")
	}
}

func TestScope_GetSymbolsByKind(t *testing.T) {
	scope := NewScope(ScopeDocument, nil, "test")
	pos := createTestPosition(1, 1)

	// Add different types of symbols
	prop := NewPropertySymbol("title", "string", pos, false, nil)
	state := NewStateSymbol("count", "number", pos, nil)
	variable := NewVariableSymbol("data", "object", pos, false)

	scope.Define(prop)
	scope.Define(state)
	scope.Define(variable)

	// Test getting properties
	properties := scope.GetSymbolsByKind(SymbolProperty)
	if len(properties) != 1 {
		t.Errorf("GetSymbolsByKind(SymbolProperty) = %d symbols, want 1", len(properties))
	}

	// Test getting states
	states := scope.GetSymbolsByKind(SymbolState)
	if len(states) != 1 {
		t.Errorf("GetSymbolsByKind(SymbolState) = %d symbols, want 1", len(states))
	}

	// Test getting variables
	variables := scope.GetSymbolsByKind(SymbolVariable)
	if len(variables) != 1 {
		t.Errorf("GetSymbolsByKind(SymbolVariable) = %d symbols, want 1", len(variables))
	}

	// Test getting non-existent kind
	functions := scope.GetSymbolsByKind(SymbolFunction)
	if len(functions) != 0 {
		t.Errorf("GetSymbolsByKind(SymbolFunction) = %d symbols, want 0", len(functions))
	}
}

func TestSymbolTable_Creation(t *testing.T) {
	st := NewSymbolTable()

	if st.global == nil {
		t.Error("SymbolTable should have global scope")
	}

	if st.current != st.global {
		t.Error("SymbolTable should start with global as current scope")
	}

	// Test built-ins are loaded
	if _, exists := st.Lookup("string"); !exists {
		t.Error("SymbolTable should have built-in 'string' type")
	}

	if _, exists := st.Lookup("number"); !exists {
		t.Error("SymbolTable should have built-in 'number' type")
	}

	if _, exists := st.Lookup("console.log"); !exists {
		t.Error("SymbolTable should have built-in 'console.log' function")
	}
}

func TestSymbolTable_ScopeManagement(t *testing.T) {
	st := NewSymbolTable()

	// Test entering scope
	docScope := st.EnterScope(ScopeDocument, "TestDoc")
	if st.current != docScope {
		t.Error("EnterScope should set current scope")
	}

	if docScope.parent != st.global {
		t.Error("New scope should have global as parent")
	}

	// Test entering nested scope
	funcScope := st.EnterScope(ScopeFunction, "TestFunc")
	if st.current != funcScope {
		t.Error("EnterScope should set current scope")
	}

	if funcScope.parent != docScope {
		t.Error("New scope should have document scope as parent")
	}

	// Test exiting scope
	err := st.ExitScope()
	if err != nil {
		t.Errorf("ExitScope() error = %v", err)
	}

	if st.current != docScope {
		t.Error("ExitScope should return to parent scope")
	}

	// Test exiting to global
	err = st.ExitScope()
	if err != nil {
		t.Errorf("ExitScope() error = %v", err)
	}

	if st.current != st.global {
		t.Error("ExitScope should return to global scope")
	}

	// Test error when trying to exit global scope
	err = st.ExitScope()
	if err == nil {
		t.Error("ExitScope() should error when trying to exit global scope")
	}
}

func TestSymbolTable_DefineAndLookup(t *testing.T) {
	st := NewSymbolTable()
	pos := createTestPosition(1, 1)

	// Test define and lookup in global scope
	variable := NewVariableSymbol("globalVar", "string", pos, false)
	err := st.Define(variable)
	if err != nil {
		t.Errorf("Define() error = %v", err)
	}

	found, exists := st.Lookup("globalVar")
	if !exists {
		t.Error("Lookup() should find global variable")
	}

	if found.Name() != "globalVar" {
		t.Errorf("Lookup() returned wrong symbol: %q", found.Name())
	}

	// Test lookup of built-in
	builtin, exists := st.Lookup("string")
	if !exists {
		t.Error("Lookup() should find built-in type")
	}

	if builtin.Kind() != SymbolBuiltIn {
		t.Errorf("Lookup() returned wrong kind for built-in: %v", builtin.Kind())
	}
}

func TestSymbolTable_ComponentManagement(t *testing.T) {
	st := NewSymbolTable()
	pos := createTestPosition(1, 1)

	// Test define component
	component := NewComponentSymbol("Button", pos, ast.DocTypeComponent)
	err := st.DefineComponent(component)
	if err != nil {
		t.Errorf("DefineComponent() error = %v", err)
	}

	// Test get component
	found, exists := st.GetComponent("Button")
	if !exists {
		t.Error("GetComponent() should find component")
	}

	if found.Name() != "Button" {
		t.Errorf("GetComponent() returned wrong component: %q", found.Name())
	}

	// Test component is also in global scope
	globalFound, exists := st.Lookup("Button")
	if !exists {
		t.Error("Component should be findable in global scope")
	}

	if globalFound.Kind() != SymbolComponent {
		t.Errorf("Component in global scope has wrong kind: %v", globalFound.Kind())
	}

	// Test get all components
	allComponents := st.GetAllComponents()
	if len(allComponents) != 1 {
		t.Errorf("GetAllComponents() returned %d components, want 1", len(allComponents))
	}

	if _, exists := allComponents["Button"]; !exists {
		t.Error("GetAllComponents() should include Button")
	}
}

func TestSymbolTable_BuiltInManagement(t *testing.T) {
	st := NewSymbolTable()

	// Test IsBuiltIn
	if !st.IsBuiltIn("string") {
		t.Error("IsBuiltIn() should return true for 'string'")
	}

	if st.IsBuiltIn("CustomType") {
		t.Error("IsBuiltIn() should return false for custom type")
	}

	// Test GetBuiltIn
	builtin, exists := st.GetBuiltIn("string")
	if !exists {
		t.Error("GetBuiltIn() should find 'string' built-in")
	}

	if builtin.Name() != "string" {
		t.Errorf("GetBuiltIn() returned wrong built-in: %q", builtin.Name())
	}

	// Test GetBuiltIn for non-existent
	_, exists = st.GetBuiltIn("NonExistent")
	if exists {
		t.Error("GetBuiltIn() should not find non-existent built-in")
	}
}

func TestSymbolTable_Debug(t *testing.T) {
	st := NewSymbolTable()
	pos := createTestPosition(1, 1)

	// Add some symbols
	st.Define(NewVariableSymbol("test", "string", pos, false))
	st.EnterScope(ScopeDocument, "TestDoc")
	st.Define(NewPropertySymbol("title", "string", pos, false, nil))

	debug := st.Debug()

	// Check that debug output contains expected elements
	if !strings.Contains(debug, "Symbol Table Debug:") {
		t.Error("Debug output should contain header")
	}

	if !strings.Contains(debug, "global") {
		t.Error("Debug output should contain global scope")
	}

	if !strings.Contains(debug, "document") {
		t.Error("Debug output should contain document scope")
	}

	if !strings.Contains(debug, "test: let test: string") {
		t.Error("Debug output should contain variable symbol")
	}

	if !strings.Contains(debug, "title: property title: string") {
		t.Error("Debug output should contain property symbol")
	}
}

func TestSymbolTable_String(t *testing.T) {
	st := NewSymbolTable()
	pos := createTestPosition(1, 1)

	// Add a component
	component := NewComponentSymbol("Button", pos, ast.DocTypeComponent)
	st.DefineComponent(component)

	str := st.String()

	if !strings.Contains(str, "SymbolTable") {
		t.Error("String() should contain 'SymbolTable'")
	}

	if !strings.Contains(str, "1 components") {
		t.Error("String() should show component count")
	}

	if !strings.Contains(str, "global") {
		t.Error("String() should show current scope")
	}
}

// Integration test
func TestSymbolTable_Integration(t *testing.T) {
	st := NewSymbolTable()
	pos := createTestPosition(1, 1)

	// Simulate processing a component
	st.EnterScope(ScopeDocument, "MyButton")

	// Define properties
	titleProp := NewPropertySymbol("title", "string", pos, false, nil)
	disabledProp := NewPropertySymbol("disabled", "boolean", pos, true, nil)

	st.Define(titleProp)
	st.Define(disabledProp)

	// Define states
	hoverState := NewStateSymbol("isHovered", "boolean", pos, nil)
	st.Define(hoverState)

	// Enter function scope
	st.EnterScope(ScopeFunction, "handleClick")

	// Define parameters
	eventParam := NewParameterSymbol("event", "MouseEvent", pos, false)
	st.Define(eventParam)

	// Test lookups
	// Should find parameter in current scope
	if _, exists := st.Lookup("event"); !exists {
		t.Error("Should find parameter in function scope")
	}

	// Should find state in parent scope
	if _, exists := st.Lookup("isHovered"); !exists {
		t.Error("Should find state in parent document scope")
	}

	// Should find property in parent scope
	if _, exists := st.Lookup("title"); !exists {
		t.Error("Should find property in parent document scope")
	}

	// Should find built-in in global scope
	if _, exists := st.Lookup("string"); !exists {
		t.Error("Should find built-in type in global scope")
	}

	// Should not find non-existent symbol
	if _, exists := st.Lookup("nonExistent"); exists {
		t.Error("Should not find non-existent symbol")
	}

	// Exit function scope
	st.ExitScope()

	// Should no longer find parameter
	if _, exists := st.Lookup("event"); exists {
		t.Error("Should not find parameter after exiting function scope")
	}

	// Should still find state
	if _, exists := st.Lookup("isHovered"); !exists {
		t.Error("Should still find state in document scope")
	}

	// Exit document scope
	st.ExitScope()

	// Should no longer find state
	if _, exists := st.Lookup("isHovered"); exists {
		t.Error("Should not find state after exiting document scope")
	}

	// Should still find built-ins
	if _, exists := st.Lookup("string"); !exists {
		t.Error("Should still find built-in type in global scope")
	}
}
