package build

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// TestDependencyGraph_AddFile tests adding files to the graph
func TestDependencyGraph_AddFile(t *testing.T) {
	dg := NewDependencyGraph()

	dg.AddFile("app/index.jml", Page, []string{"components/header.jml"})
	dg.AddFile("components/header.jml", Component, []string{})

	if len(dg.files) != 2 {
		t.Errorf("Expected 2 files in graph, got %d", len(dg.files))
	}

	indexFile := dg.files["app/index.jml"]
	if indexFile == nil {
		t.Fatal("Expected index page to be in graph")
	}

	if indexFile.Type != Page {
		t.Errorf("Expected index page to be of type Page, got %v", indexFile.Type)
	}

	if len(indexFile.Dependencies) != 1 || indexFile.Dependencies[0] != "components/header.jml" {
		t.Errorf("Expected index page to have header dependency, got %v", indexFile.Dependencies)
	}
}

// TestDependencyGraph_BuildOrder_Simple tests basic build ordering
func TestDependencyGraph_BuildOrder_Simple(t *testing.T) {
	dg := NewDependencyGraph()

	// Simple dependency chain: page -> component1 -> component2
	dg.AddFile("app/index.jml", Page, []string{"components/header.jml"})
	dg.AddFile("components/header.jml", Component, []string{"components/logo.jml"})
	dg.AddFile("components/logo.jml", Component, []string{})

	buildOrder, err := dg.BuildOrder()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := []string{"components/logo.jml", "components/header.jml", "app/index.jml"}
	if !reflect.DeepEqual(buildOrder, expected) {
		t.Errorf("Expected build order %v, got %v", expected, buildOrder)
	}
}

// TestDependencyGraph_BuildOrder_Complex tests complex dependency resolution
func TestDependencyGraph_BuildOrder_Complex(t *testing.T) {
	dg := NewDependencyGraph()

	// Complex dependency structure
	dg.AddFile("app/index.jml", Page, []string{"components/header.jml", "components/footer.jml"})
	dg.AddFile("app/about.jml", Page, []string{"components/header.jml", "components/sidebar.jml"})
	dg.AddFile("components/header.jml", Component, []string{"components/nav.jml"})
	dg.AddFile("components/footer.jml", Component, []string{})
	dg.AddFile("components/sidebar.jml", Component, []string{"components/nav.jml"})
	dg.AddFile("components/nav.jml", Component, []string{})

	buildOrder, err := dg.BuildOrder()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify that dependencies come before dependents
	positions := make(map[string]int)
	for i, file := range buildOrder {
		positions[file] = i
	}

	// Check that nav comes before header and sidebar
	if positions["components/nav.jml"] >= positions["components/header.jml"] {
		t.Error("nav.jml should come before header.jml in build order")
	}
	if positions["components/nav.jml"] >= positions["components/sidebar.jml"] {
		t.Error("nav.jml should come before sidebar.jml in build order")
	}

	// Check that components come before pages
	if positions["components/header.jml"] >= positions["app/index.jml"] {
		t.Error("header.jml should come before index.jml in build order")
	}
	if positions["components/header.jml"] >= positions["app/about.jml"] {
		t.Error("header.jml should come before about.jml in build order")
	}
}

// TestDependencyGraph_CyclicDependency tests cycle detection
func TestDependencyGraph_CyclicDependency(t *testing.T) {
	dg := NewDependencyGraph()

	// Create a cycle: a -> b -> c -> a
	dg.AddFile("pages/test.jml", Page, []string{"components/a.jml"})
	dg.AddFile("components/a.jml", Component, []string{"components/b.jml"})
	dg.AddFile("components/b.jml", Component, []string{"components/c.jml"})
	dg.AddFile("components/c.jml", Component, []string{"components/a.jml"})

	_, err := dg.BuildOrder()
	if err == nil {
		t.Fatal("Expected error for cyclic dependency, got nil")
	}

	if !strings.Contains(err.Error(), "cyclic dependency") {
		t.Errorf("Expected error to mention cyclic dependency, got: %v", err)
	}
}

