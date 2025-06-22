package build

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/page_compiler"
	"github.com/yasufadhili/jawt/internal/project"
	"os"
	"time"
)

// CompilerManager orchestrates the compilation process
type CompilerManager struct {
	project *project.Structure
}

func NewCompilerManager(project *project.Structure) *CompilerManager {
	return &CompilerManager{
		project: project,
	}
}

// CompileProject compiles the entire project
func (cm *CompilerManager) CompileProject() error {

	// Compile components first (they're dependencies)
	if err := cm.compileComponents(); err != nil {
		return fmt.Errorf("component compilation failed: %w", err)
	}

	if err := cm.compilePages(); err != nil {
		return fmt.Errorf("page compilation failed: %w", err)
	}

	if err := cm.copyAssets(); err != nil {
		return fmt.Errorf("asset copying failed: %w", err)
	}

	return nil
}

// compileComponents compiles all components in dependency order
func (cm *CompilerManager) compileComponents() error {

	for name, comp := range cm.project.Components {
		if err := cm.compileComponent(name, comp); err != nil {
			return fmt.Errorf("failed to compile component %s: %w", name, err)
		}
	}

	return nil
}

// compilePages compiles all pages
func (cm *CompilerManager) compilePages() error {

	for name, page := range cm.project.Pages {
		if err := cm.compilePage(page); err != nil {
			return fmt.Errorf("failed to compile page %s: %w", name, err)
		}
	}

	return nil
}

// compileComponent compiles a single component (placeholder)
func (cm *CompilerManager) compileComponent(name string, comp *project.ComponentInfo) error {

	// TODO: call CC (Component Compiler)

	comp.Compiled = true
	return nil
}

// compilePage compiles a single page (placeholder)
func (cm *CompilerManager) compilePage(page *project.PageInfo) error {

	compiler := page_compiler.NewPageCompiler(page)
	if err := compiler.CompilePage(); err != nil {
		return err
	}

	return nil
}

// copyAssets copies all assets to the output directory (placeholder)
func (cm *CompilerManager) copyAssets() error {

	//for _, asset := range cm.project.Assets {
	//	fmt.Printf("  📁 Copying asset: %s\n", asset)
	// TODO: actual asset copying
	//}

	return nil
}

// CompileChanged compiles only files that have changed
func (cm *CompilerManager) CompileChanged() error {

	// Check for changed components
	for name, comp := range cm.project.Components {
		if cm.hasChanged(comp.AbsolutePath, comp.LastModified) {
			if err := cm.compileComponent(name, comp); err != nil {
				return err
			}
			comp.LastModified = time.Now()
		}
	}

	// Check for changed pages
	for _, page := range cm.project.Pages {
		if cm.hasChanged(page.AbsolutePath, page.LastModified) {
			if err := cm.compilePage(page); err != nil {
				return err
			}
			page.LastModified = time.Now()
		}
	}

	return nil
}

// hasChanged checks if a file has been modified since the last compilation
func (cm *CompilerManager) hasChanged(filePath string, lastModified time.Time) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return info.ModTime().After(lastModified)
}
