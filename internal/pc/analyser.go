package pc

import "fmt"

type SemanticAnalyser struct {
	BaseVisitor // Embed the base visitor for default traversal
	Errors      []string
}

func NewSemanticAnalyser() *SemanticAnalyser {
	return &SemanticAnalyser{}
}

func (s *SemanticAnalyser) VisitProgram(node *Program) interface{} {
	fmt.Println("Semantic Analysis: Starting Page Program...")

	// continue traversal with base the visitor
	s.BaseVisitor.VisitProgram(node)

	fmt.Println("Semantic Analysis: Finished Page Program.")
	return nil
}

func (s *SemanticAnalyser) visitImportStatement(node *ImportStatement) interface{} {
	fmt.Printf("  Checking Import: %s %s from %s\n", node.Doctype, node.Identifier, node.From)
	if node.From == "" {
		s.Errors = append(s.Errors, fmt.Sprintf("Import for %s cannot have an empty 'from' path.", node.Identifier))
	}
	return nil
}

func (s *SemanticAnalyser) visitPageProperty(node *PageProperty) interface{} {
	fmt.Printf("    Checking Page Property: %s = %v (type: %T)\n", node.Key, node.Value, node.Value)

	// Example semantic check: 'version' property must be an integer
	if node.Key == "version" {
		if _, ok := node.Value.(int); !ok {
			s.Errors = append(s.Errors, fmt.Sprintf("Property '%s' expected integer value, got %T.", node.Key, node.Value))
		}
	}
	return nil
}
