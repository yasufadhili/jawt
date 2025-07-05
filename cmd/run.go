package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yasufadhili/jawt/internal/build"
	"github.com/yasufadhili/jawt/internal/core"
	"os"
)

var port int
var clearCache bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start development server with hot reload",
	Long: `Starts the development server with hot reload functionality.
Monitors your JML files for changes and automatically reloads the browser.`,
	Run: func(cmd *cobra.Command, args []string) {

		logger := core.NewDefaultLogger(core.InfoLevel)

		projectDir, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get current working directory", core.ErrorField(err))
			os.Exit(1)
		}

		cfg, err := core.LoadJawtConfig("")
		if err != nil {
			logger.Error("Failed to load JAWT configuration", core.ErrorField(err))
			os.Exit(1)
		}

		if err := cfg.Validate(); err != nil {
			logger.Error("Invalid Jawt configuration", core.ErrorField(err))
			os.Exit(1)
		}

		projectConfig, err := core.LoadProjectConfig(projectDir)
		if err != nil {
			logger.Error("Failed to load project configuration", core.ErrorField(err))
			os.Exit(1)
		}

		if err := projectConfig.Validate(); err != nil {
			logger.Error("Invalid project configuration", core.ErrorField(err))
			os.Exit(1)
		}

		paths, err := core.NewProjectPaths(projectDir, projectConfig, cfg)
		if err != nil {
			logger.Error("Failed to initialise project paths", core.ErrorField(err))
			os.Exit(1)
		}

		if err := paths.EnsureDirectories(); err != nil {
			logger.Error("Failed to create project directories", core.ErrorField(err))
			os.Exit(1)
		}

		ctx := core.NewJawtContext(cfg, projectConfig, paths, logger)

		logger.Info("Starting project",
			core.StringField("name", projectConfig.App.Name),
			core.StringField("directory", projectDir),
			core.StringField("server", projectConfig.GetServerAddress()))

		err = build.RunProject(ctx)
		if err != nil {
			logger.Error("Failed to run project", core.ErrorField(err))
			os.Exit(1)
		}

		logger.Info("JAWT project running successfully",
			core.StringField("name", projectConfig.App.Name),
			core.StringField("server", projectConfig.GetServerAddress()))

	},
}

func init() {
	runCmd.Flags().IntVarP(&port, "port", "p", 6500, "Specify custom port")
	runCmd.Flags().BoolVarP(&clearCache, "clear-cache", "c", false, "Run with cleared cache")
}
