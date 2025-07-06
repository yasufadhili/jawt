package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewProjectPaths(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	// Create a test directory structure
	err = paths.EnsureDirectories()
	if err != nil {
		t.Fatalf("failed to ensure directories: %v", err)
	}

	// Create some test TypeScript files
	testFiles := []string{
		filepath.Join(paths.ScriptsDir, "index.ts"),
		filepath.Join(paths.ScriptsDir, "utils.ts"),
		filepath.Join(paths.ScriptsDir, "components", "Button.tsx"),
		filepath.Join(paths.ScriptsDir, "services", "api.ts"), filepath.Join(paths.ScriptsDir, "not-typescript.js"), // Should be ignored
	}

	for _, file := range testFiles {
		err := os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			t.Fatalf("failed to create directory for %s: %v", file, err)
		}

		err = os.WriteFile(file, []byte("// test content"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
	}

	tsFiles, err := paths.GetTypeScriptFiles()
	if err != nil {
		t.Fatalf("failed to get TypeScript files: %v", err)
	}

	// Should find 4 TypeScript files (not the .js file)
	expectedCount := 4
	if len(tsFiles) != expectedCount {
		t.Errorf("expected %d TypeScript files, got %d", expectedCount, len(tsFiles))
	}

	// Check that all returned files have .ts or .tsx extension
	for _, file := range tsFiles {
		ext := filepath.Ext(file)
		if ext != ".ts" && ext != ".tsx" {
			t.Errorf("expected .ts or .tsx extension, got %s", ext)
		}
	}
}

func TestProjectPathsGetWatchPaths(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	watchPaths := paths.GetWatchPaths()

	expectedPaths := []string{
		paths.GetAbsolutePath(projectConfig.Paths.Pages),
		paths.GetAbsolutePath(projectConfig.Paths.Components),
		paths.GetAbsolutePath(projectConfig.Paths.Scripts),
		paths.GetAbsolutePath(projectConfig.Paths.Assets),
		paths.ProjectConfigPath,
		paths.TSConfigPath,
		paths.TailwindConfigPath,
	}

	if len(watchPaths) != len(expectedPaths) {
		t.Errorf("expected %d watch paths, got %d", len(expectedPaths), len(watchPaths))
	}

	for i, expectedPath := range expectedPaths {
		if watchPaths[i] != expectedPath {
			t.Errorf("expected watch path %s, got %s", expectedPath, watchPaths[i])
		}
	}
}

func TestProjectPathsGetTempFile(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	filename := "test-temp-file.txt"
	tempFile := paths.GetTempFile(filename)

	expectedPath := filepath.Join(paths.ProjectRoot, ".jawt", "tmp", filename)
	if tempFile != expectedPath {
		t.Errorf("expected temp file path %s, got %s", expectedPath, tempFile)
	}
}

func TestProjectPathsGetCacheFile(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	filename := "test-cache-file.json"
	cacheFile := paths.GetCacheFile(filename)

	expectedPath := filepath.Join(paths.ProjectRoot, ".jawt", "cache", filename)
	if cacheFile != expectedPath {
		t.Errorf("expected cache file path %s, got %s", expectedPath, cacheFile)
	}
}

func TestProjectPathsFindFilesWithExtension(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	// Create test directory structure
	testDir := filepath.Join(tempDir, "test-files")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	// Create test files
	testFiles := []string{
		filepath.Join(testDir, "file1.txt"),
		filepath.Join(testDir, "file2.txt"),
		filepath.Join(testDir, "subdir", "file3.txt"),
		filepath.Join(testDir, "file4.md"),
		filepath.Join(testDir, "file5.txt"),
	}

	for _, file := range testFiles {
		err := os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			t.Fatalf("failed to create directory for %s: %v", file, err)
		}

		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
	}

	// Test finding .txt files
	txtFiles, err := paths.findFilesWithExtension(testDir, ".txt")
	if err != nil {
		t.Fatalf("failed to find .txt files: %v", err)
	}

	expectedTxtCount := 4
	if len(txtFiles) != expectedTxtCount {
		t.Errorf("expected %d .txt files, got %d", expectedTxtCount, len(txtFiles))
	}

	// Test finding .md files
	mdFiles, err := paths.findFilesWithExtension(testDir, ".md")
	if err != nil {
		t.Fatalf("failed to find .md files: %v", err)
	}

	expectedMdCount := 1
	if len(mdFiles) != expectedMdCount {
		t.Errorf("expected %d .md files, got %d", expectedMdCount, len(mdFiles))
	}

	// Test non-existent directory
	nonExistentDir := filepath.Join(tempDir, "non-existent")
	files, err := paths.findFilesWithExtension(nonExistentDir, ".txt")
	if err != nil {
		t.Errorf("expected no error for non-existent directory, got %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files for non-existent directory, got %d", len(files))
	}
}

func TestProjectPathsClean(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	// Create directories
	err = paths.EnsureDirectories()
	if err != nil {
		t.Fatalf("failed to ensure directories: %v", err)
	}

	// Create some test files in the directories
	testFiles := []string{
		filepath.Join(paths.BuildDir, "test.js"),
		filepath.Join(paths.DistDir, "bundle.js"),
		filepath.Join(paths.JawtDir, "metadata.json"),
	}

	for _, file := range testFiles {
		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
	}

	// Verify directories and files exist
	dirsToCheck := []string{
		paths.JawtDir,
		paths.BuildDir,
		paths.DistDir,
	}

	for _, dir := range dirsToCheck {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("directory %s should exist before clean", dir)
		}
	}

	// Clean the directories
	err = paths.Clean()
	if err != nil {
		t.Fatalf("failed to clean directories: %v", err)
	}

	// Verify directories are removed
	for _, dir := range dirsToCheck {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Errorf("directory %s should be removed after clean", dir)
		}
	}
}

