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

// ProjectConfig represents project-specific configuration (jawt.project.json)
type ProjectConfig struct {
	// Project metadata
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	// Build configuration
	OutputDir string `json:"output_dir"`
	DistDir   string `json:"dist_dir"`
	ShadowDOM bool   `json:"shadow_dom"`

	// Development settings
	DevPort    int      `json:"dev_port"`
	EnableHMR  bool     `json:"enable_hmr"`
	WatchPaths []string `json:"watch_paths"`

	// TypeScript configuration
	TSConfigPath string `json:"ts_config_path"`

	// Tailwind configuration
	TailwindConfigPath string `json:"tailwind_config_path"`

	// Custom build scripts
	PreBuildScripts  []string `json:"pre_build_scripts"`
	PostBuildScripts []string `json:"post_build_scripts"`
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
		Name:               "jawt-project",
		Version:            "1.0.0",
		Description:        "A Jawt application",
		OutputDir:          ".jawt/build",
		DistDir:            ".jawt/dist",
		ShadowDOM:          false,
		DevPort:            6500,
		EnableHMR:          true,
		WatchPaths:         []string{"app", "components", "scripts"},
		TSConfigPath:       "tsconfig.json",
		TailwindConfigPath: "tailwind.config.js",
		PreBuildScripts:    []string{},
		PostBuildScripts:   []string{},
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
func (tc *JawtConfig) Save(configPath string) error {
	data, err := json.MarshalIndent(tc, "", "  ")
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
func (tc *JawtConfig) Validate() error {
	if tc.DefaultPort <= 0 || tc.DefaultPort > 65535 {
		return fmt.Errorf("invalid default port: %d", tc.DefaultPort)
	}

	if tc.TempDir == "" {
		return fmt.Errorf("temp directory cannot be empty")
	}

	if tc.CacheDir == "" {
		return fmt.Errorf("cache directory cannot be empty")
	}

	return nil
}

// Validate validates the project configuration
func (pc *ProjectConfig) Validate() error {
	if pc.Name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if pc.OutputDir == "" {
		return fmt.Errorf("output directory cannot be empty")
	}

	if pc.DistDir == "" {
		return fmt.Errorf("dist directory cannot be empty")
	}

	if pc.DevPort <= 0 || pc.DevPort > 65535 {
		return fmt.Errorf("invalid dev port: %d", pc.DevPort)
	}

	return nil
}
