package pc

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	parser "github.com/yasufadhili/jawt/internal/pc/parser/generated"
)

type PageCompiler struct {
	inputPath  string
	outputPath string
	filename   string
}

func NewPageCompiler(filepath string, outputPath string) *PageCompiler {
	return &PageCompiler{
		inputPath:  filepath,
		outputPath: outputPath,
	}
}

func (pc *PageCompiler) CompilePage() error {

	var allErrors []error

	// 1. Lexing
	input, err := antlr.NewFileStream(pc.inputPath)
	if err != nil {
		return err
	}
	lexer := parser.NewJMLPageLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// 2. Parsing
	p := parser.NewJMLPageParser(stream)

	syntaxErrorListener := NewErrorListener()
	p.RemoveErrorListeners()
	p.AddErrorListener(syntaxErrorListener)

	tree := p.Program()

	if len(syntaxErrorListener.Errors) > 0 {
		for _, se := range syntaxErrorListener.Errors {
			allErrors = append(allErrors, se)
		}
		return fmt.Errorf("syntax errors encountered during parsing:\n%v", formatErrors(allErrors))
	}
	fmt.Println("Parsing: OK")

	// 3. Build AST
	fmt.Println("Building AST")
	astBuilder := NewAstBuilder()
	astRoot := astBuilder.Visit(tree).(*Program)
	fmt.Printf("AST Built: %+v\n", astRoot.Doctype) // Just a quick check

	// 4. Semantic Analysis
	fmt.Println("Semantic Analysis")
	analyser := NewSemanticAnalyser()
	analyser.Visit(astRoot)

	if len(analyser.Errors) > 0 {
		for _, err := range analyser.Errors {
			fmt.Println("-", err)
		}
		return fmt.Errorf("semantic analysis failed")
	} else {
		fmt.Println("Semantic Analysis: OK")
	}

	// 5. Emit
	fmt.Println("Emitting")
	emit := NewEmitter()
	emit.Visit(astRoot)
	fmt.Println("Emit: OK")
	fmt.Println(emit.Emitted)

	return nil
}

// Helper function to format a slice of errors into a single string
func formatErrors(errors []error) string {
	if len(errors) == 0 {
		return "No errors."
	}
	s := ""
	for i, err := range errors {
		s += fmt.Sprintf("  %d. %s\n", i+1, err.Error())
	}
	return s
}
