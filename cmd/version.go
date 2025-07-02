package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of JAWT",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("JAWT %s\n", getVersion())
	},
}

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func getVersion() string {
	return fmt.Sprintf("v0.1.0-dev (commit: %s, built: %s)", commit, date)
}
