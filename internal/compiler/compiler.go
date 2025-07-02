package compiler

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"github.com/yasufadhili/jawt/internal/ast"
	"github.com/yasufadhili/jawt/internal/common"
	parser "github.com/yasufadhili/jawt/internal/compiler/parser/generated"
	"github.com/yasufadhili/jawt/internal/emitter"
	"github.com/yasufadhili/jawt/internal/project"
	"os"
	"path/filepath"
	"strings"
)

type FileCompiler struct {
	manager *Manager
	parser  *Parser
	docInfo *project.DocumentInfo
	target  common.BuildTarget
}

func NewFileCompiler(manager *Manager, docInfo *project.DocumentInfo, target common.BuildTarget) *FileCompiler {
	return &FileCompiler{
		manager: manager,
		docInfo: docInfo,
		parser:  newParser(),
		target:  target,
	}
}

func (c *FileCompiler) CompileFile() (*FileCompileResult, error) {

	input, err := antlr.NewFileStream(c.docInfo.AbsolutePath)
	if err != nil {
		return &FileCompileResult{
			Success: false,
			DocInfo: c.docInfo,
		}, fmt.Errorf("failed to read file %s: %w", c.docInfo.AbsolutePath, err)
	}

	parseResult := c.parseFile(input)

	result := &FileCompileResult{
		Success:   parseResult.Success,
		DocInfo:   c.docInfo,
		ParseTree: parseResult.Tree,
		Errors:    parseResult.Errors,
	}

	if !parseResult.Success {
		fmt.Printf("❌ Parsing failed for %s with %d errors:\n", c.docInfo.Name, len(parseResult.Errors))
		c.printErrors(parseResult.Errors)
		return result, nil
	}

	astBuilder := ast.NewAstBuilder()
	astRoot := astBuilder.Visit(parseResult.Tree).(*ast.DocumentNode)
	result.AST = astRoot

	// TODO: Symbol Collection & Semantic Analysis -> Via TypeScript too

	e := emitter.NewEmitter(astRoot, c.target)
	output := e.Emit()

	outPath := c.manager.project.OutputDir
	switch c.target {
	case common.TargetPage:
		if c.docInfo.RelativePath == "/" {
			c.docInfo.RelativePath = "index"
		} else {
			c.docInfo.RelativePath = filepath.Join(c.docInfo.RelativePath, "index")
		}
		outPath = filepath.Join(outPath, c.docInfo.RelativePath+".html")
	case common.TargetComponent:
		outPath = filepath.Join(outPath, c.docInfo.RelativePath+".js")
		outPath = strings.ReplaceAll(outPath, ".jml", "")
	}

	outDir := filepath.Dir(outPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory %s: %w", outDir, err)
	}

	osErr := os.WriteFile(outPath, []byte(output), 0644)
	if osErr != nil {
		return nil, fmt.Errorf("failed to write output file %s: %w", outPath, osErr)
	}

	c.docInfo.Compiled = true
	return result, nil
}

// FileCompileResult holds the compilation result with detailed error information
type FileCompileResult struct {
	Success   bool
	Errors    []SyntaxError
	AST       *ast.DocumentNode
	ParseTree antlr.ParseTree
	DocInfo   *project.DocumentInfo
}

// parseFile handles the actual parsing with custom error handling
func (c *FileCompiler) parseFile(input antlr.CharStream) ParseResult {
	// Reset parser state for a new file
	c.parser.errorListener.Reset()
	c.parser.errorStrategy.Reset()

	lexer := parser.NewJMLLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(c.parser.errorListener)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewJMLParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(c.parser.errorListener)
	p.SetErrorHandler(c.parser.errorStrategy)

	tree := p.Document()

	return ParseResult{
		Success: !c.parser.errorListener.HasErrors(),
		Errors:  c.parser.errorListener.GetErrors(),
		Tree:    tree,
	}
}

// ParseResult holds the result of parsing with error information
type ParseResult struct {
	Success bool
	Errors  []SyntaxError
	Tree    antlr.ParseTree
}

type Parser struct {
	errorListener *SyntaxErrorListener
	errorStrategy *ErrorStrategy
}

// newParser creates a new parser with error handling
func newParser() *Parser {
	return &Parser{
		errorListener: NewSyntaxErrorListener(),
		errorStrategy: NewErrorStrategy(3), // Allow up to 3 recovery attempts
	}
}

// printErrors displays syntax errors in a user-friendly format
func (c *FileCompiler) printErrors(errors []SyntaxError) {
	for i, err := range errors {
		fmt.Printf("  %d. Line %d:%d - %s\n", i+1, err.Line, err.Column, err.Message)
		if err.Symbol != "" && err.Symbol != "<EOF>" {
			fmt.Printf("     Near symbol: '%s'\n", err.Symbol)
		}
		if err.Context != "" {
			fmt.Printf("     Context: %s\n", err.Context)
		}

		// Add suggestions if available
		if suggestion := c.getSuggestionForError(err); suggestion != "" {
			fmt.Printf("     💡 Suggestion: %s\n", suggestion)
		}
		fmt.Println()
	}
}

// getSuggestionForError provides helpful suggestions based on error patterns
func (c *FileCompiler) getSuggestionForError(err SyntaxError) string {
	commonMistakes := map[string]string{
		"missing ')'":           "Try adding a closing parenthesis",
		"missing '}'":           "Try adding a closing brace",
		"missing '>'":           "Try adding a closing angle bracket for tag",
		"missing ';'":           "Try adding a semicolon",
		"extraneous input":      "Try removing the unexpected token",
		"mismatched input":      "Check if you're using the correct syntax",
		"no viable alternative": "Check the grammar rules for this context",
		"missing EOF":           "There might be unclosed tags or brackets",
	}

	errorMsg := err.Message
	for pattern, suggestion := range commonMistakes {
		if strings.Contains(errorMsg, pattern) {
			return suggestion
		}
	}
	return "Please check your JML page syntax"
}
