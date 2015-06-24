package iylog

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

var expstr string = "Expected '%s', got '%s'"

func Test_NewMultiLogger(t *testing.T) {
	ml := NewMultiLogger()

	// Logger shouldn't contain any loggables
	x := len(ml.loggables)
	if x > 0 {
		t.Errorf("Expected loggables to be empty, instead found %d loggables", x)
	}

	l := NewLogger(nil, DEBUG)
	ml = NewMultiLogger(l)

	// Logger should contain a single loggable
	x = len(ml.loggables)
	if x != 1 {
		t.Fatalf("Expected length of 1, instead found %d", x)
	}

	// Logger's single loggable should be the one added in the call to NewMultiLogger `lg`
	if !reflect.DeepEqual(l, ml.loggables[0]) {
		t.Errorf("Expected %v, got %v", l, ml.loggables[0])
	}
}

func Test_LoggerLevels(t *testing.T) {
	// test DEBUG logger
	lg := &testLogger{buf: &bytes.Buffer{}, lvl: DEBUG}
	l := NewMultiLogger(lg)

	// test all logging methods are called correctly
	test_allMethods(l, lg, true, true, true, true, t)

	// test all logging methods except DEBUG are called
	lg.lvl = INFO
	test_allMethods(l, lg, true, true, true, false, t)

	// test only Errorf and Warningf are called
	lg.lvl = WARNING
	test_allMethods(l, lg, true, true, false, false, t)

	// test only Errorf is called
	lg.lvl = ERROR
	test_allMethods(l, lg, true, false, false, false, t)
}

func Test_MultipleLoggers(t *testing.T) {
	lg1 := &testLogger{buf: &bytes.Buffer{}, lvl: ERROR}
	lg2 := &testLogger{buf: &bytes.Buffer{}, lvl: WARNING}
	l := NewMultiLogger(lg1, lg2)

	// test both loggers are called with correct [ERROR] message
	test_logger_func(l.Errorf, t, testCase{
		lg:     lg1,
		called: true,
		msg:    "[ERROR] test",
	}, testCase{
		lg:     lg2,
		called: true,
		msg:    "[ERROR] test",
	})

	// test only logger 2 is called with [WARNING] message
	test_logger_func(l.Warningf, t, testCase{
		lg:     lg1,
		called: false,
		msg:    "",
	}, testCase{
		lg:     lg2,
		called: true,
		msg:    "[WARNING] test",
	})

	// test neither logger is logged too on a call to Logger.Infof
	test_logger_func(l.Infof, t, testCase{
		lg:     lg1,
		called: false,
		msg:    "",
	}, testCase{
		lg:     lg2,
		called: false,
		msg:    "",
	})

	// test neither logger is logged too on a call to Logger.Debugf
	test_logger_func(l.Debugf, t, testCase{
		lg:     lg1,
		called: false,
		msg:    "",
	}, testCase{
		lg:     lg2,
		called: false,
		msg:    "",
	})

}

func Test_LogFunctionTypes(t *testing.T) {
	l := NewMultiLogger()

	// set up test logger with an empty buffer
	tl := &testLogger{buf: &bytes.Buffer{}, lvl: DEBUG}
	l.Add(tl)

	// print all log function types to buffer
	printTo(l.Error, l.Errorf, "a")

	// test all error level functions
	obt := tl.buf.String()
	exp := "[ERROR] a[ERROR] {a}"
	if obt != exp {
		t.Errorf(expstr, exp, obt)
	}
	// reset logger buffer
	tl.buf.Reset()

	// print all log function types to buffer
	printTo(l.Warning, l.Warningf, "a")

	// test all error level functions
	obt = tl.buf.String()
	exp = "[WARNING] a[WARNING] {a}"
	if obt != exp {
		t.Errorf(expstr, exp, obt)
	}
	// reset logger buffer
	tl.buf.Reset()

	// print all log function types to buffer
	printTo(l.Info, l.Infof, "a")

	// test all error level functions
	obt = tl.buf.String()
	exp = "[INFO] a[INFO] {a}"
	if obt != exp {
		t.Errorf(expstr, exp, obt)
	}
	// reset logger buffer
	tl.buf.Reset()

	// print all log function types to buffer
	printTo(l.Debug, l.Debugf, "a")

	// test all error level functions
	obt = tl.buf.String()
	exp = "[DEBUG] a[DEBUG] {a}"
	if obt != exp {
		t.Errorf(expstr, exp, obt)
	}
	// reset logger buffer
	tl.buf.Reset()
}

