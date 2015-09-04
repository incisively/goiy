package sh

import (
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/incisively/goiy/iylog"
	"github.com/stathat/go"
)

const (
	shEnvKey         = "SH_KEY" // StatHat API key environment variable.
	mb       float64 = 2 << 20
	ms               = float64(time.Millisecond)
)

var (
	std = New()
)

// count allows us to send count arguments down a channel.
type count struct {
	name string
	key  string
	n    int
}

// measure allows us to send measure arguments down a channel.
type measure struct {
	name string
	key  string
	v    float64
}

// StatHat implements the MetricsI interface for the StatHat service.
//
// StatHat ships counts and measures off to StatHat as soon as it
// receives them. To prevent blocking, StatHat maintains buffered
// channels of stats to be sent, and drops stats on the floor if they're
// full.
type StatHat struct {
	countC   chan count
	measureC chan measure
	countF   func(string, string, int)
	measureF func(string, string, float64)

	mu     sync.Mutex
	key    string
	prefix string
}

// Option is a functional option for the StatHat type.
type Option func(*StatHat)

// WithAPIKey is a functional option that sets the StatHat API key.
func WithAPIKey(k string) Option {
	return func(s *StatHat) {
		s.key = k
	}
}

// WithPrefix is a functional option that sets the prefix StatHat will
// prepend to stat names before they're shipped to the StatHat service.
func WithPrefix(p string) Option {
	return func(s *StatHat) {
		s.SetPrefix(p)
	}
}

// New returns a new StatHat type.
//
// If a key is provided via the WithAPIKey option, then it will be
// used as the StatHat API key. Otherwise, New will attempt to read the
// StatHat API key from the environment, looking for an SH_KEY variable.
func New(options ...Option) *StatHat {
	s := &StatHat{
		// TODO(edd): expose these.
		countC:   make(chan count, 10000),
		measureC: make(chan measure, 20000),
	}
	s.countF = s.sendCount
	s.measureF = s.sendMeasure

	// Apply any options.
	for _, option := range options {
		option(s)
	}

	if s.key == "" {
		s.key = os.Getenv(shEnvKey)
	}

	// Setup workers for shipping stats off to StatHat service.
	// NB copying channels so that we can switch them out in tests.
	go func(ch <-chan count) {
		for c := range ch {
			if err := stathat.PostEZCount(c.name, c.key, c.n); err != nil {
				iylog.Warning(err)
			}
		}
	}(s.countC)

	go func(ch <-chan measure) {
		for m := range ch {
			if err := stathat.PostEZValue(m.name, m.key, m.v); err != nil {
				iylog.Warning(err)
			}
		}
	}(s.measureC)
	return s
}

// SetAPIKey calls SetAPIKey on the package-level StatHat.
func SetAPIKey(k string) {
	std.SetAPIKey(k)
}

// SetAPIKey sets the StatHat service API key on the StatHat.
func (s *StatHat) SetAPIKey(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.key = k
}

// SetPrefix calls SetPrefix on the package-level StatHat.
func SetPrefix(p string) {
	std.SetPrefix(p)
}

// SetPrefix sets the prefix on the StatHat.
func (s *StatHat) SetPrefix(p string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prefix = strings.TrimSpace(p)
}

// MonitorRuntime sends runtime information via the package-level
// implementation.
func MonitorRuntime(d time.Duration) {
	std.MonitorRuntime(d)
}

// MonitorRuntime periodically sends runtime information to the StatHat
// service.
//
// Since MonitorRuntime measures various memory related statistics about
// the runtime, it currently stops the world.
//
// MonitorRuntime sends the following data:
//	goroutines  - the number of current goroutines;
//  gcpausetime - the total time the last garbage collection took (ms);
//	alloc		- the total number of bytes currently allocated (MB);
//	heapalloc   - the number of bytes allocated to the heap (MB);
//	heapobj     - the number of objects on the heap (MB);
func (s *StatHat) MonitorRuntime(d time.Duration) {
	go func() {
		prefix := "[runtime]"
		for _ = range time.Tick(d) {
			s.Measure(prefix+" "+"goroutines", float64(runtime.NumGoroutine()))

			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)

			s.Measure(prefix+" "+"gcpausetime", float64(mem.PauseNs[(mem.NumGC+255)%256])/ms)
			s.Measure(prefix+" "+"alloc", float64(mem.Alloc)/mb)
			s.Measure(prefix+" "+"heapalloc", float64(mem.HeapAlloc)/mb)
			s.Measure(prefix+" "+"heapobj", float64(mem.HeapObjects)/mb)
		}
	}()
}

// sendCount sends a count down the count channel, dropping the count on
// the floor, if the channel is full.
func (s *StatHat) sendCount(name, key string, n int) {
	select {
	case s.countC <- count{name: name, key: key, n: n}:
	default:
		iylog.Warningf("dropped count for %v", name)
	}
}

// Count calls Count on the package-level StatHat instance.
func Count(name string, n int) {
	std.Count(name, n)
}

// Count increments the stat associated with name by n.
func (s *StatHat) Count(name string, n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.key == "" {
		return
	}

	if s.prefix != "" {
		name = s.prefix + " " + name
	}
	s.countF(name, s.key, n)
}

// sendMeasure sends a measure down the measure channel, dropping the
// measure on the floor, if the channel is full.
func (s *StatHat) sendMeasure(name, key string, v float64) {
	select {
	case s.measureC <- measure{name: name, key: key, v: v}:
	default:
		iylog.Warningf("dropped measure for %v", name)
	}
}

// Measure calls Measure on the package-level StatHat instance.
func Measure(name string, v float64) {
	std.Measure(name, v)
}

// Measure sends a real-value measure to stathat.
func (s *StatHat) Measure(name string, v float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.key == "" {
		return
	}

	if s.prefix != "" {
		name = s.prefix + " " + name
	}
	s.measureF(name, s.key, v)
}

// Time calls Time on the package-level StatHat instance.
func Time(name string, t time.Time, precision time.Duration) {
	std.Time(name, t, precision)
}

// Time is a function to measure the time between `start` and when
// this function is called. The result is sent to StatHat.
// The intention for this function is to be used within a `defer`, e.g:
//
//	now := time.Now()
//	defer stat.TimeStat(now, "Timing Something", time.Millisecond)
func (s *StatHat) Time(stat string, start time.Time, precision time.Duration) {
	tms := time.Since(start) / precision
	// Stathat returns nil for PostEZValue calls anyway
	s.Measure(stat, float64(tms))
}
