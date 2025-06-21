package build

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/config"
)

type Builder struct {
	ProjectPath string
	Config      *config.Config
	discovery   *ProjectDiscovery
	project     *ProjectStructure
	compiler    *CompilerManager
	watcher     *FileWatcher
	server      *DevServer
}

// NewBuilder creates a new builder instance
func NewBuilder(rootPath string) *Builder {
	cfg, _ := config.LoadConfig(rootPath)
	discovery := NewProjectDiscovery(rootPath)

	return &Builder{
		discovery:   discovery,
		ProjectPath: rootPath,
		Config:      cfg,
	}
}

// SetConfig allows manually setting the configuration
func (b *Builder) SetConfig(cfg *config.Config) {
	b.Config = cfg
}

// Build performs a full project build
func (b *Builder) Build() error {

	project, err := b.discovery.DiscoverProject()
	if err != nil {
		return fmt.Errorf("project discovery failed: %w", err)
	}

	b.project = project

	b.compiler = NewCompilerManager(project)
	if err := b.compiler.CompileProject(); err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	return nil
}

// RunDev starts the development mode with file watching and server
func (b *Builder) RunDev() error {

	if err := b.Build(); err != nil {
		return err
	}

	b.watcher = NewFileWatcher(b.project, b.compiler)
	if err := b.watcher.Start(); err != nil {
		return fmt.Errorf("failed to start file watcher: %w", err)
	}

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
