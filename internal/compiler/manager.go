package compiler

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/yasufadhili/jawt/internal/project"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Manager struct {
	project    *project.Project
	depGraph   *DependencyGraph
	buildCache *Cache
	options    Options
}

type Options struct {
	Verbose bool
}

// Cache tracks file modification times and content hashes
type Cache struct {
	CacheFile string                `json:"-"`
	Files     map[string]FileRecord `json:"files"`
}

// FileRecord stores information about a compiled file
type FileRecord struct {
	LastModified time.Time `json:"last_modified"`
	Hash         string    `json:"hash"`
	OutputPath   string    `json:"output_path"`
	Dependencies []string  `json:"dependencies"`
}

func NewCompilerManager(project *project.Project) *Manager {
	cm := &Manager{
		project:  project,
		depGraph: NewDependencyGraph(),
	}

	cm.buildDependencyGraph()
	cm.initBuildCache()

	return cm
}

// CompileProject compiles the entire project in dependency order
func (cm *Manager) CompileProject() error {
	dir, err := os.MkdirTemp("", "jawt")
	if err != nil {
		return err
	}
	cm.project.OutputDir = dir

	// Initialise cache with a temp directory
	cm.initBuildCache()

	buildOrder, err := cm.depGraph.BuildOrder()
	if err != nil {
		return fmt.Errorf("failed to determine build order: %w", err)
	}
	if cm.options.Verbose {
		err := cm.PrintBuildPlan()
		if err != nil {
			return err
		}
	}

	for _, filePath := range buildOrder {
		if err := cm.compileFile(filePath); err != nil {
			return fmt.Errorf("failed to compile %s: %w", filePath, err)
		}

		// Update cache record after successful compilation
		if err := cm.updateCacheRecord(filePath); err != nil {
			fmt.Printf("Warning: failed to update cache for %s: %v\n", filePath, err)
		}
	}

	if err := cm.copyAssets(); err != nil {
		return fmt.Errorf("asset copying failed: %w", err)
	}

	if err := cm.saveBuildCache(); err != nil {
		fmt.Printf("Warning: failed to save build cache: %v\n", err)
	}

	return nil
}

// CompileChanged compiles only files that have changed and their dependents
func (cm *Manager) CompileChanged() error {
	filesToRecompile, err := cm.getFilesToRecompile()
	if err != nil {
		return fmt.Errorf("failed to determine files to recompile: %w", err)
	}

	if len(filesToRecompile) == 0 {
		fmt.Println("No files have changed, skipping compilation")
		return nil
	}

	fmt.Printf("Recompiling %d files...\n", len(filesToRecompile))

	for _, filePath := range filesToRecompile {
		fmt.Printf("  Compiling: %s\n", filePath)

		if err := cm.compileFile(filePath); err != nil {
			return fmt.Errorf("failed to compile %s: %w", filePath, err)
		}

		if err := cm.updateCacheRecord(filePath); err != nil {
			return fmt.Errorf("failed to update cache for %s: %w", filePath, err)
		}
	}

	if err := cm.saveBuildCache(); err != nil {
		return fmt.Errorf("failed to save build cache: %w", err)
	}

	fmt.Printf("Incremental compilation completed successfully\n")
	return nil
}

// initBuildCache initialises the build cache
func (cm *Manager) initBuildCache() {
	cacheFile := filepath.Join(cm.project.OutputDir, ".jawt_cache.json")
	cm.buildCache = &Cache{
		CacheFile: cacheFile,
		Files:     make(map[string]FileRecord),
	}

	// Try to load the existing cache
	err := cm.loadBuildCache()
	if err != nil {
		fmt.Printf("Warning: failed to load build cache: %v\n", err)
		return
	}
}

// loadBuildCache loads the build cache from disk
func (cm *Manager) loadBuildCache() error {
	data, err := os.ReadFile(cm.buildCache.CacheFile)
	if err != nil {
		// Cache file doesn't exist, start fresh
		return nil
	}

	return json.Unmarshal(data, cm.buildCache)
}

// saveBuildCache saves the build cache to disk
func (cm *Manager) saveBuildCache() error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(cm.buildCache.CacheFile), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cm.buildCache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.buildCache.CacheFile, data, 0644)
}

// hasChanged checks if a file has been modified since the last compilation
func (cm *Manager) hasChanged(filePath string) (bool, error) {
	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		return true, err // File doesn't exist, consider it changed
	}

	// Check if we have a cache record
	record, exists := cm.buildCache.Files[filePath]
	if !exists {
		return true, nil // No cache record, so file is new
	}

	// Check modification time first (quick check)
	if info.ModTime().After(record.LastModified) {
		return true, nil
	}

	// If modification time is the same, check hash for certainty
	currentHash, err := cm.calculateFileHash(filePath)
	if err != nil {
		return true, err
	}

	return currentHash != record.Hash, nil
}

