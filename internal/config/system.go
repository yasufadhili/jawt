package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func NewInstallPathsConfig(path string) (*InstallPathsConfig, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	execDir := filepath.Dir(execPath)

	return &InstallPathsConfig{
		Jawt:       execDir,
		Bin:        filepath.Join(path, "bin"),
		Config:     filepath.Join(path, "jawt.config.json"),
		Node:       filepath.Join(path, "node"),
		Npm:        filepath.Join(path, "node", "bin", "npm"),
		Tsc:        filepath.Join(path, "tsc"),
		Scrips:     filepath.Join(path, "scripts"),
		Components: filepath.Join(path, "components"),
		Modules:    filepath.Join(path, "modules"),
	}, nil
}

type InstallPathsConfig struct {
	Jawt       string
	Bin        string
	Config     string
	Node       string
	Npm        string
	Tsc        string
	Scrips     string
	Components string
	Modules    string
}

// Validate Checks if all required tools exist
func (c *InstallPathsConfig) Validate() error {
	paths := []struct{ name, path string }{
		{"node", c.Node},
		{"npm", c.Npm},
		{"tsc", c.Tsc},
	}

	for _, p := range paths {
		if _, err := os.Stat(p.path); os.IsNotExist(err) {
			return fmt.Errorf("%s not found at %s", p.name, p.path)
		}
	}
	return nil
}

// GetExecutablePath Get executable path with platform-specific extension
func (c *InstallPathsConfig) GetExecutablePath(tool string) string {
	var path string
	switch tool {
	case "node":
		path = c.Node
	case "npm":
		path = c.Npm
	case "tsc":
		path = c.Tsc
	default:
		return ""
	}

	if runtime.GOOS == "windows" {
		return path + ".exe"
	}
	return path
}
