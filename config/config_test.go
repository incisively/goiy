package config

import (
	"fmt"
	"launchpad.net/gocheck"
	"testing"
)

type ConfigSuite struct {
	config []byte
}

var _ = gocheck.Suite(&ConfigSuite{})

func Test(t *testing.T) { gocheck.TestingT(t) }

func (cs *ConfigSuite) SetUpSuite(c *gocheck.C) {
	cs.config = []byte(`
test:
    a: "Some stuff in a string"
    b: 9001
`)
}

func (cs *ConfigSuite) TestDecoderUnmarshals(c *gocheck.C) {
	// test Unmarshal works with defined env string
	var testconf, prodconf Config

	err := Unmarshal(jsondata, &testconf, "unknown")
	c.Assert(err, gocheck.Not(gocheck.IsNil))
	c.Check(err, gocheck.DeepEquals, fmt.Errorf("Cannot find env unknown in config"))

	err = Unmarshal(jsondata, &testconf, "test")
	c.Assert(err, gocheck.IsNil)
	c.Check(&testconf, gocheck.DeepEquals, &Config{
		A: "some kind of string",
		B: 100,
	})

	err = Unmarshal(jsondata, &prodconf, "production")
	c.Assert(err, gocheck.IsNil)
	c.Check(&prodconf, gocheck.DeepEquals, &Config{
		A: "production worthy string",
		B: 9001,
	})

	// test Unmarshal works for undefined env string
	var envconf EnvConfig

	err = Unmarshal(jsondata, &envconf, "")
	c.Assert(err, gocheck.IsNil)
	c.Check(&envconf, gocheck.DeepEquals, &EnvConfig{
		Test: Config{
			A: "some kind of string",
			B: 100,
		},
		Prod: Config{
			A: "production worthy string",
			B: 9001,
		},
	})

}

var jsondata []byte = []byte(`
{
	"test": {
		"a": "some kind of string",
		"b": 100
	},
	"production": {
		"a": "production worthy string",
		"b": 9001
	}
}`)

type EnvConfig struct {
	Test Config `json:"test"`
	Prod Config `json:"production"`
}

type Config struct {
	A string `json:"a"`
	B int    `json:"b"`
}
