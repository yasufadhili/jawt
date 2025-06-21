package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve production build locally",
	Long: `Serves the production build locally for previewing how your 
application will behave in production. This command requires you to 
run 'jawt build' first.`,
	Run: func(cmd *cobra.Command, args []string) {
		// FUTURE
		fmt.Println("⚠️ The 'serve' command is not yet implemented.")
		fmt.Println("Please check future versions for this functionality.")
	},
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8000, "Specify custom port")
	rootCmd.AddCommand(serveCmd)
}
