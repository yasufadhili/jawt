package printer

import (
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/tspath"
)

type SourceFileMetaDataProvider interface {
	GetSourceFileMetaData(path tspath.Path) *ast.SourceFileMetaData
}
