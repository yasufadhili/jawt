package build

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/yasufadhili/jawt/internal/project"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ProjectDiscovery handles scanning and analysing project files
type ProjectDiscovery struct {
	rootPath string
}

func NewProjectDiscovery(rootPath string) *ProjectDiscovery {
	return &ProjectDiscovery{
		rootPath: rootPath,
	}
}

// DiscoverProject scans the entire project and builds the project structure
func (pd *ProjectDiscovery) DiscoverProject() (*project.Structure, error) {
	absRoot, err := filepath.Abs(pd.rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	p := &project.Structure{
		Root:       absRoot,
		Pages:      make(map[string]*project.PageInfo),
		Components: make(map[string]*project.ComponentInfo),
		Assets:     make([]string, 0),
		BuildTime:  time.Now(),
	}

	// Load project configuration
	if err := pd.loadProjectConfig(p); err != nil {
		return nil, fmt.Errorf("failed to load project config: %w", err)
	}

	// Discover pages
	if err := pd.discoverPages(p); err != nil {
		return nil, fmt.Errorf("failed to discover pages: %w", err)
	}

	// Discover components
	if err := pd.discoverComponents(p); err != nil {
		return nil, fmt.Errorf("failed to discover components: %w", err)
	}

	// Discover assets
	if err := pd.discoverAssets(p); err != nil {
		return nil, fmt.Errorf("failed to discover assets: %w", err)
	}

	// Build dependency graph
	if err := pd.buildDependencyGraph(p); err != nil {
		return nil, fmt.Errorf("failed to build dependency graph: %w", err)
	}

	return p, nil
}

func (pd *ProjectDiscovery) loadProjectConfig(p *project.Structure) error {

	name, err := readJsonField(p.Root+"/app.json", "name")
	if err != nil {
		return err
	}

	p.Config = &project.Config{
		Name: name.(string),
	}

	return nil
}

func readJsonField(filename string, field string) (interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	file, err = os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON into a map
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	value, exists := result[field]
	if !exists {
		return nil, fmt.Errorf("field %s not found in %s", field, filename)
	}
	return value, nil
}

// discoverPages finds all page files in the app directory
func (pd *ProjectDiscovery) discoverPages(p *project.Structure) error {
	appDir := filepath.Join(p.Root, "app")

	return filepath.Walk(appDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process index.jml files
		if !info.IsDir() && info.Name() == "index.jml" {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for %s: %w", path, err)
			}

			relPath, err := filepath.Rel(appDir, filepath.Dir(path))
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", path, err)
			}
			relPath = filepath.ToSlash(relPath)
			if relPath == "." {
				relPath = "/"
			} else {
				relPath = "/" + strings.Trim(relPath, "/")
			}

			// Derive page name from directory or use "index" for root
			pageName := "index"
			if relPath != "/" {
				pageName = filepath.Base(filepath.Dir(path))
			}

			// Analyse a page file for additional metadata
			pageInfo, err := pd.analysePageFile(path, p.Root)
			if err != nil {
				return fmt.Errorf("failed to analyse page %s: %w", path, err)
			}

			pageInfo.Name = pageName
			pageInfo.RelativePath = relPath
			pageInfo.AbsolutePath = absPath
			pageInfo.Route = relPath
			pageInfo.LastModified = info.ModTime()

			p.Pages[pageName] = pageInfo
		}

		return nil
	})
}

// discoverComponents finds all component files in the components directory
func (pd *ProjectDiscovery) discoverComponents(p *project.Structure) error {
	componentsDir := filepath.Join(p.Root, "components")

	return filepath.Walk(componentsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".jml") {
			compInfo, err := pd.analyseComponentFile(path, p.Root)
			if err != nil {
				return fmt.Errorf("failed to analyse component %s: %w", path, err)
			}

			p.Components[compInfo.Name] = compInfo
		}

		return nil
	})
}

// discoverAssets finds all asset files in the assets directory
func (pd *ProjectDiscovery) discoverAssets(p *project.Structure) error {
	assetsDir := filepath.Join(p.Root, "assets")

	return filepath.Walk(assetsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(p.Root, path)
			if err != nil {
				return err
			}
			p.Assets = append(p.Assets, relPath)
		}

		return nil
	})
}

