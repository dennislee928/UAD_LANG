package common

import (
	"fmt"
	"io"
	"os"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case LogDebug:
		return "DEBUG"
	case LogInfo:
		return "INFO"
	case LogWarn:
		return "WARN"
	case LogError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging functionality
type Logger struct {
	Level      LogLevel
	Output     io.Writer
	TimeFormat string
	Prefix     string
}

// NewLogger creates a new logger with default settings
func NewLogger() *Logger {
	return &Logger{
		Level:      LogInfo,
		Output:     os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",
		Prefix:     "",
	}
}

// NewLoggerWithLevel creates a new logger with specified level
func NewLoggerWithLevel(level LogLevel) *Logger {
	logger := NewLogger()
	logger.Level = level
	return logger
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.Level = level
}

// SetOutput sets the output writer
func (l *Logger) SetOutput(w io.Writer) {
	l.Output = w
}

// SetPrefix sets the log prefix
func (l *Logger) SetPrefix(prefix string) {
	l.Prefix = prefix
}

// log writes a log message if the level is sufficient
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.Level {
		return
	}
	
	timestamp := time.Now().Format(l.TimeFormat)
	message := fmt.Sprintf(format, args...)
	
	prefix := l.Prefix
	if prefix != "" {
		prefix = prefix + " "
	}
	
	fmt.Fprintf(l.Output, "[%s] %s%s: %s\n",
		timestamp, prefix, level.String(), message)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LogDebug, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LogInfo, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LogWarn, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LogError, format, args...)
}

// Global default logger
var defaultLogger = NewLogger()

// SetDefaultLevel sets the level for the default logger
func SetDefaultLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

// SetDefaultOutput sets the output for the default logger
func SetDefaultOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

// SetDefaultPrefix sets the prefix for the default logger
func SetDefaultPrefix(prefix string) {
	defaultLogger.SetPrefix(prefix)
}

// Debug logs a debug message using the default logger
func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

// Info logs an info message using the default logger
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

// Warn logs a warning message using the default logger
func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

// Error logs an error message using the default logger
func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

