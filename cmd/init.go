package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yasufadhili/jawt/internal/build"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Create a new Jawt project",
	Long: `Creates a new Jawt project with the default structure and configuration files.
This includes app/, components/, assets/ directories and essential configuration files.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		err := build.InitProject(projectName)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error initialising project: %s\n", err)
			os.Exit(1)
		}

	},
}
