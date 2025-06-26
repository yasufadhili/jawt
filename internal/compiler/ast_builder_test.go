package compiler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ASTBuilder interface - should be implemented to build AST from source
type ASTBuilder interface {
	BuildAST(source string) (*JMLDocumentNode, error)
}

// Placeholder implementation - returns nil for now
type PlaceholderASTBuilder struct{}

func (b *PlaceholderASTBuilder) BuildAST(source string) (*JMLDocumentNode, error) {
	return nil, nil
}

type TestConfig struct {
	TestDataDir    string
	ExpectedDir    string
	ActualDir      string
	UpdateExpected bool // Set to true to update expected files instead of comparing
}

type TestCase struct {
	Name         string
	InputFile    string
	ExpectedFile string
	ShouldError  bool
}

type ASTTestSuite struct {
	config  TestConfig
	builder ASTBuilder
}

func NewASTTestSuite(config TestConfig, builder ASTBuilder) *ASTTestSuite {
	return &ASTTestSuite{
		config:  config,
		builder: builder,
	}
}

// RunAllTests Runs all tests in the test data directory
func (suite *ASTTestSuite) RunAllTests(t *testing.T) {
	testCases := suite.discoverTestCases(t)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			suite.runSingleTest(t, testCase)
		})
	}
}

// discoverTestCases Discovers test cases from the test data directory
func (suite *ASTTestSuite) discoverTestCases(t *testing.T) []TestCase {
	var testCases []TestCase

	inputPattern := filepath.Join(suite.config.TestDataDir, "*.jml")
	matches, err := filepath.Glob(inputPattern)
	if err != nil {
		t.Fatalf("Failed to find test files: %v", err)
	}

	for _, inputFile := range matches {
		baseName := strings.TrimSuffix(filepath.Base(inputFile), ".jml")
		expectedFile := filepath.Join(suite.config.ExpectedDir, baseName+".json")

		// Check if this should be an error test case
		shouldError := strings.Contains(baseName, "_error") || strings.Contains(baseName, "_invalid")

		testCases = append(testCases, TestCase{
			Name:         baseName,
			InputFile:    inputFile,
			ExpectedFile: expectedFile,
			ShouldError:  shouldError,
		})
	}

	return testCases
}

// runSingleTest Runs a single test case
func (suite *ASTTestSuite) runSingleTest(t *testing.T, testCase TestCase) {
	source, err := os.ReadFile(testCase.InputFile)
	if err != nil {
		t.Fatalf("Failed to read input file %s: %v", testCase.InputFile, err)
	}

	ast, err := suite.builder.BuildAST(string(source))

	if testCase.ShouldError {
		if err == nil {
			t.Errorf("Expected error but got none for test case %s", testCase.Name)
		}
		return // Don't check JSON for error cases
	}

	if err != nil {
		t.Errorf("Unexpected error for test case %s: %v", testCase.Name, err)
		return
	}

	// Convert AST to JSON
	actualJSON, err := suite.astToJSON(ast)
	if err != nil {
		t.Errorf("Failed to convert AST to JSON for test case %s: %v", testCase.Name, err)
		return
	}

	// Write actual result for debugging
	actualFile := filepath.Join(suite.config.ActualDir, testCase.Name+".json")
	suite.ensureDir(filepath.Dir(actualFile))
	if err := os.WriteFile(actualFile, actualJSON, 0644); err != nil {
		t.Errorf("Failed to write actual result file: %v", err)
	}

	if suite.config.UpdateExpected {
		// Update the expected file instead of comparing
		suite.ensureDir(filepath.Dir(testCase.ExpectedFile))
		if err := os.WriteFile(testCase.ExpectedFile, actualJSON, 0644); err != nil {
			t.Errorf("Failed to update expected file %s: %v", testCase.ExpectedFile, err)
		}
		return
	}

	// Compare with the expected result
	expectedJSON, err := os.ReadFile(testCase.ExpectedFile)
	if err != nil {
		t.Errorf("Failed to read expected file %s: %v", testCase.ExpectedFile, err)
		t.Logf("Actual result written to: %s", actualFile)
		return
	}

	if !suite.compareJSON(expectedJSON, actualJSON) {
		t.Errorf("AST mismatch for test case %s", testCase.Name)
		t.Logf("Expected file: %s", testCase.ExpectedFile)
		t.Logf("Actual file: %s", actualFile)
		t.Logf("Run with -update flag to update expected results")
	}
}

