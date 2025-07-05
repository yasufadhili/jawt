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

	if config.Name != "jawt-project" {
		t.Errorf("expected Name to be 'jawt-project', got %s", config.Name)
	}
	if config.Version != "1.0.0" {
		t.Errorf("expected Version to be '1.0.0', got %s", config.Version)
	}
	if config.Description != "A Jawt application" {
		t.Errorf("expected Description to be 'A Jawt application', got %s", config.Description)
	}
	if config.OutputDir != ".jawt/build" {
		t.Errorf("expected OutputDir to be '.jawt/build', got %s", config.OutputDir)
	}
	if config.DistDir != ".jawt/dist" {
		t.Errorf("expected DistDir to be '.jawt/dist', got %s", config.DistDir)
	}
	if config.ShadowDOM {
		t.Error("expected ShadowDOM to be false")
	}
	if config.DevPort != 6500 {
		t.Errorf("expected DevPort to be 6500, got %d", config.DevPort)
	}
	if !config.EnableHMR {
		t.Error("expected EnableHMR to be true")
	}
	expectedWatchPaths := []string{"app", "components", "scripts"}
	if len(config.WatchPaths) != len(expectedWatchPaths) {
		t.Errorf("expected WatchPaths length to be %d, got %d", len(expectedWatchPaths), len(config.WatchPaths))
	}
	for i, path := range expectedWatchPaths {
		if config.WatchPaths[i] != path {
			t.Errorf("expected WatchPaths[%d] to be %s, got %s", i, path, config.WatchPaths[i])
		}
	}
	if config.TSConfigPath != "tsconfig.json" {
		t.Errorf("expected TSConfigPath to be 'tsconfig.json', got %s", config.TSConfigPath)
	}
	if config.TailwindConfigPath != "tailwind.config.js" {
		t.Errorf("expected TailwindConfigPath to be 'tailwind.config.js', got %s", config.TailwindConfigPath)
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
				Name:               "test-project",
				Version:            "2.0.0",
				Description:        "Test project",
				OutputDir:          "custom-build",
				DistDir:            "custom-dist",
				ShadowDOM:          true,
				DevPort:            3000,
				EnableHMR:          false,
				WatchPaths:         []string{"src", "lib"},
				TSConfigPath:       "custom-tsconfig.json",
				TailwindConfigPath: "custom-tailwind.config.js",
				PreBuildScripts:    []string{"script1.sh"},
				PostBuildScripts:   []string{"script2.sh"},
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
				if config.Name != tt.configData.Name {
					t.Errorf("expected Name %s, got %s", tt.configData.Name, config.Name)
				}
				if config.Version != tt.configData.Version {
					t.Errorf("expected Version %s, got %s", tt.configData.Version, config.Version)
				}
				if config.DevPort != tt.configData.DevPort {
					t.Errorf("expected DevPort %d, got %d", tt.configData.DevPort, config.DevPort)
				}
				if config.ShadowDOM != tt.configData.ShadowDOM {
					t.Errorf("expected ShadowDOM %t, got %t", tt.configData.ShadowDOM, config.ShadowDOM)
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
		Name:               "test-save-project",
		Version:            "1.2.3",
		Description:        "Test save project",
		OutputDir:          "test-build",
		DistDir:            "test-dist",
		ShadowDOM:          true,
		DevPort:            4000,
		EnableHMR:          false,
		WatchPaths:         []string{"test-src"},
		TSConfigPath:       "test-tsconfig.json",
		TailwindConfigPath: "test-tailwind.config.js",
		PreBuildScripts:    []string{"pre.sh"},
		PostBuildScripts:   []string{"post.sh"},
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

	if loadedConfig.Name != config.Name {
		t.Errorf("expected Name %s, got %s", config.Name, loadedConfig.Name)
	}
	if loadedConfig.Version != config.Version {
		t.Errorf("expected Version %s, got %s", config.Version, loadedConfig.Version)
	}
	if loadedConfig.DevPort != config.DevPort {
		t.Errorf("expected DevPort %d, got %d", config.DevPort, loadedConfig.DevPort)
	}
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
			name: "empty name",
			config: &ProjectConfig{
				Name:      "",
				OutputDir: "build",
				DistDir:   "dist",
				DevPort:   8080,
			},
			expectError: true,
			errorMsg:    "project name cannot be empty",
		},
		{
			name: "empty output dir",
			config: &ProjectConfig{
				Name:      "test",
				OutputDir: "",
				DistDir:   "dist",
				DevPort:   8080,
			},
			expectError: true,
			errorMsg:    "output directory cannot be empty",
		},
		{
			name: "empty dist dir",
			config: &ProjectConfig{
				Name:      "test",
				OutputDir: "build",
				DistDir:   "",
				DevPort:   8080,
			},
			expectError: true,
			errorMsg:    "dist directory cannot be empty",
		},
		{
			name: "invalid dev port - zero",
			config: &ProjectConfig{
				Name:      "test",
				OutputDir: "build",
				DistDir:   "dist",
				DevPort:   0,
			},
			expectError: true,
			errorMsg:    "invalid dev port: " + strconv.Itoa(0),
		},
		{
			name: "invalid dev port - negative",
			config: &ProjectConfig{
				Name:      "test",
				OutputDir: "build",
				DistDir:   "dist",
				DevPort:   -1,
			},
			expectError: true,
			errorMsg:    "invalid dev port: " + strconv.Itoa(-1),
		},
		{
			name: "invalid dev port - too high",
			config: &ProjectConfig{
				Name:      "test",
				OutputDir: "build",
				DistDir:   "dist",
				DevPort:   65536,
			},
			expectError: true,
			errorMsg:    "invalid dev port: " + strconv.Itoa(65536),
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
