package bs

import (
	"errors"
	"fmt"
	"github.com/yasufadhili/jawt/internal/config"
	"os"
	"path/filepath"
)

func InitProject(projectName string) error {

	if projectName == "." {
		dir, err := getCurrentDir()
		if err != nil {
			_, _ = os.Stderr.WriteString("Error retrieving current directory: " + err.Error() + "\n")
			os.Exit(1)
		}
		_, _ = os.Stdout.WriteString("Current working directory: " + dir + "\n")
		_ = fmt.Errorf("not implemented yet")
	} else {
		err := validateFolderName(projectName)
		if err != nil {
			return err
		}

		err = createDirStructure(projectName)
		if err != nil {
			return err
		}

		err = createConfigFiles(projectName)
		if err != nil {
			return err
		}

		fmt.Printf("Project '%s' initialised\n\n", projectName)
		fmt.Printf("Run 'cd %s' to enter the project directory\n", projectName)
		fmt.Printf("Then run 'jawt run' to start the project\n")

	}

	return nil
}

func createDirStructure(parent string) error {

	subDirs := []string{"app", "components", "assets"}

	if _, err := os.Stat(parent); os.IsNotExist(err) {
		err := os.Mkdir(parent, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		// An error occurred other than "not exists"
		return err
	} else {
		return errors.New("directory already exists '" + parent + "'")
	}

	// create each subdirectory
	for _, subDir := range subDirs {
		fullPath := filepath.Join(parent, subDir)

		err := os.MkdirAll(fullPath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func createConfigFiles(dir string) error {
	err := createAppJsonFile(dir)
	if err != nil {
		return err
	}
	e := createJawtJsonFile(dir)
	if e != nil {
		return e
	}
	return nil
}

func createAppJsonFile(name string) error {

	currUser, err := getCurrentUserName()
	if err != nil {
		return err
	}

	appConfig := config.AppConfig{
		Name:        name,
		Author:      currUser,
		Version:     "1.0.0",
		Description: "Cool Jawt project",
	}

	err = createJsonFile(name, "app.json", appConfig)
	if err != nil {
		return err
	}

	return nil
}

func createJawtJsonFile(name string) error {

	projConfig := config.ProjectConfig{
		Name: name,
	}

	serverConfig := config.ServerConfig{
		Port: 6500,
	}

	jawtConfig := config.JawtConfig{
		Project: projConfig,
		Server:  serverConfig,
	}

	err := createJsonFile(name, "jawt.config.json", jawtConfig)
	if err != nil {
		return err
	}

	return nil
}
