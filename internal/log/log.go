package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var level = LevelInfo

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

func Debug(format string, v ...any) {
	levelLog(LevelDebug, format, v...)
}

func Info(format string, v ...any) {
	levelLog(LevelInfo, format, v...)
}

func Warn(format string, v ...any) {
	levelLog(LevelWarn, format, v...)
}

func Error(format string, v ...any) {
	levelLog(LevelError, format, v...)
}

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
