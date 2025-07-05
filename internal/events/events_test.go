package events

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	eventType := FileChangedEvent
	source := "test_source"

	event := NewEvent(eventType, source)

	if event.Type != eventType {
		t.Errorf("Expected event type %s, got %s", eventType, event.Type)
	}

	if event.Source != source {
		t.Errorf("Expected source %s, got %s", source, event.Source)
	}

	if event.ID == "" {
		t.Error("Expected event ID to be generated")
	}

	if event.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}

	if event.Data == nil {
		t.Error("Expected data map to be initialised")
	}
}

func TestEventWithData(t *testing.T) {
	event := NewEvent(FileChangedEvent, "test")

	// Test adding data
	event.WithData("key1", "value1")
	event.WithData("key2", 42)

	// Test retrieving data
	value1, exists1 := event.GetData("key1")
	if !exists1 || value1 != "value1" {
		t.Errorf("Expected key1 to have value 'value1', got %v (exists: %v)", value1, exists1)
	}

	value2, exists2 := event.GetData("key2")
	if !exists2 || value2 != 42 {
		t.Errorf("Expected key2 to have value 42, got %v (exists: %v)", value2, exists2)
	}

	// Test non-existent key
	_, exists3 := event.GetData("non_existent")
	if exists3 {
		t.Error("Expected non-existent key to return false")
	}
}

func TestEventString(t *testing.T) {
	event := NewEvent(FileChangedEvent, "test_source")

	str := event.String()

	// Check that the string contains expected components
	if str == "" {
		t.Error("Expected non-empty string representation")
	}

	// Basic check that it contains the event type and source
	if !contains(str, string(FileChangedEvent)) {
		t.Error("Expected string to contain event type")
	}

	if !contains(str, "test_source") {
		t.Error("Expected string to contain source")
	}
}

func TestNewEventBus(t *testing.T) {
	ctx := context.Background()
	bus := NewEventBus(ctx)

	if bus == nil {
		t.Error("Expected event bus to be created")
	}

	if bus.subscribers == nil {
		t.Error("Expected subscribers map to be initialised")
	}

	if bus.eventChan == nil {
		t.Error("Expected event channel to be initialised")
	}

	if bus.ctx == nil {
		t.Error("Expected context to be set")
	}

	// Clean up
	bus.Stop()
}

func TestEventBusStartStop(t *testing.T) {
	ctx := context.Background()
	bus := NewEventBus(ctx)

	// Test starting
	err := bus.Start()
	if err != nil {
		t.Errorf("Expected no error starting bus, got %v", err)
	}

	// Test stopping
	err = bus.Stop()
	if err != nil {
		t.Errorf("Expected no error stopping bus, got %v", err)
	}
}

func TestEventBusSubscribe(t *testing.T) {
	ctx := context.Background()
	bus := NewEventBus(ctx)
	defer bus.Stop()

	_ = false
	handler := func(event *Event) {
		_ = true
	}

	err := bus.Subscribe(FileChangedEvent, handler)
	if err != nil {
		t.Errorf("Expected no error subscribing, got %v", err)
	}

	// Check that handler was added
	bus.mu.RLock()
	handlers := bus.subscribers[FileChangedEvent]
	bus.mu.RUnlock()

	if len(handlers) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(handlers))
	}
}

func TestEventBusPublishAndHandle(t *testing.T) {
	ctx := context.Background()
	bus := NewEventBus(ctx)
	bus.Start()
	defer bus.Stop()

	receivedEvent := make(chan *Event, 1)

	// Subscribe to events
	handler := func(event *Event) {
		receivedEvent <- event
	}

	err := bus.Subscribe(FileChangedEvent, handler)
	if err != nil {
		t.Errorf("Expected no error subscribing, got %v", err)
	}

	// Publish an event
	testEvent := NewEvent(FileChangedEvent, "test")
	bus.Publish(testEvent)

	// Wait for the event to be processed
	select {
	case event := <-receivedEvent:
		if event.ID != testEvent.ID {
			t.Errorf("Expected event ID %s, got %s", testEvent.ID, event.ID)
		}
	case <-time.After(time.Second):
		t.Error("Expected to receive event within timeout")
	}
}

