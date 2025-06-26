package compiler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestJSONPrinter tests the JSONPrinter's ability to serialise an AST and write it to a file.
func TestJSONPrinter(t *testing.T) {

	ast := &JMLDocumentNode{
		Doctype: &DoctypeSpecifierNode{
			Doctype: "page",
			Name:    "TestPage",
		},
		Imports: []*ImportStatementNode{
			{
				Type:       "component",
				Identifier: "Button",
				From:       "./components/button.jml",
				IsBrowser:  false,
			},
			{
				Type:      "browser",
				IsBrowser: true,
			},
		},
		Content: &PageDefinitionNode{
			Properties: []*PropertyNode{
				{
					Name: "title",
					Value: &LiteralNode{
						Type:  "string",
						Value: "My Page",
					},
				},
			},
			Child: &ComponentElementNode{
				Name: "Button",
				Properties: []*PropertyNode{
					{
						Name: "label",
						Value: &LiteralNode{
							Type:  "string",
							Value: "Click Me",
						},
					},
				},
				Children: []*ComponentElementNode{
					{
						Name: "Text",
						Properties: []*PropertyNode{
							{
								Name: "value",
								Value: &LiteralNode{
									Type:  "string",
									Value: "Hello, World!",
								},
							},
						},
					},
				},
			},
		},
	}

	printer := NewJSONPrinter()

	jsonStr, err := printer.Print(ast)
	if err != nil {
		t.Fatalf("Failed to print AST to JSON: %v", err)
	}

	// Verify JSON is valid by parsing it
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}

	// Check some key fields in the JSON
	if jsonData["type"] != "JMLDocument" {
		t.Errorf("Expected root node type to be 'JMLDocument', got '%v'", jsonData["type"])
	}
	if doctype, ok := jsonData["doctype"].(map[string]interface{}); !ok || doctype["doctype"] != "page" {
		t.Errorf("Expected doctype to be 'page', got '%v'", doctype["doctype"])
	}
	if imports, ok := jsonData["imports"].([]interface{}); !ok || len(imports) != 2 {
		t.Errorf("Expected 2 imports, got %v", len(imports))
	}

	outputPath := filepath.Join(t.TempDir(), "ast_output.json")
	err = os.WriteFile(outputPath, []byte(jsonStr), 0644)
	if err != nil {
		t.Fatalf("Failed to write JSON to file: %v", err)
	}

	// Read the file back and verify contents
	fileContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if string(fileContent) != jsonStr {
		t.Errorf("File content does not match generated JSON")
	}

	// Verify specific content in the file
	var fileData map[string]interface{}
	if err := json.Unmarshal(fileContent, &fileData); err != nil {
		t.Fatalf("Failed to parse JSON from file: %v", err)
	}
	if fileData["type"] != "JMLDocument" {
		t.Errorf("Expected root node type in file to be 'JMLDocument', got '%v'", fileData["type"])
	}
}
