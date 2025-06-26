package build

import (
	"fmt"
	"sort"
	"strings"
)

// FileType represents the type of file in the build system
type FileType int

const (
	Page FileType = iota
	Component
)

func (ft FileType) String() string {
	return []string{"page", "component"}[ft]
}

// File represents a source file with its dependencies
type File struct {
	AbsPath      string
	RelPath      string
	Type         FileType
	Dependencies []string
}

// DependencyGraph manages the build dependency graph
type DependencyGraph struct {
	files   map[string]*File
	visited map[string]bool
	stack   map[string]bool // For cycle detection
}

func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		files:   make(map[string]*File),
		visited: make(map[string]bool),
		stack:   make(map[string]bool),
	}
}

// AddFile adds a file to the dependency graph
func (dg *DependencyGraph) AddFile(path string, fileType FileType, dependencies []string) {
	dg.files[path] = &File{
		RelPath:      path,
		Type:         fileType,
		Dependencies: dependencies,
	}
}

// BuildOrder returns the files in dependency order (topological sort)
func (dg *DependencyGraph) BuildOrder() ([]string, error) {
	var result []string
	dg.visited = make(map[string]bool)
	dg.stack = make(map[string]bool)

	pages := dg.getFilesByType(Page)

	for _, pagePath := range pages {
		if err := dg.dfsVisit(pagePath, &result); err != nil {
			return nil, err
		}
	}

	// Also include any components not reached from pages
	components := dg.getFilesByType(Component)
	for _, compPath := range components {
		if !dg.visited[compPath] {
			if err := dg.dfsVisit(compPath, &result); err != nil {
				return nil, err
			}
		}
	}

	// Reverse the result for correct build order
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

// dfsVisit performs depth-first search with cycle detection
func (dg *DependencyGraph) dfsVisit(path string, result *[]string) error {
	if dg.stack[path] {
		return fmt.Errorf("cyclic dependency detected involving: %s", path)
	}

	if dg.visited[path] {
		return nil
	}

	file, exists := dg.files[path]
	if !exists {
		return fmt.Errorf("file not found: %s", path)
	}

	dg.stack[path] = true
	dg.visited[path] = true

	// Visit all dependencies first
	for _, dep := range file.Dependencies {
		if err := dg.dfsVisit(dep, result); err != nil {
			return err
		}
	}

	dg.stack[path] = false
	*result = append(*result, path)
	return nil
}

// getFilesByType returns all files of a specific type
func (dg *DependencyGraph) getFilesByType(fileType FileType) []string {
	var files []string
	for path, file := range dg.files {
		if file.Type == fileType {
			files = append(files, path)
		}
	}
	sort.Strings(files) // For consistent ordering
	return files
}

// PrintGraph prints the dependency graph
func (dg *DependencyGraph) PrintGraph() {
	fmt.Println("=== Dependency Graph ===")

	pages := dg.getFilesByType(Page)
	components := dg.getFilesByType(Component)

	fmt.Println("\nPages:")
	for _, path := range pages {
		file := dg.files[path]
		fmt.Printf("  %s -> [%s]\n", path, strings.Join(file.Dependencies, ", "))
	}

	fmt.Println("\nComponents:")
	for _, path := range components {
		file := dg.files[path]
		fmt.Printf("  %s -> [%s]\n", path, strings.Join(file.Dependencies, ", "))
	}
}

// FindCycles detects and returns any cycles in the graph
func (dg *DependencyGraph) FindCycles() [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	stack := make(map[string]bool)
	path := make([]string, 0)

	for filePath := range dg.files {
		if !visited[filePath] {
			if cycle := dg.findCyclesDFS(filePath, visited, stack, path); cycle != nil {
				cycles = append(cycles, cycle)
			}
		}
	}

	return cycles
}

func (dg *DependencyGraph) findCyclesDFS(current string, visited, stack map[string]bool, path []string) []string {
	visited[current] = true
	stack[current] = true
	path = append(path, current)

	file, exists := dg.files[current]
	if !exists {
		return nil
	}

	for _, dep := range file.Dependencies {
		if stack[dep] {
			// Found cycle - return the cycle path
			cycleStart := -1
			for i, p := range path {
				if p == dep {
					cycleStart = i
					break
				}
			}
			if cycleStart >= 0 {
				cycle := make([]string, len(path)-cycleStart+1)
				copy(cycle, path[cycleStart:])
				cycle[len(cycle)-1] = dep
				return cycle
			}
		}

		if !visited[dep] {
			if cycle := dg.findCyclesDFS(dep, visited, stack, path); cycle != nil {
				return cycle
			}
		}
	}

	stack[current] = false
	return nil
}
