package project

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

// Config holds configuration for the project (jawt.config.json)
type Config struct {
	App          AppConfig          `json:"app"`
	Server       DevServerConfig    `json:"server"`
	Pages        DocumentConfig     `json:"pages"`
	Components   DocumentConfig     `json:"components"`
	Scripts      DocumentConfig     `json:"scripts"`
	Dependencies []DependencyConfig `json:"dependencies"`
	Build        BuildConfig        `json:"build"`
}

type AppConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Author  string `json:"author"`
	License string `json:"license"`
}

type DocumentConfig struct {
	Path  string `json:"path"`
	Alias string `json:"alias"`
}

type DependencyConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"`
}

type BuildConfig struct {
	OutputPath string `json:"outputPath"`
	Minify     bool   `json:"minify"`
}

func loadJawtConfig(projectDir string) (*Config, error) {
	v := viper.New()
	jawtConfigPath := filepath.Join(projectDir, "jawt.config.json")
	if _, err := os.Stat(jawtConfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("jawt.config.json not found in %s", projectDir)
	}

	v.SetConfigFile(jawtConfigPath)
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read jawt.config.json: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal jawt.config.json: %w", err)
	}

	// Set defaults if not specified in config
	if config.Server.Port == 0 {
		config.Server.Port = 6500
	}
	if config.Pages.Path == "" {
		config.Pages.Path = "app"
	}
	if config.Components.Path == "" {
		config.Components.Path = "components"
	}
	if config.Build.OutputPath == "" {
		config.Build.OutputPath = "dist"
	}

	return &config, nil
}

func loadAppConfig(projectDir string) (*AppConfig, error) {
	v := viper.New()

	appConfigPath := filepath.Join(projectDir, "app.json")
	if _, err := os.Stat(appConfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("app.json not found in %s", projectDir)
	}
	v.SetConfigFile(appConfigPath)
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read app.json: %w", err)
	}

	var config AppConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal app.json: %w", err)
	}

	return &config, nil
}

// DocumentInfo contains common metadata shared by all document types
type DocumentInfo struct {
	Name         string            `json:"name"`
	Title        string            `json:"title"`
	RelativePath string            `json:"relative_path"`
	AbsolutePath string            `json:"absolute_path"`
	Dependencies []string          `json:"dependencies"`
	Imports      map[string]string `json:"imports"`
	LastModified time.Time         `json:"last_modified"`
	Compiled     bool              `json:"compiled"`
}

// Document interface defines common behaviour for all document types
type Document interface {
	GetName() string
	GetTitle() string
	GetRelativePath() string
	GetAbsolutePath() string
	GetDependencies() []string
	GetImports() map[string]string
	GetLastModified() time.Time
	IsCompiled() bool
	SetCompiled(bool)
}

// PageInfo contains metadata about a page file
type PageInfo struct {
	DocumentInfo
	Route string `json:"route"`
}

// ComponentInfo contains metadata about a component file
type ComponentInfo struct {
	DocumentInfo
	UsedBy []string `json:"used_by"` // Pages/components that use this component
}

func (d *DocumentInfo) GetName() string {
	return d.Name
}

func (d *DocumentInfo) GetTitle() string {
	return d.Title
}

func (d *DocumentInfo) GetRelativePath() string {
	return d.RelativePath
}

func (d *DocumentInfo) GetAbsolutePath() string {
	return d.AbsolutePath
}

func (d *DocumentInfo) GetDependencies() []string {
	return d.Dependencies
}

func (d *DocumentInfo) GetImports() map[string]string {
	return d.Imports
}

func (d *DocumentInfo) GetLastModified() time.Time {
	return d.LastModified
}

func (d *DocumentInfo) IsCompiled() bool {
	return d.Compiled
}

func (d *DocumentInfo) SetCompiled(compiled bool) {
	d.Compiled = compiled
}

func (p *PageInfo) GetRoute() string {
	return p.Route
}

func (p *PageInfo) SetRoute(route string) {
	p.Route = route
}

func (c *ComponentInfo) GetUsedBy() []string {
	return c.UsedBy
}

func (c *ComponentInfo) AddUsedBy(name string) {
	// Avoid duplicates
	for _, existing := range c.UsedBy {
		if existing == name {
			return
		}
	}
	c.UsedBy = append(c.UsedBy, name)
}

func (c *ComponentInfo) RemoveUsedBy(name string) {
	for i, existing := range c.UsedBy {
		if existing == name {
			c.UsedBy = append(c.UsedBy[:i], c.UsedBy[i+1:]...)
			return
		}
	}
}

func (p *Project) GetAllDocuments() []Document {
	var documents []Document

	for _, page := range p.Pages {
		documents = append(documents, page)
	}

	for _, component := range p.Components {
		documents = append(documents, component)
	}

	return documents
}

func (p *Project) GetCompiledDocuments() []Document {
	var compiled []Document

	for _, doc := range p.GetAllDocuments() {
		if doc.IsCompiled() {
			compiled = append(compiled, doc)
		}
	}

	return compiled
}

func (p *Project) GetUncompiledDocuments() []Document {
	var uncompiled []Document

	for _, doc := range p.GetAllDocuments() {
		if !doc.IsCompiled() {
			uncompiled = append(uncompiled, doc)
		}
	}

	return uncompiled
}

func (p *Project) FindPageByRoute(route string) *PageInfo {
	for _, page := range p.Pages {
		if page.Route == route {
			return page
		}
	}
	return nil
}

func (p *Project) FindComponentByName(name string) *ComponentInfo {
	return p.Components[name]
}

func (p *Project) FindPageByName(name string) *PageInfo {
	return p.Pages[name]
}

// Convenience methods for accessing nested configuration values

func (c *Config) GetPort() int { return int(c.Server.Port) }

func (c *Config) GetOutputDir() string { return c.Build.OutputPath }

func (c *Config) GetProjectName() string {
	return c.App.Name
}

func (c *Config) ShouldMinify() bool {
	return c.Build.Minify
}

type DevServerConfig struct {
	Host string `json:"host"`
	Port int16  `json:"port"`
}
