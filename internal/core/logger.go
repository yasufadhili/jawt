package core

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents different log levels
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// DefaultLogger is a simple implementation of the Logger interface
type DefaultLogger struct {
	level  LogLevel
	logger *log.Logger
}

// NewDefaultLogger creates a new default logger
func NewDefaultLogger(level LogLevel) *DefaultLogger {
	return &DefaultLogger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

// Debug logs a debug message
func (l *DefaultLogger) Debug(msg string, fields ...Field) {
	if l.level <= DebugLevel {
		l.log(DebugLevel, msg, fields...)
	}
}

// Info logs an info message
func (l *DefaultLogger) Info(msg string, fields ...Field) {
	if l.level <= InfoLevel {
		l.log(InfoLevel, msg, fields...)
	}
}

// Warn logs a warning message
func (l *DefaultLogger) Warn(msg string, fields ...Field) {
	if l.level <= WarnLevel {
		l.log(WarnLevel, msg, fields...)
	}
}

// Error logs an error message
func (l *DefaultLogger) Error(msg string, fields ...Field) {
	if l.level <= ErrorLevel {
		l.log(ErrorLevel, msg, fields...)
	}
}

// Fatal logs a fatal message and exits
func (l *DefaultLogger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields...)
	os.Exit(1)
}

// log handles the actual logging
func (l *DefaultLogger) log(level LogLevel, msg string, fields ...Field) {
	timestamp := time.Now().Format("15:04:05")

	// Format the message with colour coding
	var colorCode string
	switch level {
	case DebugLevel:
		colorCode = "\033[36m" // Cyan
	case InfoLevel:
		colorCode = "\033[32m" // Green
	case WarnLevel:
		colorCode = "\033[33m" // Yellow
	case ErrorLevel:
		colorCode = "\033[31m" // Red
	case FatalLevel:
		colorCode = "\033[35m" // Magenta
	}

	resetColor := "\033[0m"

	// Build the log message
	logMsg := fmt.Sprintf("%s[%s]%s %s %s", colorCode, level.String(), resetColor, timestamp, msg)

	// Add fields if any
	if len(fields) > 0 {
		logMsg += " |"
		for _, field := range fields {
			logMsg += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}

	l.logger.Println(logMsg)
}

// SetLevel sets the log level
func (l *DefaultLogger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel returns the current log level
func (l *DefaultLogger) GetLevel() LogLevel {
	return l.level
}

// Helper functions for creating fields

func StringField(key, value string) Field {
	return Field{Key: key, Value: value}
}

func IntField(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func BoolField(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func ErrorField(err error) Field {
	return Field{Key: "error", Value: err.Error()}
}

func DurationField(key string, value time.Duration) Field {
	return Field{Key: key, Value: value.String()}
}
