iymetrics
=========

A small package to help with monitoring metrics on services, and providing implementation to push them to external services.

### MetricsI

Interface for counting / measuring metrics:

```go
type MetricsI interface {
	// Count provides a positive integer of a count at a point in time.
	// Counts tend to be summed up over time.
	Count(name string, i int) error

	// Measure provides a real-value representing an instantanious
	// measurement. Measures tend to be averaged over time.
	Measure(name string, v float64) error

	// Time is a special case of a Measure, which determines the
	// duration between two times and represents it at a time-scaled
	// defined by precision.
	Time(start time.Time, name string, precision time.Duration)
}
```

### Implementations

Currently the following implementations are available:

 - [StatHat](sh/README.md)