func TestProjectPathsEdgeCases(t *testing.T) {
	// Test with custom project config
	tempDir := t.TempDir()

	customProjectConfig := &ProjectConfig{
		App: struct {
			Name        string `json:"name"`
			Author      string `json:"author"`
			Version     string `json:"version"`
			Description string `json:"description"`
		}{
			Name:        "custom-project",
			Version:     "2.0.0",
			Description: "Custom project for testing",
		},
		Paths: struct {
			Components string `json:"components"`
			Pages      string `json:"pages"`
			Scripts    string `json:"scripts"`
			Assets     string `json:"assets"`
		}{
			Components: "custom-components",
			Pages:      "custom-app",
			Scripts:    "custom-scripts",
			Assets:     "custom-assets",
		},
		Server: struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}{
			Host: "localhost",
			Port: 6500,
		},
		Build: struct {
			OutputDir string `json:"outputDir"`
			DistDir   string `json:"distDir"`
			Minify    bool   `json:"minify"`
			ShadowDOM bool   `json:"shadowDOM"`
		}{
			OutputDir: "custom-build",
			DistDir:   "custom-dist",
			Minify:    true,
			ShadowDOM: false,
		},
		Dev: struct {
			Port       int      `json:"port"`
			EnableHMR  bool     `json:"enableHMR"`
			WatchPaths []string `json:"watchPaths"`
		}{
			Port:       9000,
			EnableHMR:  false,
			WatchPaths: []string{"src", "lib", "assets"},
		},
		Tooling: struct {
			TSConfigPath       string `json:"tsConfigPath"`
			TailwindConfigPath string `json:"tailwindConfigPath"`
		}{
			TSConfigPath:       "custom-tsconfig.json",
			TailwindConfigPath: "custom-tailwind.config.js",
		},
		Scripts: struct {
			PreBuild  []string `json:"preBuild"`
			PostBuild []string `json:"postBuild"`
		}{
			PreBuild:  []string{"prebuild.sh"},
			PostBuild: []string{"postbuild.sh"},
		},
	}

	customJawtConfig := &JawtConfig{
		TypeScriptPath:     "custom-tsc",
		TailwindPath:       "custom-tailwind",
		NodePath:           "custom-node",
		DefaultPort:        7000,
		TempDir:            "custom-tmp",
		CacheDir:           "custom-cache",
		EnableMinification: false,
		EnableSourceMaps:   false,
		EnableTreeShaking:  false,
	}

	paths, err := NewProjectPaths(tempDir, customProjectConfig, customJawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths with custom config: %v", err)
	}

	// Test that custom paths are used
	expectedBuildDir := filepath.Join(paths.ProjectRoot, customProjectConfig.Build.OutputDir)
	if paths.BuildDir != expectedBuildDir {
		t.Errorf("expected custom BuildDir %s, got %s", expectedBuildDir, paths.BuildDir)
	}

	expectedDistDir := filepath.Join(paths.ProjectRoot, customProjectConfig.Build.DistDir)
	if paths.DistDir != expectedDistDir {
		t.Errorf("expected custom DistDir %s, got %s", expectedDistDir, paths.DistDir)
	}

	// TempDir and CacheDir are now fixed to .jawt/tmp and .jawt/cache
	expectedTempDir := filepath.Join(paths.ProjectRoot, ".jawt", "tmp")
	if paths.TempDir != expectedTempDir {
		t.Errorf("expected custom TempDir %s, got %s", expectedTempDir, paths.TempDir)
	}

	expectedCacheDir := filepath.Join(paths.ProjectRoot, ".jawt", "cache")
	if paths.CacheDir != expectedCacheDir {
		t.Errorf("expected custom CacheDir %s, got %s", expectedCacheDir, paths.CacheDir)
	}

	// Test input directories
	expectedAppDir := filepath.Join(paths.ProjectRoot, customProjectConfig.Paths.Pages)
	if paths.AppDir != expectedAppDir {
		t.Errorf("expected AppDir %s, got %s", expectedAppDir, paths.AppDir)
	}

	expectedComponentsDir := filepath.Join(paths.ProjectRoot, customProjectConfig.Paths.Components)
	if paths.ComponentsDir != expectedComponentsDir {
		t.Errorf("expected ComponentsDir %s, got %s", expectedComponentsDir, paths.ComponentsDir)
	}

	expectedScriptsDir := filepath.Join(paths.ProjectRoot, customProjectConfig.Paths.Scripts)
	if paths.ScriptsDir != expectedScriptsDir {
		t.Errorf("expected ScriptsDir %s, got %s", expectedScriptsDir, paths.ScriptsDir)
	}

	expectedAssetsDir := filepath.Join(paths.ProjectRoot, customProjectConfig.Paths.Assets)
	if paths.AssetsDir != expectedAssetsDir {
		t.Errorf("expected AssetsDir %s, got %s", expectedAssetsDir, paths.AssetsDir)
	}

	// Test generated output directories
	expectedTSOutputDir := filepath.Join(paths.BuildDir, "ts")
	if paths.TypeScriptOutputDir != expectedTSOutputDir {
		t.Errorf("expected TypeScriptOutputDir %s, got %s", expectedTSOutputDir, paths.TypeScriptOutputDir)
	}

	expectedTailwindOutputDir := filepath.Join(paths.BuildDir, "styles")
	if paths.TailwindOutputDir != expectedTailwindOutputDir {
		t.Errorf("expected TailwindOutputDir %s, got %s", expectedTailwindOutputDir, paths.TailwindOutputDir)
	}

	expectedComponentsOutputDir := filepath.Join(paths.BuildDir, "components")
	if paths.ComponentsOutputDir != expectedComponentsOutputDir {
		t.Errorf("expected ComponentsOutputDir %s, got %s", expectedComponentsOutputDir, paths.ComponentsOutputDir)
	}

	// Test config file paths
	expectedTSConfigPath := filepath.Join(paths.ProjectRoot, customProjectConfig.Tooling.TSConfigPath)
	if paths.TSConfigPath != expectedTSConfigPath {
		t.Errorf("expected TSConfigPath %s, got %s", expectedTSConfigPath, paths.TSConfigPath)
	}

	expectedTailwindConfigPath := filepath.Join(paths.ProjectRoot, customProjectConfig.Tooling.TailwindConfigPath)
	if paths.TailwindConfigPath != expectedTailwindConfigPath {
		t.Errorf("expected TailwindConfigPath %s, got %s", expectedTailwindConfigPath, paths.TailwindConfigPath)
	}

	expectedProjectConfigPath := filepath.Join(paths.ProjectRoot, "jawt.project.json")
	if paths.ProjectConfigPath != expectedProjectConfigPath {
		t.Errorf("expected ProjectConfigPath %s, got %s", expectedProjectConfigPath, paths.ProjectConfigPath)
	}
}

