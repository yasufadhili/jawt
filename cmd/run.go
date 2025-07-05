package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yasufadhili/jawt/internal/execute"
)

var port int
var clearCache bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start development server with hot reload",
	Long: `Starts the development server with hot reload functionality.
Monitors your JML files for changes and automatically reloads the browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := execute.RunDev()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	runCmd.Flags().IntVarP(&port, "port", "p", 6500, "Specify custom port")
	runCmd.Flags().BoolVarP(&clearCache, "clear-cache", "c", false, "Run with cleared cache")
}
