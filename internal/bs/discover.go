package bs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func discoverPages(rootPath string) ([]string, error) {

	var pages []string

	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of root: %v", err)
	}

	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "index.jml" {
			// Get the directory path by removing the filename
			dirPath := filepath.Dir(path)
			// Convert to a relative path from the root directory
			relPath, err := filepath.Rel(absRoot, dirPath)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", path, err)
			}
			// Convert to forward-slash notation and ensure no trailing slash
			normalisedPath := strings.TrimSuffix(filepath.ToSlash(relPath), "/")
			// Append the path (or "." if the file is directly in the root)
			if normalisedPath == "." {
				normalisedPath = "/"
			}

			normalisedPath = strings.Replace(normalisedPath, "app", "/", 1)
			normalisedPath = strings.Replace(normalisedPath, "//", "/", 1)

			pages = append(pages, normalisedPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", rootPath, err)
	}

	return pages, nil
}
