package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yasufadhili/jawt/internal/build"
	"os"
)

var port int
var clearCache bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start development server",
	Long: `Starts the development server with hot reload functionality.
Monitors your Jml files for changes and automatically reloads the browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Use the configured port if not explicitly provided as a flag
		if !cmd.Flags().Changed("port") && projectManager.Project.Config != nil {
			port = int(projectManager.Project.Config.Server.Port)
		}

		builder, err := build.NewBuilder(projectManager.Project)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if clearCache {
			builder.ClearCache = true
		}

		err = builder.Run()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	},
}

func init() {
	runCmd.Flags().IntVarP(&port, "port", "p", 6500, "Specify custom port")
	runCmd.Flags().BoolVarP(&clearCache, "clear-cache", "c", false, "Run with cleared cache")
}
