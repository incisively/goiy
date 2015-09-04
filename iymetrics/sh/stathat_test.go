package sh

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/incisively/goiy/iylog"
)

func Example() {
	// Set a prefix for the package-level instance.
	SetPrefix("[service-a]")

	// Usually you would have set the SH_API key in the environment
	// but you also explicitly set it.
	SetAPIKey("ssiuhIUYDGos")

	// Monitor the process runtime every 2 minutes.
	MonitorRuntime(2 * time.Minute)

	// Sends a count of 1 to the StatHat service with the stat name
	// "[service-a] users".
	Count("users", 1)

	// Sends a measure to the StatHat service with the stat name
	// "[service-a] length".
	Measure("length", 24.22)

	// Sends a timing to the StatHat service with the name
	// "[service-a] work-ms".
	now := time.Now()
	go func() {
		// Do some work in another goroutine.
		// Work()
		defer Time("work-ms", now, time.Millisecond)
	}()
}

func ExampleStatHat() {
	// Usually you would have set the SH_API key in the environment
	// but you also explicitly set it.
	s := New(WithAPIKey("ssiuhIUYDGos"), WithPrefix("[service-a]"))

	// Monitor the process runtime every 2 minutes.
	s.MonitorRuntime(2 * time.Minute)

	// Sends a count of 1 to the StatHat service with the stat name
	// "[service-a] users".
	s.Count("users", 1)

	// Sends a measure to the StatHat service with the stat name
	// "[service-a] length".
	s.Measure("length", 24.22)

	// Sends a timing to the StatHat service with the name
	// "[service-a] work-ms".
	now := time.Now()
	go func() {
		// Do some work in another goroutine.
		// Work()
		defer s.Time("work-ms", now, time.Millisecond)
	}()
}

func TestNew(t *testing.T) {
	current := os.Getenv(shEnvKey)
	os.Setenv(shEnvKey, "helloKey")
	defer func() { os.Setenv(shEnvKey, current) }()

	s := New(WithAPIKey("foo"))
	// it overrides the environment when WithAPIKey option used.
	if s.key != "foo" {
		t.Fatalf("Expected: %q, got: %v", "foo", s.key)
	}

	// it uses environment when empty key passed in
	s = New()
	if s.key != "helloKey" {
		t.Fatalf("Expected: %q, got: %v", "helloKey", s.key)
	}
}

func TestWithPrefix(t *testing.T) {
	s := New(WithPrefix(" [service] "))

	// it sets prefix on the StatHat
	if s.prefix != "[service]" {
		t.Errorf("Expected: %v, got: %v", "[service]", s.prefix)
	}
}

func TestStatHat_SetPrefix(t *testing.T) {
	s := New()
	s.SetPrefix(" foo ")

	if s.prefix != "foo" {
		t.Errorf("expected %q, got %q", "foo", s.prefix)
	}
}

func TestStatHat_Count(t *testing.T) {
	var (
		called    bool
		name, key string
		count     int
	)
	s := New()
	s.countF = func(n string, k string, c int) {
		called = true
		name, key, count = n, k, c
	}

	// It doesn't call SH API when no key is set.
	s.Count("foo", 2)
	if called {
		t.Errorf("expected %v, got %v", false, called)
	}

	// It calls the StatHat API with the correct values when a key is
	// set.
	s.key = "fookey"
	s.prefix = "prefix"
	s.Count("stat", 10)
	if !called {
		t.Errorf("expected %v, got %v", true, called)
	}

	if name != "prefix stat" {
		t.Errorf("expected %v, got %v", "prefix stat", name)
	}

	if key != "fookey" {
		t.Errorf("expected %v, got %v", "fookey", key)
	}

	if count != 10 {
		t.Errorf("expected %v, got %v", 10, count)
	}
}

