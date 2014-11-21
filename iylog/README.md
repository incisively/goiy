Package iylog
=============

`iylog` — simple levelled logging.

`iylog` provides a `MultiLogger` type, which manages multiple levelled
loggers. The package provides a very similar API to Go's `log` package,
with the obvious extension that you can specify what level you wish to
to log things at.

Only loggers set to listen at the level logged or below will respond to
emitted messages.

As with Go's `log` package, `iylog` provides a package-level
`MultiLogger`, which is safe for use from multiple goroutines.

### Examples

#### Simple Usage

```go
package main

import (
	"log"
	"os"

	"github.com/incisively/goiy/iylog"
)

func main() {
	// initially the iylog.MultiLogger is empty
	iylog.Info("foo")
	// writes nothing

	// add a new iylog.Logger set at warning level
	iylog.Add(iylog.NewLogger(nil, iylog.WARNING))
	// and one at a lower level
	iylog.Add(iylog.NewLogger(nil, iylog.INFO))

	// “[ERROR] something went wrong” to both loggers
	iylog.Errorf("Something %s %s", "went", "wrong")

	// will not be emitted by either logger
	iylog.Debug(22)

	// will only be emitted by the first logger
	iylog.Warningln("a warning")

	// logging to a file is simply a case of setting the output of
	// a `log.Logger`, before you add it to the `iylog.MultiLogger`
	f, err := os.Create("/tmp/out.log")
	if err != nil {
		panic(err)
	}

	l := iyLog.NewLogger(log.New(f, "", log.LstdFlags))
	iylog.Add(l, iylog.DEBUG)

	// will be written to f using l, and also emitted to os.Stderr on
	// by previously added loggers.
	iylog.Debugln("some debug message")
	f.Close()
}
```

#### Custom Loggable Implementation

```go
package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/incisively/goiy/iylog"
)

type WebbyLogger struct {
	url string
	buf bytes.Buffer
	lvl iylog.Level
}

func (l *WebbyLogger) Printf(format string, v ...interface{}) {
	l.b.Reset()

	_, err := l.b.WriteString(fmt.Sprintf(format, v...))
	if err != nil {
		log.Println(err)
		return
	}

	if _, err = http.Post(l.url, "text/plain", &l.b); err != nil {
		log.Println(err)
	}
}

func (l *WebbyLogger) Level() iylog.Level {
	return l.lvl
}

func main() {
	l := &WebbyLogger{
		url: "http://some.location.com/log",
		lvl: iylog.INFO,
	}

	// add logger
	iylog.Add(l)

	// message will be POSTED to http://some.location.com/log
	iylog.Infof("2 + 2 = %v\n", 4)
}
```

