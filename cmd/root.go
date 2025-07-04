package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jawt",
	Short: "JAWT - Just Another Web Tool",
	Long: `JAWT is a tool for creating, developing, and building web applications.
It offers a streamlined workflow and unified development experience.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
	rootCmd.Flags().BoolP("version", "v", false, "Display Jawt version information")

	// rootCmd.AddCommand(initCmd)
	// rootCmd.AddCommand(runCmd)
	// rootCmd.AddCommand(buildCmd)
	// rootCmd.AddCommand(versionCmd)
	// rootCmd.AddCommand(tscCmd)
	//rootCmd.AddCommand(updateCmd) // Not yet useful
}
