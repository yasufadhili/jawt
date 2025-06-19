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

	return nil
}
