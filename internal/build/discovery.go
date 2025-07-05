package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yasufadhili/jawt/internal/core"
)

// DiscoverProjectFiles finds all JML files in the project
func DiscoverProjectFiles(ctx *core.JawtContext) ([]string, error) {
	ctx.Logger.Info("Discovering JML files in project")

	// Get paths to search
	pagesDir := ctx.ProjectConfig.GetPagesPath(ctx.Paths.ProjectRoot)
	componentsDir := ctx.ProjectConfig.GetComponentsPath(ctx.Paths.ProjectRoot)

	ctx.Logger.Debug("Searching directories",
		core.StringField("pages", pagesDir),
		core.StringField("components", componentsDir))

	// Find all .jml files
	var jmlFiles []string

	// Search pages directory
	pageFiles, err := findJMLFiles(pagesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to search pages directory: %w", err)
	}
	jmlFiles = append(jmlFiles, pageFiles...)

	// Search components directory
	componentFiles, err := findJMLFiles(componentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to search components directory: %w", err)
	}
	jmlFiles = append(jmlFiles, componentFiles...)

	ctx.Logger.Info("Jml files discovered",
		core.IntField("count", len(jmlFiles)))

	return jmlFiles, nil
}

// findJMLFiles recursively finds all .jml files in a directory
func findJMLFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".jml") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// CreateDocumentInfo creates a DocumentInfo from a JML file
func CreateDocumentInfo(path string, projectRoot string) (*DocumentInfo, error) {
	// Get file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Determine the document type based on directory
	var docType DocumentType
	if strings.Contains(path, "/pages/") || strings.Contains(path, "/app/") {
		docType = DocumentTypePage
	} else if strings.Contains(path, "/components/") {
		docType = DocumentTypeComponent
	} else {
		// Default to component if can't determine
		docType = DocumentTypeComponent
	}

	// Create document info
	relPath, err := filepath.Rel(projectRoot, path)
	if err != nil {
		relPath = path // Fallback to absolute path
	}

	// Extract name from filename (without extension)
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	return &DocumentInfo{
		Name:         name,
		RelPath:      relPath,
		AbsPath:      path,
		Type:         docType,
		Dependencies: []string{},
		DependedBy:   []string{},
		IsCompiled:   false,
		LastModified: fileInfo.ModTime(),
		Hash:         "", // TODO: Implement content hashing
	}, nil
}

// AnalyseDependencies analyses dependencies between documents
func AnalyseDependencies(docs map[string]*DocumentInfo) error {
	// TODO: Implement dependency analysis
	// 1. Parse each document to find import statements
	// 2. Update Dependencies and DependedBy fields

	return nil
}
