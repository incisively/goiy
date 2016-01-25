package iytest

import (
	"encoding/json"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

const (
	stackSize = 1024
)

func init() {
	spew.Config.SortKeys = true
}

// getTrace generates a stack trace, and snips the stack lines referring
// to functions in this package. The result is that package-users see
// the call to this package from test function as the last call in the
// stack. This makes it easier for them to quickly jump to the failing
// test.
//
// TODO(edd): This should be done properly, using runtime.Callers to get
// the correct caller counters.
func getTrace() string {
	trace := make([]byte, stackSize)
	runtime.Stack(trace, false)

	lines := strings.Split(string(trace), "\n")
	lines = append(lines[:1], lines[7:]...)
	return strings.Join(lines, "\n")
}

func equal(f func(string, ...interface{}), actual, expected interface{}, i ...int) {
	if actual != expected {
		var format string
		if len(i) > 0 {
			format = "\n[Example %d]\n got %s\nwanted %s\ntrace: %s\n"
			f(format, i[0], spew.Sdump(actual), spew.Sdump(expected), getTrace())
		} else {
			format = "\ngot %s\nwanted %s\ntrace: %s\n"
			f(format, spew.Sdump(actual), spew.Sdump(expected), getTrace())
		}
	}
}

// Equal checks actual and expected are equal, and calls t.Errorf if
// they're not.
//
// Callers can optionally pass in an example number, which will then be
// included in any failure message. Only the first value for i is used.
func Equal(t *testing.T, actual, expected interface{}, i ...int) {
	equal(t.Errorf, actual, expected, i...)
}

// Equalf is similar to Equal, but calls t.Fatalf on failure.
func Equalf(t *testing.T, actual, expected interface{}, i ...int) {
	equal(t.Fatalf, actual, expected, i...)
}

func notEqual(f func(string, ...interface{}), actual, dontwant interface{}) {
	if actual == dontwant {
		format := "\ngot %s but didn't want it\ntrace: %s\n"
		f(format, spew.Sdump(actual), getTrace())
	}
}

// NotEqual checks actual and dontwant are not equal, and calls
// t.Errorf if they are.
func NotEqual(t *testing.T, actual, dontwant interface{}) {
	notEqual(t.Errorf, actual, dontwant)
}

// NotEqualf is similar to NotEqual, but calls t.Fatalf on failure.
func NotEqualf(t *testing.T, actual, dontwant interface{}) {
	notEqual(t.Fatalf, actual, dontwant)
}

func isNil(f func(string, ...interface{}), actual interface{}, notNil bool) {
	var (
		v      = reflect.ValueOf(actual)
		format = "\nvalue is not nil:\n%s\ntrace: %s\n"
	)

	// Switch format string when checking if not nil.
	if notNil {
		format = "\nvalue is nil:\n%s\ntrace: %s\n"
	}

	var isNil bool
	if actual == nil {
		isNil = true
	} else {
		switch v.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
			isNil = v.IsNil()
		}
	}

	if notNil == isNil {
		f(format, spew.Sdump(actual), getTrace())
	}
}

// IsNil checks that actual is equivalent to nil.
func IsNil(t *testing.T, actual interface{}) {
	isNil(t.Errorf, actual, false)
}

// IsNilf is similar to IsNil, but calls t.Fatalf on failure.
func IsNilf(t *testing.T, actual interface{}) {
	isNil(t.Fatalf, actual, false)
}

// NotNil checks that actual is not equivalent to nil.
func NotNil(t *testing.T, actual interface{}) {
	isNil(t.Errorf, actual, true)
}

// NotNilf is similar to NotNil, but calls t.Fatalf on failure.
func NotNilf(t *testing.T, actual interface{}) {
	isNil(t.Fatalf, actual, true)
}

func deepEqual(f func(string, ...interface{}), actual, expected interface{}, i ...int) {
	if !reflect.DeepEqual(actual, expected) {
		var format string
		if len(i) > 0 {
			format = "\n[Example %d]\n got %s\nwanted %s\ntrace: %s\n"
			f(format, i[0], spew.Sdump(actual), spew.Sdump(expected), getTrace())
		} else {
			format = "\ngot %s\nwanted %s\ntrace: %s\n"
			f(format, spew.Sdump(actual), spew.Sdump(expected), getTrace())
		}
	}
}

// DeepEqual checks actual and expected are equal using reflect.DeepEqual,
// and calls t.Errorf if they're not.
//
// Callers can optionally pass in an example number, which will then be
// included in any failure message. Only the first value for i is used.
func DeepEqual(t *testing.T, actual, expected interface{}, i ...int) {
	deepEqual(t.Errorf, actual, expected, i...)
}

// DeepEqualf is similar to DeepEqual, but calls t.Fatalf on failure.
func DeepEqualf(t *testing.T, actual, expected interface{}, i ...int) {
	deepEqual(t.Fatalf, actual, expected, i...)
}

func jSONKeysEqual(f func(string, ...interface{}), data []byte, expected []string) {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		f("%v", err)
	}

	var keys []string
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if len(keys) != len(expected) {
		f("\ngot %#v, wanted %#v", keys, expected)
		return
	}
	sort.Strings(expected)

	for i := 0; i < len(keys); i++ {
		if keys[i] != expected[i] {
			f("\ngot %q, wanted %q\n", keys[i], expected[i])
			return
		}
	}
}

// JSONKeysEqual checks that data (which must be a single JSON object)
// contains exactly the keys in expected. If the keys differ,
// JSONKeysEqual calls t.Errorf on the first difference.
//
// expected does not need to be ordered.
func JSONKeysEqual(t *testing.T, data []byte, expected []string) {
	jSONKeysEqual(t.Errorf, data, expected)
}

// JSONKeysEqualf is similar to JSONKeysEqual, but calls t.Fatalf on
// failure.
func JSONKeysEqualf(t *testing.T, data []byte, expected []string) {
	jSONKeysEqual(t.Fatalf, data, expected)
}

// TODO(edd): DRY this up with jSONKeysEqual.
func jSONKeysEqualI(f func(string, ...interface{}), data []byte, expected []string, i int) {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		f("[Example %d] %v", i, err)
	}

	var keys []string
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if len(keys) != len(expected) {
		f("[Example %d]\ngot %#v, wanted %#v", i, keys, expected)
		return
	}
	sort.Strings(expected)

	for i := 0; i < len(keys); i++ {
		if keys[i] != expected[i] {
			f("[Example %d]\ngot %q, wanted %q\n", i, keys[i], expected[i])
			return
		}
	}
}

// JSONKeysEqualExample checks that data (which must be a single JSON
// object) contains exactly the keys in expected. If the keys differ,
// JSONKeysEqualExample calls t.Errorf on the first difference. The
// caller can provide the example number (0-based); a failiure message
// will reference it.
func JSONKeysEqualExample(t *testing.T, data []byte, expected []string, i int) {
	jSONKeysEqualI(t.Errorf, data, expected, i)
}

// JSONKeysEqualExamplef is similar to JSONKeysEqualExample, but calls
// t.Fatalf on failure.
func JSONKeysEqualExamplef(t *testing.T, data []byte, expected []string, i int) {
	jSONKeysEqualI(t.Fatalf, data, expected, i)
}
