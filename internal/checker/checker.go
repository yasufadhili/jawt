package checker

import (
	"github.com/yasufadhili/jawt/internal/ast"
	"github.com/yasufadhili/jawt/internal/diagnostic"
)

type Checker struct {
	*ast.BaseVisitor
	reporter *diagnostic.Reporter
	table    *SymbolTable
}

func NewChecker(reporter *diagnostic.Reporter) *Checker {
	t := NewSymbolTable()
	return &Checker{
		reporter: reporter,
		table:    t,
	}
}

func (c *Checker) Check(program *ast.Program) {
	program.Accept(c)
}
