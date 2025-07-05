package build

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/yasufadhili/jawt/internal/core"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// InitProject initialises a new JAWT project
func InitProject(ctx *core.JawtContext, projectName string, targetDir string) error {
	// Validate and sanitise project name
	sanitisedName, err := validateAndSanitiseProjectName(projectName)
	if err != nil {
		return fmt.Errorf("invalid project name: %w", err)
	}

	// Determine project directory
	projectDir, err := determineProjectDirectory(targetDir, sanitisedName)
	if err != nil {
		return fmt.Errorf("failed to determine project directory: %w", err)
	}

	// Check if the directory exists and handle accordingly
	if err := handleExistingDirectory(projectDir, sanitisedName); err != nil {
		return err
	}

	ctx.Logger.Info("Initialising JAWT project",
		core.Field{Key: "name", Value: sanitisedName},
		core.Field{Key: "directory", Value: projectDir})

	if err := createProjectStructure(projectDir); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	if err := generateConfigFiles(projectDir, sanitisedName); err != nil {
		return fmt.Errorf("failed to generate configuration files: %w", err)
	}

	if err := createInitialFiles(projectDir, sanitisedName); err != nil {
		return fmt.Errorf("failed to create initial files: %w", err)
	}

	// ctx.Logger.Info("Project initialised successfully", core.Field{Key: "path", Value: projectDir})

	printSuccessMessage(sanitisedName, projectDir)

	return nil
}

//go:embed templates/*
var templateFS embed.FS

// TemplateData holds the data for template rendering
type TemplateData struct {
	ProjectName string
}

// renderTemplate renders a template with the given data
func renderTemplate(templatePath string, data TemplateData) (string, error) {
	content, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// writeTemplateFile renders a template and writes it to the target file
func writeTemplateFile(templatePath, targetPath string, data TemplateData) error {
	content, err := renderTemplate(templatePath, data)
	if err != nil {
		return err
	}

	return os.WriteFile(targetPath, []byte(content), 0644)
}

// createInitialFiles creates the initial project files using embedded templates
func createInitialFiles(projectDir, projectName string) error {
	data := TemplateData{
		ProjectName: projectName,
	}

	if err := writeTemplateFile("templates/app/index.jml.tmpl",
		filepath.Join(projectDir, "app", "index.jml"), data); err != nil {
		return fmt.Errorf("failed to create home page: %w", err)
	}

	if err := writeTemplateFile("templates/components/layout.jml.tmpl",
		filepath.Join(projectDir, "components", "layout.jml"), data); err != nil {
		return fmt.Errorf("failed to create layout component: %w", err)
	}

	if err := writeTemplateFile("templates/scripts/main.ts.tmpl",
		filepath.Join(projectDir, "scripts", "main.ts"), data); err != nil {
		return fmt.Errorf("failed to create main script: %w", err)
	}

	if err := writeTemplateFile("templates/.gitignore.tmpl",
		filepath.Join(projectDir, ".gitignore"), data); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// if err := writeTemplateFile("templates/README.md.tmpl",
	// filepath.Join(projectDir, "README.md"), data); err != nil {
	//	return fmt.Errorf("failed to create README.md: %w", err)
	//}

	return nil
}

// createTSConfig creates a TypeScript configuration file using embedded template
func createTSConfig(projectDir string) error {
	content, err := templateFS.ReadFile("templates/tsconfig.json.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read TypeScript config template: %w", err)
	}

	configPath := filepath.Join(projectDir, "tsconfig.json")
	return os.WriteFile(configPath, content, 0644)
}

// createTailwindConfig creates a Tailwind CSS configuration file using embedded template
func createTailwindConfig(projectDir string) error {
	content, err := templateFS.ReadFile("templates/tailwind.config.js.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read Tailwind config template: %w", err)
	}

	configPath := filepath.Join(projectDir, "tailwind.config.js")
	return os.WriteFile(configPath, content, 0644)
}

// validateAndSanitiseProjectName validates and sanitises the project name
func validateAndSanitiseProjectName(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("project name cannot be empty")
	}

	// Remove leading/trailing whitespace
	name = strings.TrimSpace(name)

	// Convert to lowercase and replace spaces/special chars with hyphens
	name = strings.ToLower(name)
	name = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")

	// Validate final name
	if len(name) == 0 {
		return "", fmt.Errorf("project name contains no valid characters")
	}

	if len(name) > 214 {
		return "", fmt.Errorf("project name too long (max 214 characters)")
	}

	// Check for reserved names
	reservedNames := []string{
		"jawt", "node", "npm", "test", "src", "build", "dist", "public",
		"static", "assets", "components", "scripts", "styles",
	}
	for _, reserved := range reservedNames {
		if name == reserved {
			return "", fmt.Errorf("'%s' is a reserved name", reserved)
		}
	}

	// Must start with a letter
	if !regexp.MustCompile(`^[a-z]`).MatchString(name) {
		return "", fmt.Errorf("project name must start with a letter")
	}

	return name, nil
}

