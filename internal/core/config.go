package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// JawtConfig represents the global jawt configuration (jawt.config.json)
type JawtConfig struct {
	// External jawt paths
	TypeScriptPath string `json:"typescript_path"`
	TailwindPath   string `json:"tailwind_path"`
	NodePath       string `json:"node_path"`

	// Jawt-specific configurations
	DefaultPort int    `json:"default_port"`
	TempDir     string `json:"temp_dir"`
	CacheDir    string `json:"cache_dir"`

	// Build optimisation
	EnableMinification bool `json:"enable_minification"`
	EnableSourceMaps   bool `json:"enable_source_maps"`
	EnableTreeShaking  bool `json:"enable_tree_shaking"`
}

// ProjectConfig represents the new project-specific configuration structure
type ProjectConfig struct {
	App struct {
		Name        string `json:"name"`
		Author      string `json:"author"`
		Version     string `json:"version"`
		Description string `json:"description"`
	} `json:"app"`
	Components struct {
		Path  string `json:"path"`
		Alias string `json:"alias"`
	} `json:"components"`
	Pages struct {
		Path  string `json:"path"`
		Alias string `json:"alias"`
	} `json:"pages"`
	Scripts struct {
		Path  string `json:"path"`
		Alias string `json:"alias"`
	} `json:"scripts"`
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	Build struct {
		OutputDir string `json:"outputDir"`
		Minify    bool   `json:"minify"`
	} `json:"build"`
}

// DefaultJawtConfig returns a default jawt configuration
func DefaultJawtConfig() *JawtConfig {
	return &JawtConfig{
		TypeScriptPath:     "tsc",
		TailwindPath:       "tailwindcss",
		NodePath:           "node",
		DefaultPort:        6500,
		TempDir:            ".jawt/tmp",
		CacheDir:           ".jawt/cache",
		EnableMinification: true,
		EnableSourceMaps:   true,
		EnableTreeShaking:  true,
	}
}

// DefaultProjectConfig returns a default project configuration
func DefaultProjectConfig() *ProjectConfig {
	return &ProjectConfig{
		App: struct {
			Name        string `json:"name"`
			Author      string `json:"author"`
			Version     string `json:"version"`
			Description string `json:"description"`
		}{
			Name:   "jawt-project",
			Author: "",
		},
		Components: struct {
			Path  string `json:"path"`
			Alias string `json:"alias"`
		}{
			Path:  "components",
			Alias: "",
		},
		Pages: struct {
			Path  string `json:"path"`
			Alias string `json:"alias"`
		}{
			Path:  "pages",
			Alias: "",
		},
		Scripts: struct {
			Path  string `json:"path"`
			Alias string `json:"alias"`
		}{
			Path:  "scripts",
			Alias: "",
		},
		Server: struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}{
			Host: "localhost",
			Port: 6500,
		},
		Build: struct {
			OutputDir string `json:"outputDir"`
			Minify    bool   `json:"minify"`
		}{
			OutputDir: "build",
			Minify:    true,
		},
	}
}

// LoadJawtConfig loads jawt configuration from the specified path
func LoadJawtConfig(configPath string) (*JawtConfig, error) {
	// If no config path specified, use default
	if configPath == "" {
		return DefaultJawtConfig(), nil
	}

	// Check if a config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultJawtConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read jawt config: %w", err)
	}

	config := DefaultJawtConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse jawt config: %w", err)
	}

	return config, nil
}

// LoadProjectConfig loads project configuration from the specified path
func LoadProjectConfig(projectDir string) (*ProjectConfig, error) {
	configPath := filepath.Join(projectDir, "jawt.project.json")

	// Check if a config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultProjectConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read project config: %w", err)
	}

	config := DefaultProjectConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse project config: %w", err)
	}

	return config, nil
}

// Save saves the jawt configuration to the specified path
func (jc *JawtConfig) Save(configPath string) error {
	data, err := json.MarshalIndent(jc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal jawt config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// Save saves the project configuration to the project directory
func (pc *ProjectConfig) Save(projectDir string) error {
	configPath := filepath.Join(projectDir, "jawt.project.json")

	data, err := json.MarshalIndent(pc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// Validate validates the jawt configuration
func (jc *JawtConfig) Validate() error {
	if jc.DefaultPort <= 0 || jc.DefaultPort > 65535 {
		return fmt.Errorf("invalid default port: %d", jc.DefaultPort)
	}

	if jc.TempDir == "" {
		return fmt.Errorf("temp directory cannot be empty")
	}

	if jc.CacheDir == "" {
		return fmt.Errorf("cache directory cannot be empty")
	}

	return nil
}

// Validate validates the project configuration
func (pc *ProjectConfig) Validate() error {
	if pc.App.Name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if pc.Components.Path == "" {
		return fmt.Errorf("components path cannot be empty")
	}

	if pc.Pages.Path == "" {
		return fmt.Errorf("pages path cannot be empty")
	}

	if pc.Scripts.Path == "" {
		return fmt.Errorf("scripts path cannot be empty")
	}

	if pc.Build.OutputDir == "" {
		return fmt.Errorf("build output directory cannot be empty")
	}

	if pc.Server.Port <= 0 || pc.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", pc.Server.Port)
	}

	if pc.Server.Host == "" {
		return fmt.Errorf("server host cannot be empty")
	}

	return nil
}

// GetComponentsPath returns the full path to the components directory
func (pc *ProjectConfig) GetComponentsPath(projectRoot string) string {
	return filepath.Join(projectRoot, pc.Components.Path)
}

// GetPagesPath returns the full path to the pages directory
func (pc *ProjectConfig) GetPagesPath(projectRoot string) string {
	return filepath.Join(projectRoot, pc.Pages.Path)
}

// GetScriptsPath returns the full path to the scripts directory
func (pc *ProjectConfig) GetScriptsPath(projectRoot string) string {
	return filepath.Join(projectRoot, pc.Scripts.Path)
}

// GetBuildPath returns the full path to the build output directory
func (pc *ProjectConfig) GetBuildPath(projectRoot string) string {
	return filepath.Join(projectRoot, pc.Build.OutputDir)
}

// GetServerAddress returns the full server address
func (pc *ProjectConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", pc.Server.Host, pc.Server.Port)
}

// IsMinificationEnabled returns whether minification is enabled
func (pc *ProjectConfig) IsMinificationEnabled() bool {
	return pc.Build.Minify
}

// SetProjectName sets the project name
func (pc *ProjectConfig) SetProjectName(name string) {
	pc.App.Name = name
}

// SetAuthor sets the project author
func (pc *ProjectConfig) SetAuthor(author string) {
	pc.App.Author = author
}

// SetServerPort sets the server port
func (pc *ProjectConfig) SetServerPort(port int) {
	pc.Server.Port = port
}

// SetMinification enables or disables minification
func (pc *ProjectConfig) SetMinification(enabled bool) {
	pc.Build.Minify = enabled
}
