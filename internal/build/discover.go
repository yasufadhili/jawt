package build

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/core"
	"os"
	"path/filepath"
	"strings"
)

type ProjectDiscoverer struct {
	ctx *core.JawtContext
}

func NewProjectDiscoverer(ctx *core.JawtContext) ProjectDiscoverer {
	return ProjectDiscoverer{ctx: ctx}
}

func (pd *ProjectDiscoverer) DiscoverProjectFiles() ([]string, error) {
	pd.ctx.Logger.Info("Discovering JML files in project")

	// Get paths to search
	pagesDir := pd.ctx.ProjectConfig.GetPagesPath(pd.ctx.Paths.ProjectRoot)
	componentsDir := pd.ctx.ProjectConfig.GetComponentsPath(pd.ctx.Paths.ProjectRoot)

	pd.ctx.Logger.Debug("Searching directories",
		core.StringField("pages", pagesDir),
		core.StringField("components", componentsDir))

	// Find all .jml files
	var jmlFiles []string

	// Search pages directory
	pageFiles, err := pd.findJMLFiles(pagesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to search pages directory: %w", err)
	}
	jmlFiles = append(jmlFiles, pageFiles...)

	// Search components directory
	componentFiles, err := pd.findJMLFiles(componentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to search components directory: %w", err)
	}
	jmlFiles = append(jmlFiles, componentFiles...)

	pd.ctx.Logger.Info("Jml files discovered",
		core.IntField("count", len(jmlFiles)))

	return jmlFiles, nil
}

func (pd *ProjectDiscoverer) findJMLFiles(dir string) ([]string, error) {
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

func (pd *ProjectDiscoverer) CreateDocumentInfo(path string, projectRoot string) (*DocumentInfo, error) {

	// Get file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Determine the document type based on directory
	var docType DocumentType
	if strings.Contains(path, "/"+pd.ctx.Paths.AppDir+"/") {
		docType = DocumentTypePage
	} else if strings.Contains(path, "/"+pd.ctx.Paths.ComponentsDir+"/") {
		docType = DocumentTypeComponent
	} else {
		// Default to component if can't determine
		docType = DocumentTypeComponent
	}

	// Create document info
	relPath, err := filepath.Rel(projectRoot, path)
	if err != nil {
		relPath = path // Fallback to the absolute path
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
		Hash:         "",
	}, nil
}

// AnalyseDependencies analyses dependencies between documents
func AnalyseDependencies(docs map[string]*DocumentInfo) error {
	// TODO: Implement dependency analysis
	// 1. Parse each document to find import statements
	// 2. Update Dependencies and DependedBy fields

	return nil
}
