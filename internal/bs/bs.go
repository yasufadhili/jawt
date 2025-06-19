package bs

import (
	"github.com/yasufadhili/jawt/internal/cc"
	"github.com/yasufadhili/jawt/internal/pc"
)

type BuildSystem struct {
	ProjectRoot string
	Pages       []pc.Page
	Components  []cc.Component
}

func NewBuildSystem() *BuildSystem {
	return &BuildSystem{}
}
