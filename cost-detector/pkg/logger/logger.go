package logger

import (
	"fmt"
	"time"
)

// Logger is a simple logger
type Logger struct {
	Level string // "debug", "info", "warning", "error"
}

// NewLogger creates a new logger
func NewLogger(level string) *Logger {
	return &Logger{
		Level: level,
	}
}

// Info logs info message
func (l *Logger) Info(msg string) {
	fmt.Printf("[%s] INFO: %s\n", time.Now().Format("15:04:05"), msg)
}

// Error logs error message
func (l *Logger) Error(msg string) {
	fmt.Printf("[%s] ERROR: %s\n", time.Now().Format("15:04:05"), msg)
}

// Debug logs debug message
func (l *Logger) Debug(msg string) {
	if l.Level == "debug" {
		fmt.Printf("[%s] DEBUG: %s\n", time.Now().Format("15:04:05"), msg)
	}
}
