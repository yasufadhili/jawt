package build

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
)

// ProjectInitializer handles creating new JAWT projects
type ProjectInitializer struct {
	targetPath string
	config     *InitConfig
}

// InitConfig holds configuration for project initialization
type InitConfig struct {
	ProjectName string `json:"name"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Template    string `json:"template"`
	Port        int    `json:"port"`
}

// InitProject Package-level convenience function for external use
func InitProject(targetPath, projectName string) error {
	initializer := NewProjectInitializer(targetPath)
	return initializer.InitProject(projectName)
}

// NewProjectInitializer creates a new project initializer
func NewProjectInitializer(targetPath string) *ProjectInitializer {
	return &ProjectInitializer{
		targetPath: targetPath,
	}
}

// InitProject initializes a new JAWT project
func (pi *ProjectInitializer) InitProject(projectName string) error {
	// Validate and prepare configuration
	config, err := pi.prepareConfig(projectName)
	if err != nil {
		return fmt.Errorf("failed to prepare configuration: %w", err)
	}

	pi.config = config

	// Create directory structure
	if err := pi.createDirectoryStructure(); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Generate configuration files
	if err := pi.generateConfigFiles(); err != nil {
		return fmt.Errorf("failed to generate configuration files: %w", err)
	}

	// Generate template files
	if err := pi.generateTemplateFiles(); err != nil {
		return fmt.Errorf("failed to generate template files: %w", err)
	}

	// Verify project structure
	if err := pi.verifyProject(); err != nil {
		return fmt.Errorf("project verification failed: %w", err)
	}

	pi.printSuccessMessage()
	return nil
}

// prepareConfig validates input and prepares initialization configuration
func (pi *ProjectInitializer) prepareConfig(projectName string) (*InitConfig, error) {
	var actualProjectName string
	//var targetDir string

	if projectName == "." {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}

		actualProjectName = filepath.Base(currentDir)
		//targetDir = currentDir
	} else {
		if err := validateProjectName(projectName); err != nil {
			return nil, err
		}

		actualProjectName = projectName
		//targetDir = filepath.Join(pi.targetPath, projectName)
	}

	currentUser, err := getCurrentUser()
	if err != nil {
		currentUser = "Unknown"
	}

	return &InitConfig{
		ProjectName: actualProjectName,
		Author:      currentUser,
		Version:     "1.0.0",
		Description: fmt.Sprintf("A JAWT project called %s", actualProjectName),
		Template:    "default",
		Port:        6500,
	}, nil
}

// createDirectoryStructure creates the project directory structure
func (pi *ProjectInitializer) createDirectoryStructure() error {
	// Determine target directory

	d, e := os.Getwd()
	if e != nil {
		return fmt.Errorf("failed to get current directory: %w", e)
	}
	fmt.Println(d)

	var targetDir string
	if pi.config.ProjectName == filepath.Base(pi.targetPath) {
		// Initializing in the current directory
		targetDir = pi.targetPath
	} else {
		// Creating a new directory
		targetDir = filepath.Join(pi.targetPath, pi.config.ProjectName)
	}

	// Check if the directory already exists and is not empty
	if info, err := os.Stat(targetDir); err == nil {
		if info.IsDir() {
			// Check if the directory is empty
			entries, err := os.ReadDir(targetDir)
			if err != nil {
				return fmt.Errorf("failed to read directory %s: %w", targetDir, err)
			}

			if len(entries) > 0 {
				return fmt.Errorf("directory %s already exists and is not empty", targetDir)
			}
		} else {
			return fmt.Errorf("path %s exists but is not a directory", targetDir)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check directory %s: %w", targetDir, err)
	}

	// Create directory structure
	directories := []string{
		"app",
		"components",
		"assets",
		".dist", // Build output directory
	}

	// Create the main directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create subdirectories
	for _, dir := range directories {
		fullPath := filepath.Join(targetDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
		}
	}

	// Update a target path for later operations
	pi.targetPath = targetDir
	return nil
}

// generateConfigFiles creates the project configuration files
func (pi *ProjectInitializer) generateConfigFiles() error {
	// Generate app.json
	appConfig := map[string]interface{}{
		"name":        pi.config.ProjectName,
		"author":      pi.config.Author,
		"version":     pi.config.Version,
		"description": pi.config.Description,
	}

	if err := pi.writeJSONFile("app.json", appConfig); err != nil {
		return fmt.Errorf("failed to create app.json: %w", err)
	}

	// Generate jawt.config.json
	jawtConfig := map[string]interface{}{
		"project": map[string]interface{}{
			"name": pi.config.ProjectName,
		},
		"server": map[string]interface{}{
			"port": pi.config.Port,
		},
		"build": map[string]interface{}{
			"output": "dist",
			"minify": true,
		},
	}

	if err := pi.writeJSONFile("jawt.config.json", jawtConfig); err != nil {
		return fmt.Errorf("failed to create jawt.config.json: %w", err)
	}

	// Generate .gitignore
	gitignoreContent := `# Build output
.dist/

# Development files
.jawt-cache/
*.log

# OS files
.DS_Store
Thumbs.db

# Editor files
.vscode/
.idea/
*.swp
*.swo
*~
`

	if err := pi.writeTextFile(".gitignore", gitignoreContent); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}

func (pi *ProjectInitializer) generateTemplateFiles() error {

	// Generate app/index.jml
	indexContent := fmt.Sprintf(`_doctype page index

import Layout from "components/layout"

Page {
  title: "%s"
  description: "Welcome to %s - Built with JAWT"
  
  Layout {
    content: "Welcome to your new JAWT project!"
  }
}
`, pi.config.ProjectName, pi.config.ProjectName)

	if err := pi.writeTextFile("app/index.jml", indexContent); err != nil {
		return fmt.Errorf("failed to create app/index.jml: %w", err)
	}

	// Generate components/layout.jml
	layoutContent := `_doctype component layout

Container {
  style: "min-h-screen bg-gray-50 flex flex-col"
  
  Header {
    style: "bg-white shadow-sm border-b px-6 py-4"
    
    Title {
      style: "text-2xl font-bold text-gray-900"
      text: "JAWT Application"
    }
  }
  
  Main {
    style: "flex-1 container mx-auto px-6 py-8"
    
    Section {
      style: "max-w-2xl mx-auto text-center"
      
      Heading {
        style: "text-4xl font-bold text-gray-900 mb-4"
        text: props.content || "Hello, JAWT!"
      }
      
      Paragraph {
        style: "text-lg text-gray-600 mb-8"
        text: "Start building your web application with JAWT's declarative approach."
      }
      
      Card {
        style: "bg-white rounded-lg shadow-md p-6 text-left"
        
        CardTitle {
          style: "text-xl font-semibold text-gray-800 mb-3"
          text: "Getting Started"
        }
        
        List {
          style: "space-y-2 text-gray-600"
          
          ListItem {
            text: "Edit app/index.jml to modify this page"
          }
          ListItem {
            text: "Create new components in the components/ directory"
          }
          ListItem {
            text: "Add assets to the assets/ directory"
          }
          ListItem {
            text: "Run 'jawt build' to build for production"
          }
        }
      }
    }
  }
  
  Footer {
    style: "bg-white border-t px-6 py-4 text-center text-gray-500"
    text: "Built with JAWT"
  }
}
`

	if err := pi.writeTextFile("components/layout.jml", layoutContent); err != nil {
		return fmt.Errorf("failed to create components/layout.jml: %w", err)
	}

	// Generate a sample favicon
	faviconSVG := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="32" height="32">
  <rect width="32" height="32" fill="#3B82F6"/>
  <text x="16" y="22" font-family="Arial, sans-serif" font-size="18" font-weight="bold" text-anchor="middle" fill="white">J</text>
</svg>`

	if err := pi.writeTextFile("assets/favicon.svg", faviconSVG); err != nil {
		return fmt.Errorf("failed to create assets/favicon.svg: %w", err)
	}

	//readmeContent := fmt.Sprintf(`# JAWT Application`, pi.config.ProjectName, pi.config.Description, pi.config.Port, pi.config.Port, pi.config.ProjectName)
	readmeContent := fmt.Sprintf("")

	if err := pi.writeTextFile("README.md", readmeContent); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	return nil
}

