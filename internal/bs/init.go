package bs

import "os"

func InitProject(projectName string) error {

	if projectName == "." {
		dir, err := getCurrentDir()
		if err != nil {
			_, _ = os.Stderr.WriteString("Error retrieving current directory: " + err.Error() + "\n")
			os.Exit(1)
		}
		_, _ = os.Stdout.WriteString("Current working directory: " + dir + "\n")
	} else {
		err := validateFolderName(projectName)
		if err != nil {
			return err
		}
	}

	return nil
}