// TestDependencyGraph_FindCycles tests cycle detection functionality
func TestDependencyGraph_FindCycles(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() *DependencyGraph
		wantCycle bool
	}{
		{
			name: "No cycles",
			setupFunc: func() *DependencyGraph {
				dg := NewDependencyGraph()
				dg.AddFile("app/home.jml", Page, []string{"components/header.jml"})
				dg.AddFile("components/header.jml", Component, []string{"components/logo.jml"})
				dg.AddFile("components/logo.jml", Component, []string{})
				return dg
			},
			wantCycle: false,
		},
		{
			name: "Simple cycle",
			setupFunc: func() *DependencyGraph {
				dg := NewDependencyGraph()
				dg.AddFile("components/a.jml", Component, []string{"components/b.jml"})
				dg.AddFile("components/b.jml", Component, []string{"components/a.jml"})
				return dg
			},
			wantCycle: true,
		},
		{
			name: "Complex cycle",
			setupFunc: func() *DependencyGraph {
				dg := NewDependencyGraph()
				dg.AddFile("pages/test.jml", Page, []string{"components/a.jml"})
				dg.AddFile("components/a.jml", Component, []string{"components/b.jml"})
				dg.AddFile("components/b.jml", Component, []string{"components/c.jml"})
				dg.AddFile("components/c.jml", Component, []string{"components/a.jml"})
				dg.AddFile("components/d.jml", Component, []string{}) // Unrelated component
				return dg
			},
			wantCycle: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dg := tt.setupFunc()
			cycles := dg.FindCycles()

			hasCycle := len(cycles) > 0
			if hasCycle != tt.wantCycle {
				t.Errorf("Expected cycle detection to be %v, got %v", tt.wantCycle, hasCycle)
			}
		})
	}
}

// TestDependencyGraph_MissingDependency tests handling of missing files
func TestDependencyGraph_MissingDependency(t *testing.T) {
	dg := NewDependencyGraph()

	dg.AddFile("app/home.jml", Page, []string{"components/missing.jml"})

	_, err := dg.BuildOrder()
	if err == nil {
		t.Fatal("Expected error for missing dependency, got nil")
	}

	if !strings.Contains(err.Error(), "file not found") {
		t.Errorf("Expected error to mention missing file, got: %v", err)
	}
}

// TestDependencyGraph_GetFilesByType tests file type filtering
func TestDependencyGraph_GetFilesByType(t *testing.T) {
	dg := NewDependencyGraph()

	dg.AddFile("app/home.jml", Page, []string{})
	dg.AddFile("app/about.jml", Page, []string{})
	dg.AddFile("components/header.jml", Component, []string{})
	dg.AddFile("components/footer.jml", Component, []string{})

	pages := dg.getFilesByType(Page)
	components := dg.getFilesByType(Component)

	expectedPages := []string{"app/about.jml", "app/home.jml"}
	expectedComponents := []string{"components/footer.jml", "components/header.jml"}

	sort.Strings(pages)
	sort.Strings(components)

	if !reflect.DeepEqual(pages, expectedPages) {
		t.Errorf("Expected pages %v, got %v", expectedPages, pages)
	}

	if !reflect.DeepEqual(components, expectedComponents) {
		t.Errorf("Expected components %v, got %v", expectedComponents, components)
	}
}

// TestDependencyGraph_EmptyGraph tests behaviour with empty graph
func TestDependencyGraph_EmptyGraph(t *testing.T) {
	dg := NewDependencyGraph()

	buildOrder, err := dg.BuildOrder()
	if err != nil {
		t.Fatalf("Unexpected error for empty graph: %v", err)
	}

	if len(buildOrder) != 0 {
		t.Errorf("Expected empty build order, got %v", buildOrder)
	}

	cycles := dg.FindCycles()
	if len(cycles) != 0 {
		t.Errorf("Expected no cycles in empty graph, got %v", cycles)
	}
}

// TestDependencyGraph_SelfReference tests self-referencing files
func TestDependencyGraph_SelfReference(t *testing.T) {
	dg := NewDependencyGraph()

	// File that depends on itself
	dg.AddFile("components/self.jml", Component, []string{"components/self.jml"})

	_, err := dg.BuildOrder()
	if err == nil {
		t.Fatal("Expected error for self-referencing file, got nil")
	}

	cycles := dg.FindCycles()
	if len(cycles) == 0 {
		t.Error("Expected to find cycle for self-referencing file")
	}
}

