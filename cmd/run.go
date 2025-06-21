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
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		builder := build.NewBuilder(dir)

		if clearCache {
			fmt.Println("🧹 Clearing cache...")
			// TODO: Implement cache clearing in builder
		}

		fmt.Printf("🚀 Starting development server on port %d...\n", port)

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
