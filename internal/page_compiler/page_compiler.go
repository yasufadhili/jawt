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
