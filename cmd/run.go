package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yasufadhili/jawt/internal/build"
)

var port int
var clearCache bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start development server with hot reload",
	Long: `Starts the development server with hot reload functionality.
Monitors your JML files for changes and automatically reloads the browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Use the configured port if not explicitly provided as a flag
		if !cmd.Flags().Changed("port") && projectConfig != nil {
			port = projectConfig.Jawt.Port
		}

		fmt.Printf("🚀 Starting development server for %s on port %d...\n",
			projectConfig.App.Name, port)

		builder, err := build.NewBuilder(projectDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building project: %v\n", err)
			os.Exit(1)
		}

		if clearCache {
			fmt.Println("🧹 Clearing cache...")
			// TODO: Implement cache clearing in builder
		}

		err = builder.RunDev()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	runCmd.Flags().IntVarP(&port, "port", "p", 6500, "Specify custom port")
	runCmd.Flags().BoolVarP(&clearCache, "clear-cache", "c", false, "Run with cleared cache")
}
