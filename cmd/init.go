package cmd

import (
	"fmt"
	"os"

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

		fmt.Printf("✅ Successfully created new project: %s\n", projectName)
	},
}
