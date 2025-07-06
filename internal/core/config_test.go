package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestDefaultJawtConfig(t *testing.T) {
	config := DefaultJawtConfig()

	if config.TypeScriptPath != "tsc" {
		t.Errorf("expected TypeScriptPath to be 'tsc', got %s", config.TypeScriptPath)
	}
	if config.TailwindPath != "tailwindcss" {
		t.Errorf("expected TailwindPath to be 'tailwindcss', got %s", config.TailwindPath)
	}
	if config.NodePath != "node" {
		t.Errorf("expected NodePath to be 'node', got %s", config.NodePath)
	}
	if config.DefaultPort != 6500 {
		t.Errorf("expected DefaultPort to be 6500, got %d", config.DefaultPort)
	}
	if config.TempDir != ".jawt/tmp" {
		t.Errorf("expected TempDir to be '.jawt/tmp', got %s", config.TempDir)
	}
	if config.CacheDir != ".jawt/cache" {
		t.Errorf("expected CacheDir to be '.jawt/cache', got %s", config.CacheDir)
	}
	if !config.EnableMinification {
		t.Error("expected EnableMinification to be true")
	}
	if !config.EnableSourceMaps {
		t.Error("expected EnableSourceMaps to be true")
	}
	if !config.EnableTreeShaking {
		t.Error("expected EnableTreeShaking to be true")
	}
}

func TestDefaultProjectConfig(t *testing.T) {
	config := DefaultProjectConfig()

	if config.App.Name != "jawt-project" {
		t.Errorf("expected Name to be 'jawt-project', got %s", config.App.Name)
	}
	if config.App.Version != "0.1.0" {
		t.Errorf("expected App.Version to be '0.1.0', got %s", config.App.Version)
	}
	if config.App.Description != "A Jawt application" {
		t.Errorf("expected App.Description to be 'A Jawt application', got %s", config.App.Description)
	}
	if config.Paths.Components != "components" {
		t.Errorf("expected Paths.Components to be 'components', got %s", config.Paths.Components)
	}
	if config.Paths.Pages != "app" {
		t.Errorf("expected Paths.Pages to be 'app', got %s", config.Paths.Pages)
	}
	if config.Paths.Scripts != "scripts" {
		t.Errorf("expected Paths.Scripts to be 'scripts', got %s", config.Paths.Scripts)
	}
	if config.Paths.Assets != "assets" {
		t.Errorf("expected Paths.Assets to be 'assets', got %s", config.Paths.Assets)
	}
	if config.Server.Host != "localhost" {
		t.Errorf("expected Server.Host to be 'localhost', got %s", config.Server.Host)
	}
	if config.Server.Port != 6500 {
		t.Errorf("expected Server.Port to be 6500, got %d", config.Server.Port)
	}
	if config.Build.OutputDir != ".jawt/build" {
		t.Errorf("expected Build.OutputDir to be '.jawt/build', got %s", config.Build.OutputDir)
	}
	if config.Build.DistDir != ".jawt/dist" {
		t.Errorf("expected Build.DistDir to be '.jawt/dist', got %s", config.Build.DistDir)
	}
	if !config.Build.Minify {
		t.Error("expected Build.Minify to be true")
	}
	if config.Build.ShadowDOM {
		t.Error("expected Build.ShadowDOM to be false")
	}
	if config.Dev.Port != 6500 {
		t.Errorf("expected Dev.Port to be 6500, got %d", config.Dev.Port)
	}
	if !config.Dev.EnableHMR {
		t.Error("expected Dev.EnableHMR to be true")
	}
	expectedWatchPaths := []string{"app", "components", "scripts", "assets"}
	if len(config.Dev.WatchPaths) != len(expectedWatchPaths) {
		t.Errorf("expected Dev.WatchPaths length to be %d, got %d", len(expectedWatchPaths), len(config.Dev.WatchPaths))
	}
	for i, path := range expectedWatchPaths {
		if config.Dev.WatchPaths[i] != path {
			t.Errorf("expected Dev.WatchPaths[%d] to be %s, got %s", i, path, config.Dev.WatchPaths[i])
		}
	}
	if config.Tooling.TSConfigPath != "tsconfig.json" {
		t.Errorf("expected Tooling.TSConfigPath to be 'tsconfig.json', got %s", config.Tooling.TSConfigPath)
	}
	if config.Tooling.TailwindConfigPath != "tailwind.config.js" {
		t.Errorf("expected Tooling.TailwindConfigPath to be 'tailwind.config.js', got %s", config.Tooling.TailwindConfigPath)
	}
	if len(config.Scripts.PreBuild) != 0 {
		t.Errorf("expected Scripts.PreBuild to be empty, got %v", config.Scripts.PreBuild)
	}
	if len(config.Scripts.PostBuild) != 0 {
		t.Errorf("expected Scripts.PostBuild to be empty, got %v", config.Scripts.PostBuild)
	}
}

