package build

import "time"

// DocumentType represents the type of document (page or component)
type DocumentType int

const (
	DocumentTypePage DocumentType = iota
	DocumentTypeComponent
)

// DocumentInfo represents common information for all document types
type DocumentInfo struct {
	Name         string
	RelPath      string
	AbsPath      string
	Type         DocumentType
	Dependencies []string
	DependedBy   []string
	IsCompiled   bool
	LastModified time.Time
	Hash         string
}

// ComponentInfo represents information about a component
type ComponentInfo struct {
	DocumentInfo
	Props map[string]string
}

// PageInfo represents information about a page
type PageInfo struct {
	DocumentInfo
	Route string
}
