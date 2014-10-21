package config

import (
	"encoding/json"
	"fmt"
	"io"
)

// Unmarshal wraps json.Unmarshal
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Unmarshal parses out a json configuration file which is segmented
// by an environment string.
func UnmarshalEnv(data []byte, v interface{}, env string) error {
	// construct a new map of string to json raw message type
	var envs map[string]json.RawMessage

	// initially unmarshal environments
	if err := json.Unmarshal(data, &envs); err != nil {
		return err
	}

	// if the environment desired is not present return an error
	conf, ok := envs[env]
	if !ok {
		return EnvNotFoundError{env}
	}

	// unmarshal the data in to the provided interface `v`
	if err := json.Unmarshal(conf, v); err != nil {
		return err
	}

	return nil
}

// Decoder unmarshals config data from an io.Reader
// into a target struct type
type Decoder struct {
	*json.Decoder
}

// NewDecoder takes an io.Reader to unmarshal
// and return a pointer to a new Decoder
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		Decoder: json.NewDecoder(r),
	}
}

// Decode performs the unmarshaling on the contents of the io.Reader
// in to the target interface `v` for a given environment `env`.
func (dec *Decoder) DecodeEnv(v interface{}, env string) error {
	return dec.decode(v, env)
}

// decode performs the decoding when `env` is not the empty string.
// If the env key is not found in the unmarshalled result it will
// return an error.
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
		return EnvNotFoundError{env}
	}

	// unmarshal the data in to the provided interface `v`
	if err = json.Unmarshal(conf, v); err != nil {
		return
	}

	return
}

// EnvNotFoundError is returned when an environment requested
// to be unmarshalled is not found in the provided data.
type EnvNotFoundError struct {
	env string
}

// Error returns a description of the unknown environment.
func (e EnvNotFoundError) Error() string {
	return fmt.Sprintf("Cannot find env '%s' in config", e.env)
}