func printTo(log func(...interface{}), logf func(string, ...interface{}), msg string) {
	log(msg)
	logf("{%s}", msg)
}

func Test_CapturePanic(t *testing.T) {
	// set log output to provided buffer

	tl := &testLogger{buf: &bytes.Buffer{}, lvl: ERROR}
	l := NewMultiLogger(tl)

	defer func() {
		// check recovered message exists
		rec := recover()
		if rec == nil {
			t.Errorf("Expected panic 'test', got nil")
		}

		// check it matches original panic string
		if rec.(string) != "test" {
			t.Errorf("Expected panic 'test', got '%s'", rec)
		}

		// ensure correct log message was sent
		msg := tl.buf.String()
		exp := ".*\\[ERROR\\] test"
		if ok, _ := regexp.MatchString(exp, msg); !ok {
			t.Errorf("Expected match %s, got %s", exp, msg)
		}
	}()

	// defer panic capture
	defer l.CapturePanic()
	// panic with string 'test'
	panic("test")
}

// test_allMethods checks that when a single loggable of a
// fixed level is embedded within a Logger struct, that the
// correct behaviour occurs when each of the Logger.<Level>f
// formatted logging functions are called.
func test_allMethods(l *MultiLogger, lg *testLogger, err, wrn, inf, dbg bool, t *testing.T) {
	// test call to Errorf logs as expected, if at all
	test_logger_func(l.Errorf, t, testCase{
		lg:     lg,
		called: err,
		msg:    "[ERROR] test",
	})
	// test call to Warningf logs as expected, if at all
	test_logger_func(l.Warningf, t, testCase{
		lg:     lg,
		called: wrn,
		msg:    "[WARNING] test",
	})
	// test call to Infof logs as expected, if at all
	test_logger_func(l.Infof, t, testCase{
		lg:     lg,
		called: inf,
		msg:    "[INFO] test",
	})

	// test call to Debugf logs as expected, if at all
	test_logger_func(l.Debugf, t, testCase{
		lg:     lg,
		called: dbg,
		msg:    "[DEBUG] test",
	})

}

// testCase is used for checking that a testLogger
// is interacted with correctly, when embedded within
// a iylog.Logger.
type testCase struct {
	lg     *testLogger
	called bool
	msg    string
}

// test_logger_func checks to see when a provided printf
// function is called that all the provided testCase
// structures criteria
func test_logger_func(logf func(f string, v ...interface{}), t *testing.T, tests ...testCase) {
	logf("test")

	for _, test := range tests {
		if test.called {
			resp := test.lg.buf.String()
			if resp != test.msg {
				t.Errorf(expstr, test.msg, resp)
			}
		} else {
			// Loggable should not have been called
			if test.lg.called {
				t.Errorf("Loggable.Printf was called!")
			}
		}

		test.lg.called = false
		test.lg.buf.Reset()
	}
}

// testLogger implements iylog.Logger
//
// It contains a buffer, logging Level and a boolean
// which is set to true when a call to Printf is made
type testLogger struct {
	buf    *bytes.Buffer
	lvl    Level
	called bool
}

// Printf prints the formatted message to the underlying buffer
// and sets testLogger.called to true
func (t *testLogger) Printf(f string, v ...interface{}) {
	t.called = true
	fmt.Fprintf(t.buf, f, v...)
}

//Level return the testLogger.lvl
func (t *testLogger) Level() Level {
	return t.lvl
}