// calculateFileHash calculates SHA256 hash of a file
func (cm *Manager) calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// getChangedFiles returns a list of files that have changed
func (cm *Manager) getChangedFiles() ([]string, error) {
	var changedFiles []string

	for _, comp := range cm.project.Components {
		changed, err := cm.hasChanged(comp.AbsolutePath)
		if err != nil {
			return nil, fmt.Errorf("error checking component %s: %w", comp.RelativePath, err)
		}
		if changed {
			changedFiles = append(changedFiles, comp.AbsolutePath)
		}
	}

	for _, page := range cm.project.Pages {
		changed, err := cm.hasChanged(page.AbsolutePath)
		if err != nil {
			return nil, fmt.Errorf("error checking page %s: %w", page.RelativePath, err)
		}
		if changed {
			changedFiles = append(changedFiles, page.AbsolutePath)
		}
	}

	return changedFiles, nil
}

// getDependentFiles returns files that depend on the given file
func (cm *Manager) getDependentFiles(filePath string) []string {
	var dependents []string

	for path, node := range cm.depGraph.files {
		for _, dep := range node.Dependencies {
			if dep == filePath {
				dependents = append(dependents, path)
				break
			}
		}
	}

	return dependents
}

// getFilesToRecompile returns all files that need recompilation
func (cm *Manager) getFilesToRecompile() ([]string, error) {
	changedFiles, err := cm.getChangedFiles()
	if err != nil {
		return nil, err
	}

	if len(changedFiles) == 0 {
		return []string{}, nil
	}

	// Build a set of files to recompile
	toRecompile := make(map[string]bool)

	for _, file := range changedFiles {
		toRecompile[file] = true
	}

	// Add files that depend on changed files (cascade compilation)
	var queue []string
	queue = append(queue, changedFiles...)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		dependents := cm.getDependentFiles(current)
		for _, dep := range dependents {
			if !toRecompile[dep] {
				toRecompile[dep] = true
				queue = append(queue, dep)
			}
		}
	}

	// Convert the map to slice and sort by dependency order
	var result []string
	buildOrder, err := cm.depGraph.BuildOrder()
	if err != nil {
		return nil, err
	}

	for _, file := range buildOrder {
		if toRecompile[file] {
			result = append(result, file)
		}
	}

	return result, nil
}

// updateCacheRecord updates the cache record for a file
func (cm *Manager) updateCacheRecord(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	hash, err := cm.calculateFileHash(filePath)
	if err != nil {
		return err
	}

	var dependencies []string
	if node, exists := cm.depGraph.files[filePath]; exists {
		dependencies = node.Dependencies
	}

	cm.buildCache.Files[filePath] = FileRecord{
		LastModified: info.ModTime(),
		Hash:         hash,
		OutputPath:   "", // Could be populated with the actual output path
		Dependencies: dependencies,
	}

	return nil
}

// ClearCache removes the build cache
func (cm *Manager) ClearCache() error {
	cm.buildCache.Files = make(map[string]FileRecord)
	return os.Remove(cm.buildCache.CacheFile)
}

// GetCacheStats returns statistics about the build cache
func (cm *Manager) GetCacheStats() CacheStats {
	return CacheStats{
		CachedFiles:  len(cm.buildCache.Files),
		CacheSize:    cm.calculateCacheSize(),
		LastModified: cm.getLastCacheModification(),
	}
}

// CacheStats holds build cache statistics
type CacheStats struct {
	CachedFiles  int
	CacheSize    int64
	LastModified time.Time
}

// calculateCacheSize calculates the size of the cache file
func (cm *Manager) calculateCacheSize() int64 {
	info, err := os.Stat(cm.buildCache.CacheFile)
	if err != nil {
		return 0
	}
	return info.Size()
}