func TestEventBusMultipleSubscribers(t *testing.T) {
	ctx := context.Background()
	bus := NewEventBus(ctx)
	bus.Start()
	defer bus.Stop()

	var wg sync.WaitGroup
	eventCount := 3
	receivedEvents := make([]chan *Event, eventCount)

	// Create multiple subscribers
	for i := 0; i < eventCount; i++ {
		receivedEvents[i] = make(chan *Event, 1)
		wg.Add(1)

		handler := func(ch chan *Event) EventHandler {
			return func(event *Event) {
				ch <- event
				wg.Done()
			}
		}(receivedEvents[i])

		err := bus.Subscribe(FileChangedEvent, handler)
		if err != nil {
			t.Errorf("Expected no error subscribing handler %d, got %v", i, err)
		}
	}

	// Publish an event
	testEvent := NewEvent(FileChangedEvent, "test")
	bus.Publish(testEvent)

	// Wait for all handlers to be called
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		// All handlers called successfully
	case <-time.After(time.Second):
		t.Error("Expected all handlers to be called within timeout")
	}
}

func TestEventBusUnsubscribe(t *testing.T) {
	ctx := context.Background()
	bus := NewEventBus(ctx)
	defer bus.Stop()

	handler := func(event *Event) {}

	// Subscribe
	err := bus.Subscribe(FileChangedEvent, handler)
	if err != nil {
		t.Errorf("Expected no error subscribing, got %v", err)
	}

	// Verify subscription exists
	bus.mu.RLock()
	initialCount := len(bus.subscribers[FileChangedEvent])
	bus.mu.RUnlock()

	if initialCount != 1 {
		t.Errorf("Expected 1 subscriber, got %d", initialCount)
	}

	// Unsubscribe
	err = bus.Unsubscribe(FileChangedEvent, handler)
	if err != nil {
		t.Errorf("Expected no error unsubscribing, got %v", err)
	}

	// Verify subscription removed
	bus.mu.RLock()
	finalCount := len(bus.subscribers[FileChangedEvent])
	bus.mu.RUnlock()

	if finalCount != 0 {
		t.Errorf("Expected 0 subscribers after unsubscribe, got %d", finalCount)
	}
}

func TestEventBusUnsubscribeNonExistent(t *testing.T) {
	ctx := context.Background()
	bus := NewEventBus(ctx)
	defer bus.Stop()

	handler := func(event *Event) {}

	// Try to unsubscribe from non-existent event type
	err := bus.Unsubscribe(FileChangedEvent, handler)
	if err == nil {
		t.Error("Expected error when unsubscribing from non-existent event type")
	}
}

func TestEventBusContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	bus := NewEventBus(ctx)
	bus.Start()

	// Cancel context
	cancel()

	// Stop should complete quickly
	done := make(chan bool)
	go func() {
		bus.Stop()
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-time.After(time.Second):
		t.Error("Expected bus to stop within timeout after context cancellation")
	}
}

func TestEventBusHandlerPanic(t *testing.T) {
	ctx := context.Background()
	bus := NewEventBus(ctx)
	bus.Start()
	defer bus.Stop()

	panicHandler := func(event *Event) {
		panic("test panic")
	}

	normalHandler := func(event *Event) {
		// This should still execute despite the panic
	}

	err := bus.Subscribe(FileChangedEvent, panicHandler)
	if err != nil {
		t.Errorf("Expected no error subscribing panic handler, got %v", err)
	}

	err = bus.Subscribe(FileChangedEvent, normalHandler)
	if err != nil {
		t.Errorf("Expected no error subscribing normal handler, got %v", err)
	}

	// Publish event - this should not crash the bus
	testEvent := NewEvent(FileChangedEvent, "test")
	bus.Publish(testEvent)

	// Give time for handlers to execute
	time.Sleep(100 * time.Millisecond)

	// Bus should still be operational
	anotherEvent := NewEvent(FileChangedEvent, "test2")
	bus.Publish(anotherEvent)
}

func TestEventBusBufferedChannel(t *testing.T) {
	ctx := context.Background()
	bus := NewEventBus(ctx)
	// Don't start the bus to test buffering

	// Publish events up to buffer size
	for i := 0; i < 100; i++ {
		event := NewEvent(FileChangedEvent, "test")
		bus.Publish(event)
	}

	// Channel should have 100 events
	if len(bus.eventChan) != 100 {
		t.Errorf("Expected 100 events in channel, got %d", len(bus.eventChan))
	}

	// One more event should be dropped (channel full)
	event := NewEvent(FileChangedEvent, "test")
	bus.Publish(event)

	// Channel should still have 100 events
	if len(bus.eventChan) != 100 {
		t.Errorf("Expected 100 events in channel after overflow, got %d", len(bus.eventChan))
	}

	bus.Stop()
}

// Test event creation helpers

