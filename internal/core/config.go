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
	Paths struct {
		Components string `json:"components"`
		Pages      string `json:"pages"`
		Scripts    string `json:"scripts"`
		Assets     string `json:"assets"`
	} `json:"paths"`
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	Build struct {
		OutputDir string `json:"outputDir"`
		DistDir   string `json:"distDir"`
		Minify    bool   `json:"minify"`
		ShadowDOM bool   `json:"shadowDOM"`
	} `json:"build"`
	Dev struct {
		Port       int      `json:"port"`
		EnableHMR  bool     `json:"enableHMR"`
		WatchPaths []string `json:"watchPaths"`
	} `json:"dev"`
	Tooling struct {
		TSConfigPath       string `json:"tsConfigPath"`
		TailwindConfigPath string `json:"tailwindConfigPath"`
	} `json:"tooling"`
	Scripts struct {
		PreBuild  []string `json:"preBuild"`
		PostBuild []string `json:"postBuild"`
	} `json:"scripts"`
}

// BuildOptions represents build-time options and detected features
type BuildOptions struct {
	UsesTailwindCSS bool
}

// NewBuildOptions creates a new BuildOptions instance
func NewBuildOptions() *BuildOptions {
	return &BuildOptions{
		UsesTailwindCSS: false,
	}
}

// DefaultJawtConfig returns a default jawt configuration
func DefaultJawtConfig() *JawtConfig {
	return &JawtConfig{
		TypeScriptPath:     "tsc",
		TailwindPath:       "tailwindcss",
		NodePath:           "node",
		DefaultPort:        6500,
		TempDir:            filepath.Join(".jawt", "temp"),
		CacheDir:           filepath.Join(".jawt", "cache"),
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
			Name:        "jawt-project",
			Author:      "",
			Version:     "0.1.0",
			Description: "A Jawt application",
		},
		Paths: struct {
			Components string `json:"components"`
			Pages      string `json:"pages"`
			Scripts    string `json:"scripts"`
			Assets     string `json:"assets"`
		}{
			Components: "components",
			Pages:      "app",
			Scripts:    "scripts",
			Assets:     "assets",
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
			DistDir   string `json:"distDir"`
			Minify    bool   `json:"minify"`
			ShadowDOM bool   `json:"shadowDOM"`
		}{
			OutputDir: filepath.Join(".jawt", "build"),
			DistDir:   filepath.Join(".jawt", "dist"),
			Minify:    true,
			ShadowDOM: false,
		},
		Dev: struct {
			Port       int      `json:"port"`
			EnableHMR  bool     `json:"enableHMR"`
			WatchPaths []string `json:"watchPaths"`
		}{
			Port:       6500,
			EnableHMR:  true,
			WatchPaths: []string{"app", "components", "scripts", "assets"},
		},
		Tooling: struct {
			TSConfigPath       string `json:"tsConfigPath"`
			TailwindConfigPath string `json:"tailwindConfigPath"`
		}{
			TSConfigPath:       "tsconfig.json",
			TailwindConfigPath: "tailwind.config.js",
		},
		Scripts: struct {
			PreBuild  []string `json:"preBuild"`
			PostBuild []string `json:"postBuild"`
		}{
			PreBuild:  []string{},
			PostBuild: []string{},
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
		return fmt.Errorf("app name cannot be empty")
	}

	if pc.Paths.Components == "" {
		return fmt.Errorf("components path cannot be empty")
	}

	if pc.Paths.Pages == "" {
		return fmt.Errorf("pages path cannot be empty")
	}

	if pc.Paths.Scripts == "" {
		return fmt.Errorf("scripts path cannot be empty")
	}

	if pc.Paths.Assets == "" {
		return fmt.Errorf("assets path cannot be empty")
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

	if pc.Dev.Port <= 0 || pc.Dev.Port > 65535 {
		return fmt.Errorf("invalid dev server port: %d", pc.Dev.Port)
	}

	return nil
}

// GetComponentsPath returns the full path to the components directory
func (pc *ProjectConfig) GetComponentsPath(projectRoot string) string {
	return filepath.Join(projectRoot, pc.Paths.Components)
}

// GetPagesPath returns the full path to the pages directory
func (pc *ProjectConfig) GetPagesPath(projectRoot string) string {
	return filepath.Join(projectRoot, pc.Paths.Pages)
}

// GetScriptsPath returns the full path to the scripts directory
func (pc *ProjectConfig) GetScriptsPath(projectRoot string) string {
	return filepath.Join(projectRoot, pc.Paths.Scripts)
}

// GetAssetsPath returns the full path to the assets directory
func (pc *ProjectConfig) GetAssetsPath(projectRoot string) string {
	return filepath.Join(projectRoot, pc.Paths.Assets)
}

// GetBuildOutputDir returns the full path to the build output directory
func (pc *ProjectConfig) GetBuildOutputDir(projectRoot string) string {
	return filepath.Join(projectRoot, pc.Build.OutputDir)
}

// GetDistDir returns the full path to the distribution directory
func (pc *ProjectConfig) GetDistDir(projectRoot string) string {
	return filepath.Join(projectRoot, pc.Build.DistDir)
}

// GetServerAddress returns the full server address
func (pc *ProjectConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", pc.Server.Host, pc.Server.Port)
}

// GetDevServerAddress returns the full dev server address
func (pc *ProjectConfig) GetDevServerAddress() string {
	return fmt.Sprintf("%s:%d", pc.Server.Host, pc.Dev.Port)
}

// IsMinificationEnabled returns whether minification is enabled
func (pc *ProjectConfig) IsMinificationEnabled() bool {
	return pc.Build.Minify
}

// IsShadowDOMEnabled returns whether Shadow DOM is enabled
func (pc *ProjectConfig) IsShadowDOMEnabled() bool {
	return pc.Build.ShadowDOM
}

// IsHMR enabled returns whether HMR is enabled
func (pc *ProjectConfig) IsHMRenabled() bool {
	return pc.Dev.EnableHMR
}

// GetWatchPaths returns the paths to watch for changes
func (pc *ProjectConfig) GetWatchPaths() []string {
	return pc.Dev.WatchPaths
}

// GetTSConfigPath returns the path to the TypeScript config file
func (pc *ProjectConfig) GetTSConfigPath(projectRoot string) string {
	return filepath.Join(projectRoot, pc.Tooling.TSConfigPath)
}

// GetTailwindConfigPath returns the path to the Tailwind CSS config file
func (pc *ProjectConfig) GetTailwindConfigPath(projectRoot string) string {
	return filepath.Join(projectRoot, pc.Tooling.TailwindConfigPath)
}

// GetPreBuildScripts returns the pre-build scripts
func (pc *ProjectConfig) GetPreBuildScripts() []string {
	return pc.Scripts.PreBuild
}

// GetPostBuildScripts returns the post-build scripts
func (pc *ProjectConfig) GetPostBuildScripts() []string {
	return pc.Scripts.PostBuild
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

// SetDevServerPort sets the dev server port
func (pc *ProjectConfig) SetDevServerPort(port int) {
	pc.Dev.Port = port
}

// SetMinification enables or disables minification
func (pc *ProjectConfig) SetMinification(enabled bool) {
	pc.Build.Minify = enabled
}

// SetShadowDOM enables or disables Shadow DOM
func (pc *ProjectConfig) SetShadowDOM(enabled bool) {
	pc.Build.ShadowDOM = enabled
}

// SetHMR enables or disables HMR
func (pc *ProjectConfig) SetHMR(enabled bool) {
	pc.Dev.EnableHMR = enabled
}
