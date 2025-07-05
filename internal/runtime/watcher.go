package runtime

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yasufadhili/jawt/internal/core"
)

// FileWatcher manages file system watching
type FileWatcher struct {
	ctx      context.Context
	cancel   context.CancelFunc
	logger   core.Logger
	eventBus core.EventBus

	watcher *fsnotify.Watcher
	paths   []string

	mu            sync.RWMutex
	debounceMap   map[string]time.Time
	debounceDelay time.Duration

	// File patterns to watch
	watchPatterns  []string
	ignorePatterns []string

	wg sync.WaitGroup
}

func NewFileWatcher(ctx context.Context, logger core.Logger, eventBus core.EventBus) (*FileWatcher, error) {
	watcherCtx, cancel := context.WithCancel(ctx)

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		cancel()
		return nil, err
	}

	return &FileWatcher{
		ctx:           watcherCtx,
		cancel:        cancel,
		logger:        logger,
		eventBus:      eventBus,
		watcher:       fsWatcher,
		debounceMap:   make(map[string]time.Time),
		debounceDelay: 100 * time.Millisecond,
		watchPatterns: []string{".jml", ".ts", ".tsx", ".js", ".jsx", ".json", ".css"},
		ignorePatterns: []string{
			".git/", "node_modules/", ".jawt/", "dist/", "build/",
			".DS_Store", "*.tmp", "*.swp", "*.swo",
		},
	}, nil
}

func (fw *FileWatcher) Start() error {
	fw.logger.Info("Starting file watcher")

	fw.wg.Add(1)
	go fw.watchLoop()

	return nil
}

func (fw *FileWatcher) Stop() error {
	fw.logger.Info("Stopping file watcher")

	fw.cancel()

	if err := fw.watcher.Close(); err != nil {
		fw.logger.Error("Failed to close file watcher", core.ErrorField(err))
	}

	fw.wg.Wait()
	return nil
}

func (fw *FileWatcher) AddPath(path string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	// Check if path is already watched
	for _, p := range fw.paths {
		if p == path {
			return nil
		}
	}

	// Add to fsnotify watcher
	if err := fw.watcher.Add(path); err != nil {
		return err
	}

	fw.paths = append(fw.paths, path)
	fw.logger.Debug("Added watch path", core.StringField("path", path))

	return nil
}

func (fw *FileWatcher) RemovePath(path string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	// Remove from fsnotify watcher
	if err := fw.watcher.Remove(path); err != nil {
		return err
	}

	// Remove from our paths slice
	for i, p := range fw.paths {
		if p == path {
			fw.paths = append(fw.paths[:i], fw.paths[i+1:]...)
			break
		}
	}

	fw.logger.Debug("Removed watch path", core.StringField("path", path))
	return nil
}

// AddPathsRecursive adds paths recursively
func (fw *FileWatcher) AddPathsRecursive(paths []string) error {
	for _, path := range paths {
		if err := fw.addPathRecursive(path); err != nil {
			fw.logger.Error("Failed to add path recursively",
				core.StringField("path", path),
				core.ErrorField(err))
			return err
		}
	}
	return nil
}

// addPathRecursive adds a path and all its subdirectories
func (fw *FileWatcher) addPathRecursive(path string) error {
	return filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if directory should be ignored
		if info.IsDir() && fw.shouldIgnoreDirectory(walkPath) {
			return filepath.SkipDir
		}

		// Only add directories to watcher
		if info.IsDir() {
			if err := fw.AddPath(walkPath); err != nil {
				fw.logger.Warn("Failed to add directory to watcher",
					core.StringField("path", walkPath),
					core.ErrorField(err))
			}
		}

		return nil
	})
}

// SetWatchPatterns sets the file patterns to watch
func (fw *FileWatcher) SetWatchPatterns(patterns []string) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	fw.watchPatterns = patterns
	fw.logger.Debug("Set watch patterns", core.StringField("patterns", strings.Join(patterns, ", ")))
}

func (fw *FileWatcher) SetIgnorePatterns(patterns []string) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	fw.ignorePatterns = patterns
	fw.logger.Debug("Set ignore patterns", core.StringField("patterns", strings.Join(patterns, ", ")))
}

