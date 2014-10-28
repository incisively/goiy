Package `goiy/iylog`
=================

`goiy/iylog` is a logging package, with similarities to go stdlib `log`.
It contains the `Logger` type, the `Loggable` interface and the `Level` enum.
Just importing the package is enough to register loggables to the packages standard `Logger` using `iylog.Add(loggables ...Loggable)`.

### Example Usage

```go
package main

import (
    "log"

    "github.com/incisively/goiy/iylog"
)

type Logger struct {
    logger *log.Logger
    level  iylog.Level
}

func NewLogger(out io.Writer, level iylog.Level) *Logger {
    return &Logger{
        logger: log.New(out, "", log.LstdFlags),
        level: level,
    }
}

func (s *Logger) Printf(format string, v ...interface{}) {
    return s.logger.Printf(format, v...)
}

func (s *Logger) Level() log.Level {
    return s.level
}

func main() {
    iylog.Add(NewLogger(os.Stdout, iylog.WARNING))
    iylog.Add(NewLogger(os.Stderr, iylog.INFO))

    // “[ERROR] Something went wrong” to both Stdout and Stderr
    iylog.Errorf("Something %s %s", "went", "wrong")
    // prints nowhere as level is too high
    iylog.Debugf("stuff is happening")
    // “[INFO] Some info” to Stderr, but not Stdout
    iylog.Infof("Some info")
}

```
