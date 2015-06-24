package iylog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Level int

const (
	DEBUG Level = 1 << iota
	INFO
	WARNING
	ERROR
)

var std *MultiLogger = NewMultiLogger()

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

// Logger implements the Loggable interface.
//
// It turns a log.Logger into something that can be used with a
// MultiLogger.
type Logger struct {
	logger *log.Logger
	level  Level
}

// NewLogger returns a new Logger instance.
//
// If l is nil, NewLogger uses a log.Logger which writes to os.Stderr
func NewLogger(logger *log.Logger, level Level) *Logger {
	if logger == nil {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return &Logger{logger: logger, level: level}
}

// NewLoggerFromWriter returns a new Logger instance using a standard
// log.Logger with the provided writer.
//
// To create a Logger equivalent to the default Go
// log.Logger, call NewLogger with a nil logger. NewLoggerFromWriter
// panics if w is nil.
func NewLoggerFromWriter(w io.Writer, level Level) *Logger {
	if w == nil {
		panic("the io.Writer must not be nil")
	}
	return NewLogger(log.New(w, "", log.LstdFlags), level)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *Logger) Level() Level { return l.level }

// Loggable describes the set of type
// which can be used for logging within
// the iylog.Logger.
type Loggable interface {
	Printf(format string, v ...interface{})
	Level() Level
}

// MultiLogger maintains multiple Loggable objects.
//
// Each logging operation on a MultiLogger will be passed on to each
// Loggable object in turn, where it will be logged if the Loggable's
// level is less than or equal to the level of the message.
//
// A MultiLogger can be used simultaneously from multiple goroutines.
type MultiLogger struct {
	loggables []Loggable
	mu        *sync.RWMutex
}

// NewMultiLogger returns a ready to use MultiLogger
func NewMultiLogger(loggables ...Loggable) *MultiLogger {
	l := &MultiLogger{
		loggables: make([]Loggable, 0),
		mu:        &sync.RWMutex{},
	}
	l.Add(loggables...)
	return l
}

// Add adds the loggables to the MultiLogger.
func (m *MultiLogger) Add(loggables ...Loggable) {
	m.mu.Lock()
	m.loggables = append(m.loggables, loggables...)
	m.mu.Unlock()
}

// Add adds the loggables to package-level MultiLogger.
func Add(loggables ...Loggable) {
	std.Add(loggables...)
}

// CapturePanic logs panics with a level ERROR
func (m *MultiLogger) CapturePanic() {
	if rec := recover(); rec != nil {
		m.Error(rec)
		panic(rec)
	}
}

// CapturePanic calls CapturePanic on std Logger
func CapturePanic() {
	std.CapturePanic()
}

// Errorf prints to all loggers with a level of ERROR or above
func (m *MultiLogger) Errorf(format string, v ...interface{}) {
	m.printf(ERROR, format, v...)
}

// Errorf prints to all loggers registered within the iylog
// package standard logger, with a level of ERROR or above.
func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

// Error prints to all loggers with a level of ERROR or above
func (m *MultiLogger) Error(v ...interface{}) {
	m.print(ERROR, v...)
}

// Error prints to all loggers registered within the iylog
// package standard logger, with a level of ERROR or above.
func Error(v ...interface{}) {
	std.Error(v...)
}

// Warningf prints to all loggers with a level of WARNING or above
func (m *MultiLogger) Warningf(format string, v ...interface{}) {
	m.printf(WARNING, format, v...)
}

// Warningf prints to all loggers registered within the iylog
// package standard logger, with a level of WARNING or above.
func Warningf(format string, v ...interface{}) {
	std.Warningf(format, v...)
}

// Warning prints to all loggers with a level of WARNING or above
func (m *MultiLogger) Warning(v ...interface{}) {
	m.print(WARNING, v...)
}

// Warning prints to all loggers registered within the iylog
// package standard logger, with a level of WARNING or above.
func Warning(v ...interface{}) {
	std.Warning(v...)
}

// Infof prints to all loggers with a level of INFO or above
func (m *MultiLogger) Infof(format string, v ...interface{}) {
	m.printf(INFO, format, v...)
}

// Infof prints to all loggers registered within the iylog
// package standard logger, with a level of INFO or above.
func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

// Info prints to all loggers with a level of INFO or above
func (m *MultiLogger) Info(v ...interface{}) {
	m.print(INFO, v...)
}

// Info prints to all loggers registered within the iylog
// package standard logger, with a level of INFO or above.
func Info(v ...interface{}) {
	std.Info(v...)
}

// Debugf prints to all loggers with a level of DEBUG
func (m *MultiLogger) Debugf(format string, v ...interface{}) {
	m.printf(DEBUG, format, v...)
}

// Debugf prints to all loggers registered within the iylog
// package standard logger, with a level of DEBUG.
func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

// Debug prints to all loggers with a level of DEBUG or above
func (m *MultiLogger) Debug(v ...interface{}) {
	m.print(DEBUG, v...)
}

// Debug prints to all loggers registered within the iylog
// package standard logger, with a level of DEBUG or above.
func Debug(v ...interface{}) {
	std.Debug(v...)
}

// print performs the printing of levels and messages
// to the Logger's set of loggers.
func (m *MultiLogger) print(level Level, v ...interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	format := fmt.Sprintf("[%s] ", level)
	for i := 0; i < len(v); i++ {
		format += "%v"
		if i < len(v)-1 {
			format += " "
		}
	}

	for _, l := range m.loggables {
		if level >= l.Level() {
			l.Printf(format, v...)
		}
	}
}

// printf performs the printing and formatting of levels and messages
// to the Loggers set of loggers.
func (m *MultiLogger) printf(level Level, format string, v ...interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, l := range m.loggables {
		if level >= l.Level() {
			l.Printf(fmt.Sprintf("[%s] ", level)+format, v...)
		}
	}
}
