package logger

import (
	"fmt"
	"time"
)

// Logger provides structured logging
// WHY: Structured logs are easier to parse, search, and alert on
// Better than printf for production systems
type Logger struct {
	Level string // "debug", "info", "warning", "error"
}

// NewLogger creates a new logger
func NewLogger(level string) *Logger {
	return &Logger{
		Level: level,
	}
}

// Info logs an informational message
func (l *Logger) Info(msg string) {
	fmt.Printf("[%s] INFO: %s\n", time.Now().Format("15:04:05"), msg)
}

// Warning logs a warning message
func (l *Logger) Warning(msg string) {
	fmt.Printf("[%s] WARN: %s\n", time.Now().Format("15:04:05"), msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	fmt.Printf("[%s] ERROR: %s\n", time.Now().Format("15:04:05"), msg)
}

// Debug logs a debug message (only if level is "debug")
func (l *Logger) Debug(msg string) {
	if l.Level == "debug" {
		fmt.Printf("[%s] DEBUG: %s\n", time.Now().Format("15:04:05"), msg)
	}
}