func TestCreateFileChangeEvent(t *testing.T) {
	filePath := "/path/to/file.txt"
	source := "file_watcher"

	event := CreateFileChangeEvent(source, filePath)

	if event.Type != FileChangedEvent {
		t.Errorf("Expected event type %s, got %s", FileChangedEvent, event.Type)
	}

	if event.Source != source {
		t.Errorf("Expected source %s, got %s", source, event.Source)
	}

	path, exists := event.GetData("file_path")
	if !exists || path != filePath {
		t.Errorf("Expected file_path %s, got %v (exists: %v)", filePath, path, exists)
	}
}

func TestCreateCompilationStartEvent(t *testing.T) {
	source := "compiler"

	event := CreateCompilationStartEvent(source)

	if event.Type != CompilationStartedEvent {
		t.Errorf("Expected event type %s, got %s", CompilationStartedEvent, event.Type)
	}

	if event.Source != source {
		t.Errorf("Expected source %s, got %s", source, event.Source)
	}
}

func TestCreateCompilationFinishEvent(t *testing.T) {
	source := "compiler"
	success := true
	duration := time.Second * 5

	event := CreateCompilationFinishEvent(source, success, duration)

	if event.Type != CompilationFinishedEvent {
		t.Errorf("Expected event type %s, got %s", CompilationFinishedEvent, event.Type)
	}

	successData, exists := event.GetData("success")
	if !exists || successData != success {
		t.Errorf("Expected success %v, got %v (exists: %v)", success, successData, exists)
	}

	durationData, exists := event.GetData("duration")
	if !exists || durationData != duration.String() {
		t.Errorf("Expected duration %s, got %v (exists: %v)", duration.String(), durationData, exists)
	}
}

func TestCreateCompilationErrorEvent(t *testing.T) {
	source := "compiler"
	testErr := errors.New("compilation failed")

	event := CreateCompilationErrorEvent(source, testErr)

	if event.Type != CompilationErrorEvent {
		t.Errorf("Expected event type %s, got %s", CompilationErrorEvent, event.Type)
	}

	errorData, exists := event.GetData("error")
	if !exists || errorData != testErr.Error() {
		t.Errorf("Expected error %s, got %v (exists: %v)", testErr.Error(), errorData, exists)
	}
}

func TestCreateProcessStartEvent(t *testing.T) {
	source := "process_manager"
	processName := "test_process"
	pid := 1234

	event := CreateProcessStartEvent(source, processName, pid)

	if event.Type != ProcessStartedEvent {
		t.Errorf("Expected event type %s, got %s", ProcessStartedEvent, event.Type)
	}

	nameData, exists := event.GetData("process_name")
	if !exists || nameData != processName {
		t.Errorf("Expected process_name %s, got %v (exists: %v)", processName, nameData, exists)
	}

	pidData, exists := event.GetData("pid")
	if !exists || pidData != pid {
		t.Errorf("Expected pid %d, got %v (exists: %v)", pid, pidData, exists)
	}
}

func TestCreateProcessStopEvent(t *testing.T) {
	source := "process_manager"
	processName := "test_process"
	pid := 1234

	event := CreateProcessStopEvent(source, processName, pid)

	if event.Type != ProcessStoppedEvent {
		t.Errorf("Expected event type %s, got %s", ProcessStoppedEvent, event.Type)
	}

	nameData, exists := event.GetData("process_name")
	if !exists || nameData != processName {
		t.Errorf("Expected process_name %s, got %v (exists: %v)", processName, nameData, exists)
	}

	pidData, exists := event.GetData("pid")
	if !exists || pidData != pid {
		t.Errorf("Expected pid %d, got %v (exists: %v)", pid, pidData, exists)
	}
}

func TestCreateProcessOutputEvent(t *testing.T) {
	source := "process_manager"
	processName := "test_process"
	output := "Hello, World!"

	event := CreateProcessOutputEvent(source, processName, output)

	if event.Type != ProcessOutputEvent {
		t.Errorf("Expected event type %s, got %s", ProcessOutputEvent, event.Type)
	}

	nameData, exists := event.GetData("process_name")
	if !exists || nameData != processName {
		t.Errorf("Expected process_name %s, got %v (exists: %v)", processName, nameData, exists)
	}

	outputData, exists := event.GetData("output")
	if !exists || outputData != output {
		t.Errorf("Expected output %s, got %v (exists: %v)", output, outputData, exists)
	}
}

