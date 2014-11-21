package iylog

import (
	"fmt"
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
		m.Errorln(rec)
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
	m.Errorf(fmt.Sprint(v...))
}

// Error prints to all loggers registered within the iylog
// package standard logger, with a level of ERROR or above.
func Error(v ...interface{}) {
	std.Error(v...)
}

// Errorln prints to all loggers with a level of ERROR or above
func (m *MultiLogger) Errorln(v ...interface{}) {
	m.Error((append(v, "\n"))...)
}

// Errorln prints with a line break to all loggers registered
// within the iylog package standard logger, with a level of
// ERROR or above.
func Errorln(v ...interface{}) {
	std.Errorln(v...)
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
	m.Warningf(fmt.Sprint(v...))
}

// Warning prints to all loggers registered within the iylog
// package standard logger, with a level of WARNING or above.
func Warning(v ...interface{}) {
	std.Warning(v...)
}

// Warningln prints to all loggers with a level of WARNING or above
func (m *MultiLogger) Warningln(v ...interface{}) {
	m.Warning((append(v, "\n"))...)
}

// Warningln prints with a line break to all loggers registered
// within the iylog package standard logger, with a level of
// WARNING or above.
func Warningln(v ...interface{}) {
	std.Warningln(v...)
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
	m.Infof(fmt.Sprint(v...))
}

// Info prints to all loggers registered within the iylog
// package standard logger, with a level of INFO or above.
func Info(v ...interface{}) {
	std.Info(v...)
}

// Infoln prints to all loggers with a level of INFO or above
func (m *MultiLogger) Infoln(v ...interface{}) {
	m.Info((append(v, "\n"))...)
}

// Infoln prints with a line break to all loggers registered
// within the iylog package standard logger, with a level of
// INFO or above.
func Infoln(v ...interface{}) {
	std.Infoln(v...)
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
	m.Debugf(fmt.Sprint(v...))
}

// Debug prints to all loggers registered within the iylog
// package standard logger, with a level of DEBUG or above.
func Debug(v ...interface{}) {
	std.Debug(v...)
}

// Debugln prints to all loggers with a level of DEBUG or above
func (m *MultiLogger) Debugln(v ...interface{}) {
	m.Debug((append(v, "\n"))...)
}

// Debugln prints with a line break to all loggers registered
// within the iylog package standard logger, with a level of
// DEBUG or above.
func Debugln(v ...interface{}) {
	std.Debugln(v...)
}

// printf performs the printing and formatting of levels and messages
// to the Loggers set of loggers. It uses a mutex to ensure
// routine safety.
func (m *MultiLogger) printf(level Level, format string, v ...interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, l := range m.loggables {
		if level >= l.Level() {
			l.Printf(fmt.Sprintf("[%s] %s", level, format), v...)
		}
	}
}
