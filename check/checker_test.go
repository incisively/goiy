package check

import (
	"launchpad.net/gocheck"
	"testing"
)

type CheckSuite struct{}

var _ = gocheck.Suite(&CheckSuite{})

func Test(t *testing.T) { gocheck.TestingT(t) }

func (cs *CheckSuite) Test_JsonHasKey_HandlesBadArgs(c *gocheck.C) {
	// build json checker
	HasKey := jsonHasKey{}

	// check info is as expected
	c.Check(HasKey.Info(), gocheck.DeepEquals, &gocheck.CheckerInfo{
		Name:   "JsonHasKey",
		Params: []string{"obtained", "key"},
	})

	// A call to check with bad arguments should error

	// first argument should be string/[]byte
	res, err := HasKey.Check([]interface{}{1}, []string{})
	c.Assert(res, gocheck.Equals, false)
	c.Check(err, gocheck.Equals, "First argument must be either a string or []byte not int")

	// second argument should be a string
	res, err = HasKey.Check([]interface{}{"one", 1}, []string{})
	c.Assert(res, gocheck.Equals, false)
	c.Check(err, gocheck.Equals, "Second argument must be a string not int")

	// first argument must be valid JSON
	res, err = HasKey.Check([]interface{}{"one", "two"}, []string{})
	c.Assert(res, gocheck.Equals, false)
	c.Check(err, gocheck.Equals, "invalid character 'o' looking for beginning of value")
}

func (cs *CheckSuite) Test_JsonHasKey(c *gocheck.C) {
	// build json checker
	HasKey := jsonHasKey{}

	// check it returns false for a missing key
	res, err := HasKey.Check([]interface{}{jsondata, "three"}, []string{})
	c.Check(res, gocheck.Equals, false)
	c.Check(err, gocheck.Equals, "")

	// check it return true for a present key
	res, err = HasKey.Check([]interface{}{jsondata, "one"}, []string{})
	c.Check(res, gocheck.Equals, true)
	c.Check(err, gocheck.Equals, "")

	// check it behaves as a Checker
	c.Check(jsondata, JsonHasKey, "one")
	c.Check(jsondata, gocheck.Not(JsonHasKey), "three")
}

var jsondata string = `
{
    "one": "",
    "two": "",
    "four": ""
}`
