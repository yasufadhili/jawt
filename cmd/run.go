package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yasufadhili/jawt/internal/core"
	"github.com/yasufadhili/jawt/internal/runtime"
	"os"
)

var port int
var clearCache bool
var verbose bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start development server with hot reload",
	Long: `Starts the development server with hot reload functionality.
Monitors your JML files for changes and automatically reloads the browser.`,
	Run: func(cmd *cobra.Command, args []string) {

		var logLevel core.LogLevel
		if verbose {
			logLevel = core.DebugLevel
		} else {
			logLevel = core.WarnLevel
		}
		logger := core.NewDefaultLogger(logLevel)

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

		buildOptions := core.NewBuildOptions()
		// Check for tailwind.config.js
		tailwindConfigPath := paths.GetAbsolutePath(projectConfig.TailwindConfigPath)
		if _, err := os.Stat(tailwindConfigPath); err == nil {
			buildOptions.UsesTailwindCSS = true
		} else if !os.IsNotExist(err) {
			logger.Error("Failed to check for tailwind config file", core.ErrorField(err))
			os.Exit(1)
		}

		ctx := core.NewJawtContext(cfg, projectConfig, paths, logger, buildOptions)

		logger.Info("Starting project",
			core.StringField("name", projectConfig.App.Name),
			core.StringField("directory", projectDir),
			core.StringField("server", projectConfig.GetServerAddress()))

		// Create and start the orchestrator
		orchestrator, err := runtime.NewOrchestrator(cmd.Context(), logger, ctx)
		if err != nil {
			logger.Error("Failed to create orchestrator", core.ErrorField(err))
			os.Exit(1)
		}

		if err := orchestrator.StartAll(); err != nil {
			logger.Error("Failed to start orchestrator", core.ErrorField(err))
			os.Exit(1)
		}

		// Wait for interruption
		<-cmd.Context().Done()

		// Stop the orchestrator
		if err := orchestrator.StopAll(); err != nil {
			logger.Error("Failed to stop orchestrator", core.ErrorField(err))
		}

		logger.Info("JAWT project stopped",
			core.StringField("name", projectConfig.App.Name))
	},
}

func init() {
	runCmd.Flags().IntVarP(&port, "port", "p", 6500, "Specify custom port")
	runCmd.Flags().BoolVarP(&clearCache, "clear-cache", "c", false, "Run with cleared cache")
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}
