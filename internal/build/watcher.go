package build

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/yasufadhili/jawt/internal/project"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileWatcher struct {
	project   *project.Project
	watcher   *fsnotify.Watcher
	isRunning bool
	stopCh    chan struct{}
	mu        sync.RWMutex

	// Event debouncing
	debounceMap   map[string]*time.Timer
	debounceMu    sync.Mutex
	debounceDelay time.Duration

	// Event handlers
	errorHandler  func(error)
	changeHandler func(string)

	// File filters
	watchedExts map[string]bool
	ignoredDirs map[string]bool
}

// NewFileWatcher creates a new file watcher instance
func NewFileWatcher(project *project.Project) *FileWatcher {
	return &FileWatcher{
		project:       project,
		stopCh:        make(chan struct{}),
		debounceMap:   make(map[string]*time.Timer),
		debounceDelay: 300 * time.Millisecond,
		watchedExts: map[string]bool{
			".jml":  true,
			".ts":   true,
			".css":  true,
			".json": true,
		},
		ignoredDirs: map[string]bool{
			".git":        true,
			".jawt":       true,
			"dist":        true,
			"build":       true,
			".cache":      true,
			".jawt-cache": true,
		},
	}
}

// SetErrorHandler sets the error handler function
func (fw *FileWatcher) SetErrorHandler(handler func(error)) {
	fw.errorHandler = handler
}

// SetChangeHandler sets the file change handler function
func (fw *FileWatcher) SetChangeHandler(handler func(string)) {
	fw.changeHandler = handler
}

// SetDebounceDelay sets the debounce delay for file events
func (fw *FileWatcher) SetDebounceDelay(delay time.Duration) {
	fw.debounceDelay = delay
}

// AddWatchedExtension adds a file extension to watch
func (fw *FileWatcher) AddWatchedExtension(ext string) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	fw.watchedExts[ext] = true
}

// RemoveWatchedExtension removes a file extension from watching
func (fw *FileWatcher) RemoveWatchedExtension(ext string) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	delete(fw.watchedExts, ext)
}

// shouldWatchFile determines if a file should be watched
func (fw *FileWatcher) shouldWatchFile(filePath string) bool {
	// Check if a file extension is watched
	ext := filepath.Ext(filePath)
	if !fw.watchedExts[ext] {
		return false
	}

	// Check if any parent directory is ignored
	dir := filepath.Dir(filePath)
	for dir != "." && dir != "/" {
		if fw.ignoredDirs[filepath.Base(dir)] {
			return false
		}
		dir = filepath.Dir(dir)
	}

	return true
}

// shouldWatchDir determines if a directory should be watched
func (fw *FileWatcher) shouldWatchDir(dirPath string) bool {
	dirName := filepath.Base(dirPath)
	return !fw.ignoredDirs[dirName]
}

// Start begins watching for file changes
func (fw *FileWatcher) Start() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if fw.isRunning {
		return fmt.Errorf("file watcher is already running")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file system watcher: %w", err)
	}

	fw.watcher = watcher
	fw.isRunning = true

	// Add project root and all subdirectories to the watcher
	if err := fw.addWatchPaths(); err != nil {
		fw.watcher.Close()
		fw.isRunning = false
		return fmt.Errorf("failed to add watch paths: %w", err)
	}

	// Start an event processing goroutine
	go fw.processEvents()

	return nil
}

// addWatchPaths adds all relevant paths to the watcher
func (fw *FileWatcher) addWatchPaths() error {
	if err := fw.watcher.Add(fw.project.RootPath); err != nil {
		return err
	}

	// Walk through the project structure and add directories
	return filepath.Walk(fw.project.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
		}

		if info.IsDir() && fw.shouldWatchDir(path) {
			if err := fw.watcher.Add(path); err != nil {
				// Log error but continue - some directories might not be accessible
				if fw.errorHandler != nil {
					fw.errorHandler(fmt.Errorf("failed to watch directory %s: %w", path, err))
				}
			}
		}

		return nil
	})
}

// processEvents processes file system events
func (fw *FileWatcher) processEvents() {
	defer fw.watcher.Close()

	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}
			fw.handleEvent(event)

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			if fw.errorHandler != nil {
				fw.errorHandler(fmt.Errorf("file watcher error: %w", err))
			}

		case <-fw.stopCh:
			return
		}
	}
}

