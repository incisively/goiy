package iyconfig

import (
	"io/ioutil"
	"strings"
	"testing"

	"launchpad.net/gocheck"
)

type ConfigSuite struct{}

var _ = gocheck.Suite(&ConfigSuite{})

func Test(t *testing.T) { gocheck.TestingT(t) }

func (cs *ConfigSuite) TestUnmarshal(c *gocheck.C) {
	// test Unmarshal works with defined env string
	var testconf, prodconf Config

	err := UnmarshalEnv(jsondata, &testconf, "unknown")
	c.Assert(err, gocheck.Not(gocheck.IsNil))
	c.Check(err.Error(), gocheck.DeepEquals, "Cannot find env 'unknown' in config")

	err = UnmarshalEnv(jsondata, &testconf, "test")
	c.Assert(err, gocheck.IsNil)
	c.Check(&testconf, gocheck.DeepEquals, &Config{
		A: "some kind of string",
		B: 100,
	})

	err = UnmarshalEnv(jsondata, &prodconf, "production")
	c.Assert(err, gocheck.IsNil)
	c.Check(&prodconf, gocheck.DeepEquals, &Config{
		A: "production worthy string",
		B: 9001,
	})

	// test Unmarshal works for undefined env string
	var envconf EnvConfig

	err = UnmarshalEnv(jsondata, &envconf, "")
	c.Assert(err, gocheck.Not(gocheck.IsNil))
	c.Check(err.Error(), gocheck.DeepEquals, "Cannot find env '' in config")

	err = Unmarshal(jsondata, &envconf)
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
	err := dec.DecodeEnv(&testconf, "unknown")
	c.Assert(err, gocheck.Not(gocheck.IsNil))
	c.Check(err.Error(), gocheck.DeepEquals, "Cannot find env 'unknown' in config")

	// test Decoder decodes test config
	reader = strings.NewReader(string(jsondata))
	dec = NewDecoder(reader)
	err = dec.DecodeEnv(&testconf, "test")
	c.Assert(err, gocheck.IsNil)
	c.Check(&testconf, gocheck.DeepEquals, &Config{
		A: "some kind of string",
		B: 100,
	})

	// test Decoder decodes production config
	reader = strings.NewReader(string(jsondata))
	dec = NewDecoder(reader)
	err = dec.DecodeEnv(&prodconf, "production")
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
	err = dec.DecodeEnv(&envconf, "")
	c.Assert(err, gocheck.NotNil)
	c.Check(err.Error(), gocheck.Equals, "Cannot find env '' in config")

	dec = NewDecoder(strings.NewReader(string(jsondata)))
	err = dec.Decode(&envconf)
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

func (cs *ConfigSuite) TestDecodeFromFileP(c *gocheck.C) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}

	if _, err := f.Write(jsondata); err != nil {
		panic(err)
	}

	conf := Config{}
	DecodeFromFileP(f.Name(), &conf, "production")
	c.Check(conf, gocheck.DeepEquals, Config{
		A: "production worthy string",
		B: 9001,
	})
}

var jsondata = []byte(`
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