// SetDebounceDelay sets the debounce delay for file events
func (fw *FileWatcher) SetDebounceDelay(delay time.Duration) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	fw.debounceDelay = delay
	fw.logger.Debug("Set debounce delay", core.DurationField("delay", delay))
}

// watchLoop is the main watch loop
func (fw *FileWatcher) watchLoop() {
	defer fw.wg.Done()

	for {
		select {
		case <-fw.ctx.Done():
			return
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}
			fw.handleEvent(event)
		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			fw.logger.Error("File watcher error", core.ErrorField(err))
		}
	}
}

// handleEvent handles a file system event
func (fw *FileWatcher) handleEvent(event fsnotify.Event) {
	// Check if file should be ignored
	if fw.shouldIgnoreFile(event.Name) {
		return
	}

	// Check if file matches watch patterns
	if !fw.shouldWatchFile(event.Name) {
		return
	}

	// Debounce the event
	if fw.isDebouncedEvent(event.Name) {
		return
	}

	fw.logger.Debug("File event",
		core.StringField("file", event.Name),
		core.StringField("operation", event.Op.String()))

	// Determine event type and publish
	var eventType EventType
	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		eventType = FileCreatedEvent
	case event.Op&fsnotify.Write == fsnotify.Write:
		eventType = FileChangedEvent
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		eventType = FileDeletedEvent
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		eventType = FileDeletedEvent // Treat rename as delete
	default:
		eventType = FileChangedEvent
	}

	// Create and publish event
	runtimeEvent := NewEvent(eventType, "file_watcher").
		WithData("file_path", event.Name).
		WithData("operation", event.Op.String())

	fw.eventBus.Publish(runtimeEvent)
}

// shouldIgnoreFile checks if a file should be ignored
func (fw *FileWatcher) shouldIgnoreFile(filePath string) bool {
	fw.mu.RLock()
	defer fw.mu.RUnlock()

	// Get relative path and filename
	fileName := filepath.Base(filePath)
	relativePath := filepath.Clean(filePath)

	// Check against ignore patterns
	for _, pattern := range fw.ignorePatterns {
		// Check if pattern matches filename
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return true
		}

		// Check if pattern matches path
		if strings.Contains(relativePath, pattern) {
			return true
		}
	}

	return false
}

// shouldIgnoreDirectory checks if a directory should be ignored
func (fw *FileWatcher) shouldIgnoreDirectory(dirPath string) bool {
	fw.mu.RLock()
	defer fw.mu.RUnlock()

	dirName := filepath.Base(dirPath)
	relativePath := filepath.Clean(dirPath)

	// Check against ignore patterns
	for _, pattern := range fw.ignorePatterns {
		// Remove trailing slash for comparison
		cleanPattern := strings.TrimSuffix(pattern, "/")

		// Check if pattern matches directory name
		if matched, _ := filepath.Match(cleanPattern, dirName); matched {
			return true
		}

		// Check if pattern matches path
		if strings.Contains(relativePath, cleanPattern) {
			return true
		}
	}

	return false
}

// shouldWatchFile checks if a file matches watch patterns
func (fw *FileWatcher) shouldWatchFile(filePath string) bool {
	fw.mu.RLock()
	defer fw.mu.RUnlock()

	// If no patterns specified, watch all files
	if len(fw.watchPatterns) == 0 {
		return true
	}

	fileName := filepath.Base(filePath)
	fileExt := filepath.Ext(filePath)

	// Check against watch patterns
	for _, pattern := range fw.watchPatterns {
		// Check if pattern matches file extension
		if pattern == fileExt {
			return true
		}

		// Check if pattern matches filename
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return true
		}
	}

	return false
}

// isDebouncedEvent checks if an event should be debounced
func (fw *FileWatcher) isDebouncedEvent(filePath string) bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()
	lastEvent, exists := fw.debounceMap[filePath]

	if exists && now.Sub(lastEvent) < fw.debounceDelay {
		return true
	}

	fw.debounceMap[filePath] = now
	return false
}
