package sourcemap

import "github.com/yasufadhili/jawt/internal/tsc/core"

type Source interface {
	Text() string
	FileName() string
	LineMap() []core.TextPos
}
