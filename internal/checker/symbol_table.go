package checker

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/ast"
	"strings"
)

// SymbolKind represents the different types of symbols
type SymbolKind int

const (
	SymbolProperty SymbolKind = iota
	SymbolState
	SymbolComponent
	SymbolFunction
	SymbolVariable
	SymbolParameter
	SymbolImport
	SymbolBuiltIn
)

func (s SymbolKind) String() string {
	switch s {
	case SymbolProperty:
		return "property"
	case SymbolState:
		return "state"
	case SymbolComponent:
		return "component"
	case SymbolFunction:
		return "function"
	case SymbolVariable:
		return "variable"
	case SymbolParameter:
		return "parameter"
	case SymbolImport:
		return "import"
	case SymbolBuiltIn:
		return "built-in"
	default:
		return "unknown"
	}
}

// ScopeKind represents different types of scopes
type ScopeKind int

const (
	ScopeGlobal ScopeKind = iota
	ScopeDocument
	ScopeFunction
	ScopeBlock
	ScopeElement
)

func (s ScopeKind) String() string {
	switch s {
	case ScopeGlobal:
		return "global"
	case ScopeDocument:
		return "document"
	case ScopeFunction:
		return "function"
	case ScopeBlock:
		return "block"
	case ScopeElement:
		return "element"
	default:
		return "unknown"
	}
}

// Symbol represents a symbol in the symbol table
type Symbol interface {
	Name() string
	Kind() SymbolKind
	Type() string
	Position() ast.Position
	String() string
}

// BaseSymbol provides common functionality for all symbols
type BaseSymbol struct {
	name     string
	kind     SymbolKind
	typeStr  string
	position ast.Position
}

func (b *BaseSymbol) Name() string           { return b.name }
func (b *BaseSymbol) Kind() SymbolKind       { return b.kind }
func (b *BaseSymbol) Type() string           { return b.typeStr }
func (b *BaseSymbol) Position() ast.Position { return b.position }
func (b *BaseSymbol) String() string {
	return fmt.Sprintf("%s %s: %s", b.kind, b.name, b.typeStr)
}

// PropertySymbol represents a property declaration
type PropertySymbol struct {
	BaseSymbol
	IsOptional   bool
	DefaultValue ast.Expression
}

func NewPropertySymbol(name, typeStr string, pos ast.Position, optional bool, defaultValue ast.Expression) *PropertySymbol {
	return &PropertySymbol{
		BaseSymbol: BaseSymbol{
			name:     name,
			kind:     SymbolProperty,
			typeStr:  typeStr,
			position: pos,
		},
		IsOptional:   optional,
		DefaultValue: defaultValue,
	}
}

func (p *PropertySymbol) String() string {
	optional := ""
	if p.IsOptional {
		optional = "?"
	}
	return fmt.Sprintf("property %s%s: %s", p.name, optional, p.typeStr)
}

// StateSymbol represents a state declaration
type StateSymbol struct {
	BaseSymbol
	DefaultValue ast.Expression
}

func NewStateSymbol(name, typeStr string, pos ast.Position, defaultValue ast.Expression) *StateSymbol {
	return &StateSymbol{
		BaseSymbol: BaseSymbol{
			name:     name,
			kind:     SymbolState,
			typeStr:  typeStr,
			position: pos,
		},
		DefaultValue: defaultValue,
	}
}

func (s *StateSymbol) String() string {
	return fmt.Sprintf("state %s: %s", s.name, s.typeStr)
}

// ComponentSymbol represents a component (document) symbol
type ComponentSymbol struct {
	BaseSymbol
	DocType    ast.DocType
	Properties map[string]*PropertySymbol
	States     map[string]*StateSymbol
	Imports    map[string]*ImportSymbol
}

func NewComponentSymbol(name string, pos ast.Position, docType ast.DocType) *ComponentSymbol {
	return &ComponentSymbol{
		BaseSymbol: BaseSymbol{
			name:     name,
			kind:     SymbolComponent,
			typeStr:  docType.String(),
			position: pos,
		},
		DocType:    docType,
		Properties: make(map[string]*PropertySymbol),
		States:     make(map[string]*StateSymbol),
		Imports:    make(map[string]*ImportSymbol),
	}
}

func (c *ComponentSymbol) String() string {
	return fmt.Sprintf("component %s (%s)", c.name, c.DocType)
}

// FunctionSymbol represents a function declaration
type FunctionSymbol struct {
	BaseSymbol
	Parameters []*ParameterSymbol
	ReturnType string
}

func NewFunctionSymbol(name, returnType string, pos ast.Position, params []*ParameterSymbol) *FunctionSymbol {
	return &FunctionSymbol{
		BaseSymbol: BaseSymbol{
			name:     name,
			kind:     SymbolFunction,
			typeStr:  returnType,
			position: pos,
		},
		Parameters: params,
		ReturnType: returnType,
	}
}

