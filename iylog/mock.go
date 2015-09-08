package iylog

import (
	"fmt"
	"reflect"
	"sync"
)

// MemLogger is an in-memory implementation of a Logger.
//
// A MemLogger can be used to check if log calls were made. By default
// a MemLogger will capture any logged messages at level DEBUG or
// above. This can be changed using the SetLevel method.
//
// A MemLogger is safe for use by multiple goroutines.
type MemLogger struct {
	i int

	mu       sync.Mutex
	messages []Message
	level    Level
}

type Message struct {
	Format string
	Args   []interface{}
}

func (m Message) String() string {
	return fmt.Sprintf(m.Format, m.Args...)
}

func NewMemLogger() *MemLogger {
	return &MemLogger{level: DEBUG}
}

func (l *MemLogger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages, l.i = []Message{}, 0
}

func (l *MemLogger) Printf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, Message{Format: format, Args: v})
}

func (l *MemLogger) Level() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

func (l *MemLogger) Messages() []Message {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.messages
}

func (l *MemLogger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *MemLogger) Called() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.messages) > 0
}

func (l *MemLogger) CalledN() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.messages)
}

func (l *MemLogger) CalledWith(format string, v ...interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.messages) == 0 {
		panic("MemLogger was never called")
	}
	for _, msg := range l.messages {
		result := format == msg.Format && reflect.DeepEqual(v, msg.Args)
		if result {
			return result
		}
	}
	return false
}

func (l *MemLogger) NextCalledWith(format string, v ...interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.messages) <= l.i {
		panic("No more messages")
	}
	result := format == l.messages[l.i].Format && reflect.DeepEqual(v, l.messages[l.i].Args)
	l.i++
	return result
}

func (l *MemLogger) LastCalledWith(format string, v ...interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.messages) == 0 {
		panic("MemLogger was never called")
	}
	return format == l.messages[len(l.messages)-1].Format && reflect.DeepEqual(v, l.messages[len(l.messages)-1].Args)
}
