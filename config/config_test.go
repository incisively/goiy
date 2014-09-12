package config

import (
	"launchpad.net/gocheck"
	"strings"
	"testing"
)

type ConfigSuite struct{}

var _ = gocheck.Suite(&ConfigSuite{})

func Test(t *testing.T) { gocheck.TestingT(t) }

func (cs *ConfigSuite) TestUnmarshal(c *gocheck.C) {
	// test Unmarshal works with defined env string
	var testconf, prodconf Config

	err := Unmarshal(jsondata, &testconf, "unknown")
	c.Assert(err, gocheck.Not(gocheck.IsNil))
	c.Check(err.Error(), gocheck.DeepEquals, "Cannot find env unknown in config")

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

func (cs *ConfigSuite) TestDecoderUnmarshals(c *gocheck.C) {
	// test Unmarshal works with defined env string
	var testconf, prodconf Config

	// test decoder fails on unrecognised env
	reader := strings.NewReader(string(jsondata))
	dec := NewDecoder(reader)
	err := dec.Decode(&testconf, "unknown")
	c.Assert(err, gocheck.Not(gocheck.IsNil))
	c.Check(err.Error(), gocheck.DeepEquals, "Cannot find env unknown in config")

	// test Decoder decodes test config
	reader = strings.NewReader(string(jsondata))
	dec = NewDecoder(reader)
	err = dec.Decode(&testconf, "test")
	c.Assert(err, gocheck.IsNil)
	c.Check(&testconf, gocheck.DeepEquals, &Config{
		A: "some kind of string",
		B: 100,
	})

	// test Decoder decodes production config
	reader = strings.NewReader(string(jsondata))
	dec = NewDecoder(reader)
	err = dec.Decode(&prodconf, "production")
	c.Assert(err, gocheck.IsNil)
	c.Check(&prodconf, gocheck.DeepEquals, &Config{
		A: "production worthy string",
		B: 9001,
	})

	// test Decoder works for undefined env string
	var envconf EnvConfig

	// test Decoder decodes entire config on empty env string
	reader = strings.NewReader(string(jsondata))
	dec = NewDecoder(reader)
	err = dec.Decode(&envconf, "")
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