// handleEvent processes a single file system event
func (fw *FileWatcher) handleEvent(event fsnotify.Event) {
	// Only watch files we care about
	if !fw.shouldWatchFile(event.Name) {
		return
	}

	switch {
	case event.Op&fsnotify.Write == fsnotify.Write:
		fw.debounceFileChange(event.Name)
	case event.Op&fsnotify.Create == fsnotify.Create:
		fw.debounceFileChange(event.Name)
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		fw.debounceFileChange(event.Name)
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		fw.debounceFileChange(event.Name)
	}
}

// debounceFileChange debounces file change events to avoid rapid rebuilds
func (fw *FileWatcher) debounceFileChange(filePath string) {
	fw.debounceMu.Lock()
	defer fw.debounceMu.Unlock()

	// Cancel the existing timer for this file
	if timer, exists := fw.debounceMap[filePath]; exists {
		timer.Stop()
	}

	// Create new timer
	fw.debounceMap[filePath] = time.AfterFunc(fw.debounceDelay, func() {
		fw.triggerFileChange(filePath)

		// Clean up the timer
		fw.debounceMu.Lock()
		delete(fw.debounceMap, filePath)
		fw.debounceMu.Unlock()
	})
}

// triggerFileChange triggers the file change handler
func (fw *FileWatcher) triggerFileChange(filePath string) {
	if fw.changeHandler != nil {
		// Convert to a relative path for cleaner output
		relPath, err := filepath.Rel(fw.project.RootPath, filePath)
		if err != nil {
			relPath = filePath
		}
		fw.changeHandler(relPath)
	}
}

// Stop stops the file watcher
func (fw *FileWatcher) Stop() {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if !fw.isRunning {
		return
	}

	fw.isRunning = false
	close(fw.stopCh)

	// Cancel all pending debounce timers
	fw.debounceMu.Lock()
	for filePath, timer := range fw.debounceMap {
		timer.Stop()
		delete(fw.debounceMap, filePath)
	}
	fw.debounceMu.Unlock()
}

// IsRunning returns whether the watcher is currently running
func (fw *FileWatcher) IsRunning() bool {
	fw.mu.RLock()
	defer fw.mu.RUnlock()
	return fw.isRunning
}

// GetWatchedPaths returns a list of currently watched paths
func (fw *FileWatcher) GetWatchedPaths() []string {
	if fw.watcher == nil {
		return []string{}
	}

	return fw.watcher.WatchList()
}

// GetStats returns watcher statistics
func (fw *FileWatcher) GetStats() WatcherStats {
	fw.mu.RLock()
	defer fw.mu.RUnlock()

	stats := WatcherStats{
		IsRunning:      fw.isRunning,
		WatchedPaths:   len(fw.GetWatchedPaths()),
		WatchedExts:    make([]string, 0, len(fw.watchedExts)),
		IgnoredDirs:    make([]string, 0, len(fw.ignoredDirs)),
		DebounceDelay:  fw.debounceDelay,
		PendingChanges: len(fw.debounceMap),
	}

	for ext := range fw.watchedExts {
		stats.WatchedExts = append(stats.WatchedExts, ext)
	}

	for dir := range fw.ignoredDirs {
		stats.IgnoredDirs = append(stats.IgnoredDirs, dir)
	}

	return stats
}

// WatcherStats holds file watcher statistics
type WatcherStats struct {
	IsRunning      bool
	WatchedPaths   int
	WatchedExts    []string
	IgnoredDirs    []string
	DebounceDelay  time.Duration
	PendingChanges int
}

// String returns a formatted string of watcher stats
func (ws WatcherStats) String() string {
	status := "stopped"
	if ws.IsRunning {
		status = "running"
	}

	return fmt.Sprintf(
		"Watcher: %s, %d paths, %d extensions, %dms debounce, %d pending",
		status, ws.WatchedPaths, len(ws.WatchedExts),
		ws.DebounceDelay.Milliseconds(), ws.PendingChanges,
	)
}

// AddIgnoredDirectory adds a directory to the ignore list
func (fw *FileWatcher) AddIgnoredDirectory(dirName string) {
	fw.ignoredDirs[dirName] = true
}

// RemoveIgnoredDirectory removes a directory from the ignore list
func (fw *FileWatcher) RemoveIgnoredDirectory(dirName string) {
	delete(fw.ignoredDirs, dirName)
}

// Restart restarts the file watcher (useful after configuration changes)
func (fw *FileWatcher) Restart() error {
	if fw.IsRunning() {
		fw.Stop()
		// Give it a moment to clean up
		time.Sleep(100 * time.Millisecond)
	}
	return fw.Start()
}
