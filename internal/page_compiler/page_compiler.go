package page_compiler

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	parser "github.com/yasufadhili/jawt/internal/pc/parser/generated"
	"github.com/yasufadhili/jawt/internal/project"
)

type PageCompiler struct {
	pageInfo *project.PageInfo
}

func NewPageCompiler(pageInfo *project.PageInfo) *PageCompiler {
	return &PageCompiler{
		pageInfo: pageInfo,
	}
}

func (pc *PageCompiler) CompilePage() error {

	fmt.Printf("Compiling page %s:%s\n", pc.pageInfo.Name, pc.pageInfo.AbsolutePath)

	input, err := antlr.NewFileStream(pc.pageInfo.AbsolutePath)
	if err != nil {
		return err
	}
	lexer := parser.NewJMLPageLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewJMLPageParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))

	tree := p.Page()

	fmt.Println(tree.ToStringTree(nil, p))

	// Build AST
	astBuilder := NewAstBuilder()
	astRoot := astBuilder.Visit(tree).(*Page)
	if astRoot.Doctype != nil {
		fmt.Printf("AST Built: %+v\n", astRoot.Doctype) // Just a quick check
	}

	pc.pageInfo.Compiled = true
	return nil
}

func (pc *PageCompiler) lexPage() (*antlr.CommonTokenStream, error) {

	input, err := antlr.NewFileStream(pc.pageInfo.AbsolutePath)
	if err != nil {
		return nil, err
	}
	lexer := parser.NewJMLPageLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	return stream, nil
}

func (pc *PageCompiler) parsePage() (parser.IPageContext, []error) {
	var errors []error
	stream, err := pc.lexPage()
	if err != nil {
		errors = append(errors, err)
		return nil, errors
	}

	p := parser.NewJMLPageParser(stream)

	tree := p.Page()

	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))

	return tree, nil
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
