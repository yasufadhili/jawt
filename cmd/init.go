package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yasufadhili/jawt/internal/build"
	"github.com/yasufadhili/jawt/internal/core"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Create a new Jawt project",
	Long: `Creates a new Jawt project with the default structure and configuration files.
This includes app/, components/, assets/ directories and essential configuration files.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		logger := core.NewDefaultLogger(core.InfoLevel)

		cfg, err := core.LoadJawtConfig("")
		if err != nil {
			logger.Error("Failed to load JAWT configuration", core.ErrorField(err))
			os.Exit(1)
		}

		if err := cfg.Validate(); err != nil {
			logger.Error("Invalid JAWT configuration", core.ErrorField(err))
			os.Exit(1)
		}

		projectConfig := core.DefaultProjectConfig()
		projectConfig.SetProjectName(projectName)

		// Determine project directory (use current directory or create new one based on project name)
		projectDir, err := filepath.Abs(projectName)
		if err != nil {
			logger.Error("Failed to determine project directory", core.ErrorField(err))
			os.Exit(1)
		}

		paths, err := core.NewProjectPaths(projectDir, projectConfig, cfg)
		if err != nil {
			logger.Error("Failed to initialise project paths", core.ErrorField(err))
			os.Exit(1)
		}

		ctx := core.NewJawtContext(cfg, projectConfig, paths, logger, nil)

		err = build.InitProject(ctx, projectName, projectDir)
		if err != nil {
			logger.Error("Failed to initialise project", core.ErrorField(err))
			os.Exit(1)
		}

		// Ensure all necessary directories are created
		if err := paths.EnsureDirectories(); err != nil {
			logger.Error("Failed to create project directories", core.ErrorField(err))
			os.Exit(1)
		}

		// Save project configuration
		if err := projectConfig.Save(projectDir); err != nil {
			logger.Error("Failed to save project configuration", core.ErrorField(err))
			os.Exit(1)
		}

		logger.Info("JAWT project initialised successfully",
			core.StringField("name", projectName),
			core.StringField("path", projectDir))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
