package project

import "time"

type Structure struct {
	Root       string                    `json:"root"`
	Config     *Config                   `json:"config"`
	Pages      map[string]*PageInfo      `json:"pages"`
	Components map[string]*ComponentInfo `json:"components"`
	Assets     []string                  `json:"assets"`
	BuildTime  time.Time                 `json:"build_time"`
	TempDir    string
}

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

type ServerConfig struct {
	Port int `json:"port"`
}