// TestDependencyGraph_MultipleCycles tests detection of multiple cycles
func TestDependencyGraph_MultipleCycles(t *testing.T) {
	dg := NewDependencyGraph()

	// First cycle: a -> b -> a
	dg.AddFile("components/a.jml", Component, []string{"components/b.jml"})
	dg.AddFile("components/b.jml", Component, []string{"components/a.jml"})

	// Second cycle: x -> y -> z -> x
	dg.AddFile("components/x.jml", Component, []string{"components/y.jml"})
	dg.AddFile("components/y.jml", Component, []string{"components/z.jml"})
	dg.AddFile("components/z.jml", Component, []string{"components/x.jml"})

	cycles := dg.FindCycles()
	if len(cycles) < 2 {
		t.Errorf("Expected at least 2 cycles, got %d", len(cycles))
	}
}

// TestDependencyGraph_DiamondDependency tests diamond dependency pattern
func TestDependencyGraph_DiamondDependency(t *testing.T) {
	dg := NewDependencyGraph()

	// Diamond pattern: page -> (header, footer) -> shared
	dg.AddFile("app/home.jml", Page, []string{"components/header.jml", "components/footer.jml"})
	dg.AddFile("components/header.jml", Component, []string{"components/shared.jml"})
	dg.AddFile("components/footer.jml", Component, []string{"components/shared.jml"})
	dg.AddFile("components/shared.jml", Component, []string{})

	buildOrder, err := dg.BuildOrder()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify that shared comes before both header and footer
	positions := make(map[string]int)
	for i, file := range buildOrder {
		positions[file] = i
	}

	if positions["components/shared.jml"] >= positions["components/header.jml"] {
		t.Error("shared.jml should come before header.jml in build order")
	}
	if positions["components/shared.jml"] >= positions["components/footer.jml"] {
		t.Error("shared.jml should come before footer.jml in build order")
	}
}

// TestDependencyGraph_OrphanedComponents tests components not referenced by pages
func TestDependencyGraph_OrphanedComponents(t *testing.T) {
	dg := NewDependencyGraph()

	dg.AddFile("app/home.jml", Page, []string{"components/header.jml"})
	dg.AddFile("components/header.jml", Component, []string{})
	dg.AddFile("components/orphaned.jml", Component, []string{}) // Not referenced by any page

	buildOrder, err := dg.BuildOrder()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// All files should be included in build order
	if len(buildOrder) != 3 {
		t.Errorf("Expected 3 files in build order, got %d", len(buildOrder))
	}

	// Check that orphaned component is included
	found := false
	for _, file := range buildOrder {
		if file == "components/orphaned.jml" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Orphaned component should be included in build order")
	}
}

// BenchmarkDependencyGraph_BuildOrder benchmarks the build order calculation
func BenchmarkDependencyGraph_BuildOrder(b *testing.B) {
	dg := NewDependencyGraph()

	// Create a large dependency graph
	for i := 0; i < 100; i++ {
		pageName := fmt.Sprintf("pages/page%d.jml", i)
		deps := []string{
			fmt.Sprintf("components/header%d.jml", i%10),
			fmt.Sprintf("components/footer%d.jml", i%10),
		}
		dg.AddFile(pageName, Page, deps)
	}

	for i := 0; i < 50; i++ {
		headerName := fmt.Sprintf("components/header%d.jml", i)
		footerName := fmt.Sprintf("components/footer%d.jml", i)
		baseName := fmt.Sprintf("components/base%d.jml", i%5)

		dg.AddFile(headerName, Component, []string{baseName})
		dg.AddFile(footerName, Component, []string{baseName})
		dg.AddFile(baseName, Component, []string{})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := dg.BuildOrder()
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// TestFileType_String tests the string representation of file types
func TestFileType_String(t *testing.T) {
	tests := []struct {
		fileType FileType
		expected string
	}{
		{Page, "Page"},
		{Component, "Component"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.fileType.String(); got != tt.expected {
				t.Errorf("FileType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
