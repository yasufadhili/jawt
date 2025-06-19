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

	fmt.Println(wDir)

	return nil
}
