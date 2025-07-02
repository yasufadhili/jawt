package cmd

import (
	"github.com/spf13/cobra"
)

var outputDir string

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build optimised production bundle",
	Long: `Compiles your Jawt application into production-ready web standard files.
Generates optimised web standard output.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	buildCmd.Flags().StringVarP(&outputDir, "output", "o", "dist", "Specify custom output directory")
}
