package core

import (
	"os"
	"path/filepath"
)

// ProjectPaths manages all file paths for the project
type ProjectPaths struct {
	// Base directories
	ProjectRoot string
	WorkingDir  string

	// Input directories
	AppDir        string
	ComponentsDir string
	ScriptsDir    string
	AssetsDir     string

	// Output directories
	JawtDir  string
	BuildDir string
	DistDir  string
	TempDir  string
	CacheDir string

	// Generated files
	TypeScriptOutputDir string
	TailwindOutputDir   string
	ComponentsOutputDir string

	// Config files
	TSConfigPath       string
	TailwindConfigPath string
	ProjectConfigPath  string
}

// NewProjectPaths creates a new ProjectPaths instance
func NewProjectPaths(projectRoot string, projectConfig *ProjectConfig, jawtConfig *JawtConfig) (*ProjectPaths, error) {
	absProjectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, err
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	paths := &ProjectPaths{
		ProjectRoot: absProjectRoot,
		WorkingDir:  workingDir,
	}

	// Set up input directories
	paths.AppDir = filepath.Join(absProjectRoot, "app")
	paths.ComponentsDir = filepath.Join(absProjectRoot, "components")
	paths.ScriptsDir = filepath.Join(absProjectRoot, "scripts")
	paths.AssetsDir = filepath.Join(absProjectRoot, "assets")

	// Set up output directories based on config
	paths.JawtDir = filepath.Join(absProjectRoot, ".jawt")
	paths.BuildDir = filepath.Join(paths.JawtDir, "build")
	paths.DistDir = filepath.Join(paths.JawtDir, "dist")
	paths.TempDir = filepath.Join(paths.JawtDir, "temp")
	paths.CacheDir = filepath.Join(paths.JawtDir, "cache")

	// Set up generated output directories
	paths.TypeScriptOutputDir = filepath.Join(paths.BuildDir, "ts")
	paths.TailwindOutputDir = filepath.Join(paths.BuildDir, "styles")
	paths.ComponentsOutputDir = filepath.Join(paths.BuildDir, "components")

	// Set up config file paths
	paths.TSConfigPath = filepath.Join(absProjectRoot, projectConfig.TSConfigPath)
	paths.TailwindConfigPath = filepath.Join(absProjectRoot, projectConfig.TailwindConfigPath)
	paths.ProjectConfigPath = filepath.Join(absProjectRoot, "jawt.project.json")

	return paths, nil
}

// EnsureDirectories creates all necessary directories
func (p *ProjectPaths) EnsureDirectories() error {
	dirs := []string{
		p.JawtDir,
		p.BuildDir,
		p.DistDir,
		p.TempDir,
		p.CacheDir,
		p.TypeScriptOutputDir,
		p.TailwindOutputDir,
		p.ComponentsOutputDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// GetRelativePath returns a path relative to the project root
func (p *ProjectPaths) GetRelativePath(path string) string {
	rel, err := filepath.Rel(p.ProjectRoot, path)
	if err != nil {
		return path
	}
	return rel
}

// GetAbsolutePath returns an absolute path from a relative path
func (p *ProjectPaths) GetAbsolutePath(relativePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	return filepath.Join(p.ProjectRoot, relativePath)
}

// GetJMLFiles returns all JML files in the project
func (p *ProjectPaths) GetJMLFiles() ([]string, error) {
	var jmlFiles []string

	// Check app directory
	appFiles, err := p.findFilesWithExtension(p.AppDir, ".jml")
	if err == nil {
		jmlFiles = append(jmlFiles, appFiles...)
	}

	// Check components directory
	componentFiles, err := p.findFilesWithExtension(p.ComponentsDir, ".jml")
	if err == nil {
		jmlFiles = append(jmlFiles, componentFiles...)
	}

	return jmlFiles, nil
}

// GetTypeScriptFiles returns all TypeScript files in the project
func (p *ProjectPaths) GetTypeScriptFiles() ([]string, error) {
	var tsFiles []string

	// Check scripts directory
	scriptFiles, err := p.findFilesWithExtension(p.ScriptsDir, ".ts")
	if err == nil {
		tsFiles = append(tsFiles, scriptFiles...)
	}

	// Also check for .tsx files
	tsxFiles, err := p.findFilesWithExtension(p.ScriptsDir, ".tsx")
	if err == nil {
		tsFiles = append(tsFiles, tsxFiles...)
	}

	return tsFiles, nil
}

// GetWatchPaths returns all paths that should be watched for changes
func (p *ProjectPaths) GetWatchPaths() []string {
	return []string{
		p.AppDir,
		p.ComponentsDir,
		p.ScriptsDir,
		p.AssetsDir,
		p.ProjectConfigPath,
		p.TSConfigPath,
		p.TailwindConfigPath,
	}
}

// GetTempFile returns a temporary file path
func (p *ProjectPaths) GetTempFile(filename string) string {
	return filepath.Join(p.TempDir, filename)
}

// GetCacheFile returns a cache file path
func (p *ProjectPaths) GetCacheFile(filename string) string {
	return filepath.Join(p.CacheDir, filename)
}

// findFilesWithExtension recursively finds files with the given extension
func (p *ProjectPaths) findFilesWithExtension(dir, ext string) ([]string, error) {
	var files []string

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return files, nil
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ext {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// Clean removes all generated directories
func (p *ProjectPaths) Clean() error {
	dirs := []string{
		p.JawtDir,
		p.BuildDir,
		p.DistDir,
	}

	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}

	return nil
}
