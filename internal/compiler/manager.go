package compiler

import (
	"github.com/yasufadhili/jawt/internal/project"
	"time"
)

type Manager struct {
	project *project.Project
	//depGraph   *DependencyGraph
	buildCache *Cache
}

// Cache tracks file modification times and content hashes
type Cache struct {
	CacheFile string                `json:"-"`
	Files     map[string]FileRecord `json:"files"`
}

// FileRecord stores information about a compiled file
type FileRecord struct {
	LastModified time.Time `json:"last_modified"`
	Hash         string    `json:"hash"`
	OutputPath   string    `json:"output_path"`
	Dependencies []string  `json:"dependencies"`
}

func NewCompilerManager(project *project.Project) *Manager {
	cm := &Manager{
		project: project,
		//depGraph: NewDependencyGraph(),
	}

	//cm.buildDependencyGraph()
	//cm.initBuildCache()

	return cm
}

func (cm *Manager) CompileChanged() error {
	return nil
}
