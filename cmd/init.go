package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yasufadhili/jawt/internal/build"
	"github.com/yasufadhili/jawt/internal/core"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Create a new Jawt project",
	Long: `Creates a new Jawt project with the default structure and configuration files.
This includes app/, components/, assets/ directories and essential configuration files.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		projectName := args[0]

		cfg, err := core.LoadJawtConfig("")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		ctx := core.NewJawtContext(cfg, nil, nil, nil)

		err = build.InitProject(ctx, projectName, projectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initialising project: %v\n", err)
			os.Exit(1)
		}

	},
}
