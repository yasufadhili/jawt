package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Create a new Jawt project",
	Long: `Creates a new Jawt project with the default structure and configuration files.
This includes essential directories and configuration files.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚠️ The 'init' command is not yet implemented.")
		fmt.Println("Please check future versions for this functionality.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