func TestLoadJawtConfig(t *testing.T) {
	tests := []struct {
		name         string
		configPath   string
		createConfig bool
		configData   *JawtConfig
		expectError  bool
	}{
		{
			name:         "empty config path returns default",
			configPath:   "",
			createConfig: false,
			expectError:  false,
		},
		{
			name:         "non-existent config returns default",
			configPath:   "non-existent-config.json",
			createConfig: false,
			expectError:  false,
		},
		{
			name:         "valid config file",
			configPath:   "test-config.json",
			createConfig: true,
			configData: &JawtConfig{
				TypeScriptPath:     "custom-tsc",
				TailwindPath:       "custom-tailwind",
				NodePath:           "custom-node",
				DefaultPort:        8080,
				TempDir:            "custom-tmp",
				CacheDir:           "custom-cache",
				EnableMinification: false,
				EnableSourceMaps:   false,
				EnableTreeShaking:  false,
			},
			expectError: false,
		},
		{
			name:         "invalid json",
			configPath:   "invalid-config.json",
			createConfig: true,
			configData:   nil, // Will create invalid JSON
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createConfig {
				if tt.configData == nil {
					// Create invalid JSON
					err := os.WriteFile(tt.configPath, []byte("invalid json"), 0644)
					if err != nil {
						t.Fatalf("failed to create invalid config file: %v", err)
					}
				} else {
					// Create valid JSON
					data, err := json.MarshalIndent(tt.configData, "", "  ")
					if err != nil {
						t.Fatalf("failed to marshal config data: %v", err)
					}
					err = os.WriteFile(tt.configPath, data, 0644)
					if err != nil {
						t.Fatalf("failed to create config file: %v", err)
					}
				}
				defer os.Remove(tt.configPath)
			}

			config, err := LoadJawtConfig(tt.configPath)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if config == nil {
				t.Error("expected config but got nil")
				return
			}

			if tt.configData != nil {
				if config.TypeScriptPath != tt.configData.TypeScriptPath {
					t.Errorf("expected TypeScriptPath %s, got %s", tt.configData.TypeScriptPath, config.TypeScriptPath)
				}
				if config.DefaultPort != tt.configData.DefaultPort {
					t.Errorf("expected DefaultPort %d, got %d", tt.configData.DefaultPort, config.DefaultPort)
				}
			}
		})
	}
}

