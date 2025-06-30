package build

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/yasufadhili/jawt/internal/compiler"
	"github.com/yasufadhili/jawt/internal/project"
	"github.com/yasufadhili/jawt/internal/server"
	"sync"
	"time"
)

type Builder struct {
	project         *project.Project
	ClearCache      bool
	discovery       *project.Discovery
	watcher         *FileWatcher
	server          *server.DevServer
	compilerManager *compiler.Manager

	// Error state management
	errorState *ErrorState
	isRunning  bool
	stopChan   chan struct{}
	mu         sync.RWMutex

	options Options
}

type Options struct {
	Verbose bool
}

// NewBuilder creates a new builder instance
func NewBuilder(p *project.Project) (*Builder, error) {
	discovery := project.NewProjectDiscovery(p)
	cm := compiler.NewCompilerManager(p)
	return &Builder{
		project:         p,
		discovery:       discovery,
		compilerManager: cm,
		ClearCache:      false,
		errorState:      &ErrorState{},
		stopChan:        make(chan struct{}),
	}, nil
}

// Build performs a full project build with error state management
func (b *Builder) Build() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	p, err := b.discovery.DiscoverProject()
	if err != nil {
		buildErr := fmt.Errorf("project discovery failed: %w", err)
		if b.errorState.shouldShowError(buildErr) {
			b.printError("Project Discovery", buildErr)
		}
		return buildErr
	}
	b.project = p
	b.compilerManager = compiler.NewCompilerManager(p)

	if err := b.compilerManager.ValidateDependencies(); err != nil {
		buildErr := fmt.Errorf("dependency validation failed: %w", err)
		if b.errorState.shouldShowError(buildErr) {
			b.printError("Dependency Validation", buildErr)
		}
		return buildErr
	}

	if b.options.Verbose {
		if err := b.compilerManager.PrintBuildPlan(); err != nil {
			buildErr := fmt.Errorf("failed to print build plan: %w", err)
			if b.errorState.shouldShowError(buildErr) {
				b.printError("Build Plan", buildErr)
			}
			return buildErr
		}
	}

	// Compile in dependency order
	if err := b.compilerManager.CompileProject(); err != nil {
		buildErr := err
		if b.errorState.shouldShowError(buildErr) {
			b.printError("Compilation", buildErr)
		}
		return buildErr
	}

	if b.errorState.shouldShowError(nil) {
		b.printSuccess()
	}

	return nil
}

// SetConfig allows manually setting the configuration
func (b *Builder) SetConfig(cfg *project.Config) {
	b.project.Config = cfg
}

// ErrorState tracks the current error state to prevent spam
type ErrorState struct {
	mu                sync.RWMutex
	lastError         error
	lastErrorHash     string
	lastErrorTime     time.Time
	errorCount        int
	hasShownError     bool
	successAfterError bool
}

// hashError creates a simple hash of the error for comparison
func (es *ErrorState) hashError(err error) string {
	if err == nil {
		return ""
	}
	h := md5.Sum([]byte(err.Error()))
	return hex.EncodeToString(h[:])
}

// shouldShowError determines if we should show this error
func (es *ErrorState) shouldShowError(err error) bool {
	es.mu.Lock()
	defer es.mu.Unlock()

	if err == nil {
		// Clear error state on success
		if es.hasShownError && !es.successAfterError {
			es.successAfterError = true
			return true // Show a success message
		}
		es.reset()
		return false
	}

	currentHash := es.hashError(err)
	now := time.Now()

	// If it's a new error, or we haven't shown this error yet
	if currentHash != es.lastErrorHash || !es.hasShownError {
		es.lastError = err
		es.lastErrorHash = currentHash
		es.lastErrorTime = now
		es.errorCount = 1
		es.hasShownError = true
		es.successAfterError = false
		return true
	}

	// Same error - increment count but don't show
	es.errorCount++

	// Show the error again after a longer period (5 minutes) as a reminder
	if now.Sub(es.lastErrorTime) > 5*time.Minute {
		es.lastErrorTime = now
		return true
	}

	return false
}