// analysePageFile extracts metadata from a page file
func (pd *ProjectDiscovery) analysePageFile(filePath, rootPath string) (*project.PageInfo, error) {
	relPath, err := filepath.Rel(rootPath, filePath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	// Extract name from the filename
	name := strings.TrimSuffix(filepath.Base(filePath), ".jml")

	// Extract title from the first non-empty line
	title, err := pd.extractTitleFromFile(filePath)
	if err != nil {
		title = name // fallback to name if title extraction fails
	}

	// Generate route from directory structure
	route := pd.generateRoute(relPath)

	// Parse imports (placeholder)
	imports, dependencies := pd.parseImports(filePath)

	docInfo := project.DocumentInfo{
		Name:         name,
		Title:        title,
		RelativePath: relPath,
		AbsolutePath: filePath,
		Dependencies: dependencies,
		Imports:      imports,
		LastModified: info.ModTime(),
		Compiled:     false,
	}

	return &project.PageInfo{
		DocumentInfo: docInfo,
		Route:        route,
	}, nil
}

// analyseComponentFile extracts metadata from a component file
func (pd *ProjectDiscovery) analyseComponentFile(filePath, rootPath string) (*project.ComponentInfo, error) {
	relPath, err := filepath.Rel(rootPath, filePath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	// Extract name from the filename
	name := strings.TrimSuffix(filepath.Base(filePath), ".jml")

	// Extract title from the first non-empty line
	title, err := pd.extractTitleFromFile(filePath)
	if err != nil {
		title = name // fallback to name if title extraction fails
	}

	// Parse imports (placeholder)
	imports, dependencies := pd.parseImports(filePath)

	docInfo := project.DocumentInfo{
		Name:         name,
		Title:        title,
		RelativePath: relPath,
		AbsolutePath: filePath,
		Dependencies: dependencies,
		Imports:      imports,
		LastModified: info.ModTime(),
		Compiled:     false,
	}

	return &project.ComponentInfo{
		DocumentInfo: docInfo,
		UsedBy:       make([]string, 0),
	}, nil
}

// extractTitleFromFile reads the first non-empty line to extract title
func (pd *ProjectDiscovery) extractTitleFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "/*") {
			// Look for a doctype declaration or component name
			if strings.HasPrefix(line, "_doctype") {
				parts := strings.Fields(line)
				if len(parts) >= 3 {
					return parts[2], nil // Return the identifier after "page" or "component"
				}
			}
			return line, nil
		}
	}

	return "", scanner.Err()
}

// generateRoute creates a route from the file path
func (pd *ProjectDiscovery) generateRoute(relPath string) string {
	// Convert a file path to a route
	route := filepath.Dir(relPath)
	route = strings.ReplaceAll(route, "\\", "/")
	route = strings.TrimPrefix(route, "app")

	if route == "." || route == "" {
		return "/"
	}

	if !strings.HasPrefix(route, "/") {
		route = "/" + route
	}

	return route
}

// parseImports extracts import statements from the file (placeholder)
func (pd *ProjectDiscovery) parseImports(filePath string) (map[string]string, []string) {
	// TODO: implement actual parsing
	// For now, return empty maps
	return make(map[string]string), make([]string, 0)
}

// buildDependencyGraph creates the dependency relationships
func (pd *ProjectDiscovery) buildDependencyGraph(p *project.Structure) error {
	// Build reverse dependencies for components
	for pageName, pageInfo := range p.Pages {
		for _, dep := range pageInfo.Dependencies {
			if comp, exists := p.Components[dep]; exists {
				comp.UsedBy = append(comp.UsedBy, pageName)
			}
		}
	}

	for compName, compInfo := range p.Components {
		for _, dep := range compInfo.Dependencies {
			if comp, exists := p.Components[dep]; exists {
				comp.UsedBy = append(comp.UsedBy, compName)
			}
		}
	}

	return nil
}