// Convert AST to pretty-printed JSON
func (suite *ASTTestSuite) astToJSON(ast *JMLDocumentNode) ([]byte, error) {
	// Convert AST to a JSON-serializable structure
	jsonData := suite.astToMap(ast)
	return json.MarshalIndent(jsonData, "", "  ")
}

// Convert AST node to map for JSON serialisation
func (suite *ASTTestSuite) astToMap(node interface{}) interface{} {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *JMLDocumentNode:
		return map[string]interface{}{
			"type":    "JMLDocument",
			"doctype": suite.astToMap(n.Doctype),
			"imports": suite.sliceToMap(n.Imports),
			"content": suite.astToMap(n.Content),
		}

	case *DoctypeSpecifierNode:
		return map[string]interface{}{
			"type":    "DoctypeSpecifier",
			"doctype": n.Doctype,
			"name":    n.Name,
		}

	case *ImportStatementNode:
		return map[string]interface{}{
			"type":       "ImportStatement",
			"importType": n.Type,
			"identifier": n.Identifier,
			"from":       n.From,
			"isBrowser":  n.IsBrowser,
		}

	case *PageDefinitionNode:
		return map[string]interface{}{
			"type":       "PageDefinition",
			"properties": suite.sliceToMap(n.Properties),
			"child":      suite.astToMap(n.Child),
		}

	case *ComponentDefinitionNode:
		return map[string]interface{}{
			"type":    "ComponentDefinition",
			"element": suite.astToMap(n.Element),
		}

	case *ComponentElementNode:
		return map[string]interface{}{
			"type":       "ComponentElement",
			"name":       n.Name,
			"properties": suite.sliceToMap(n.Properties),
			"children":   suite.sliceToMap(n.Children),
		}

	case *PropertyNode:
		return map[string]interface{}{
			"type":  "Property",
			"name":  n.Name,
			"value": suite.astToMap(n.Value),
		}

	case *LiteralNode:
		return map[string]interface{}{
			"type":      "Literal",
			"valueType": n.Type,
			"value":     n.Value,
		}

	default:
		return fmt.Sprintf("Unknown node type: %T", node)
	}
}

// sliceToMap Converts a slice to a map representation
func (suite *ASTTestSuite) sliceToMap(slice interface{}) []interface{} {
	switch s := slice.(type) {
	case []*ImportStatementNode:
		result := make([]interface{}, len(s))
		for i, item := range s {
			result[i] = suite.astToMap(item)
		}
		return result

	case []*PropertyNode:
		result := make([]interface{}, len(s))
		for i, item := range s {
			result[i] = suite.astToMap(item)
		}
		return result

	case []*ComponentElementNode:
		result := make([]interface{}, len(s))
		for i, item := range s {
			result[i] = suite.astToMap(item)
		}
		return result

	default:
		return []interface{}{}
	}
}

// compareJSON Compares two JSON byte arrays (ignoring formatting)
func (suite *ASTTestSuite) compareJSON(expected, actual []byte) bool {
	var expectedData, actualData interface{}

	if err := json.Unmarshal(expected, &expectedData); err != nil {
		return false
	}

	if err := json.Unmarshal(actual, &actualData); err != nil {
		return false
	}

	expectedJSON, _ := json.Marshal(expectedData)
	actualJSON, _ := json.Marshal(actualData)

	return string(expectedJSON) == string(actualJSON)
}

// Ensure directory exists
func (suite *ASTTestSuite) ensureDir(dir string) {
	os.MkdirAll(dir, 0755)
}

// Test setup
func TestASTBuilder(t *testing.T) {
	config := TestConfig{
		TestDataDir:    "testdata/input",
		ExpectedDir:    "testdata/expected",
		ActualDir:      "testdata/actual",
		UpdateExpected: false, // Set to true to update expected files
	}

	builder := &PlaceholderASTBuilder{}
	suite := NewASTTestSuite(config, builder)

	suite.RunAllTests(t)
}

// Helper function to run tests with an update flag
func TestASTBuilderUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping update test in short mode")
	}

	config := TestConfig{
		TestDataDir:    "testdata/input",
		ExpectedDir:    "testdata/expected",
		ActualDir:      "testdata/actual",
		UpdateExpected: true,
	}

	builder := &PlaceholderASTBuilder{}
	suite := NewASTTestSuite(config, builder)

	suite.RunAllTests(t)
}

func BenchmarkASTBuilder(b *testing.B) {
	builder := &PlaceholderASTBuilder{}
	source := `_doctype page test
	
Page {
  title: "Test Page"
  Layout {
    content: "Hello World"
  }
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := builder.BuildAST(source)
		if err != nil {
			b.Fatalf("AST building failed: %v", err)
		}
	}
}