// reset clears the error state
func (es *ErrorState) reset() {
	es.lastError = nil
	es.lastErrorHash = ""
	es.errorCount = 0
	es.hasShownError = false
	es.successAfterError = false
}

// getErrorSummary returns a summary of the current error state
func (es *ErrorState) getErrorSummary() string {
	es.mu.RLock()
	defer es.mu.RUnlock()

	if es.lastError == nil {
		return ""
	}

	if es.errorCount > 1 {
		return fmt.Sprintf("(occurred %d times, last at %s)",
			es.errorCount, es.lastErrorTime.Format("15:04:05"))
	}
	return ""
}

// printError prints error messages with context
func (b *Builder) printError(context string, err error) {
	timestamp := time.Now().Format("15:04:05")
	summary := b.errorState.getErrorSummary()

	fmt.Printf("\n🔴 [%s] %s Error:\n", timestamp, context)
	fmt.Printf("   %s\n", err.Error())
	if summary != "" {
		fmt.Printf("   %s\n", summary)
	}
	fmt.Printf("   Watching for changes...\n\n")
}

// printSuccess prints success messages
func (b *Builder) printSuccess() {
	timestamp := time.Now().Format("15:04:05")
	//stats := b.compiler.GetCompilationStats()

	fmt.Printf("\n✅ [%s] Build completed successfully!\n", timestamp)
	//fmt.Printf("   %s\n", stats.String())
	fmt.Printf("   Watching for changes...\n\n")
}

// printIncrementalSuccess prints incremental build success
func (b *Builder) printIncrementalSuccess() {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("✅ [%s] Incremental build completed\n", timestamp)
}

// handleWatcherError handles errors from the file watcher
func (b *Builder) handleWatcherError(err error) {
	if b.errorState.shouldShowError(err) {
		b.printError("File Watcher", err)
	}
}

// handleFileChange handles file change events
func (b *Builder) handleFileChange(filePath string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] File changed: %s\n", timestamp, filePath)

	// Attempt incremental build
	if err := b.BuildIncremental(); err != nil {
		fmt.Printf("   Incremental build failed, trying full rebuild...\n")
		if err := b.Build(); err != nil {
			// Full build also failed - error already handled by Build()
			return
		}
	}
}

// GetProject returns the current project structure
func (b *Builder) GetProject() *project.Project {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.project
}

// IsRunning returns whether the builder is currently running
func (b *Builder) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.isRunning
}

// GetErrorState returns the current error state (for testing/debugging)
func (b *Builder) GetErrorState() *ErrorState {
	return b.errorState
}

// ClearErrorState manually clears the error state
func (b *Builder) ClearErrorState() {
	b.errorState.mu.Lock()
	defer b.errorState.mu.Unlock()
	b.errorState.reset()
}

// GetStats returns current build statistics
func (b *Builder) GetStats() Stats {
	b.mu.RLock()
	defer b.mu.RUnlock()

	stats := Stats{
		IsRunning:  b.isRunning,
		HasErrors:  b.errorState.lastError != nil,
		ErrorCount: b.errorState.errorCount,
	}

	//if b.compiler != nil {
	//	compStats := b.compiler.GetCompilationStats()
	//	stats.CompilationStats = &compStats
	//}

	return stats
}

// Stats holds overall build statistics
type Stats struct {
	IsRunning  bool
	HasErrors  bool
	ErrorCount int
	//CompilationStats *CompilationStats
}

// String returns a formatted string of build stats
func (s Stats) String() string {
	status := "stopped"
	if s.IsRunning {
		status = "running"
	}

	errorInfo := ""
	if s.HasErrors {
		errorInfo = fmt.Sprintf(", %d errors", s.ErrorCount)
	}

	compInfo := ""
	//if s.CompilationStats != nil {
	//	compInfo = fmt.Sprintf(", %s", s.CompilationStats.String())
	//}

	return fmt.Sprintf("Status: %s%s%s", status, errorInfo, compInfo)
}
