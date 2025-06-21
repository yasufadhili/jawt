package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of JAWT",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("JAWT v0.0.1")
		// TODO: Implement proper version handling
	},
}
