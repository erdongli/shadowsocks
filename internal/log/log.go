package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Level int

const (
	// Fine-grained informational events that are most useful to
	// debug an application.
	LevelDebug Level = iota

	// Informational messages that highlight the progress of the
	// application at coarse-grained level.
	LevelInfo

	// Potentially harmful situations.
	LevelWarn

	// Error events that might still allow the application to continue
	// running.
	LevelError

	// Severe error events that will presumably lead the application
	// to abort.
	LevelFatal
)

var level = LevelInfo

// String returns a human readable string of the log level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// SetLevel sets the global log level specifying which message levels will be logged.
func SetLevel(l string) {
	switch strings.ToLower(l) {
	case "debug":
		level = LevelDebug
	case "info":
		level = LevelInfo
	case "warn":
		level = LevelWarn
	case "error":
		level = LevelError
	case "fatal":
		level = LevelFatal
	default:
		level = LevelInfo
	}
}

// Debug logs a message with the DEBUG level.
func Debug(format string, v ...any) {
	levelLog(LevelDebug, format, v...)
}

// Info logs a message with the INFO level.
func Info(format string, v ...any) {
	levelLog(LevelInfo, format, v...)
}

// Warn logs a message with the WARN level.
func Warn(format string, v ...any) {
	levelLog(LevelWarn, format, v...)
}

// Error logs a message with the ERROR level.
func Error(format string, v ...any) {
	levelLog(LevelError, format, v...)
}

// Fatal logs a message with the FATAL level, and aborts the applicatin with status code 1.
func Fatal(format string, v ...any) {
	levelLog(LevelFatal, format, v...)
	os.Exit(1)
}

func levelLog(l Level, format string, v ...any) {
	if l < level {
		return
	}

	msg := fmt.Sprintf(format, v...)
	log.Printf("[%s] %s", l, msg)
}
