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
		// Use a configured output directory if not explicitly provided as a flag
		if !cmd.Flags().Changed("output") && projectConfig != nil {
			outputDir = projectConfig.Jawt.OutputDir
		}

		fmt.Printf("🔨 Building %s for production...\n", projectConfig.App.Name)
		fmt.Printf("📁 Output directory: %s\n", outputDir)

		builder, err := build.NewBuilder(projectDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building project: %v\n", err)
			os.Exit(1)
		}

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