func (f *FunctionSymbol) String() string {
	var params []string
	for _, param := range f.Parameters {
		params = append(params, param.String())
	}
	return fmt.Sprintf("function %s(%s): %s", f.name, strings.Join(params, ", "), f.ReturnType)
}

// ParameterSymbol represents a function parameter
type ParameterSymbol struct {
	BaseSymbol
	IsOptional bool
}

func NewParameterSymbol(name, typeStr string, pos ast.Position, optional bool) *ParameterSymbol {
	return &ParameterSymbol{
		BaseSymbol: BaseSymbol{
			name:     name,
			kind:     SymbolParameter,
			typeStr:  typeStr,
			position: pos,
		},
		IsOptional: optional,
	}
}

func (p *ParameterSymbol) String() string {
	optional := ""
	if p.IsOptional {
		optional = "?"
	}
	return fmt.Sprintf("%s%s: %s", p.name, optional, p.typeStr)
}

// VariableSymbol represents a variable declaration
type VariableSymbol struct {
	BaseSymbol
	IsConst bool
}

func NewVariableSymbol(name, typeStr string, pos ast.Position, isConst bool) *VariableSymbol {
	return &VariableSymbol{
		BaseSymbol: BaseSymbol{
			name:     name,
			kind:     SymbolVariable,
			typeStr:  typeStr,
			position: pos,
		},
		IsConst: isConst,
	}
}

func (v *VariableSymbol) String() string {
	kind := "let"
	if v.IsConst {
		kind = "const"
	}
	return fmt.Sprintf("%s %s: %s", kind, v.name, v.typeStr)
}

// ImportSymbol represents an import statement
type ImportSymbol struct {
	BaseSymbol
	Path       string
	ImportKind ast.ImportKind
}

func NewImportSymbol(name, path string, pos ast.Position, kind ast.ImportKind) *ImportSymbol {
	return &ImportSymbol{
		BaseSymbol: BaseSymbol{
			name:     name,
			kind:     SymbolImport,
			typeStr:  kind.String(),
			position: pos,
		},
		Path:       path,
		ImportKind: kind,
	}
}

func (i *ImportSymbol) String() string {
	return fmt.Sprintf("import %s %s from %s", i.ImportKind, i.name, i.Path)
}

// BuiltInSymbol represents built-in functions and types
type BuiltInSymbol struct {
	BaseSymbol
	Description string
}

func NewBuiltInSymbol(name, typeStr, description string) *BuiltInSymbol {
	return &BuiltInSymbol{
		BaseSymbol: BaseSymbol{
			name:     name,
			kind:     SymbolBuiltIn,
			typeStr:  typeStr,
			position: ast.Position{}, // Built-ins don't have positions
		},
		Description: description,
	}
}

func (b *BuiltInSymbol) String() string {
	return fmt.Sprintf("built-in %s: %s", b.name, b.typeStr)
}

// Scope represents a lexical scope
type Scope struct {
	kind     ScopeKind
	parent   *Scope
	children []*Scope
	symbols  map[string]Symbol
	name     string // Optional name for debugging
}

func NewScope(kind ScopeKind, parent *Scope, name string) *Scope {
	scope := &Scope{
		kind:     kind,
		parent:   parent,
		children: make([]*Scope, 0),
		symbols:  make(map[string]Symbol),
		name:     name,
	}

	if parent != nil {
		parent.children = append(parent.children, scope)
	}

	return scope
}

func (s *Scope) Kind() ScopeKind { return s.kind }
func (s *Scope) Parent() *Scope  { return s.parent }
func (s *Scope) Name() string    { return s.name }

// Define adds a symbol to the current scope
func (s *Scope) Define(symbol Symbol) error {
	name := symbol.Name()
	if _, exists := s.symbols[name]; exists {
		return fmt.Errorf("symbol '%s' already defined in scope", name)
	}
	s.symbols[name] = symbol
	return nil
}

// Lookup searches for a symbol in the current scope
func (s *Scope) Lookup(name string) (Symbol, bool) {
	symbol, exists := s.symbols[name]
	return symbol, exists
}

// LookupRecursive searches for a symbol in the current scope and parent scopes
func (s *Scope) LookupRecursive(name string) (Symbol, bool) {
	if symbol, exists := s.symbols[name]; exists {
		return symbol, true
	}
	if s.parent != nil {
		return s.parent.LookupRecursive(name)
	}
	return nil, false
}

// GetAllSymbols returns all symbols in the current scope
func (s *Scope) GetAllSymbols() map[string]Symbol {
	result := make(map[string]Symbol)
	for name, symbol := range s.symbols {
		result[name] = symbol
	}
	return result
}

// GetSymbolsByKind returns all symbols of a specific kind in the current scope
func (s *Scope) GetSymbolsByKind(kind SymbolKind) []Symbol {
	var result []Symbol
	for _, symbol := range s.symbols {
		if symbol.Kind() == kind {
			result = append(result, symbol)
		}
	}
	return result
}

