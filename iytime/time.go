package iytime

import (
	"encoding/json"
	"fmt"
	"time"
)

// A Duration is a drop-in replacement for a time.Duration, and
// represents the elapsed time between two instants as an int64
// nanosecond count.
//
// Duration implements the json.Marshaler and json.Unmarshaler
// interfaces, as well as providing method wrappers for the underlying
// time.Duration methods.
type Duration time.Duration

// MarshalJSON implements the json.Marshaler.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprint(d))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(parsed)
	return nil
}

// Hours returns the duration as a floating point number of hours.
func (d Duration) Hours() float64 { return time.Duration(d).Hours() }

// Minutes returns the duration as a floating point number of minutes.
func (d Duration) Minutes() float64 { return time.Duration(d).Minutes() }

// Nanoseconds returns the duration as an integer nanosecond count.
func (d Duration) Nanoseconds() int64 { return time.Duration(d).Nanoseconds() }

// Seconds returns the duration as a floating point number of seconds.
func (d Duration) Seconds() float64 { return time.Duration(d).Seconds() }

// String returns a string representing the duration in the form "72h3m0.5s".
// Leading zero units are omitted.  As a special case, durations less than one
// second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure
// that the leading digit is non-zero.  The zero duration formats as 0,
// with no unit.
func (d Duration) String() string { return time.Duration(d).String() }