func TestLoadProjectConfig(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		createConfig bool
		configData   *ProjectConfig
		expectError  bool
	}{
		{
			name:         "no config file returns default",
			createConfig: false,
			expectError:  false,
		},
		{
			name:         "valid config file",
			createConfig: true,
			configData: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name:        "test-project",
					Version:     "2.0.0",
					Description: "Test project",
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
					Host: "custom-host",
					Port: 8080,
				},
				Build: struct {
					OutputDir string `json:"outputDir"`
					DistDir   string `json:"distDir"`
					Minify    bool   `json:"minify"`
					ShadowDOM bool   `json:"shadowDOM"`
				}{
					OutputDir: "custom-build",
					DistDir:   "custom-dist",
					Minify:    false,
					ShadowDOM: true,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       3000,
					EnableHMR:  false,
					WatchPaths: []string{"src", "lib"},
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
					PreBuild:  []string{"script1.sh"},
					PostBuild: []string{"script2.sh"},
				},
			},
			expectError: false,
		},
		{
			name:         "invalid json",
			createConfig: true,
			configData:   nil, // Will create invalid JSON
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectDir := filepath.Join(tempDir, tt.name)
			err := os.MkdirAll(projectDir, 0755)
			if err != nil {
				t.Fatalf("failed to create project directory: %v", err)
			}

			if tt.createConfig {
				configPath := filepath.Join(projectDir, "jawt.project.json")
				if tt.configData == nil {
					// Create invalid JSON
					err := os.WriteFile(configPath, []byte("invalid json"), 0644)
					if err != nil {
						t.Fatalf("failed to create invalid config file: %v", err)
					}
				} else {
					// Create valid JSON
					data, err := json.MarshalIndent(tt.configData, "", "  ")
					if err != nil {
						t.Fatalf("failed to marshal config data: %v", err)
					}
					err = os.WriteFile(configPath, data, 0644)
					if err != nil {
						t.Fatalf("failed to create config file: %v", err)
					}
				}
			}

			config, err := LoadProjectConfig(projectDir)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if config == nil {
				t.Error("expected config but got nil")
				return
			}

			if tt.configData != nil {
				if config.App.Name != tt.configData.App.Name {
					t.Errorf("expected Name %s, got %s", tt.configData.App.Name, config.App.Name)
				}
				if config.App.Version != tt.configData.App.Version {
					t.Errorf("expected Version %s, got %s", tt.configData.App.Version, config.App.Version)
				}
				if config.Server.Port != tt.configData.Server.Port {
					t.Errorf("expected DevPort %d, got %d", tt.configData.Server.Port, config.Server.Port)
				}
				if config.Build.ShadowDOM != tt.configData.Build.ShadowDOM {
					t.Errorf("expected ShadowDOM %t, got %t", tt.configData.Build.ShadowDOM, config.Build.ShadowDOM)
				}
			}
		})
	}
}

func TestJawtConfigSave(t *testing.T) {
	config := &JawtConfig{
		TypeScriptPath:     "test-tsc",
		TailwindPath:       "test-tailwind",
		NodePath:           "test-node",
		DefaultPort:        9000,
		TempDir:            "test-tmp",
		CacheDir:           "test-cache",
		EnableMinification: true,
		EnableSourceMaps:   true,
		EnableTreeShaking:  true,
	}

	tempFile := filepath.Join(t.TempDir(), "test-config.json")

	err := config.Save(tempFile)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Load the saved config and verify
	loadedConfig, err := LoadJawtConfig(tempFile)
	if err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}

	if loadedConfig.TypeScriptPath != config.TypeScriptPath {
		t.Errorf("expected TypeScriptPath %s, got %s", config.TypeScriptPath, loadedConfig.TypeScriptPath)
	}
	if loadedConfig.DefaultPort != config.DefaultPort {
		t.Errorf("expected DefaultPort %d, got %d", config.DefaultPort, loadedConfig.DefaultPort)
	}
}

func TestProjectConfigSave(t *testing.T) {
	config := &ProjectConfig{
		App: struct {
			Name        string `json:"name"`
			Author      string `json:"author"`
			Version     string `json:"version"`
			Description string `json:"description"`
		}{
			Name:        "test-save-project",
			Version:     "1.2.3",
			Description: "Test save project",
		},
	}

	tempDir := t.TempDir()

	err := config.Save(tempDir)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Load the saved config and verify
	loadedConfig, err := LoadProjectConfig(tempDir)
	if err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}

	if loadedConfig.App.Name != config.App.Name {
		t.Errorf("expected Name %s, got %s", config.App.Name, loadedConfig.App.Name)
	}
	if loadedConfig.App.Version != config.App.Version {
		t.Errorf("expected Version %s, got %s", config.App.Version, loadedConfig.App.Version)
	}
	if loadedConfig.Server.Port != config.Server.Port {
		t.Errorf("expected DevPort %d, got %d", config.Server.Port, loadedConfig.Server.Port)
	}
	// if loadedConfig.HasTailwindConfig != config.HasTailwindConfig {
	//	t.Errorf("expected HasTailwindConfig %t, got %t", config.HasTailwindConfig, loadedConfig.HasTailwindConfig)
	// }
}

func TestJawtConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *JawtConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid config",
			config:      DefaultJawtConfig(),
			expectError: false,
		},
		{
			name: "invalid port - zero",
			config: &JawtConfig{
				DefaultPort: 0,
				TempDir:     "temp",
				CacheDir:    "cache",
			},
			expectError: true,
			errorMsg:    "invalid default port: " + strconv.Itoa(0),
		},
		{
			name: "invalid port - negative",
			config: &JawtConfig{
				DefaultPort: -1,
				TempDir:     "temp",
				CacheDir:    "cache",
			},
			expectError: true,
			errorMsg:    "invalid default port: " + strconv.Itoa(-1),
		},
		{
			name: "invalid port - too high",
			config: &JawtConfig{
				DefaultPort: 65536,
				TempDir:     "temp",
				CacheDir:    "cache",
			},
			expectError: true,
			errorMsg:    "invalid default port: " + strconv.Itoa(65536),
		},
		{
			name: "empty temp dir",
			config: &JawtConfig{
				DefaultPort: 8080,
				TempDir:     "",
				CacheDir:    "cache",
			},
			expectError: true,
			errorMsg:    "temp directory cannot be empty",
		},
		{
			name: "empty cache dir",
			config: &JawtConfig{
				DefaultPort: 8080,
				TempDir:     "temp",
				CacheDir:    "",
			},
			expectError: true,
			errorMsg:    "cache directory cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestProjectConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *ProjectConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid config",
			config:      DefaultProjectConfig(),
			expectError: false,
		},
		{
			name: "empty app name",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "components",
					Pages:      "app",
					Scripts:    "scripts",
					Assets:     "assets",
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
					OutputDir: "build",
					Minify:    true,
					ShadowDOM: false,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       6500,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "app name cannot be empty",
		},
		{
			name: "empty components path",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "test",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "",
					Pages:      "app",
					Scripts:    "scripts",
					Assets:     "assets",
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
					OutputDir: "build",
					Minify:    true,
					ShadowDOM: false,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       6500,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "components path cannot be empty",
		},
		{
			name: "empty pages path",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "test",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "components",
					Pages:      "",
					Scripts:    "scripts",
					Assets:     "assets",
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
				}(struct {
					OutputDir string `json:"outputDir"`
					DistDir   string `json:"dist"`
					Minify    bool   `json:"minify"`
					ShadowDOM bool   `json:"shadowDOM"`
				}{
					OutputDir: "build",
					Minify:    true,
					ShadowDOM: false,
				}),
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       6500,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "pages path cannot be empty",
		},
		{
			name: "empty scripts path",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "test",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "components",
					Pages:      "app",
					Scripts:    "",
					Assets:     "assets",
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
					OutputDir: "build",
					Minify:    true,
					ShadowDOM: false,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       6500,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "scripts path cannot be empty",
		},
		{
			name: "empty assets path",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "test",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "components",
					Pages:      "app",
					Scripts:    "scripts",
					Assets:     "",
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
					OutputDir: "build",
					Minify:    true,
					ShadowDOM: false,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       6500,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "assets path cannot be empty",
		},
		{
			name: "empty output dir",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "test",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "components",
					Pages:      "app",
					Scripts:    "scripts",
					Assets:     "assets",
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
					OutputDir: "",
					Minify:    true,
					ShadowDOM: false,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       6500,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "build output directory cannot be empty",
		},
		{
			name: "invalid server port - zero",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "test",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "components",
					Pages:      "app",
					Scripts:    "scripts",
					Assets:     "assets",
				},
				Server: struct {
					Host string `json:"host"`
					Port int    `json:"port"`
				}{
					Host: "localhost",
					Port: 0,
				},
				Build: struct {
					OutputDir string `json:"outputDir"`
					DistDir   string `json:"distDir"`
					Minify    bool   `json:"minify"`
					ShadowDOM bool   `json:"shadowDOM"`
				}{
					OutputDir: "build",
					Minify:    true,
					ShadowDOM: false,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       6500,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "invalid server port: " + strconv.Itoa(0),
		},
		{
			name: "invalid server port - negative",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "test",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "components",
					Pages:      "app",
					Scripts:    "scripts",
					Assets:     "assets",
				},
				Server: struct {
					Host string `json:"host"`
					Port int    `json:"port"`
				}{
					Host: "localhost",
					Port: -1,
				},
				Build: struct {
					OutputDir string `json:"outputDir"`
					DistDir   string `json:"distDir"`
					Minify    bool   `json:"minify"`
					ShadowDOM bool   `json:"shadowDOM"`
				}{
					OutputDir: "build",
					Minify:    true,
					ShadowDOM: false,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       6500,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "invalid server port: " + strconv.Itoa(-1),
		},
		{
			name: "invalid server port - too high",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "test",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "components",
					Pages:      "app",
					Scripts:    "scripts",
					Assets:     "assets",
				},
				Server: struct {
					Host string `json:"host"`
					Port int    `json:"port"`
				}{
					Host: "localhost",
					Port: 65536,
				},
				Build: struct {
					OutputDir string `json:"outputDir"`
					DistDir   string `json:"distDir"`
					Minify    bool   `json:"minify"`
					ShadowDOM bool   `json:"shadowDOM"`
				}{
					OutputDir: "build",
					Minify:    true,
					ShadowDOM: false,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       6500,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "invalid server port: " + strconv.Itoa(65536),
		},
		{
			name: "empty server host",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "test",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "components",
					Pages:      "app",
					Scripts:    "scripts",
					Assets:     "assets",
				},
				Server: struct {
					Host string `json:"host"`
					Port int    `json:"port"`
				}{
					Host: "",
					Port: 6500,
				},
				Build: struct {
					OutputDir string `json:"outputDir"`
					DistDir   string `json:"distDir"`
					Minify    bool   `json:"minify"`
					ShadowDOM bool   `json:"shadowDOM"`
				}{
					OutputDir: "build",
					Minify:    true,
					ShadowDOM: false,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       6500,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "server host cannot be empty",
		},
		{
			name: "invalid dev port - zero",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "test",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "components",
					Pages:      "app",
					Scripts:    "scripts",
					Assets:     "assets",
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
					OutputDir: "build",
					Minify:    true,
					ShadowDOM: false,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       0,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "invalid dev server port: " + strconv.Itoa(0),
		},
		{
			name: "invalid dev port - negative",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "test",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "components",
					Pages:      "app",
					Scripts:    "scripts",
					Assets:     "assets",
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
					OutputDir: "build",
					Minify:    true,
					ShadowDOM: false,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       -1,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "invalid dev server port: " + strconv.Itoa(-1),
		},
		{
			name: "invalid dev port - too high",
			config: &ProjectConfig{
				App: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     string `json:"version"`
					Description string `json:"description"`
				}{
					Name: "test",
				},
				Paths: struct {
					Components string `json:"components"`
					Pages      string `json:"pages"`
					Scripts    string `json:"scripts"`
					Assets     string `json:"assets"`
				}{
					Components: "components",
					Pages:      "app",
					Scripts:    "scripts",
					Assets:     "assets",
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
					OutputDir: "build",
					Minify:    true,
					ShadowDOM: false,
				},
				Dev: struct {
					Port       int      `json:"port"`
					EnableHMR  bool     `json:"enableHMR"`
					WatchPaths []string `json:"watchPaths"`
				}{
					Port:       65536,
					EnableHMR:  true,
					WatchPaths: []string{"app", "components", "scripts", "assets"},
				},
				Tooling: struct {
					TSConfigPath       string `json:"tsConfigPath"`
					TailwindConfigPath string `json:"tailwindConfigPath"`
				}{
					TSConfigPath:       "tsconfig.json",
					TailwindConfigPath: "tailwind.config.js",
				},
				Scripts: struct {
					PreBuild  []string `json:"preBuild"`
					PostBuild []string `json:"postBuild"`
				}{
					PreBuild:  []string{},
					PostBuild: []string{},
				},
			},
			expectError: true,
			errorMsg:    "invalid dev server port: " + strconv.Itoa(65536),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
