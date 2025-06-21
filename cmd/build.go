package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yasufadhili/jawt/internal/build"
)

var outputDir string

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build optimised production bundle",
	Long: `Compiles your JAWT application into production-ready web standard files.
Generates optimised HTML, CSS, JavaScript and WebAssembly output.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		builder := build.NewBuilder(dir)

		// TODO: Pass output directory to builder when supported

		err = builder.Build()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building project: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Build completed successfully!")
	},
}

func init() {
	buildCmd.Flags().StringVarP(&outputDir, "output", "o", "dist", "Specify custom output directory")
}