// determineProjectDirectory determines the target directory for the project
func determineProjectDirectory(targetDir, projectName string) (string, error) {
	var projectDir string

	if targetDir == "." {
		// Initialise in current directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		projectDir = cwd
	} else {
		// Create new directory
		if filepath.IsAbs(targetDir) {
			projectDir = targetDir
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return "", fmt.Errorf("failed to get current directory: %w", err)
			}
			projectDir = filepath.Join(cwd, targetDir)
		}
	}

	return projectDir, nil
}

// handleExistingDirectory checks if directory exists and handles accordingly
func handleExistingDirectory(projectDir, projectName string) error {
	info, err := os.Stat(projectDir)
	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to check project directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("target path exists but is not a directory")
	}

	// Check if the directory is empty
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return fmt.Errorf("failed to read project directory: %w", err)
	}

	if len(entries) > 0 {
		// Directory is not empty, check for an existing JAWT project
		jawtConfigPath := filepath.Join(projectDir, "jawt.project.json")
		if _, err := os.Stat(jawtConfigPath); err == nil {
			return fmt.Errorf("directory already contains a JAWT project")
		}

		// Allow initialisation in non-empty directory if no conflicting files
		conflictingFiles := []string{
			"jawt.project.json", "tsconfig.json", "tailwind.config.js",
			"app", "components", "scripts",
		}

		for _, file := range conflictingFiles {
			if _, err := os.Stat(filepath.Join(projectDir, file)); err == nil {
				return fmt.Errorf("directory contains conflicting file/directory: %s", file)
			}
		}
	}

	return nil
}

// createProjectStructure creates the basic directory structure
func createProjectStructure(projectDir string) error {
	directories := []string{
		"app",
		"components",
		"scripts",
		"assets",
	}

	for _, dir := range directories {
		fullPath := filepath.Join(projectDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// ProjectConfig represents the new project configuration structure
type ProjectConfig struct {
	App struct {
		Name   string `json:"name"`
		Author string `json:"author"`
	} `json:"app"`
	Components struct {
		Path  string `json:"path"`
		Alias string `json:"alias"`
	} `json:"components"`
	Pages struct {
		Path  string `json:"path"`
		Alias string `json:"alias"`
	} `json:"pages"`
	Scripts struct {
		Path  string `json:"path"`
		Alias string `json:"alias"`
	} `json:"scripts"`
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	Build struct {
		OutputDir string `json:"outputDir"`
		Minify    bool   `json:"minify"`
	} `json:"build"`
}

// generateConfigFiles creates the configuration files
func generateConfigFiles(projectDir, projectName string) error {
	// Create jawt.project.json
	projectConfig := ProjectConfig{
		App: struct {
			Name   string `json:"name"`
			Author string `json:"author"`
		}{
			Name:   projectName,
			Author: "",
		},
		Components: struct {
			Path  string `json:"path"`
			Alias string `json:"alias"`
		}{
			Path:  "components",
			Alias: "",
		},
		Pages: struct {
			Path  string `json:"path"`
			Alias string `json:"alias"`
		}{
			Path:  "app",
			Alias: "",
		},
		Scripts: struct {
			Path  string `json:"path"`
			Alias string `json:"alias"`
		}{
			Path:  "scripts",
			Alias: "",
		},
		Server: struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}{
			Host: "localhost",
			Port: 6500,
		},
		Build: struct {
			OutputDir string `json:"outputDir"`
			Minify    bool   `json:"minify"`
		}{
			OutputDir: "build",
			Minify:    true,
		},
	}

	if err := saveProjectConfig(projectDir, &projectConfig); err != nil {
		return fmt.Errorf("failed to save project config: %w", err)
	}

	// if err := createTSConfig(projectDir); err != nil {
	//	return fmt.Errorf("failed to create TypeScript config: %w", err)
	// }

	// if err:= createTailwindConfig(projectDir); err != nil {
	// 	return fmt.Errorf("failed to create Tailwind config: %w", err)
	// }

	return nil
}

// saveProjectConfig saves the project configuration
func saveProjectConfig(projectDir string, config *ProjectConfig) error {
	configPath := filepath.Join(projectDir, "jawt.project.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// printSuccessMessage prints a success message with next steps
func printSuccessMessage(projectName, projectDir string) {
	fmt.Printf(`
🎉 Successfully created JAWT project '%s'!

📂 Project created at: %s

🚀 Next steps:
   cd %s
   jawt run

📚 Learn more:
   • Edit app/index.jml to modify your home page
   • Create components in components/
   • Add TypeScript functionality in scripts/
   • Run 'jawt --help' for all available commands

Happy coding! 🎯
`, projectName, projectDir, filepath.Base(projectDir))
}
