package stat

import (
	"log"
	"time"

	"github.com/e-dard/gev"
	"github.com/stathat/go"
)

// StatHat implements the Statter interface for the Stat Hat service
type StatHat struct {
	key string `env:"SH_KEY"`
}

var std Statter = NewStatHat("")

// NewStatHat returns a new StatHat type.
//
// If key is not the empty string, then it will be used as the Stat Hat
// API key. Otherwise, NewStatHat will attempt to read the Stat Hat API
// key from the environment, looking for a SH_KEY variable.
//
// NewStatHat panics if there is a problem reading this variable, though
// it won't panic if the variable is missing from the environment.
func NewStatHat(key string) (s StatHat) {
	if key != "" {
		s.key = key
		return
	}

	str := struct {
		Key string `env:"SH_KEY"`
	}{}

	if err := gev.Unmarshal(&str); err != nil {
		panic(err)
	}

	s.key = str.Key

	return
}

// Count increments a Stat Hat counter by n. It's threadsafe, and will
// not make a call if the Stat Hat API key is not present.
func (s StatHat) Count(stat string, n int) error {
	if s.key == "" {
		return nil
	}
	return stathat.PostEZCount(stat, s.key, n)
}

// Measure sends a value to a Stat Hat value. It's threadsafe, and will
// not make a call if the Stat Hat API key is not present.
func (s StatHat) Measure(stat string, v float64) error {
	if s.key == "" {
		return nil
	}
	return stathat.PostEZValue(stat, s.key, v)
}

// Time is a function to measure the time between `start` and when
// this function is called. The result is sent to StatHat.
// The intention for this function is to be used within a `defer`, e.g.
//	now := time.Now()
//	defer stat.TimeStat(now, "Timing Something", time.Millisecond)
func (s StatHat) Time(start time.Time, stat string, dur time.Duration) {
	tms := time.Since(start) / dur
	go func() {
		if err := s.Measure(stat, float64(tms)); err != nil {
			log.Println(err)
		}
	}()
}

func Count(stat string, n int) error {
	return std.Count(stat, n)
}

func Measure(stat string, v float64) error {
	return std.Measure(stat, v)
}

func Time(start time.Time, stat string, dur time.Duration) {
	std.Time(start, stat, dur)
}
