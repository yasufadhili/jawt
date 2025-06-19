package main

import (
	"fmt"
	cmd "github.com/yasufadhili/jawt/cmd"
	"github.com/yasufadhili/jawt/internal/build"
	"os"
)

func main() {
	command, err := cmd.ParseArgs()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		cmd.PrintUsage()
		os.Exit(1)
	}

	switch c := command.(type) {

	case cmd.InitCommand:

		currentDir, err := os.Getwd()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		targetPath := currentDir

		err = build.InitProject(targetPath, c.ProjectName)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error initialising project: %v\n", err)
			os.Exit(1)
		}

	case cmd.RunCommand:

		dir, err := os.Getwd()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		builder := build.NewBuilder(dir)

		if c.ClearCache {
			fmt.Println("🧹 Clearing cache...")
			// TODO: Implement cache clearing in builder
		}

		err = builder.RunDev()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case cmd.BuildCommand:
		dir, err := os.Getwd()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		builder := build.NewBuilder(dir)
		err = builder.Build()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error building project: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Build completed successfully!")

	case cmd.VersionCommand:
		fmt.Println("JAWT version 0.1.0")
		fmt.Println("A minimal web application builder") // TODO: Implement proper version handling
		// TODO Call version implementation

	case cmd.HelpCommand:
		cmd.PrintUsage()

	default:
		fmt.Printf("Unknown command: %T\n", c)
		cmd.PrintUsage()
		os.Exit(1)

	}

}