func TestCreateHotReloadEvent(t *testing.T) {
	source := "hot_reload"
	files := []string{"file1.go", "file2.go", "file3.go"}

	event := CreateHotReloadEvent(source, files)

	if event.Type != HotReloadEvent {
		t.Errorf("Expected event type %s, got %s", HotReloadEvent, event.Type)
	}

	filesData, exists := event.GetData("files")
	if !exists {
		t.Error("Expected files data to exist")
	}

	// Type assertion to check the slice
	if filesSlice, ok := filesData.([]string); ok {
		if len(filesSlice) != len(files) {
			t.Errorf("Expected %d files, got %d", len(files), len(filesSlice))
		}
		for i, file := range files {
			if filesSlice[i] != file {
				t.Errorf("Expected file %s at index %d, got %s", file, i, filesSlice[i])
			}
		}
	} else {
		t.Error("Expected files data to be []string")
	}
}

func TestCreateProcessErrorEvent(t *testing.T) {
	source := "process_manager"
	processName := "test_process"
	testErr := errors.New("process error")

	event := CreateProcessErrorEvent(source, processName, testErr)

	if event.Type != ErrorEvent {
		t.Errorf("Expected event type %s, got %s", ErrorEvent, event.Type)
	}

	errorData, exists := event.GetData("error")
	if !exists || errorData != testErr.Error() {
		t.Errorf("Expected error %s, got %v (exists: %v)", testErr.Error(), errorData, exists)
	}
}

func TestGenerateEventID(t *testing.T) {
	id1 := generateEventID()
	id2 := generateEventID()

	if id1 == "" {
		t.Error("Expected non-empty event ID")
	}

	if id2 == "" {
		t.Error("Expected non-empty event ID")
	}

	if id1 == id2 {
		t.Error("Expected unique event IDs")
	}
}

func TestEventTypeConstants(t *testing.T) {
	// Test that all event type constants are defined
	eventTypes := []EventType{
		FileChangedEvent,
		FileCreatedEvent,
		FileDeletedEvent,
		CompilationStartedEvent,
		CompilationFinishedEvent,
		CompilationErrorEvent,
		ProcessStartedEvent,
		ProcessStoppedEvent,
		ProcessErrorEvent,
		ProcessOutputEvent,
		ServerStartedEvent,
		ServerStoppedEvent,
		ClientConnectedEvent,
		ClientDisconnectedEvent,
		HotReloadEvent,
		ShutdownEvent,
		ErrorEvent,
	}

	for _, eventType := range eventTypes {
		if string(eventType) == "" {
			t.Errorf("Event type %s should not be empty", eventType)
		}
	}
}

// Benchmark tests

func BenchmarkNewEvent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewEvent(FileChangedEvent, "test")
	}
}

func BenchmarkEventBusPublish(b *testing.B) {
	ctx := context.Background()
	bus := NewEventBus(ctx)
	bus.Start()
	defer bus.Stop()

	event := NewEvent(FileChangedEvent, "test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Publish(event)
	}
}

func BenchmarkEventBusSubscribe(b *testing.B) {
	ctx := context.Background()
	bus := NewEventBus(ctx)
	defer bus.Stop()

	handler := func(event *Event) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Subscribe(FileChangedEvent, handler)
	}
}

// Helper functions

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Race condition tests

func TestEventBusConcurrentPublish(t *testing.T) {
	ctx := context.Background()
	bus := NewEventBus(ctx)
	bus.Start()
	defer bus.Stop()

	var wg sync.WaitGroup
	goroutines := 10
	eventsPerGoroutine := 100

	// Subscribe to count events
	var eventCount int64
	var mu sync.Mutex

	handler := func(event *Event) {
		mu.Lock()
		eventCount++
		mu.Unlock()
	}

	bus.Subscribe(FileChangedEvent, handler)

	// Launch concurrent publishers
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				event := NewEvent(FileChangedEvent, "test")
				bus.Publish(event)
			}
		}()
	}

	wg.Wait()

	// Give time for all events to be processed
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	finalCount := eventCount
	mu.Unlock()

	expectedCount := int64(goroutines * eventsPerGoroutine)
	if finalCount != expectedCount {
		t.Errorf("Expected %d events to be processed, got %d", expectedCount, finalCount)
	}
}

func TestEventBusConcurrentSubscribe(t *testing.T) {
	ctx := context.Background()
	bus := NewEventBus(ctx)
	defer bus.Stop()

	var wg sync.WaitGroup
	goroutines := 10

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			handler := func(event *Event) {}
			bus.Subscribe(FileChangedEvent, handler)
		}()
	}

	wg.Wait()

	// Check that all subscriptions were added
	bus.mu.RLock()
	handlerCount := len(bus.subscribers[FileChangedEvent])
	bus.mu.RUnlock()

	if handlerCount != goroutines {
		t.Errorf("Expected %d handlers, got %d", goroutines, handlerCount)
	}
}