func (s *Scope) String() string {
	symbolCount := len(s.symbols)
	return fmt.Sprintf("Scope{%s, %d symbols}", s.kind, symbolCount)
}

// SymbolTable manages the symbol table and scopes
type SymbolTable struct {
	global     *Scope
	current    *Scope
	components map[string]*ComponentSymbol
	builtIns   map[string]*BuiltInSymbol
}

func NewSymbolTable() *SymbolTable {
	global := NewScope(ScopeGlobal, nil, "global")
	st := &SymbolTable{
		global:     global,
		current:    global,
		components: make(map[string]*ComponentSymbol),
		builtIns:   make(map[string]*BuiltInSymbol),
	}

	// Add built-in types and functions
	st.addBuiltIns()

	return st
}

func (st *SymbolTable) addBuiltIns() {
	// Built-in types
	builtInTypes := []struct {
		name, description string
	}{
		{"string", "String type"},
		{"number", "Number type"},
		{"boolean", "Boolean type"},
		{"object", "Object type"},
		{"array", "Array type"},
		{"void", "Void type"},
		{"any", "Any type"},
	}

	for _, builtin := range builtInTypes {
		symbol := NewBuiltInSymbol(builtin.name, "type", builtin.description)
		st.builtIns[builtin.name] = symbol
		st.global.Define(symbol)
	}

	// Built-in functions
	builtInFunctions := []struct {
		name, signature, description string
	}{
		{"console.log", "(...args: any[]) => void", "Console logging function"},
		{"Math.floor", "(x: number) => number", "Math floor function"},
		{"Math.ceil", "(x: number) => number", "Math ceil function"},
		{"Math.round", "(x: number) => number", "Math round function"},
	}

	for _, builtin := range builtInFunctions {
		symbol := NewBuiltInSymbol(builtin.name, builtin.signature, builtin.description)
		st.builtIns[builtin.name] = symbol
		st.global.Define(symbol)
	}
}

// EnterScope creates a new scope and makes it current
func (st *SymbolTable) EnterScope(kind ScopeKind, name string) *Scope {
	newScope := NewScope(kind, st.current, name)
	st.current = newScope
	return newScope
}

// ExitScope returns to the parent scope
func (st *SymbolTable) ExitScope() error {
	if st.current.parent == nil {
		return fmt.Errorf("cannot exit global scope")
	}
	st.current = st.current.parent
	return nil
}

// CurrentScope returns the current scope
func (st *SymbolTable) CurrentScope() *Scope {
	return st.current
}

// GlobalScope returns the global scope
func (st *SymbolTable) GlobalScope() *Scope {
	return st.global
}

// Define adds a symbol to the current scope
func (st *SymbolTable) Define(symbol Symbol) error {
	return st.current.Define(symbol)
}

// Lookup searches for a symbol starting from the current scope
func (st *SymbolTable) Lookup(name string) (Symbol, bool) {
	return st.current.LookupRecursive(name)
}

// DefineComponent adds a component symbol to the global scope
func (st *SymbolTable) DefineComponent(component *ComponentSymbol) error {
	st.components[component.Name()] = component
	return st.global.Define(component)
}

// GetComponent returns a component symbol by name
func (st *SymbolTable) GetComponent(name string) (*ComponentSymbol, bool) {
	component, exists := st.components[name]
	return component, exists
}

// GetAllComponents returns all registered components
func (st *SymbolTable) GetAllComponents() map[string]*ComponentSymbol {
	result := make(map[string]*ComponentSymbol)
	for name, component := range st.components {
		result[name] = component
	}
	return result
}

// IsBuiltIn checks if a symbol is a built-in
func (st *SymbolTable) IsBuiltIn(name string) bool {
	_, exists := st.builtIns[name]
	return exists
}

// GetBuiltIn returns a built-in symbol by name
func (st *SymbolTable) GetBuiltIn(name string) (*BuiltInSymbol, bool) {
	builtin, exists := st.builtIns[name]
	return builtin, exists
}

// String returns a string representation of the symbol table
func (st *SymbolTable) String() string {
	return fmt.Sprintf("SymbolTable{%d components, current: %s}",
		len(st.components), st.current.String())
}

// Debug prints the symbol table structure for debugging
func (st *SymbolTable) Debug() string {
	var builder strings.Builder
	builder.WriteString("Symbol Table Debug:\n")
	st.debugScope(st.global, 0, &builder)
	return builder.String()
}

func (st *SymbolTable) debugScope(scope *Scope, depth int, builder *strings.Builder) {
	indent := strings.Repeat("  ", depth)
	builder.WriteString(fmt.Sprintf("%s%s\n", indent, scope.String()))

	for name, symbol := range scope.symbols {
		builder.WriteString(fmt.Sprintf("%s  %s: %s\n", indent, name, symbol.String()))
	}

	for _, child := range scope.children {
		st.debugScope(child, depth+1, builder)
	}
}
