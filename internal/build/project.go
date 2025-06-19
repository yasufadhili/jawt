package build

import "time"

type ProjectStructure struct {
	Root       string                    `json:"root"`
	Config     *ProjectConfig            `json:"config"`
	Pages      map[string]*PageInfo      `json:"pages"`
	Components map[string]*ComponentInfo `json:"components"`
	Assets     []string                  `json:"assets"`
	BuildTime  time.Time                 `json:"build_time"`
}

// PageInfo contains metadata about a page file
type PageInfo struct {
	Name         string            `json:"name"`
	Title        string            `json:"title"`
	RelativePath string            `json:"relative_path"`
	AbsolutePath string            `json:"absolute_path"`
	Route        string            `json:"route"`
	Dependencies []string          `json:"dependencies"`
	Imports      map[string]string `json:"imports"`
	LastModified time.Time         `json:"last_modified"`
	Compiled     bool              `json:"compiled"`
}

// ComponentInfo contains metadata about a component file
type ComponentInfo struct {
	Name         string            `json:"name"`
	Title        string            `json:"title"`
	RelativePath string            `json:"relative_path"`
	AbsolutePath string            `json:"absolute_path"`
	Dependencies []string          `json:"dependencies"`
	Imports      map[string]string `json:"imports"`
	LastModified time.Time         `json:"last_modified"`
	Compiled     bool              `json:"compiled"`
	UsedBy       []string          `json:"used_by"` // Pages/components that use this component
}

// ProjectConfig holds configuration for the build system
type ProjectConfig struct {
	Name        string       `json:"name"`
	Version     string       `json:"version"`
	Author      string       `json:"author"`
	Description string       `json:"description"`
	Server      ServerConfig `json:"server"`
}

type ServerConfig struct {
	Port int `json:"port"`
}
