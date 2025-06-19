package bs

import (
	"fmt"
	"os"
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
