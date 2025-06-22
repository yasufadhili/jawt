package page_compiler

import (
	"fmt"
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
	fmt.Println("Compiling page...")
	fmt.Println(pc.pageInfo)
	pc.pageInfo.Compiled = true
	return nil
}
