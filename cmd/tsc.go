package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var tscCmd = &cobra.Command{
	Use:   "tsc [typescript-args...]",
	Short: "Run TypeScript compiler",
	Long: `Run the TypeScript compiler within Jawt with the provided arguments.
All arguments are passed directly to tsc.`,
	Run: func(cmd *cobra.Command, args []string) {
		// FUTURE
		fmt.Println("⚠️ The 'tsc' command is not yet implemented.")
		fmt.Println("Please check future versions for this functionality.")
	},
}
