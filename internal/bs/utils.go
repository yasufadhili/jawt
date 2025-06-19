package bs

import (
	"os"
)

func getCurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}