// verifyProject checks that the project was created correctly
func (pi *ProjectInitializer) verifyProject() error {
	// Use the project discovery system to verify structure
	discovery := NewProjectDiscovery(pi.targetPath)
	project, err := discovery.DiscoverProject()
	if err != nil {
		return fmt.Errorf("project verification failed: %w", err)
	}

	// Check that essential files exist
	requiredFiles := []string{
		"app.json",
		"jawt.config.json",
		"app/index.jml",
		"components/layout.jml",
	}

	for _, file := range requiredFiles {
		fullPath := filepath.Join(pi.targetPath, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("required file missing: %s", file)
		}
	}

	// Verify the project structure matches expected
	if len(project.Pages) == 0 {
		return fmt.Errorf("no pages discovered in project")
	}

	if len(project.Components) == 0 {
		return fmt.Errorf("no components discovered in project")
	}

	return nil
}

// printSuccessMessage displays success information
func (pi *ProjectInitializer) printSuccessMessage() {
	fmt.Printf("✅ Project '%s' initialised successfully!\n\n", pi.config.ProjectName)

	if pi.config.ProjectName != filepath.Base(pi.targetPath) {
		fmt.Printf("📁 Run 'cd %s' to enter the project directory\n", pi.config.ProjectName)
	}

	fmt.Printf("🚀 Run 'jawt run' to start the development server\n")
	fmt.Printf("🏗️  Run 'jawt build' to build for production\n")
	fmt.Printf("📖 Check README.md for more information\n\n")

	fmt.Printf("🎉 Happy coding with JAWT!\n")
}

