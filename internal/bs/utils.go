package bs

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
)

func getCurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// validateFolderName checks if the provided folder or project name is valid.
// It returns an error if the name is invalid, nil otherwise.
func validateFolderName(name string) error {
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

func createJsonFile(parentDir string, fileName string, data interface{}) error {
	// Ensure the parent directory exists
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", parentDir, err)
	}

	filePath := filepath.Join(parentDir, fileName)

	// Marshal the data to JSON with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	// Create or overwrite the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	// Write JSON data to the file
	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write to file %s: %v", filePath, err)
	}

	return nil
}

func getCurrentUserName() (string, error) {
	currUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %v", err)
	}
	return currUser.Username, nil
}

func createFile(parentDir string, fileName string, data []byte) error {

	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", parentDir, err)
	}

	filePath := filepath.Join(parentDir, fileName)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func readJsonField(filename string, field string) (interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	file, err = os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON into a map
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	value, exists := result[field]
	if !exists {
		return nil, fmt.Errorf("field %s not found in %s", field, filename)
	}
	return value, nil
}
