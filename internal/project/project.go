package project

import "time"

// Structure represents the entire project structure
type Structure struct {
	Root       string                    `json:"root"`
	Config     *Config                   `json:"config"`
	Pages      map[string]*PageInfo      `json:"pages"`
	Components map[string]*ComponentInfo `json:"components"`
	Assets     []string                  `json:"assets"`
	BuildTime  time.Time                 `json:"build_time"`
	TempDir    string                    `json:"temp_dir"`
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

// Config holds configuration for the project to be used by the build system
type Config struct {
	Name        string       `json:"name"`
	Version     string       `json:"version"`
	Author      string       `json:"author"`
	Description string       `json:"description"`
	Server      ServerConfig `json:"server"`
}

// ServerConfig contains server-specific configuration
type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
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

func (s *Structure) GetAllDocuments() []Document {
	var documents []Document

	for _, page := range s.Pages {
		documents = append(documents, page)
	}

	for _, component := range s.Components {
		documents = append(documents, component)
	}

	return documents
}

func (s *Structure) GetCompiledDocuments() []Document {
	var compiled []Document

	for _, doc := range s.GetAllDocuments() {
		if doc.IsCompiled() {
			compiled = append(compiled, doc)
		}
	}

	return compiled
}

func (s *Structure) GetUncompiledDocuments() []Document {
	var uncompiled []Document

	for _, doc := range s.GetAllDocuments() {
		if !doc.IsCompiled() {
			uncompiled = append(uncompiled, doc)
		}
	}

	return uncompiled
}

func (s *Structure) FindPageByRoute(route string) *PageInfo {
	for _, page := range s.Pages {
		if page.Route == route {
			return page
		}
	}
	return nil
}

func (s *Structure) FindComponentByName(name string) *ComponentInfo {
	return s.Components[name]
}

func (s *Structure) FindPageByName(name string) *PageInfo {
	return s.Pages[name]
}
