package log

import (
	"fmt"
	"log"
	"strings"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

var level = Info

func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func SetLevel(l string) {
	switch strings.ToLower(l) {
	case "debug":
		level = Debug
	case "info":
		level = Info
	case "warn":
		level = Warn
	case "error":
		level = Error
	default:
		level = Info
	}
}

func Printf(l Level, format string, v ...any) {
	if l < level {
		return
	}

	msg := fmt.Sprintf(format, v...)
	log.Printf("[%s] %s", l, msg)
}