// writeJSONFile writes data as JSON to a file
func (pi *ProjectInitializer) writeJSONFile(filename string, data interface{}) error {
	fullPath := filepath.Join(pi.targetPath, filename)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(fullPath, jsonData, 0644)
}

// writeTextFile writes text content to a file
func (pi *ProjectInitializer) writeTextFile(filename, content string) error {
	fullPath := filepath.Join(pi.targetPath, filename)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.WriteFile(fullPath, []byte(content), 0644)
}

// validateFolderName checks if the provided folder or project name is valid.
// It returns an error if the name is invalid, nil otherwise.
func validateProjectName(name string) error {
	// Regular expression for a valid folder / project name:
	// - Starts with a letter,
	// - Contains only letters, numbers, underscores,
	// - Maximum length of 255 characters
	// - Does not contain reserved characters or sequences
	const maxLength = 255
	validNamePattern := `^[a-zA-Z][a-zA-Z0-9_]{0,254}$`

	if name == "" {
		return fmt.Errorf("folder name cannot be empty")
	}

	if len(name) > maxLength {
		return fmt.Errorf("folder name exceeds maximum length of %d characters", maxLength)
	}

	reservedNames := []string{".", "..", "CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	for _, reserved := range reservedNames {
		if name == reserved {
			return fmt.Errorf("folder name '%s' is a reserved name", name)
		}
	}

	// Check for invalid characters using regex
	matched, err := regexp.MatchString(validNamePattern, name)
	if err != nil {
		return fmt.Errorf("error validating folder name: %v", err)
	}
	if !matched {
		return fmt.Errorf("folder name contains invalid characters or format")
	}

	// Check for consecutive dots or hyphens
	if regexp.MustCompile(`[.-]{2,}`).MatchString(name) {
		return fmt.Errorf("folder name cannot contain consecutive dots or hyphens")
	}

	// Check for spaces
	if regexp.MustCompile(`\s`).MatchString(name) {
		return fmt.Errorf("folder name cannot contain spaces")
	}

	return nil
}

// getCurrentUser gets the current system user
func getCurrentUser() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	// Try to get the display name, fall back to username
	if currentUser.Name != "" {
		return currentUser.Name, nil
	}

	return currentUser.Username, nil
}
