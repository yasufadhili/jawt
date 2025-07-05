package events

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// EventType represents the type of event
type EventType string

const (
	// File system events
	FileChangedEvent EventType = "file_changed"
	FileCreatedEvent EventType = "file_created"
	FileDeletedEvent EventType = "file_deleted"

	// Compilation events
	CompilationStartedEvent  EventType = "compilation_started"
	CompilationFinishedEvent EventType = "compilation_finished"
	CompilationErrorEvent    EventType = "compilation_error"

	// Process events
	ProcessStartedEvent EventType = "process_started"
	ProcessStoppedEvent EventType = "process_stopped"
	ProcessErrorEvent   EventType = "process_error"
	ProcessOutputEvent  EventType = "process_output"

	// Development server events
	ServerStartedEvent      EventType = "server_started"
	ServerStoppedEvent      EventType = "server_stopped"
	ClientConnectedEvent    EventType = "client_connected"
	ClientDisconnectedEvent EventType = "client_disconnected"

	// Hot reload events
	HotReloadEvent EventType = "hot_reload"

	// System events
	ShutdownEvent EventType = "shutdown"
	ErrorEvent    EventType = "error"
)

// Event represents an event in the system
type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// NewEvent creates a new event
func NewEvent(eventType EventType, source string) *Event {
	return &Event{
		ID:        generateEventID(),
		Type:      eventType,
		Timestamp: time.Now(),
		Source:    source,
		Data:      make(map[string]interface{}),
	}
}

// WithData adds data to the event
func (e *Event) WithData(key string, value interface{}) *Event {
	e.Data[key] = value
	return e
}

// GetData retrieves data from the event
func (e *Event) GetData(key string) (interface{}, bool) {
	value, exists := e.Data[key]
	return value, exists
}

// String returns a string representation of the event
func (e *Event) String() string {
	return fmt.Sprintf("[%s] %s from %s at %s", e.ID, e.Type, e.Source, e.Timestamp.Format("15:04:05"))
}

// EventHandler handles events
type EventHandler func(event *Event)

// EventBus manages event publishing and subscription
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[EventType][]EventHandler
	eventChan   chan *Event
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

// NewEventBus creates a new event bus
func NewEventBus(ctx context.Context) *EventBus {
	busCtx, cancel := context.WithCancel(ctx)

	return &EventBus{
		subscribers: make(map[EventType][]EventHandler),
		eventChan:   make(chan *Event, 100), // Buffered channel
		ctx:         busCtx,
		cancel:      cancel,
	}
}

// Start starts the event bus
func (eb *EventBus) Start() error {
	eb.wg.Add(1)
	go eb.eventLoop()
	return nil
}

// Stop stops the event bus
func (eb *EventBus) Stop() error {
	eb.cancel()
	eb.wg.Wait()
	return nil
}

// Publish publishes an event
func (eb *EventBus) Publish(event *Event) {
	select {
	case eb.eventChan <- event:
		// Event published successfully
	case <-eb.ctx.Done():
		// Context cancelled, stop publishing
	default:
		// Channel full, drop the event (could log this)
	}
}

// Subscribe subscribes to events of a specific type
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
	return nil
}

// Unsubscribe removes a handler from the subscription list
func (eb *EventBus) Unsubscribe(eventType EventType, handler EventHandler) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	handlers, exists := eb.subscribers[eventType]
	if !exists {
		return fmt.Errorf("no subscribers for event type %s", eventType)
	}

	// Remove the handler (basic implementation)
	// TODO: need to compare function pointers
	// or use a subscription ID system
	for i, h := range handlers {
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
			eb.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	return nil
}

// eventLoop processes events
func (eb *EventBus) eventLoop() {
	defer eb.wg.Done()

	for {
		select {
		case event := <-eb.eventChan:
			eb.handleEvent(event)
		case <-eb.ctx.Done():
			return
		}
	}
}

// handleEvent dispatches an event to all subscribers
func (eb *EventBus) handleEvent(event *Event) {
	eb.mu.RLock()
	handlers, exists := eb.subscribers[event.Type]
	eb.mu.RUnlock()

	if !exists {
		return
	}

	// Execute handlers concurrently
	var wg sync.WaitGroup
	for _, handler := range handlers {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					// Log panic but don't crash the event bus
					fmt.Printf("Event handler panic: %v\n", r)
				}
			}()
			h(event)
		}(handler)
	}

	wg.Wait()
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Common event creators

func CreateFileChangeEvent(source, filePath string) *Event {
	return NewEvent(FileChangedEvent, source).WithData("file_path", filePath)
}

func CreateCompilationStartEvent(source string) *Event {
	return NewEvent(CompilationStartedEvent, source)
}

func CreateCompilationFinishEvent(source string, success bool, duration time.Duration) *Event {
	return NewEvent(CompilationFinishedEvent, source).
		WithData("success", success).
		WithData("duration", duration.String())
}

func CreateCompilationErrorEvent(source string, err error) *Event {
	return NewEvent(CompilationErrorEvent, source).WithData("error", err.Error())
}

func CreateProcessStartEvent(source, processName string, pid int) *Event {
	return NewEvent(ProcessStartedEvent, source).
		WithData("process_name", processName).
		WithData("pid", pid)
}

func CreateProcessStopEvent(source, processName string, pid int) *Event {
	return NewEvent(ProcessStoppedEvent, source).
		WithData("process_name", processName).
		WithData("pid", pid)
}

func CreateProcessOutputEvent(source, processName, output string) *Event {
	return NewEvent(ProcessOutputEvent, source).
		WithData("process_name", processName).
		WithData("output", output)
}

func CreateHotReloadEvent(source string, files []string) *Event {
	return NewEvent(HotReloadEvent, source).WithData("files", files)
}

func CreateProcessErrorEvent(source string, processName string, err error) *Event {
	return NewEvent(ErrorEvent, source).WithData("error", err.Error())
}
