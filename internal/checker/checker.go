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
		BaseVisitor: &ast.BaseVisitor{},
		reporter:    reporter,
		table:       t,
	}
}
