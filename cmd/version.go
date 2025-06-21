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
	// TODO: Implement proper version handling
	if version == "dev" {
		return "v0.0.1-dev"
	}
	return version
}