func TestProjectPathsWithSymlinks(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	// Create a symlink to the temp directory
	symlinkPath := filepath.Join(tempDir, "symlink-project")
	targetPath := filepath.Join(tempDir, "actual-project")

	err := os.MkdirAll(targetPath, 0755)
	if err != nil {
		t.Fatalf("failed to create target directory: %v", err)
	}

	// Create symlink (skip test if symlinks are not supported)
	err = os.Symlink(targetPath, symlinkPath)
	if err != nil {
		t.Skipf("symlinks not supported on this platform: %v", err)
	}

	paths, err := NewProjectPaths(symlinkPath, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths with symlink: %v", err)
	}

	// Ensure the project root resolves to absolute path
	if !filepath.IsAbs(paths.ProjectRoot) {
		t.Errorf("expected absolute project root, got %s", paths.ProjectRoot)
	}

	// Test that the absolute path is resolved correctly
	expectedAbsPath, _ := filepath.Abs(targetPath)

	if paths.ProjectRoot != expectedAbsPath {
		t.Errorf("expected ProjectRoot %s, got %s", expectedAbsPath, paths.ProjectRoot)
	}

	// Test input directories
	expectedAppDir := filepath.Join(expectedAbsPath, projectConfig.Paths.Pages)
	if paths.AppDir != expectedAppDir {
		t.Errorf("expected AppDir %s, got %s", expectedAppDir, paths.AppDir)
	}

	expectedComponentsDir := filepath.Join(expectedAbsPath, projectConfig.Paths.Components)
	if paths.ComponentsDir != expectedComponentsDir {
		t.Errorf("expected ComponentsDir %s, got %s", expectedComponentsDir, paths.ComponentsDir)
	}

	expectedScriptsDir := filepath.Join(expectedAbsPath, projectConfig.Paths.Scripts)
	if paths.ScriptsDir != expectedScriptsDir {
		t.Errorf("expected ScriptsDir %s, got %s", expectedScriptsDir, paths.ScriptsDir)
	}

	expectedAssetsDir := filepath.Join(expectedAbsPath, projectConfig.Paths.Assets)
	if paths.AssetsDir != expectedAssetsDir {
		t.Errorf("expected AssetsDir %s, got %s", expectedAssetsDir, paths.AssetsDir)
	}

	// Test output directories
	expectedBuildDir := filepath.Join(expectedAbsPath, projectConfig.Build.OutputDir)
	if paths.BuildDir != expectedBuildDir {
		t.Errorf("expected BuildDir %s, got %s", expectedBuildDir, paths.BuildDir)
	}

	expectedDistDir := filepath.Join(expectedAbsPath, projectConfig.Build.DistDir)
	if paths.DistDir != expectedDistDir {
		t.Errorf("expected DistDir %s, got %s", expectedDistDir, paths.DistDir)
	}

	expectedTempDir := filepath.Join(expectedAbsPath, ".jawt", "tmp")
	if paths.TempDir != expectedTempDir {
		t.Errorf("expected TempDir %s, got %s", expectedTempDir, paths.TempDir)
	}

	expectedCacheDir := filepath.Join(expectedAbsPath, ".jawt", "cache")
	if paths.CacheDir != expectedCacheDir {
		t.Errorf("expected CacheDir %s, got %s", expectedCacheDir, paths.CacheDir)
	}

	// Test generated output directories
	expectedTSOutputDir := filepath.Join(paths.BuildDir, "ts")
	if paths.TypeScriptOutputDir != expectedTSOutputDir {
		t.Errorf("expected TypeScriptOutputDir %s, got %s", expectedTSOutputDir, paths.TypeScriptOutputDir)
	}

	expectedTailwindOutputDir := filepath.Join(paths.BuildDir, "styles")
	if paths.TailwindOutputDir != expectedTailwindOutputDir {
		t.Errorf("expected TailwindOutputDir %s, got %s", expectedTailwindOutputDir, paths.TailwindOutputDir)
	}

	expectedComponentsOutputDir := filepath.Join(paths.BuildDir, "components")
	if paths.ComponentsOutputDir != expectedComponentsOutputDir {
		t.Errorf("expected ComponentsOutputDir %s, got %s", expectedComponentsOutputDir, paths.ComponentsOutputDir)
	}

	// Test config file paths
	expectedTSConfigPath := filepath.Join(expectedAbsPath, projectConfig.Tooling.TSConfigPath)
	if paths.TSConfigPath != expectedTSConfigPath {
		t.Errorf("expected TSConfigPath %s, got %s", expectedTSConfigPath, paths.TSConfigPath)
	}

	expectedTailwindConfigPath := filepath.Join(expectedAbsPath, projectConfig.Tooling.TailwindConfigPath)
	if paths.TailwindConfigPath != expectedTailwindConfigPath {
		t.Errorf("expected TailwindConfigPath %s, got %s", expectedTailwindConfigPath, paths.TailwindConfigPath)
	}

	expectedProjectConfigPath := filepath.Join(expectedAbsPath, "jawt.project.json")
	if paths.ProjectConfigPath != expectedProjectConfigPath {
		t.Errorf("expected ProjectConfigPath %s, got %s", expectedProjectConfigPath, paths.ProjectConfigPath)
	}
}

