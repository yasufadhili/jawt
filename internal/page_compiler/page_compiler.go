package page_compiler

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"github.com/yasufadhili/jawt/internal/error_handler"
	parser "github.com/yasufadhili/jawt/internal/page_compiler/parser/generated"
	"github.com/yasufadhili/jawt/internal/project"
	"os"
	"path/filepath"
)

type PageCompiler struct {
	pageInfo   *project.PageInfo
	parser     *Parser
	outputPath string
}

func NewPageCompiler(pageInfo *project.PageInfo, outPath string) *PageCompiler {
	return &PageCompiler{
		pageInfo:   pageInfo,
		parser:     newParser(),
		outputPath: outPath,
	}
}

// CompileResult holds the compilation result with detailed error information
type CompileResult struct {
	Success   bool
	Errors    []error_handler.SyntaxError
	AST       *Page
	ParseTree antlr.ParseTree
	PageInfo  *project.PageInfo
}

func (pc *PageCompiler) CompilePage() (*CompileResult, error) {
	fmt.Printf("Compiling page %s:%s\n", pc.pageInfo.Name, pc.pageInfo.AbsolutePath)

	input, err := antlr.NewFileStream(pc.pageInfo.AbsolutePath)
	if err != nil {
		return &CompileResult{
			Success:  false,
			PageInfo: pc.pageInfo,
		}, fmt.Errorf("failed to read file %s: %w", pc.pageInfo.AbsolutePath, err)
	}

	parseResult := pc.parseFile(input)

	// Create compile result
	result := &CompileResult{
		Success:   parseResult.Success,
		Errors:    parseResult.Errors,
		ParseTree: parseResult.Tree,
		PageInfo:  pc.pageInfo,
	}

	// If parsing failed, return early with errors
	if !parseResult.Success {
		fmt.Printf("❌ Parsing failed for %s with %d errors:\n", pc.pageInfo.Name, len(parseResult.Errors))
		pc.printErrors(parseResult.Errors)
		return result, nil
	}

	astBuilder := NewAstBuilder()
	astRoot := astBuilder.Visit(parseResult.Tree).(*Page)
	result.AST = astRoot

	if astRoot.Doctype != nil {
		//fmt.Printf("AST Built: %+v\n", astRoot.Doctype)
	}

	emitConfig := &EmitterConfig{
		TailwindCDN: "https://cdn.tailwindcss.com",
		CustomCSS: []string{
			"https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap",
		},
		BodyClasses: []string{
			"bg-gradient-to-br",
			"from-blue-50",
			"to-indigo-100",
			"min-h-screen",
			"font-inter",
		},
		WrapperClass: "max-w-4xl mx-auto px-6 py-12",
	}

	configurableEmitter := NewConfigurableHTMLEmitter(emitConfig)
	customHTML := configurableEmitter.EmitHTML(astRoot)

	outPath := filepath.Join(pc.outputPath, pc.pageInfo.RelativePath)

	err = os.WriteFile(outPath+"/index.html", []byte(customHTML), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write HTML file: %w", err)
	}

	pc.pageInfo.Compiled = true
	return result, nil
}

// parseFile handles the actual parsing with custom error handling
func (pc *PageCompiler) parseFile(input antlr.CharStream) ParseResult {
	// Reset parser state for a new file
	pc.parser.errorListener.Reset()
	pc.parser.errorStrategy.Reset()

	lexer := parser.NewJMLPageLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(pc.parser.errorListener)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewJMLPageParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(pc.parser.errorListener)
	p.SetErrorHandler(pc.parser.errorStrategy)

	tree := p.Page()

	return ParseResult{
		Success: !pc.parser.errorListener.HasErrors(),
		Errors:  pc.parser.errorListener.GetErrors(),
		Tree:    tree,
	}
}

