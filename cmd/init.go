package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/yasufadhili/jawt/internal/config"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yasufadhili/jawt/internal/build"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Create a new JAWT project",
	Long: `Creates a new JAWT project with the default structure and configuration files.
This includes app/, components/, assets/ directories and essential configuration files.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		targetPath := currentDir

		err = build.InitProject(targetPath, projectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initialising project: %v\n", err)
			os.Exit(1)
		}

		err = createDefaultConfigFiles(filepath.Join(targetPath, projectName), projectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating configuration files: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Successfully created new project: %s\n", projectName)
		fmt.Printf("📂 Project location: %s/%s\n", targetPath, projectName)
		fmt.Println("\nNext steps:")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Println("  jawt run")
	},
}

// createDefaultConfigFiles creates the default app.json and jawt.config.json files
func createDefaultConfigFiles(projectDir, projectName string) error {
	// Create app.json
	appConfig := config.AppConfig{
		Name:         projectName,
		Description:  "A JAWT application",
		Version:      "0.1.0",
		Author:       "",
		License:      "MIT",
		Dependencies: []string{},
	}

	appJSON, err := json.MarshalIndent(appConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling app.json: %w", err)
	}

	err = os.WriteFile(filepath.Join(projectDir, "app.json"), appJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing app.json: %w", err)
	}

	// Create jawt.config.json
	jawtConfig := config.JawtConfig{
		Project: config.ProjectConfig{
			Name: projectName,
		},
		Server: config.ServerConfig{
			Port: 6500,
		},
		Build: config.BuildConfig{
			Output: "dist",
			Minify: true,
		},
	}

	jawtConfigJSON, err := json.MarshalIndent(jawtConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling jawt.config.json: %w", err)
	}

	err = os.WriteFile(filepath.Join(projectDir, "jawt.config.json"), jawtConfigJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing jawt.config.json: %w", err)
	}

	return nil
}