func TestProjectPathsEnsureDirectories(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	err = paths.EnsureDirectories()
	if err != nil {
		t.Fatalf("failed to ensure directories: %v", err)
	}

	// Check that all directories were created
	dirsToCheck := []string{
		paths.JawtDir,
		paths.BuildDir,
		paths.DistDir,
		paths.TempDir,
		paths.CacheDir,
		paths.TypeScriptOutputDir,
		paths.TailwindOutputDir,
		paths.ComponentsOutputDir,
	}

	for _, dir := range dirsToCheck {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("directory %s was not created", dir)
		}
	}
}

func TestProjectPathsGetRelativePath(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	tests := []struct {
		name         string
		absolutePath string
		expected     string
	}{
		{
			name:         "app directory",
			absolutePath: paths.AppDir,
			expected:     projectConfig.Paths.Pages,
		},
		{
			name:         "components directory",
			absolutePath: paths.ComponentsDir,
			expected:     projectConfig.Paths.Components,
		},
		{
			name:         "build directory",
			absolutePath: paths.BuildDir,
			expected:     projectConfig.Build.OutputDir,
		},
		{
			name:         "file in app directory",
			absolutePath: filepath.Join(paths.AppDir, "index.jml"),
			expected:     filepath.Join(projectConfig.Paths.Pages, "index.jml"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := paths.GetRelativePath(tt.absolutePath)
			if result != tt.expected {
				t.Errorf("expected relative path %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestProjectPathsGetAbsolutePath(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	tests := []struct {
		name         string
		relativePath string
		expected     string
	}{
		{
			name:         "app directory",
			relativePath: projectConfig.Paths.Pages,
			expected:     paths.AppDir,
		},
		{
			name:         "components directory",
			relativePath: projectConfig.Paths.Components,
			expected:     paths.ComponentsDir,
		},
		{
			name:         "file in app directory",
			relativePath: filepath.Join(projectConfig.Paths.Pages, "index.jml"),
			expected:     filepath.Join(paths.AppDir, "index.jml"),
		},
		{
			name:         "already absolute path",
			relativePath: paths.AppDir,
			expected:     paths.AppDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := paths.GetAbsolutePath(tt.relativePath)
			if result != tt.expected {
				t.Errorf("expected absolute path %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestProjectPathsGetJMLFiles(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	// Create test directory structure
	err = paths.EnsureDirectories()
	if err != nil {
		t.Fatalf("failed to ensure directories: %v", err)
	}

	// Create some test JML files
	testFiles := []string{
		filepath.Join(paths.AppDir, "index.jml"),
		filepath.Join(paths.AppDir, "about.jml"),
		filepath.Join(paths.AppDir, "subdir", "nested.jml"),
		filepath.Join(paths.ComponentsDir, "header.jml"),
		filepath.Join(paths.ComponentsDir, "footer.jml"),
		filepath.Join(paths.AppDir, "not-jml.txt"), // Should be ignored
	}

	for _, file := range testFiles {
		err := os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			t.Fatalf("failed to create directory for %s: %v", file, err)
		}

		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
	}

	jmlFiles, err := paths.GetJMLFiles()
	if err != nil {
		t.Fatalf("failed to get JML files: %v", err)
	}

	// Should find 5 JML files (not the .txt file)
	expectedCount := 5
	if len(jmlFiles) != expectedCount {
		t.Errorf("expected %d JML files, got %d", expectedCount, len(jmlFiles))
	}

	// Check that all returned files have .jml extension
	for _, file := range jmlFiles {
		if filepath.Ext(file) != ".jml" {
			t.Errorf("expected .jml extension, got %s", filepath.Ext(file))
		}
	}
}

func TestProjectPathsGetTypeScriptFiles(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	// Create test directory structure
	err = paths.EnsureDirectories()
	if err != nil {
		t.Fatalf("failed to ensure directories: %v", err)
	}

	// Create some test TypeScript files
	testFiles := []string{
		filepath.Join(paths.ScriptsDir, "index.ts"),
		filepath.Join(paths.ScriptsDir, "utils.ts"),
		filepath.Join(paths.ScriptsDir, "components", "Button.tsx"),
		filepath.Join(paths.ScriptsDir, "services", "api.ts"),
		filepath.Join(paths.ScriptsDir, "not-typescript.js"), // Should be ignored
	}

	for _, file := range testFiles {
		err := os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			t.Fatalf("failed to create directory for %s: %v", file, err)
		}

		err = os.WriteFile(file, []byte("// test content"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
	}

	tsFiles, err := paths.GetTypeScriptFiles()
	if err != nil {
		t.Fatalf("failed to get TypeScript files: %v", err)
	}

	// Should find 4 TypeScript files (not the .js file)
	expectedCount := 4
	if len(tsFiles) != expectedCount {
		t.Errorf("expected %d TypeScript files, got %d", expectedCount, len(tsFiles))
	}

	// Check that all returned files have .ts or .tsx extension
	for _, file := range tsFiles {
		ext := filepath.Ext(file)
		if ext != ".ts" && ext != ".tsx" {
			t.Errorf("expected .ts or .tsx extension, got %s", ext)
		}
	}
}

func TestProjectPathsGetJMLFiles2(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	// Create test directory structure
	err = paths.EnsureDirectories()
	if err != nil {
		t.Fatalf("failed to ensure directories: %v", err)
	}

	// Create some test JML files
	testFiles := []string{
		filepath.Join(paths.AppDir, "index.jml"),
		filepath.Join(paths.AppDir, "about.jml"),
		filepath.Join(paths.AppDir, "subdir", "nested.jml"),
		filepath.Join(paths.ComponentsDir, "header.jml"),
		filepath.Join(paths.ComponentsDir, "footer.jml"),
		filepath.Join(paths.AppDir, "not-jml.txt"), // Should be ignored
	}

	for _, file := range testFiles {
		err := os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			t.Fatalf("failed to create directory for %s: %v", file, err)
		}

		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
	}

	jmlFiles, err := paths.GetJMLFiles()
	if err != nil {
		t.Fatalf("failed to get JML files: %v", err)
	}

	// Should find 5 JML files (not the .txt file)
	expectedCount := 5
	if len(jmlFiles) != expectedCount {
		t.Errorf("expected %d JML files, got %d", expectedCount, len(jmlFiles))
	}

	// Check that all returned files have .jml extension
	for _, file := range jmlFiles {
		if filepath.Ext(file) != ".jml" {
			t.Errorf("expected .jml extension, got %s", filepath.Ext(file))
		}
	}
}

func TestProjectPathsGetWatchPaths2(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	watchPaths := paths.GetWatchPaths()

	expectedPaths := []string{
		paths.GetAbsolutePath(projectConfig.Paths.Pages),
		paths.GetAbsolutePath(projectConfig.Paths.Components),
		paths.GetAbsolutePath(projectConfig.Paths.Scripts),
		paths.GetAbsolutePath(projectConfig.Paths.Assets),
		paths.ProjectConfigPath,
		paths.TSConfigPath,
		paths.TailwindConfigPath,
	}

	if len(watchPaths) != len(expectedPaths) {
		t.Errorf("expected %d watch paths, got %d", len(expectedPaths), len(watchPaths))
	}

	for i, expectedPath := range expectedPaths {
		if watchPaths[i] != expectedPath {
			t.Errorf("expected watch path %s, got %s", expectedPath, watchPaths[i])
		}
	}
}

func TestProjectPathsGetTempFile2(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	filename := "test-temp-file.txt"
	tempFile := paths.GetTempFile(filename)

	expectedPath := filepath.Join(paths.ProjectRoot, ".jawt", "tmp", filename)
	if tempFile != expectedPath {
		t.Errorf("expected cache file path %s, got %s", expectedPath, tempFile)
	}
}

func TestProjectPathsGetCacheFile2(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	filename := "test-cache-file.json"
	cacheFile := paths.GetCacheFile(filename)

	expectedPath := filepath.Join(paths.ProjectRoot, ".jawt", "cache", filename)
	if cacheFile != expectedPath {
		t.Errorf("expected cache file path %s, got %s", expectedPath, cacheFile)
	}
}

func TestProjectPathsEnsureDirectories2(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	err = paths.EnsureDirectories()
	if err != nil {
		t.Fatalf("failed to ensure directories: %v", err)
	}

	// Check that all directories were created
	dirsToCheck := []string{
		paths.JawtDir,
		paths.BuildDir,
		paths.DistDir,
		filepath.Join(paths.ProjectRoot, ".jawt", "tmp"),
		filepath.Join(paths.ProjectRoot, ".jawt", "cache"),
		paths.TypeScriptOutputDir,
		paths.TailwindOutputDir,
		paths.ComponentsOutputDir,
	}

	for _, dir := range dirsToCheck {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("directory %s was not created", dir)
		}
	}
}

func TestProjectPathsGetRelativePath2(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	tests := []struct {
		name         string
		absolutePath string
		expected     string
	}{
		{
			name:         "app directory",
			absolutePath: paths.AppDir,
			expected:     projectConfig.Paths.Pages,
		},
		{
			name:         "components directory",
			absolutePath: paths.ComponentsDir,
			expected:     projectConfig.Paths.Components,
		},
		{
			name:         "build directory",
			absolutePath: paths.BuildDir,
			expected:     projectConfig.Build.OutputDir,
		},
		{
			name:         "file in app directory",
			absolutePath: filepath.Join(paths.AppDir, "index.jml"),
			expected:     filepath.Join(projectConfig.Paths.Pages, "index.jml"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := paths.GetRelativePath(tt.absolutePath)
			if result != tt.expected {
				t.Errorf("expected relative path %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestProjectPathsGetAbsolutePath2(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	tests := []struct {
		name         string
		relativePath string
		expected     string
	}{
		{
			name:         "app directory",
			relativePath: projectConfig.Paths.Pages,
			expected:     paths.AppDir,
		},
		{
			name:         "components directory",
			relativePath: projectConfig.Paths.Components,
			expected:     paths.ComponentsDir,
		},
		{
			name:         "file in app directory",
			relativePath: filepath.Join(projectConfig.Paths.Pages, "index.jml"),
			expected:     filepath.Join(paths.AppDir, "index.jml"),
		},
		{
			name:         "already absolute path",
			relativePath: paths.AppDir,
			expected:     paths.AppDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := paths.GetAbsolutePath(tt.relativePath)
			if result != tt.expected {
				t.Errorf("expected absolute path %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestProjectPathsFindFilesWithExtension2(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	// Create test directory structure
	testDir := filepath.Join(tempDir, "test-files")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	// Create test files
	testFiles := []string{
		filepath.Join(testDir, "file1.txt"),
		filepath.Join(testDir, "file2.txt"),
		filepath.Join(testDir, "subdir", "file3.txt"),
		filepath.Join(testDir, "file4.md"),
		filepath.Join(testDir, "file5.txt"),
	}

	for _, file := range testFiles {
		err := os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			t.Fatalf("failed to create directory for %s: %v", file, err)
		}

		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
	}

	// Test finding .txt files
	txtFiles, err := paths.findFilesWithExtension(testDir, ".txt")
	if err != nil {
		t.Fatalf("failed to find .txt files: %v", err)
	}

	expectedTxtCount := 4
	if len(txtFiles) != expectedTxtCount {
		t.Errorf("expected %d .txt files, got %d", expectedTxtCount, len(txtFiles))
	}

	// Test finding .md files
	mdFiles, err := paths.findFilesWithExtension(testDir, ".md")
	if err != nil {
		t.Fatalf("failed to find .md files: %v", err)
	}

	expectedMdCount := 1
	if len(mdFiles) != expectedMdCount {
		t.Errorf("expected %d .md files, got %d", expectedMdCount, len(mdFiles))
	}

	// Test non-existent directory
	nonExistentDir := filepath.Join(tempDir, "non-existent")
	files, err := paths.findFilesWithExtension(nonExistentDir, ".txt")
	if err != nil {
		t.Errorf("expected no error for non-existent directory, got %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files for non-existent directory, got %d", len(files))
	}
}

func TestProjectPathsClean2(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	paths, err := NewProjectPaths(tempDir, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths: %v", err)
	}

	// Create directories
	err = paths.EnsureDirectories()
	if err != nil {
		t.Fatalf("failed to ensure directories: %v", err)
	}

	// Create some test files in the directories
	testFiles := []string{
		filepath.Join(paths.BuildDir, "test.js"),
		filepath.Join(paths.DistDir, "bundle.js"),
		filepath.Join(paths.JawtDir, "metadata.json"),
	}

	for _, file := range testFiles {
		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
	}

	// Verify directories and files exist
	dirsToCheck := []string{
		paths.JawtDir,
		paths.BuildDir,
		paths.DistDir,
	}

	for _, dir := range dirsToCheck {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("directory %s should exist before clean", dir)
		}
	}

	// Clean the directories
	err = paths.Clean()
	if err != nil {
		t.Fatalf("failed to clean directories: %v", err)
	}

	// Verify directories are removed
	for _, dir := range dirsToCheck {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Errorf("directory %s should be removed after clean", dir)
		}
	}
}

func TestProjectPathsEdgeCases2(t *testing.T) {
	// Test with custom project config
	tempDir := t.TempDir()

	customProjectConfig := &ProjectConfig{
		App: struct {
			Name        string `json:"name"`
			Author      string `json:"author"`
			Version     string `json:"version"`
			Description string `json:"description"`
		}{
			Name:        "custom-project",
			Version:     "2.0.0",
			Description: "Custom project for testing",
		},
		Paths: struct {
			Components string `json:"components"`
			Pages      string `json:"pages"`
			Scripts    string `json:"scripts"`
			Assets     string `json:"assets"`
		}{
			Components: "custom-components",
			Pages:      "custom-app",
			Scripts:    "custom-scripts",
			Assets:     "custom-assets",
		},
		Server: struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}{
			Host: "localhost",
			Port: 6500,
		},
		Build: struct {
			OutputDir string `json:"outputDir"`
			DistDir   string `json:"distDir"`
			Minify    bool   `json:"minify"`
			ShadowDOM bool   `json:"shadowDOM"`
		}{
			OutputDir: "custom-build",
			DistDir:   "custom-dist",
			Minify:    true,
			ShadowDOM: false,
		},
		Dev: struct {
			Port       int      `json:"port"`
			EnableHMR  bool     `json:"enableHMR"`
			WatchPaths []string `json:"watchPaths"`
		}{
			Port:       9000,
			EnableHMR:  false,
			WatchPaths: []string{"src", "lib", "assets"},
		},
		Tooling: struct {
			TSConfigPath       string `json:"tsConfigPath"`
			TailwindConfigPath string `json:"tailwindConfigPath"`
		}{
			TSConfigPath:       "custom-tsconfig.json",
			TailwindConfigPath: "custom-tailwind.config.js",
		},
		Scripts: struct {
			PreBuild  []string `json:"preBuild"`
			PostBuild []string `json:"postBuild"`
		}{
			PreBuild:  []string{"prebuild.sh"},
			PostBuild: []string{"postbuild.sh"},
		},
	}

	customJawtConfig := &JawtConfig{
		TypeScriptPath:     "custom-tsc",
		TailwindPath:       "custom-tailwind",
		NodePath:           "custom-node",
		DefaultPort:        7000,
		TempDir:            "custom-tmp",
		CacheDir:           "custom-cache",
		EnableMinification: false,
		EnableSourceMaps:   false,
		EnableTreeShaking:  false,
	}

	paths, err := NewProjectPaths(tempDir, customProjectConfig, customJawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths with custom config: %v", err)
	}

	// Test that custom paths are used
	expectedBuildDir := filepath.Join(paths.ProjectRoot, customProjectConfig.Build.OutputDir)
	if paths.BuildDir != expectedBuildDir {
		t.Errorf("expected custom BuildDir %s, got %s", expectedBuildDir, paths.BuildDir)
	}

	expectedDistDir := filepath.Join(paths.ProjectRoot, customProjectConfig.Build.DistDir)
	if paths.DistDir != expectedDistDir {
		t.Errorf("expected custom DistDir %s, got %s", expectedDistDir, paths.DistDir)
	}

	// TempDir and CacheDir are now fixed to .jawt/tmp and .jawt/cache
	expectedTempDir := filepath.Join(paths.ProjectRoot, ".jawt", "tmp")
	if paths.TempDir != expectedTempDir {
		t.Errorf("expected custom TempDir %s, got %s", expectedTempDir, paths.TempDir)
	}

	expectedCacheDir := filepath.Join(paths.ProjectRoot, ".jawt", "cache")
	if paths.CacheDir != expectedCacheDir {
		t.Errorf("expected custom CacheDir %s, got %s", expectedCacheDir, paths.CacheDir)
	}

	// Test input directories
	expectedAppDir := filepath.Join(paths.ProjectRoot, customProjectConfig.Paths.Pages)
	if paths.AppDir != expectedAppDir {
		t.Errorf("expected AppDir %s, got %s", expectedAppDir, paths.AppDir)
	}

	expectedComponentsDir := filepath.Join(paths.ProjectRoot, customProjectConfig.Paths.Components)
	if paths.ComponentsDir != expectedComponentsDir {
		t.Errorf("expected ComponentsDir %s, got %s", expectedComponentsDir, paths.ComponentsDir)
	}

	expectedScriptsDir := filepath.Join(paths.ProjectRoot, customProjectConfig.Paths.Scripts)
	if paths.ScriptsDir != expectedScriptsDir {
		t.Errorf("expected ScriptsDir %s, got %s", expectedScriptsDir, paths.ScriptsDir)
	}

	expectedAssetsDir := filepath.Join(paths.ProjectRoot, customProjectConfig.Paths.Assets)
	if paths.AssetsDir != expectedAssetsDir {
		t.Errorf("expected AssetsDir %s, got %s", expectedAssetsDir, paths.AssetsDir)
	}

	// Test generated output directories
	expectedTSOutputDir := filepath.Join(paths.BuildDir, "ts")
	if paths.TypeScriptOutputDir != expectedTSOutputDir {
		t.Errorf("expected TypeScriptOutputDir %s, got %s", expectedTSOutputDir, paths.TypeScriptOutputDir)
	}

	expectedTailwindOutputDir := filepath.Join(paths.BuildDir, "styles")
	if paths.TailwindOutputDir != expectedTailwindOutputDir {
		t.Errorf("expected TailwindOutputDir %s, got %s", expectedTailwindOutputDir, paths.TailwindOutputDir)
	}

	expectedComponentsOutputDir := filepath.Join(paths.BuildDir, "components")
	if paths.ComponentsOutputDir != expectedComponentsOutputDir {
		t.Errorf("expected ComponentsOutputDir %s, got %s", expectedComponentsOutputDir, paths.ComponentsOutputDir)
	}

	// Test config file paths
	expectedTSConfigPath := filepath.Join(paths.ProjectRoot, customProjectConfig.Tooling.TSConfigPath)
	if paths.TSConfigPath != expectedTSConfigPath {
		t.Errorf("expected TSConfigPath %s, got %s", expectedTSConfigPath, paths.TSConfigPath)
	}

	expectedTailwindConfigPath := filepath.Join(paths.ProjectRoot, customProjectConfig.Tooling.TailwindConfigPath)
	if paths.TailwindConfigPath != expectedTailwindConfigPath {
		t.Errorf("expected TailwindConfigPath %s, got %s", expectedTailwindConfigPath, paths.TailwindConfigPath)
	}

	expectedProjectConfigPath := filepath.Join(paths.ProjectRoot, "jawt.project.json")
	if paths.ProjectConfigPath != expectedProjectConfigPath {
		t.Errorf("expected ProjectConfigPath %s, got %s", expectedProjectConfigPath, paths.ProjectConfigPath)
	}
}

func TestProjectPathsWithSymlinks2(t *testing.T) {
	tempDir := t.TempDir()
	projectConfig := DefaultProjectConfig()
	jawtConfig := DefaultJawtConfig()

	// Create a symlink to the temp directory
	symlinkPath := filepath.Join(tempDir, "symlink-project")
	targetPath := filepath.Join(tempDir, "actual-project")

	err := os.MkdirAll(targetPath, 0755)
	if err != nil {
		t.Fatalf("failed to create target directory: %v", err)
	}

	// Create symlink (skip test if symlinks are not supported)
	err = os.Symlink(targetPath, symlinkPath)
	if err != nil {
		t.Skipf("symlinks not supported on this platform: %v", err)
	}

	paths, err := NewProjectPaths(symlinkPath, projectConfig, jawtConfig)
	if err != nil {
		t.Fatalf("failed to create project paths with symlink: %v", err)
	}

	// Ensure the project root resolves to absolute path
	if !filepath.IsAbs(paths.ProjectRoot) {
		t.Errorf("expected absolute project root, got %s", paths.ProjectRoot)
	}

	// Test that the absolute path is resolved correctly
	expectedAbsPath, _ := filepath.Abs(targetPath)

	if paths.ProjectRoot != expectedAbsPath {
		t.Errorf("expected ProjectRoot %s, got %s", expectedAbsPath, paths.ProjectRoot)
	}

	// Test input directories
	expectedAppDir := filepath.Join(expectedAbsPath, projectConfig.Paths.Pages)
	if paths.AppDir != expectedAppDir {
		t.Errorf("expected AppDir %s, got %s", expectedAppDir, paths.AppDir)
	}

	expectedComponentsDir := filepath.Join(expectedAbsPath, projectConfig.Paths.Components)
	if paths.ComponentsDir != expectedComponentsDir {
		t.Errorf("expected ComponentsDir %s, got %s", expectedComponentsDir, paths.ComponentsDir)
	}

	expectedScriptsDir := filepath.Join(expectedAbsPath, projectConfig.Paths.Scripts)
	if paths.ScriptsDir != expectedScriptsDir {
		t.Errorf("expected ScriptsDir %s, got %s", expectedScriptsDir, paths.ScriptsDir)
	}

	expectedAssetsDir := filepath.Join(expectedAbsPath, projectConfig.Paths.Assets)
	if paths.AssetsDir != expectedAssetsDir {
		t.Errorf("expected AssetsDir %s, got %s", expectedAssetsDir, paths.AssetsDir)
	}

	// Test output directories
	expectedBuildDir := filepath.Join(expectedAbsPath, projectConfig.Build.OutputDir)
	if paths.BuildDir != expectedBuildDir {
		t.Errorf("expected BuildDir %s, got %s", expectedBuildDir, paths.BuildDir)
	}

	expectedDistDir := filepath.Join(expectedAbsPath, projectConfig.Build.DistDir)
	if paths.DistDir != expectedDistDir {
		t.Errorf("expected DistDir %s, got %s", expectedDistDir, paths.DistDir)
	}

	expectedTempDir := filepath.Join(expectedAbsPath, ".jawt", "tmp")
	if paths.TempDir != expectedTempDir {
		t.Errorf("expected TempDir %s, got %s", expectedTempDir, paths.TempDir)
	}

	expectedCacheDir := filepath.Join(expectedAbsPath, ".jawt", "cache")
	if paths.CacheDir != expectedCacheDir {
		t.Errorf("expected CacheDir %s, got %s", expectedCacheDir, paths.CacheDir)
	}

	// Test generated output directories
	expectedTSOutputDir := filepath.Join(paths.BuildDir, "ts")
	if paths.TypeScriptOutputDir != expectedTSOutputDir {
		t.Errorf("expected TypeScriptOutputDir %s, got %s", expectedTSOutputDir, paths.TypeScriptOutputDir)
	}

	expectedTailwindOutputDir := filepath.Join(paths.BuildDir, "styles")
	if paths.TailwindOutputDir != expectedTailwindOutputDir {
		t.Errorf("expected TailwindOutputDir %s, got %s", expectedTailwindOutputDir, paths.TailwindOutputDir)
	}

	expectedComponentsOutputDir := filepath.Join(paths.BuildDir, "components")
	if paths.ComponentsOutputDir != expectedComponentsOutputDir {
		t.Errorf("expected ComponentsOutputDir %s, got %s", expectedComponentsOutputDir, paths.ComponentsOutputDir)
	}

	// Test config file paths
	expectedTSConfigPath := filepath.Join(expectedAbsPath, projectConfig.Tooling.TSConfigPath)
	if paths.TSConfigPath != expectedTSConfigPath {
		t.Errorf("expected TSConfigPath %s, got %s", expectedTSConfigPath, paths.TSConfigPath)
	}

	expectedTailwindConfigPath := filepath.Join(expectedAbsPath, projectConfig.Tooling.TailwindConfigPath)
	if paths.TailwindConfigPath != expectedTailwindConfigPath {
		t.Errorf("expected TailwindConfigPath %s, got %s", expectedTailwindConfigPath, paths.TailwindConfigPath)
	}

	expectedProjectConfigPath := filepath.Join(expectedAbsPath, "jawt.project.json")
	if paths.ProjectConfigPath != expectedProjectConfigPath {
		t.Errorf("expected ProjectConfigPath %s, got %s", expectedProjectConfigPath, paths.ProjectConfigPath)
	}
}
