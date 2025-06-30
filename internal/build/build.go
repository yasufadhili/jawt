package build

import "github.com/yasufadhili/jawt/internal/project"

type Builder struct {
	Project    *project.Project
	ClearCache bool
}

func NewBuilder(project *project.Project) (*Builder, error) {
	return &Builder{
		Project:    project,
		ClearCache: false,
	}, nil
}
