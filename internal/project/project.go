package project

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Manager struct {
	Project *Project
	Options DevOptions
}

type DevOptions struct {
	Verbose bool
}

func NewProjectManager(project *Project) *Manager {
	return &Manager{
		Project: project,
	}
}

// LoadProject initialises a project from a root directory
func (m *Manager) LoadProject() (*Project, error) {
	return nil, nil
}

// ValidateProject checks project structure and configuration
func (m *Manager) ValidateProject() []Error {
	return nil
}

// GetProjectConfig returns resolved configuration with cascading
func (m *Manager) GetProjectConfig() *Config {
	return m.Project.Config
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
	OutputDir  string                    `json:"output_dir"`
}

func NewProject(rootPath string) (*Project, error) {
	config, err := LoadConfig(rootPath)
	if err != nil {
		return nil, err
	}
	return &Project{
		RootPath: rootPath,
		Config:   config,
	}, nil
}

// LoadConfig loads both app.json and jawt.config.json from the specified directory
func LoadConfig(projectPath string) (*Config, error) {
	config := &Config{}

	// Load app.json
	appConfig, err := loadAppConfig(projectPath)
	if err != nil {
		return nil, err
	}
	config.App = *appConfig

	// Load jawt.config.json
	jawtConfig, err := loadJawtConfig(projectPath)
	if err != nil {
		return nil, err
	}
	config = jawtConfig
	jawtConfig.App = *appConfig

	if config.App.Name == "" {
		return nil, fmt.Errorf("project name is required in app.json")
	}

	return config, nil
}

type Error struct {
	Message string `json:"message"`
}

type Metadata struct{}

// IsJawtProject checks if the current directory is a JAWT project
// by checking for the existence of app.json and jawt.config.json
func IsJawtProject(dir string) bool {
	appConfigPath := filepath.Join(dir, "app.json")
	jawtConfigPath := filepath.Join(dir, "jawt.config.json")

	_, appErr := os.Stat(appConfigPath)
	_, jawtErr := os.Stat(jawtConfigPath)

	return !os.IsNotExist(appErr) && !os.IsNotExist(jawtErr)
}

type ChangeCallback func(event *FileChangeEvent) error

type FileChangeEvent struct {
	Type      ChangeType
	FilePath  string
	Timestamp time.Time
}

type ChangeType int
