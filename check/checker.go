package check

import (
	"encoding/json"
	"fmt"

	"launchpad.net/gocheck"
)

// Checker var JsonHasKey
var JsonHasKey gocheck.Checker = &jsonHasKey{}

// jsonHasKey implements gocheck.Checker. It is used
// to check whether or not a string/[]byte when unmarshalled
// from JSON has a particular key present or not
type jsonHasKey struct{}

// Check checks that the argument provided are as expected.
// Unmarshals params[0] from a string or a []byte JSON blob in to a map.
// It then checks whether the key in params[1] is present within
// the map.
func (j *jsonHasKey) Check(params []interface{}, names []string) (bool, string) {
	var data []byte

	switch d := params[0].(type) {
	case string:
		data = []byte(d)
	case []byte:
		data = d
	default:
		return false, fmt.Sprintf("First argument must be either a string or []byte not %T", d)
	}

	key, ok := params[1].(string)
	if !ok {
		return false, fmt.Sprintf("Second argument must be a string not %T", params[1])
	}

	var mapj map[string]interface{}
	if err := json.Unmarshal(data, &mapj); err != nil {
		return false, err.Error()
	}

	_, ok = mapj[key]
	return ok, ""
}

// Info returns necessary information for gocheck to
// present a useful human-readable message when something
// doesn't check out.
func (j *jsonHasKey) Info() *gocheck.CheckerInfo {
	return &gocheck.CheckerInfo{
		Name:   "JsonHasKey",
		Params: []string{"obtained", "key"},
	}
}
