package core

import (
	"context"
	"sync"
)

// JawtContext represents the global context passed to every subsystem
type JawtContext struct {
	ctx    context.Context
	cancel context.CancelFunc

	// Core configuration
	JawtConfig    *JawtConfig
	ProjectConfig *ProjectConfig
	Paths         *ProjectPaths
	Logger        Logger
	EventBus      EventBus

	// Runtime state
	mu       sync.RWMutex
	metadata map[string]interface{}
}

// NewJawtContext creates a new jawt context with the given configurations
func NewJawtContext(jawtConfig *JawtConfig, projectConfig *ProjectConfig, paths *ProjectPaths, logger Logger, eventBus EventBus) *JawtContext {
	ctx, cancel := context.WithCancel(context.Background())

	return &JawtContext{
		ctx:           ctx,
		cancel:        cancel,
		JawtConfig:    jawtConfig,
		ProjectConfig: projectConfig,
		Paths:         paths,
		Logger:        logger,
		EventBus:      eventBus,
		metadata:      make(map[string]interface{}),
	}
}

// Context returns the underlying context.Context
func (tc *JawtContext) Context() context.Context {
	return tc.ctx
}

// Cancel cancels the context and shuts down all subsystems
func (tc *JawtContext) Cancel() {
	tc.cancel()
}

// SetMetadata stores arbitrary metadata in the context
func (tc *JawtContext) SetMetadata(key string, value interface{}) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.metadata[key] = value
}

// GetMetadata retrieves metadata from the context
func (tc *JawtContext) GetMetadata(key string) (interface{}, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	value, exists := tc.metadata[key]
	return value, exists
}

// EventBus interface for pub-sub communication
type EventBus interface {
	Publish(event Event)
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(eventType string, handler EventHandler) error
	Start() error
	Stop() error
}

// EventHandler handles events from the event bus
type EventHandler func(event Event)

// Logger interface for structured logging
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}
