## StatHat

The StatHat implementation ships counts and measures off to StatHat. The StatHat implementation will look for an environment variable `SH_KEY` if a StatHat API key is not provided to `NewStatHat`.

You can either use the package-level instance, which works much in the same way as [log.Logger](http://golang.org/pkg/log/) does, or you can create your own `StatHat` instance.

### Monitoring Runtime

With `StatHat` you can also setup automatic monitoring of certain aspects of the runtime.
Simply make a call to `MonitorRuntime` and pass in a duration with which the `StatHat` should send the metrics.

### Example Usage

Using the package-level implementation:

```go
package main

import (
	"time"

	"github.com/incisively/goiy/iymetrics/sh"
)

func main() {
	// Set a prefix for the package-level instance.
	sh.SetPrefix("[service-a]")

	// If you have the SH_API key in the environment you don't
	// need to call this.
	sh.SetAPIKey("ssiuhIUYDGos")

	// Monitor the process runtime every 2 minutes.
	sh.MonitorRuntime(2 * time.Minute)

	// Sends a count of 1 to the StatHat service with the stat name
	// "[service-a] users".
	sh.Count("users", 1)

	// Sends a measure to the StatHat service with the stat name
	// "[service-a] length".
	sh.Measure("length", 24.22)

	// Sends a timing to the StatHat service with the name
	// "[service-a] work-ms".
	now := time.Now()
	go func() {
		// Do some work in another goroutine.
		// Work()
		defer sh.Time("work-ms", now, time.Millisecond)
	}()
}
```
