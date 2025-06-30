package project

import "time"

type Manager struct {
	Project *Project
}

func NewProjectManager(project *Project) *Manager {
	return &Manager{
		Project: project,
	}
}

// LoadProject initialises a project from a root directory
func (m *Manager) LoadProject(rootPath string) (*Project, error) {
	return nil, nil
}

// ValidateProject checks project structure and configuration
func (m *Manager) ValidateProject() []Error {
	return nil
}

// GetProjectConfig returns resolved configuration with cascading
func (m *Manager) GetProjectConfig() (*Config, error) {
	return nil, nil
}

// WatchProject sets up the file system watching for development mode
func (m *Manager) WatchProject() error {
	return nil
}

type Project struct {
	RootPath   string                    `json:"root_path"`
	Config     *Config                   `json:"config"`
	Metadata   *Metadata                 `json:"metadata"`
	Pages      map[string]*PageInfo      `json:"pages"`
	Components map[string]*ComponentInfo `json:"components"`
	Assets     []string                  `json:"assets"`
	BuildTime  time.Time                 `json:"build_time"`
	TempDir    string                    `json:"temp_dir"`
}

func NewProject(rootPath string) *Project {
	return &Project{
		RootPath: rootPath,
	}
}

type Error struct {
	Message string `json:"message"`
}

type Metadata struct{}
