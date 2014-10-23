package iylog

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Level enumeration is used to describe
// a Loggables logging level.
type Level int

const (
	ERROR Level = iota
	WARNING
	INFO
	DEBUG
)

// String returns the string representation
// of a Level
func (l Level) String() string {
	switch l {
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARNING"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	}
	return "UNKNOWN"
}

// LevelFromString returns the corresponding
// Level based on the provided string `level`
func LevelFromString(level string) Level {
	switch level {
	case "ERROR":
		return ERROR
	case "WARNING":
		return WARNING
	case "INFO":
		return INFO
	default:
		return DEBUG
	}
}

// logger implements the Loggable interface for
// purposes of providing a default Loggable type.
type logger struct {
	l *log.Logger
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.l.Printf(format, v...)
}

func (l *logger) Level() Level { return DEBUG }

var std *Logger = NewLogger(&logger{l: log.New(os.Stderr, "", log.LstdFlags)})

// Loggable interface describes the set of
// types which can be used for logging within
// the iylog.Logger
type Loggable interface {
	Printf(format string, v ...interface{})
	Level() Level
}

// Logger represents multiple active Loggable objects.
//
// Each logging operation will be passed on to each Loggable
// object in turn, where it will be outputted depending on the
// level of the operation.
//
// A Logger can be used simultaneously from multiple goroutines.
type Logger struct {
	loggables []Loggable
	mu        sync.Mutex
}

// NewLogger returns a ready to use Logger and adds
// all of the provided Loggable objects.
func NewLogger(loggables ...Loggable) *Logger {
	logger := &Logger{
		loggables: make([]Loggable, 0),
	}
	logger.Add(loggables...)
	return logger
}

// Add adds the provided loggables
// to the Logger set of loggables.
func (m *Logger) Add(loggables ...Loggable) {
	m.mu.Lock()
	m.loggables = append(m.loggables, loggables...)
	m.mu.Unlock()
}

// Add calls logger.AddLoggables on the
// iylog package standard logger.
func Add(loggables ...Loggable) {
	std.Add(loggables...)
}

// Errorf prints to all loggers with a level of ERROR or above
func (m *Logger) Errorf(format string, v ...interface{}) {
	m.printf(ERROR, format, v...)
}

// Errorf prints to all loggers registered within the iylog
// package standard logger, with a level of ERROR or above.
func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

// Warningf prints to all loggers with a level of WARNING or above
func (m *Logger) Warningf(format string, v ...interface{}) {
	m.printf(WARNING, format, v...)
}

// Warningf prints to all loggers registered within the iylog
// package standard logger, with a level of WARNING or above.
func Warningf(format string, v ...interface{}) {
	std.Warningf(format, v...)
}

// Infof prints to all loggers with a level of INFO or above
func (m *Logger) Infof(format string, v ...interface{}) {
	m.printf(INFO, format, v...)
}

// Infof prints to all loggers registered within the iylog
// package standard logger, with a level of INFO or above.
func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

// Debugf prints to all loggers with a level of DEBUG
func (m *Logger) Debugf(format string, v ...interface{}) {
	m.printf(DEBUG, format, v...)
}

// Debugf prints to all loggers registered within the iylog
// package standard logger, with a level of DEBUG.
func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

// printf performs the printing and formating of levels and messages
// to the Loggers set of loggers. It uses a mutex to ensure
// routine safety.
func (m *Logger) printf(level Level, format string, v ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, l := range m.loggables {
		if level <= l.Level() {
			l.Printf(fmt.Sprintf("[%s] %s", level, format), v...)
		}
	}
}
