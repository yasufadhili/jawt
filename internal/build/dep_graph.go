package build

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/core"
	"regexp"
)

// ExtractDependencies extracts component and script dependencies from a JML file.
func ExtractDependencies(content string) []string {
	re := regexp.MustCompile(`import\s+(component|script)\s+\w+\s+from\s+"([^"]+)"`)
	matches := re.FindAllStringSubmatch(content, -1)

	var dependencies []string
	for _, match := range matches {
		if len(match) == 3 {
			dependencies = append(dependencies, match[2])
		}
	}

	return dependencies
}

type DependencyGraph interface {
	// Core graph operations
	AddNode(path string, docType DocumentType) error
	RemoveNode(path string) error
	AddDependency(from, to string) error
	RemoveDependency(from, to string) error

	// Query operations
	GetDependencies(path string) []string
	GetDependents(path string) []string
	GetAllNodes() []string

	// Analysis operations
	HasCycle() bool
	GetCycles() [][]string
	GetTopologicalOrder() ([]string, error)
	GetCompilationOrder() ([]string, error)

	// Validation
	ValidateGraph() error

	// Utilities
	IsConnected(from, to string) bool
	GetShortestPath(from, to string) []string

	// Additional methods for better build system integration
	GetNodesByType(docType DocumentType) []string
	GetTransitiveDependencies(path string) []string
	GetTransitiveDependents(path string) []string
}

type dependencyGraph struct {
	nodes        map[string]*GraphNode
	edges        map[string][]string // Adjacency list: node -> dependencies
	reverseEdges map[string][]string // Reverse adjacency list: node -> dependents
	logger       core.Logger
}

type GraphNode struct {
	Path     string
	DocType  DocumentType
	Metadata map[string]interface{} // For future extensibility
}

func NewDependencyGraph() DependencyGraph {
	return &dependencyGraph{
		nodes:        make(map[string]*GraphNode),
		edges:        make(map[string][]string),
		reverseEdges: make(map[string][]string),
	}
}

func (dg *dependencyGraph) AddNode(path string, docType DocumentType) error {
	if _, exists := dg.nodes[path]; exists {
		// Node already exists, update it
		dg.nodes[path].DocType = docType
		return nil
	}

	dg.nodes[path] = &GraphNode{
		Path:     path,
		DocType:  docType,
		Metadata: make(map[string]interface{}),
	}

	// Initialise empty dependency lists
	if _, exists := dg.edges[path]; !exists {
		dg.edges[path] = []string{}
	}
	if _, exists := dg.reverseEdges[path]; !exists {
		dg.reverseEdges[path] = []string{}
	}

	return nil
}

func (dg *dependencyGraph) RemoveNode(path string) error {
	if _, exists := dg.nodes[path]; !exists {
		return fmt.Errorf("node %s does not exist", path)
	}

	// Remove all dependencies from this node
	for _, dep := range dg.edges[path] {
		dg.RemoveDependency(path, dep)
	}

	// Remove all dependencies to this node
	for _, dependent := range dg.reverseEdges[path] {
		dg.RemoveDependency(dependent, path)
	}

	// Remove the node
	delete(dg.nodes, path)
	delete(dg.edges, path)
	delete(dg.reverseEdges, path)

	return nil
}

func (dg *dependencyGraph) AddDependency(from, to string) error {
	// Ensure both nodes exist
	if _, exists := dg.nodes[from]; !exists {
		return fmt.Errorf("source node %s does not exist", from)
	}
	if _, exists := dg.nodes[to]; !exists {
		return fmt.Errorf("target node %s does not exist", to)
	}

	if from == to {
		return fmt.Errorf("cannot add self-dependency for node %s", from)
	}

	// Check if dependency already exists
	for _, dep := range dg.edges[from] {
		if dep == to {
			return nil // Dependency already exists
		}
	}

	if dg.wouldCreateCycle(from, to) {
		return fmt.Errorf("adding dependency %s -> %s would create a cycle", from, to)
	}

	// Add dependency
	dg.edges[from] = append(dg.edges[from], to)
	dg.reverseEdges[to] = append(dg.reverseEdges[to], from)

	return nil
}

func (dg *dependencyGraph) RemoveDependency(from, to string) error {
	// Remove from edges
	if deps, exists := dg.edges[from]; exists {
		for i, dep := range deps {
			if dep == to {
				dg.edges[from] = append(deps[:i], deps[i+1:]...)
				break
			}
		}
	}

	// Remove from reverse edges
	if dependents, exists := dg.reverseEdges[to]; exists {
		for i, dependent := range dependents {
			if dependent == from {
				dg.reverseEdges[to] = append(dependents[:i], dependents[i+1:]...)
				break
			}
		}
	}

	return nil
}

