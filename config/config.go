package config

import (
	"encoding/json"
	"fmt"
	"io"
)

// Unmarshal parses out a json configuration file which is aggregated
// by an environment string. If no env string is provided it unmarshals
// the entire data []byte in to the provided interface v.
func Unmarshal(data []byte, v interface{}, env string) error {
	// if no env string is provided parse as normal json
	if env == "" {
		return json.Unmarshal(data, v)
	}
	// if env string provided use env lookup approach
	return unmarshal(data, v, env)
}

// unmarshal a configuration json file, but only the top level
// object with the key `env`
func unmarshal(data []byte, v interface{}, env string) (err error) {
	// construct a new map of string to json raw message type
	var envs map[string]json.RawMessage

	// initially unmarshal environments
	if err = json.Unmarshal(data, &envs); err != nil {
		return
	}

	// if the environment desired is not present return an error
	conf, ok := envs[env]
	if !ok {
		return fmt.Errorf("Cannot find env %s in config", env)
	}

	// unmarshal the data in to the provided interface `v`
	if err = json.Unmarshal(conf, v); err != nil {
		return
	}

	return
}

type Decoder struct {
	*json.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		Decoder: json.NewDecoder(r),
	}
}

func (dec *Decoder) Decode(v interface{}, env string) error {
	if env == "" {
		return dec.Decoder.Decode(v)
	}
	return dec.decode(v, env)
}

func (dec *Decoder) decode(v interface{}, env string) (err error) {
	// construct a new map of string to json raw message type
	var envs map[string]json.RawMessage

	// initially unmarshal environments
	if err = dec.Decoder.Decode(&envs); err != nil {
		return
	}

	// if the environment desired is not present return an error
	conf, ok := envs[env]
	if !ok {
		return fmt.Errorf("Cannot find env %s in config", env)
	}

	// unmarshal the data in to the provided interface `v`
	if err = json.Unmarshal(conf, v); err != nil {
		return
	}

	return
}
