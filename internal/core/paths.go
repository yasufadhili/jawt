package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ProjectPaths manages all file paths for the project
type ProjectPaths struct {
	// Base directories
	ProjectRoot string
	WorkingDir  string

	// Configuration references
	ProjectConfig *ProjectConfig
	JawtConfig    *JawtConfig

	// --- User-Facing Directories ---
	AppDir        string
	ComponentsDir string
	ScriptsDir    string
	AssetsDir     string
	DistDir       string // Final output for the user

	// --- Managed Workspace (.jawt) ---
	JawtDir         string // Root .jawt directory
	BuildDir        string // Intermediate build artifacts (.jawt/build)
	CacheDir        string // General-purpose cache (.jawt/cache)
	SrcDir          string // "Virtual" source root for compilers (.jawt/src)
	UserSrcDir      string // User code copied to the workspace (.jawt/src/user)
	InternalSrcDir  string // Jawt's internal bundled code (.jawt/src/internal)
	NodeModulesDir  string // Managed node_modules (.jawt/node_modules)
	ToolsDir        string // Managed tools like node/tsc (.jawt/tools)
	GeneratedDir    string // For generated code like routes, manifests (.jawt/generated)
	TailwindCSSPath string // Path to the output tailwind.css file

	// --- Key File Paths ---
	ProjectConfigPath  string
	TSConfigPath       string // Path to the generated tsconfig.json in .jawt
	TailwindConfigPath string // Path to the generated tailwind.config.js in .jawt
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
		ProjectRoot:   absProjectRoot,
		WorkingDir:    workingDir,
		ProjectConfig: projectConfig,
		JawtConfig:    jawtConfig,
	}

	// --- User-Facing Directories ---
	paths.AppDir = filepath.Join(absProjectRoot, projectConfig.Paths.Pages)
	paths.ComponentsDir = filepath.Join(absProjectRoot, projectConfig.Paths.Components)
	paths.ScriptsDir = filepath.Join(absProjectRoot, projectConfig.Paths.Scripts)
	paths.AssetsDir = filepath.Join(absProjectRoot, projectConfig.Paths.Assets)
	paths.DistDir = projectConfig.GetDistDir(absProjectRoot)

	// --- Managed Workspace (.jawt) ---
	paths.JawtDir = filepath.Join(absProjectRoot, ".jawt")
	paths.BuildDir = filepath.Join(paths.JawtDir, "build")
	paths.CacheDir = filepath.Join(paths.JawtDir, "cache")
	paths.SrcDir = filepath.Join(paths.JawtDir, "src")
	paths.UserSrcDir = filepath.Join(paths.SrcDir, "user")
	paths.InternalSrcDir = filepath.Join(paths.SrcDir, "internal")
	paths.NodeModulesDir = filepath.Join(paths.JawtDir, "node_modules")
	paths.ToolsDir = filepath.Join(paths.JawtDir, "tools")
	paths.GeneratedDir = filepath.Join(paths.JawtDir, "generated")
	paths.TailwindCSSPath = filepath.Join(paths.BuildDir, "tailwind.css")

	// --- Key File Paths ---
	paths.ProjectConfigPath = filepath.Join(absProjectRoot, "jawt.project.json")
	paths.TSConfigPath = filepath.Join(paths.JawtDir, "jawt.tsconfig.json")
	paths.TailwindConfigPath = filepath.Join(paths.JawtDir, "tailwind.config.js")

	return paths, nil
}

// ResolveExecutablePath resolves the absolute path to an executable.
// It checks:
// 1. If the provided path is already an absolute executable path.
// 2. If the command exists in the system's PATH.
// 3. If the command exists relative to the current JAWT executable.
func ResolveExecutablePath(cmd string) (string, error) {
	// 1. Check if the provided path is already an absolute executable path
	if filepath.IsAbs(cmd) {
		if _, err := os.Stat(cmd); err == nil {
			return cmd, nil
		}
	}

	// 2. Check if the command exists in the system's PATH
	if path, err := exec.LookPath(cmd); err == nil {
		return path, nil
	}

	// 3. Check if the command exists relative to the current JAWT executable
	jawtExec, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get current executable path: %w", err)
	}
	jawtDir := filepath.Dir(jawtExec)
	localPath := filepath.Join(jawtDir, cmd)
	if runtime.GOOS == "windows" {
		localPath += ".exe"
	}
	if _, err := os.Stat(localPath); err == nil {
		return localPath, nil
	}

	return "", fmt.Errorf("executable '%s' not found in PATH or relative to JAWT executable", cmd)
}

// EnsureDirectories creates all necessary directories
func (p *ProjectPaths) EnsureDirectories() error {
	// Create all necessary directories
	// User-facing directories are not created automatically, except for DistDir
	dirs := []string{
		p.DistDir,
		p.JawtDir,
		p.BuildDir,
		p.CacheDir,
		p.SrcDir,
		p.UserSrcDir,
		p.InternalSrcDir,
		p.NodeModulesDir,
		p.ToolsDir,
		p.GeneratedDir,
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
	watchPaths := make([]string, len(p.ProjectConfig.Dev.WatchPaths))
	for i, path := range p.ProjectConfig.Dev.WatchPaths {
		watchPaths[i] = p.GetAbsolutePath(path)
	}

	// Always watch config files
	watchPaths = append(watchPaths,
		p.ProjectConfigPath,
		p.TSConfigPath,
		p.TailwindConfigPath,
	)

	return watchPaths
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