// TODO(edd): DRY this up with TestStatHat_Count.
func TestStatHat_Measure(t *testing.T) {
	var (
		called    bool
		name, key string
		value     float64
	)
	s := New()
	s.measureF = func(n string, k string, v float64) {
		called = true
		name, key, value = n, k, v
	}

	// It doesn't call SH API when no key is set.
	s.Measure("foo", 2.4)
	if called {
		t.Errorf("expected %v, got %v", false, called)
	}

	// It calls the StatHat API with the correct values when a key is
	// set.
	s.key = "fookey"
	s.prefix = "prefix"
	s.Measure("stat", 12.3)
	if !called {
		t.Errorf("expected %v, got %v", true, called)
	}

	if name != "prefix stat" {
		t.Errorf("expected %v, got %v", "prefix stat", name)
	}

	if key != "fookey" {
		t.Errorf("expected %v, got %v", "fookey", key)
	}

	if value != 12.3 {
		t.Errorf("expected %v, got %v", 12.3, value)
	}
}

func TestStatHat_Time(t *testing.T) {
	var (
		name  string
		value float64
	)
	s := New(WithAPIKey("fookey"))
	s.measureF = func(n string, _ string, v float64) {
		name, value = n, v
	}

	// It measures the time in milliseconds between the start time and
	// the current time.
	now := time.Now()
	time.Sleep(20 * time.Millisecond)
	s.Time("timing stat", now, time.Millisecond)

	if name != "timing stat" {
		t.Errorf("expected %v, got %v", "timing stat", name)
	}

	diff := math.Abs(value - 20.0)
	if diff > 2 {
		t.Errorf("expected value to be within %vms, was %v", 2.0, diff)
	}
}

func TestStatHat_sendCount(t *testing.T) {
	s := New()
	ml := iylog.NewMockLogger()
	iylog.Add(ml)
	defer iylog.Reset()

	c := make(chan count, 1)
	s.countC = c

	// When there is room in the buffer is sends the right count down the
	// channel.
	// fmt.Println(len(c))
	s.sendCount("foo", "key", 2)
	select {
	case v := <-c:
		if v.name != "foo" {
			t.Errorf("expected %v, got %v", "foo", v.name)
		}

		if v.key != "key" {
			t.Errorf("expected %v, got %v", "key", v.key)
		}

		if v.n != 2 {
			t.Errorf("expected %v, got %v", 2, v.n)
		}
	case <-time.After(time.Second):
		fmt.Println(ml.Messages())
		t.Error("timed out waiting for count")
	}

	if ml.Called() {
		t.Error("default case should not have been triggered")
	}
	ml.Reset()

	// When there is no room in the buffer is drops the count on the
	// floor, and logs a warning.
	s.sendCount("foo", "key", 2)
	s.sendCount("foo", "key", 2)
	if !ml.Called() {
		t.Error("default case should have been triggered")
	}
}

// TODO(edd): DRY this up with TestStatHat_sendCount
func TestStatHat_sendMeasure(t *testing.T) {
	s := New()
	ml := iylog.NewMockLogger()
	iylog.Add(ml)
	defer iylog.Reset()

	c := make(chan measure, 1)
	s.measureC = c

	// When there is room in the buffer is sends the right count down the
	// channel.
	s.sendMeasure("foo", "key", 2.3)

	select {
	case v := <-c:
		if v.name != "foo" {
			t.Errorf("expected %v, got %v", "foo", v.name)
		}

		if v.key != "key" {
			t.Errorf("expected %v, got %v", "key", v.key)
		}

		if v.v != 2.3 {
			t.Errorf("expected %v, got %v", 2.3, v.v)
		}
	case <-time.After(time.Second):
		t.Error("timed out waiting for count")
	}

	if ml.Called() {
		t.Error("default case should not have been triggered")
	}
	ml.Reset()

	// When there is no room in the buffer is drops the count on the
	// floor, and logs a warning.
	s.sendMeasure("foo", "key", 2.3)
	s.sendMeasure("foo", "key", 2.3)
	if !ml.Called() {
		t.Error("default case should have been triggered")
	}
}
