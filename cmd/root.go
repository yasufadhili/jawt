package cmd

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/config"
	"os"

	"github.com/spf13/cobra"
)

// Global configuration for use across commands
var (
	projectConfig *config.Config
	projectDir    string
)

var rootCmd = &cobra.Command{
	Use:   "jawt",
	Short: "JAWT - Just Another Web Tool",
	Long: `JAWT is a tool for creating, developing, and building minimal web applications.
It offers a streamlined workflow and unified development experience.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip config loading for commands that don't require it
		if cmd.Name() == "init" || cmd.Name() == "version" {
			return
		}

		var err error
		projectDir, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		if !config.IsJawtProject(projectDir) {
			fmt.Fprintf(os.Stderr, "Error: Current directory is not a JAWT project.\n")
			fmt.Fprintf(os.Stderr, "Run 'jawt init <project-name>' to create a new project.\n")
			os.Exit(1)
		}

		projectConfig, err = config.LoadConfig(projectDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading project configuration: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Display help information for the command")
	rootCmd.Flags().BoolP("version", "v", false, "Display JAWT version information")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(tscCmd)
	//rootCmd.AddCommand(updateCmd) // Not yet useful
}
