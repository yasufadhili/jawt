package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

type InitCommand struct {
	ProjectName string
}

type RunCommand struct {
	Port       int
	ClearCache bool
}

type BuildCommand struct{}

type VersionCommand struct{}

type HelpCommand struct{}

type ParsedCommand interface{}

func ParseArgs() (ParsedCommand, error) {
	args := os.Args

	if len(args) < 1 {
		return nil, fmt.Errorf("no program name provided")
	}

	programName := args[0]

	// If no command provided, show help
	if len(args) < 2 {
		return HelpCommand{}, nil
	}

	command := args[1]
	commandArgs := args[2:]

	switch strings.ToLower(command) {
	case "init":
		return parseInitCommand(commandArgs)
	case "run":
		return parseRunCommand(commandArgs)
	case "build":
		return parseBuildCommand(commandArgs)
	case "version":
		return parseVersionCommand(commandArgs)
	case "help":
		return parseHelpCommand(commandArgs)
	default:
		return nil, fmt.Errorf("unknown command: %s\n\nUse '%s help' for usage information", command, programName)
	}
}

func parseInitCommand(args []string) (InitCommand, error) {
	if len(args) == 0 {
		return InitCommand{}, fmt.Errorf("init command requires a project name or '.' for current directory")
	}

	if len(args) > 1 {
		return InitCommand{}, fmt.Errorf("init command takes only one argument (project name or '.')")
	}

	projectName := args[0]
	if projectName != "." && strings.Contains(projectName, "/") {
		return InitCommand{}, fmt.Errorf("invalid project name: %s", projectName)
	}

	return InitCommand{
		ProjectName: projectName,
	}, nil
}

func parseRunCommand(args []string) (RunCommand, error) {
	cmd := RunCommand{
		Port:       6500, // default port
		ClearCache: false,
	}

	i := 0
	for i < len(args) {
		arg := args[i]

		switch arg {
		case "-c", "--clear-cache":
			cmd.ClearCache = true
		case "-p", "--port":
			if i+1 >= len(args) {
				return cmd, fmt.Errorf("port flag requires a value")
			}
			i++
			port, err := strconv.Atoi(args[i])
			if err != nil {
				return cmd, fmt.Errorf("invalid port number: %s", args[i])
			}
			if port < 1 || port > 65535 {
				return cmd, fmt.Errorf("port number must be between 1 and 65535")
			}
			cmd.Port = port
		default:
			// Try to parse as port number if it's a number
			if port, err := strconv.Atoi(arg); err == nil {
				if port < 1 || port > 65535 {
					return cmd, fmt.Errorf("port number must be between 1 and 65535")
				}
				cmd.Port = port
			} else {
				return cmd, fmt.Errorf("unknown argument: %s", arg)
			}
		}
		i++
	}

	return cmd, nil
}

func parseBuildCommand(args []string) (BuildCommand, error) {
	if len(args) > 0 {
		return BuildCommand{}, fmt.Errorf("build command does not accept any arguments")
	}

	return BuildCommand{}, nil
}

func parseVersionCommand(args []string) (VersionCommand, error) {
	if len(args) > 0 {
		return VersionCommand{}, fmt.Errorf("version command does not accept any arguments")
	}

	return VersionCommand{}, nil
}

func parseHelpCommand(args []string) (HelpCommand, error) {
	if len(args) > 0 {
		return HelpCommand{}, fmt.Errorf("help command does not accept any arguments")
	}

	return HelpCommand{}, nil
}

// GetProgramName returns the program name from os.Args
func GetProgramName() string {
	if len(os.Args) > 0 {
		return os.Args[0]
	}
	return "program"
}

// PrintUsage prints the usage information
func PrintUsage() {
	programName := GetProgramName()
	fmt.Printf(`Usage: %s <command> [options]

Commands:
  init <name>    Initialize a project with the given name or '.' for current directory
  run [port]     Run the project (optionally specify port, default: 8080)
                 Options:
                   -c, --clear-cache    Clear cache before running
                   -p, --port <port>    Specify port number
  build          Build the project
  version        Display version information
  help           Display this help message

Examples:
  %s init myproject
  %s init .
  %s run
  %s run 3000
  %s run -c
  %s run -p 3000 -c
  %s build
  %s version
  %s help
`, programName, programName, programName, programName, programName, programName, programName, programName, programName, programName)
}
