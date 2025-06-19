package build

import "fmt"

type Builder struct {
	discovery *ProjectDiscovery
	project   *ProjectStructure
	compiler  *CompilerManager
	watcher   *FileWatcher
	server    *DevServer
}

// NewBuilder creates a new builder instance
func NewBuilder(rootPath string) *Builder {
	discovery := NewProjectDiscovery(rootPath)

	return &Builder{
		discovery: discovery,
	}
}

// Build performs a full project build
func (b *Builder) Build() error {
	//fmt.Println("🏗️  Starting JAWT build process...")

	// Discover project structure
	project, err := b.discovery.DiscoverProject()
	if err != nil {
		return fmt.Errorf("project discovery failed: %w", err)
	}

	b.project = project

	/*
		fmt.Printf("📊 Project Summary:\n")
		fmt.Printf("   📄 Pages: %d\n", len(project.Pages))
		fmt.Printf("   📦 Components: %d\n", len(project.Components))
		fmt.Printf("   📁 Assets: %d\n", len(project.Assets))
	*/

	// Compile project
	b.compiler = NewCompilerManager(project)
	if err := b.compiler.CompileProject(); err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	return nil
}

// RunDev starts the development mode with file watching and server
func (b *Builder) RunDev() error {

	// Build the project first
	if err := b.Build(); err != nil {
		return err
	}

	// Start file watcher
	b.watcher = NewFileWatcher(b.project, b.compiler)
	if err := b.watcher.Start(); err != nil {
		return fmt.Errorf("failed to start file watcher: %w", err)
	}

	// Start development server
	b.server = NewDevServer(b.project)
	if err := b.server.Start(); err != nil {
		return fmt.Errorf("failed to start development server: %w", err)
	}

	return nil
}

// Stop stops all running services
func (b *Builder) Stop() {
	if b.watcher != nil {
		b.watcher.Stop()
	}
	if b.server != nil {
		b.server.Stop()
	}
}

// GetProject returns the current project structure
func (b *Builder) GetProject() *ProjectStructure {
	return b.project
}
