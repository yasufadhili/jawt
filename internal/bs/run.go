package bs

import (
	"fmt"
	"os"
)

func RunProject(clearCache bool) error {

	wDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if !fileExists("jawt.config.json") {
		return fmt.Errorf("jawt.config.json not found in %s", wDir)
	}

	if !fileExists("app.json") {
		return fmt.Errorf("app.json not found in %s", wDir)
	}

	name, err := readJsonField("app.json", "name")
	if err != nil {
		return err
	}

	// TODO: Implement Build Process
	// Parse all page and component names
	// Resolve dependencies
	// Compile each page or component accordingly

	fmt.Printf("Running project '%s'\n", name)
	fmt.Printf("Visit localhost:%d to view the project\n", 6500) // TODO: Allow using port passed by user

	return nil
}
