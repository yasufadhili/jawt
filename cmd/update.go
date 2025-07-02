package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Jawt and all dependencies to the latest stable version",
	Long:  `Downloads and installs Jawt, and it's dependencies to use with the toolchain.`,
	Run: func(cmd *cobra.Command, args []string) {

		// FUTURE
		fmt.Println("⚠️ The 'serve' command is not yet implemented.")
		fmt.Println("Please check future versions for this functionality.")

	},
}