func (dg *dependencyGraph) GetDependencies(path string) []string {
	if deps, exists := dg.edges[path]; exists {
		// Return a copy to prevent external modification
		result := make([]string, len(deps))
		copy(result, deps)
		return result
	}
	return []string{}
}

func (dg *dependencyGraph) GetDependents(path string) []string {
	if dependents, exists := dg.reverseEdges[path]; exists {
		// Return a copy to prevent external modification
		result := make([]string, len(dependents))
		copy(result, dependents)
		return result
	}
	return []string{}
}

func (dg *dependencyGraph) GetAllNodes() []string {
	nodes := make([]string, 0, len(dg.nodes))
	for path := range dg.nodes {
		nodes = append(nodes, path)
	}
	return nodes
}

// Analysis operations

func (dg *dependencyGraph) HasCycle() bool {
	// DFS-based cycle detection using three colours
	// White (0): unvisited, Grey (1): visiting, Black (2): visited
	colour := make(map[string]int)

	for node := range dg.nodes {
		colour[node] = 0 // White
	}

	for node := range dg.nodes {
		if colour[node] == 0 {
			if dg.hasCycleDFS(node, colour) {
				return true
			}
		}
	}
	return false
}

func (dg *dependencyGraph) hasCycleDFS(node string, colour map[string]int) bool {
	colour[node] = 1 // Grey (visiting)

	for _, dep := range dg.edges[node] {
		if colour[dep] == 1 {
			// Back edge found - cycle detected
			return true
		}
		if colour[dep] == 0 && dg.hasCycleDFS(dep, colour) {
			return true
		}
	}

	colour[node] = 2 // Black (visited)
	return false
}

func (dg *dependencyGraph) GetCycles() [][]string {
	var cycles [][]string
	colour := make(map[string]int)
	parent := make(map[string]string)

	for node := range dg.nodes {
		colour[node] = 0 // White
		parent[node] = ""
	}

	for node := range dg.nodes {
		if colour[node] == 0 {
			dg.findCyclesDFS(node, colour, parent, &cycles)
		}
	}

	return cycles
}

func (dg *dependencyGraph) findCyclesDFS(node string, colour map[string]int, parent map[string]string, cycles *[][]string) {
	colour[node] = 1 // Grey

	for _, dep := range dg.edges[node] {
		if colour[dep] == 1 {
			// Back edge found - extract cycle
			cycle := dg.extractCycle(node, dep, parent)
			*cycles = append(*cycles, cycle)
		} else if colour[dep] == 0 {
			parent[dep] = node
			dg.findCyclesDFS(dep, colour, parent, cycles)
		}
	}

	colour[node] = 2 // Black
}

func (dg *dependencyGraph) extractCycle(start, end string, parent map[string]string) []string {
	var cycle []string
	current := start

	for {
		cycle = append(cycle, current)
		if current == end {
			break
		}
		current = parent[current]
		if current == "" {
			break
		}
	}

	// Reverse to get correct order
	for i, j := 0, len(cycle)-1; i < j; i, j = i+1, j-1 {
		cycle[i], cycle[j] = cycle[j], cycle[i]
	}

	return cycle
}

func (dg *dependencyGraph) GetTopologicalOrder() ([]string, error) {
	if dg.HasCycle() {
		return nil, fmt.Errorf("cannot create topological order: graph has cycles")
	}

	// Kahn's algorithm for topological sorting
	inDegree := make(map[string]int)
	for node := range dg.nodes {
		inDegree[node] = len(dg.reverseEdges[node])
	}

	queue := make([]string, 0)
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	result := make([]string, 0, len(dg.nodes))
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		for _, dep := range dg.edges[node] {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				queue = append(queue, dep)
			}
		}
	}

	if len(result) != len(dg.nodes) {
		return nil, fmt.Errorf("topological sort failed: possible cycle detected")
	}

	return result, nil
}

func (dg *dependencyGraph) GetCompilationOrder() ([]string, error) {
	// For compilation, dependencies compiled first,
	// This is the reverse of topological order
	topoOrder, err := dg.GetTopologicalOrder()
	if err != nil {
		return nil, err
	}

	// Reverse the order
	compilationOrder := make([]string, len(topoOrder))
	for i, node := range topoOrder {
		compilationOrder[len(topoOrder)-1-i] = node
	}

	return compilationOrder, nil
}

