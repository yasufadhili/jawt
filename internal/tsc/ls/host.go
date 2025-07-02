package ls

import (
	"github.com/yasufadhili/jawt/internal/tsc/compiler"
	"github.com/yasufadhili/jawt/internal/tsc/lsp/lsproto"
)

type Host interface {
	GetProgram() *compiler.Program
	GetPositionEncoding() lsproto.PositionEncodingKind
	GetLineMap(fileName string) *LineMap
}