// CompilePageWithRecovery attempts compilation with error recovery
func (pc *PageCompiler) CompilePageWithRecovery() (*CompileResult, error) {
	fmt.Printf("Compiling page with recovery %s:%s\n", pc.pageInfo.Name, pc.pageInfo.AbsolutePath)

	result, err := pc.CompilePage()
	if err != nil {
		return result, err
	}

	// If there were syntax errors, but we want to attempt partial compilation
	if !result.Success && len(result.Errors) > 0 {
		fmt.Printf("⚠️  Attempting partial compilation despite %d syntax errors\n", len(result.Errors))

		// Try to build AST even with errors (if a parse tree was created)
		if result.ParseTree != nil {
			astBuilder := NewAstBuilder()
			// Use a try-catch equivalent for AST building
			if astRoot := pc.tryBuildAST(astBuilder, result.ParseTree); astRoot != nil {
				result.AST = astRoot
				fmt.Println("✅ Partial AST built successfully")
			}
		}
	}

	return result, nil
}

// tryBuildAST attempts to build AST with error recovery
func (pc *PageCompiler) tryBuildAST(astBuilder *AstBuilder, tree antlr.ParseTree) *Page {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("⚠️  AST building failed: %v\n", r)
		}
	}()

	if tree == nil {
		return nil
	}

	astRoot := astBuilder.Visit(tree)
	if astRoot == nil {
		return nil
	}

	if page, ok := astRoot.(*Page); ok {
		return page
	}

	return nil
}

// printErrors displays syntax errors in a user-friendly format
func (pc *PageCompiler) printErrors(errors []error_handler.SyntaxError) {
	for i, err := range errors {
		fmt.Printf("  %d. Line %d:%d - %s\n", i+1, err.Line, err.Column, err.Message)
		if err.Symbol != "" && err.Symbol != "<EOF>" {
			fmt.Printf("     Near symbol: '%s'\n", err.Symbol)
		}
		if err.Context != "" {
			fmt.Printf("     Context: %s\n", err.Context)
		}

		// Add suggestions if available
		if suggestion := pc.getSuggestionForError(err); suggestion != "" {
			fmt.Printf("     💡 Suggestion: %s\n", suggestion)
		}
		fmt.Println()
	}
}

// getSuggestionForError provides helpful suggestions based on error patterns
func (pc *PageCompiler) getSuggestionForError(err error_handler.SyntaxError) string {
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
		if contains(errorMsg, pattern) {
			return suggestion
		}
	}
	return "Please check your JML page syntax"
}

// contains is a helper function for case-insensitive string matching
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				anySubstring(s, substr)))
}

func anySubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ValidatePage performs syntax validation without full compilation
func (pc *PageCompiler) ValidatePage() (bool, []error_handler.SyntaxError, error) {
	input, err := antlr.NewFileStream(pc.pageInfo.AbsolutePath)
	if err != nil {
		return false, nil, fmt.Errorf("failed to read file %s: %w", pc.pageInfo.AbsolutePath, err)
	}

	parseResult := pc.parseFile(input)
	return parseResult.Success, parseResult.Errors, nil
}

// ParseResult holds the result of parsing with error information
type ParseResult struct {
	Success bool
	Errors  []error_handler.SyntaxError
	Tree    antlr.ParseTree
}

// Parser wraps ANTLR parser with error handling
type Parser struct {
	errorListener *error_handler.SyntaxErrorListener
	errorStrategy *error_handler.ErrorStrategy
}

// newParser creates a new parser with error handling
func newParser() *Parser {
	return &Parser{
		errorListener: error_handler.NewSyntaxErrorListener(),
		errorStrategy: error_handler.NewErrorStrategy(3), // Allow up to 3 recovery attempts
	}
}

// CompilePageLegacy Legacy method for backwards compatibility - now returns error information
func (pc *PageCompiler) CompilePageLegacy() error {
	result, err := pc.CompilePage()
	if err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("compilation failed with %d syntax errors", len(result.Errors))
	}

	return nil
}

// GetCompilerStats returns statistics about the compilation process
func (pc *PageCompiler) GetCompilerStats() map[string]interface{} {
	return map[string]interface{}{
		"page_name":      pc.pageInfo.Name,
		"page_path":      pc.pageInfo.AbsolutePath,
		"is_compiled":    pc.pageInfo.Compiled,
		"error_recovery": pc.parser.errorStrategy != nil,
		"max_recovery":   3,
	}
}