// Validation

func (dg *dependencyGraph) ValidateGraph() error {
	// Check for invalid dependencies (references to non-existent nodes)
	for node, deps := range dg.edges {
		for _, dep := range deps {
			if _, exists := dg.nodes[dep]; !exists {
				return fmt.Errorf("node %s has dependency on non-existent node %s", node, dep)
			}
		}
	}

	// Check consistency between edges and reverseEdges
	for node, deps := range dg.edges {
		for _, dep := range deps {
			found := false
			for _, reverseDep := range dg.reverseEdges[dep] {
				if reverseDep == node {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("inconsistency: edge %s->%s not found in reverse edges", node, dep)
			}
		}
	}

	// Check reverse consistency
	for node, dependents := range dg.reverseEdges {
		for _, dependent := range dependents {
			found := false
			for _, dep := range dg.edges[dependent] {
				if dep == node {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("inconsistency: reverse edge %s->%s not found in forward edges", dependent, node)
			}
		}
	}

	return nil
}

// Utilities

func (dg *dependencyGraph) IsConnected(from, to string) bool {
	// Check if there's a path from 'from' to 'to'
	visited := make(map[string]bool)
	return dg.dfsSearch(from, to, visited)
}

func (dg *dependencyGraph) dfsSearch(current, target string, visited map[string]bool) bool {
	if current == target {
		return true
	}

	visited[current] = true
	for _, dep := range dg.edges[current] {
		if !visited[dep] {
			if dg.dfsSearch(dep, target, visited) {
				return true
			}
		}
	}
	return false
}

func (dg *dependencyGraph) GetShortestPath(from, to string) []string {
	// BFS to find the shortest path
	if from == to {
		return []string{from}
	}

	queue := [][]string{{from}}
	visited := make(map[string]bool)
	visited[from] = true

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		current := path[len(path)-1]

		for _, dep := range dg.edges[current] {
			if dep == to {
				return append(path, dep)
			}

			if !visited[dep] {
				visited[dep] = true
				newPath := make([]string, len(path)+1)
				copy(newPath, path)
				newPath[len(path)] = dep
				queue = append(queue, newPath)
			}
		}
	}

	return nil // No path found
}

// GetNodesByType returns all nodes of a specific document type
func (dg *dependencyGraph) GetNodesByType(docType DocumentType) []string {
	var nodes []string
	for path, node := range dg.nodes {
		if node.DocType == docType {
			nodes = append(nodes, path)
		}
	}
	return nodes
}

// GetTransitiveDependencies returns all transitive dependencies of a node
func (dg *dependencyGraph) GetTransitiveDependencies(path string) []string {
	visited := make(map[string]bool)
	var result []string

	dg.collectTransitiveDeps(path, visited, &result)

	return result
}

func (dg *dependencyGraph) collectTransitiveDeps(node string, visited map[string]bool, result *[]string) {
	if visited[node] {
		return
	}

	visited[node] = true

	for _, dep := range dg.edges[node] {
		// TODO: Make sure not to duplicate dependencies
		// Hack exists in tests to deduplicate
		// only add nodes to the result slice when they haven't been visited yet
		if !visited[dep] {
			*result = append(*result, dep)
		}
		dg.collectTransitiveDeps(dep, visited, result)
	}
}

// GetTransitiveDependents returns all transitive dependents of a node
func (dg *dependencyGraph) GetTransitiveDependents(path string) []string {
	visited := make(map[string]bool)
	var result []string

	dg.collectTransitiveDependents(path, visited, &result)

	return result
}

func (dg *dependencyGraph) collectTransitiveDependents(node string, visited map[string]bool, result *[]string) {
	if visited[node] {
		return
	}

	visited[node] = true

	for _, dependent := range dg.reverseEdges[node] {
		*result = append(*result, dependent)
		dg.collectTransitiveDependents(dependent, visited, result)
	}
}

// Helper method to check if adding a dependency would create a cycle
func (dg *dependencyGraph) wouldCreateCycle(from, to string) bool {
	// If we can reach 'from' from 'to', then adding from->to would create a cycle
	visited := make(map[string]bool)
	return dg.canReach(to, from, visited)
}

func (dg *dependencyGraph) canReach(start, target string, visited map[string]bool) bool {
	if start == target {
		return true
	}

	if visited[start] {
		return false
	}

	visited[start] = true

	for _, dep := range dg.edges[start] {
		if dg.canReach(dep, target, visited) {
			return true
		}
	}

	return false
}
