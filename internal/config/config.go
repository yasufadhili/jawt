package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config combines all configuration types
type Config struct {
	App  AppConfig
	Jawt JawtConfig
}

// JawtConfig holds the JAWT configuration from jawt.config.json
// Fixed to match the nested structure created by the initialiser
type JawtConfig struct {
	Project ProjectConfig `mapstructure:"project"`
	Server  ServerConfig  `mapstructure:"server"`
	Build   BuildConfig   `mapstructure:"build"`
}

// ProjectConfig holds project-specific settings
type ProjectConfig struct {
	Name string `mapstructure:"name"`
}

// ServerConfig holds server-specific settings
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// BuildConfig holds build-specific settings
type BuildConfig struct {
	Output string `mapstructure:"output"`
	Minify bool   `mapstructure:"minify"`
}

// AppConfig holds the application configuration from app.json
// This structure matches what the initialiser creates
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	Version     string `mapstructure:"version"`
	Author      string `mapstructure:"author"`
	// Optional fields that might be added later
	License      string   `mapstructure:"license,omitempty"`
	Dependencies []string `mapstructure:"dependencies,omitempty"`
}

// LoadConfig loads both app.json and jawt.config.json from the specified directory
func LoadConfig(projectDir string) (*Config, error) {
	config := &Config{}

	// Load app.json
	appConfig, err := loadAppConfig(projectDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load app.json: %w", err)
	}
	config.App = *appConfig

	// Load jawt.config.json
	jawtConfig, err := loadJawtConfig(projectDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load jawt.config.json: %w", err)
	}
	config.Jawt = *jawtConfig

	return config, nil
}

// loadAppConfig loads the app.json file
func loadAppConfig(projectDir string) (*AppConfig, error) {
	v := viper.New()

	appConfigPath := filepath.Join(projectDir, "app.json")
	if _, err := os.Stat(appConfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("app.json not found in %s", projectDir)
	}

	v.SetConfigFile(appConfigPath)
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading app.json: %w", err)
	}

	var appConfig AppConfig
	if err := v.Unmarshal(&appConfig); err != nil {
		return nil, fmt.Errorf("error unmarshalling app.json: %w", err)
	}

	return &appConfig, nil
}

// loadJawtConfig loads the jawt.config.json file
func loadJawtConfig(projectDir string) (*JawtConfig, error) {
	v := viper.New()

	jawtConfigPath := filepath.Join(projectDir, "jawt.config.json")
	if _, err := os.Stat(jawtConfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("jawt.config.json not found in %s", projectDir)
	}

	v.SetConfigFile(jawtConfigPath)
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading jawt.config.json: %w", err)
	}

	var jawtConfig JawtConfig
	if err := v.Unmarshal(&jawtConfig); err != nil {
		return nil, fmt.Errorf("error unmarshalling jawt.config.json: %w", err)
	}

	// Set defaults if not specified in config
	if jawtConfig.Server.Port == 0 {
		jawtConfig.Server.Port = 6500
	}
	if jawtConfig.Build.Output == "" {
		jawtConfig.Build.Output = "dist"
	}

	return &jawtConfig, nil
}

// IsJawtProject checks if the current directory is a JAWT project
// by checking for the existence of app.json and jawt.config.json
func IsJawtProject(dir string) bool {
	appPath := filepath.Join(dir, "app.json")
	jawtConfigPath := filepath.Join(dir, "jawt.config.json")

	_, appErr := os.Stat(appPath)
	_, jawtErr := os.Stat(jawtConfigPath)

	return !os.IsNotExist(appErr) && !os.IsNotExist(jawtErr)
}

// GetPort Convenience methods for accessing nested configuration values
func (c *Config) GetPort() int {
	return c.Jawt.Server.Port
}

func (c *Config) GetOutputDir() string {
	return c.Jawt.Build.Output
}

func (c *Config) GetProjectName() string {
	return c.App.Name
}

func (c *Config) ShouldMinify() bool {
	return c.Jawt.Build.Minify
}
