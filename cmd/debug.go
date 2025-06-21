package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var debugPort int

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Start JAWT debugger",
	Long: `Starts the JAWT debugger, which provides debugging tools and insights
into your application's compilation and runtime behaviour.`,
	Run: func(cmd *cobra.Command, args []string) {
		// FUTURE
		fmt.Printf("⚠️ The 'debug' command is not yet implemented.\n")
		fmt.Printf("Please check future versions for this functionality.\n")
	},
}

func init() {
	debugCmd.Flags().IntVarP(&debugPort, "port", "p", 6501, "Specify custom port")
	rootCmd.AddCommand(debugCmd)
}
