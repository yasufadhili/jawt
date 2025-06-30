package project

type Manager struct {
	Project *Project
	TempDir string
}

type Project struct {
	RootPath string    `json:"root_path"`
	Config   *Config   `json:"config"`
	Metadata *Metadata `json:"metadata"`
}

func NewProjectManager() *Manager {
	return &Manager{}
}

// LoadProject initialises a project from a root directory
func (p *Manager) LoadProject(rootPath string) (*Project, error) {
	return nil, nil
}

// ValidateProject checks project structure and configuration
func (p *Manager) ValidateProject() []Error {
	return nil
}

// GetProjectConfig returns resolved configuration with cascading
func (p *Manager) GetProjectConfig() (*Config, error) {
	return nil, nil
}

// WatchProject sets up the file system watching for development mode
func (p *Manager) WatchProject() error {
	return nil
}

type Error struct {
	Message string `json:"message"`
}

type Metadata struct{}
