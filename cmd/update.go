package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const (
	NodeVersion = "v20.11.0"
	NodeBaseURL = "https://nodejs.org/dist"
)

type InstallPathsConfig struct {
	Jawt   string
	Config string
	Node   string
	Npm    string
	Tsc    string
}

func NewInstallPathsConfig() (*InstallPathsConfig, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	execDir := filepath.Dir(execPath)

	// Check if necessary directories exist in the current directory
	if pathsExist(execDir) {
		return buildConfig(execDir), nil
	}

	// Check parent directory
	parentDir := filepath.Dir(execDir)
	if pathsExist(parentDir) {
		return buildConfig(parentDir), nil
	}

	// Default to executable directory
	return buildConfig(execDir), nil
}

func pathsExist(basePath string) bool {
	// Check if this looks like our toolchain directory structure
	configPath := filepath.Join(basePath, "jawt.config.json")
	_, err := os.Stat(configPath)
	return err == nil
}

func buildConfig(basePath string) *InstallPathsConfig {
	nodeDir := filepath.Join(basePath, "node")

	return &InstallPathsConfig{
		Jawt:   basePath,
		Config: filepath.Join(basePath, "jawt.config.json"),
		Node:   getNodeExecutable(nodeDir),
		Npm:    getNpmExecutable(nodeDir),
		Tsc:    getTscExecutable(nodeDir),
	}
}

func getNodeExecutable(nodeDir string) string {
	nodePath := filepath.Join(nodeDir, "bin", "node")
	if runtime.GOOS == "windows" {
		return filepath.Join(nodeDir, "node.exe")
	}
	return nodePath
}

func getNpmExecutable(nodeDir string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(nodeDir, "npm.cmd")
	}
	return filepath.Join(nodeDir, "bin", "npm")
}

func getTscExecutable(nodeDir string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(nodeDir, "node_modules", ".bin", "tsc.cmd")
	}
	return filepath.Join(nodeDir, "bin", "tsc")
}

func downloadAndExtractNode(config *InstallPathsConfig) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("node download only supported on Linux")
	}

	arch := "x64"
	if runtime.GOARCH == "arm64" {
		arch = "arm64"
	}

	filename := fmt.Sprintf("node-%s-linux-%s.tar.gz", NodeVersion, arch)
	url := fmt.Sprintf("%s/%s/%s", NodeBaseURL, NodeVersion, filename)

	fmt.Printf("Downloading Node.js %s for Linux %s...\n", NodeVersion, arch)

	// Download
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download Node.js: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download Node.js: HTTP %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "node-*.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Copy download to the temp file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save download: %w", err)
	}

	// Extract to node directory
	nodeDir := filepath.Dir(config.Node)
	if err := os.RemoveAll(nodeDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing node directory: %w", err)
	}

	if err := os.MkdirAll(nodeDir, 0755); err != nil {
		return fmt.Errorf("failed to create node directory: %w", err)
	}

	// Extract tar.gz
	tmpFile.Seek(0, 0)
	if err := extractTarGz(tmpFile, nodeDir); err != nil {
		return fmt.Errorf("failed to extract Node.js: %w", err)
	}

	fmt.Println("Node.js extracted successfully")
	return nil
}

func extractTarGz(src io.Reader, destDir string) error {
	gzr, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Remove the top-level directory from the path
		parts := strings.Split(header.Name, "/")
		if len(parts) <= 1 {
			continue
		}
		relativePath := filepath.Join(parts[1:]...)

		target := filepath.Join(destDir, relativePath)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(file, tr); err != nil {
				file.Close()
				return err
			}
			file.Close()
		}
	}

	return nil
}

// Install TypeScript globally using our portable npm
func installTypeScript(config *InstallPathsConfig) error {
	fmt.Println("Installing TypeScript compiler...")

	cmd := exec.Command(config.Npm, "install", "-g", "typescript")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set NODE_PATH to our portable installation
	nodeModulesPath := filepath.Join(filepath.Dir(config.Node), "lib", "node_modules")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("NODE_PATH=%s", nodeModulesPath),
		fmt.Sprintf("PATH=%s:%s", filepath.Dir(config.Node), os.Getenv("PATH")),
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install TypeScript: %w", err)
	}

	fmt.Println("TypeScript installed successfully")
	return nil
}

// CallTSC Calls TypeScript compiler with arguments
func CallTSC(config *InstallPathsConfig, args []string) error {
	if _, err := os.Stat(config.Tsc); os.IsNotExist(err) {
		return fmt.Errorf("TypeScript compiler not found at %s. Run update first", config.Tsc)
	}

	cmd := exec.Command(config.Node, config.Tsc)
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Set environment for portable setup
	nodeModulesPath := filepath.Join(filepath.Dir(config.Node), "lib", "node_modules")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("NODE_PATH=%s", nodeModulesPath),
		fmt.Sprintf("PATH=%s:%s", filepath.Dir(config.Node), os.Getenv("PATH")),
	)

	return cmd.Run()
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Jawt and all dependencies to the latest version",
	Long: `Downloads and installs Jawt, Node.js and TypeScript to use with the toolchain.
This command will:
1. Download Jawt it's dependencies from wherever it's called
2. Setup Node.js for the current Jawt installation
3. Install TypeScript within the Jawt environment`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := NewInstallPathsConfig()
		if err != nil {
			return fmt.Errorf("failed to determine install paths: %w", err)
		}

		fmt.Printf("Installing to: %s\n", config.Jawt)

		// Download and extract Node.js (Linux only for now)
		if err := downloadAndExtractNode(config); err != nil {
			return err
		}

		if err := installTypeScript(config); err != nil {
			return err
		}

		fmt.Println("Update completed successfully!")
		fmt.Printf("Node.js: %s\n", config.Node)
		fmt.Printf("npm: %s\n", config.Npm)
		fmt.Printf("tsc: %s\n", config.Tsc)

		return nil
	},
}

var tscCmd = &cobra.Command{
	Use:   "tsc [typescript-args...]",
	Short: "Run TypeScript compiler",
	Long: `Run the TypeScript compiler within Jawt with the provided arguments.
All arguments are passed directly to tsc.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := NewInstallPathsConfig()
		if err != nil {
			return fmt.Errorf("failed to determine install paths: %w", err)
		}

		return CallTSC(config, args)
	},
}