// getLastCacheModification returns when the cache was last modified
func (cm *Manager) getLastCacheModification() time.Time {
	info, err := os.Stat(cm.buildCache.CacheFile)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

// buildDependencyGraph populates the dependency graph from the project structure
func (cm *Manager) buildDependencyGraph() {
	for _, comp := range cm.project.Components {
		cm.depGraph.AddFile(comp.AbsolutePath, Component, comp.Dependencies)
	}

	for _, page := range cm.project.Pages {
		cm.depGraph.AddFile(page.AbsolutePath, Page, page.Dependencies)
	}
}

// compileFile compiles a single file based on its type
func (cm *Manager) compileFile(filePath string) error {
	if comp := cm.findComponentByPath(filePath); comp != nil {
		return cm.compileComponent(comp)
	}

	if page := cm.findPageByPath(filePath); page != nil {
		return cm.compilePage(page)
	}

	return fmt.Errorf("unknown file type for %s", filePath)
}

// findComponentByPath finds a component by its file path
func (cm *Manager) findComponentByPath(filePath string) *project.ComponentInfo {
	for _, comp := range cm.project.Components {
		if comp.AbsolutePath == filePath {
			return comp
		}
	}
	return nil
}

// compileComponent compiles a single component
func (cm *Manager) compileComponent(comp *project.ComponentInfo) error {
	c := NewFileCompiler(cm, &comp.DocumentInfo)

	res, err := c.CompileFile()
	if err != nil {
		return err
	}
	if !res.Success {
		fmt.Printf("Found %d syntax errors\n", len(res.Errors))
		for _, err := range res.Errors {
			_ = fmt.Errorf(err.Error())
		}
		return fmt.Errorf("compilation failed")
	}
	return nil
}

// compilePage compiles a single page
func (cm *Manager) compilePage(page *project.PageInfo) error {
	c := NewFileCompiler(cm, &page.DocumentInfo)

	res, err := c.CompileFile()
	if err != nil {
		return err
	}

	if !res.Success {
		fmt.Printf("Found %d syntax errors\n", len(res.Errors))
		for _, err := range res.Errors {
			_ = fmt.Errorf(err.Error())
		}
		return fmt.Errorf("compilation failed")
	}

	return nil
}

// copyAssets copies all assets to the output directory (placeholder)
func (cm *Manager) copyAssets() error {
	//for _, asset := range cm.project.Assets {
	// fmt.Printf("  📁 Copying asset: %s\n", asset)
	// TODO: actual asset copying
	//}
	return nil
}

// findPageByPath finds a page by its file path
func (cm *Manager) findPageByPath(filePath string) *project.PageInfo {
	for _, page := range cm.project.Pages {
		if page.AbsolutePath == filePath {
			return page
		}
	}
	return nil
}

// ValidateDependencies checks for dependency issues
func (cm *Manager) ValidateDependencies() error {
	// Check for cycles
	if cycles := cm.depGraph.FindCycles(); len(cycles) > 0 {
		return fmt.Errorf("circular dependencies detected: %v", cycles)
	}

	// Check for missing dependencies
	for filePath, node := range cm.depGraph.files {
		for _, dep := range node.Dependencies {
			if _, exists := cm.depGraph.files[dep]; !exists {
				return fmt.Errorf("missing dependency %s for file %s", dep, filePath)
			}
		}
	}

	return nil
}

// GetBuildOrder returns the order files should be compiled in
func (cm *Manager) GetBuildOrder() ([]string, error) {
	return cm.depGraph.BuildOrder()
}

// PrintBuildPlan outputs the planned build order
func (cm *Manager) PrintBuildPlan() error {
	buildOrder, err := cm.GetBuildOrder()
	if err != nil {
		return err
	}

	fmt.Println("Build Plan:")
	for i, filePath := range buildOrder {
		node := cm.depGraph.files[filePath]
		fmt.Printf("%d. %s (%s)\n", i+1, filePath, node.Type)
		if len(node.Dependencies) > 0 {
			fmt.Printf("   Dependencies: %v\n", node.Dependencies)
		}
	}

	return nil
}

// GetCompilationStats returns statistics about the compilation
func (cm *Manager) GetCompilationStats() CompilationStats {
	stats := CompilationStats{
		TotalComponents: len(cm.project.Components),
		TotalPages:      len(cm.project.Pages),
		TotalAssets:     len(cm.project.Assets),
	}

	for _, comp := range cm.project.Components {
		if comp.IsCompiled() {
			stats.CompiledComponents++
		}
	}

	for _, page := range cm.project.Pages {
		if page.IsCompiled() {
			stats.CompiledPages++
		}
	}

	return stats
}

// CompilationStats holds compilation statistics
type CompilationStats struct {
	TotalComponents    int
	CompiledComponents int
	TotalPages         int
	CompiledPages      int
	TotalAssets        int
}

// String returns a formatted string of compilation stats
func (cs CompilationStats) String() string {
	return fmt.Sprintf(
		"Components: %d/%d, Pages: %d/%d, Assets: %d",
		cs.CompiledComponents, cs.TotalComponents,
		cs.CompiledPages, cs.TotalPages,
		cs.TotalAssets,
	)
}
