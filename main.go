package main

import (
	"fmt"
	cmd "github.com/yasufadhili/jawt/cmd"
	"github.com/yasufadhili/jawt/internal/bs"
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
		fmt.Printf("Initialising project %s\n", c.ProjectName)
		err := bs.InitProject(c.ProjectName)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case cmd.RunCommand:
		if c.ClearCache {
			err := bs.RunProject(true)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}
		err := bs.RunProject(false)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case cmd.BuildCommand:
		fmt.Println("Building project...")
		// TODO: Handle build

	case cmd.VersionCommand:
		fmt.Println("Version 0.1.0") // TODO: Implement proper version handling
	// TODO Call version implementation

	case cmd.HelpCommand:
		cmd.PrintUsage()

	default:
		fmt.Printf("Unknown command: %T\n", c)
		os.Exit(1)

	}

}
