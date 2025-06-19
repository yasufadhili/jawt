package bs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type page struct {
	Name    string
	RelPath string
	AbsPath string
}

func discoverPages(rootPath string) ([]page, error) {
	var pages []page

	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of root: %v", err)
	}

	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "index.jml" {
			dirPath := filepath.Dir(path)
			relPath, err := filepath.Rel(absRoot, dirPath)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", path, err)
			}
			relPath = strings.TrimSuffix(filepath.ToSlash(relPath), "/")
			if relPath == "." {
				relPath = "/"
			}

			// Derive page name from the directory or use "index" for root
			pageName := "index"
			if relPath != "/" {
				pageName = filepath.Base(dirPath)
			}

			// Use the absolute path of the index.jml file
			absPath, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for %s: %w", path, err)
			}

			pages = append(pages, page{
				Name:    pageName,
				RelPath: relPath,
				AbsPath: absPath,
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", rootPath, err)
	}

	return pages, nil
}

type component struct {
	Name    string
	RelPath string
	AbsPath string
}

func discoverComponents(rootPath string) ([]component, error) {
	var components []component

	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of directory: %v", err)
	}

	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".jml") {
			// TODO: Read first non-empty line of file and get component name
			// Get the component name (filename without .jml extension)
			componentName := strings.TrimSuffix(info.Name(), ".jml")

			// Get the relative path from the root directory
			relPath, err := filepath.Rel(absRoot, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", path, err)
			}

			absPath, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for %s: %w", path, err)
			}

			// Convert to forward-slash notation and remove .jml extension
			relPath = strings.TrimSuffix(filepath.ToSlash(relPath), ".jml")

			components = append(components, component{
				Name:    componentName,
				RelPath: relPath,
				AbsPath: absPath,
			})

		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", rootPath, err)
	}

	return components, nil
}
