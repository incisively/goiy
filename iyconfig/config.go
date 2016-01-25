package iyconfig

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Unmarshal wraps json.Unmarshal.
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// UnmarshalEnv parses a JSON configuration file where multiple
// environments are specified.
//
// Specifically, it expects a file with a layout similar to the
// following:
//    {
//      "production": {},
//      "dev": {},
//      "local": {}
//    }
//
// UnmarshalEnv will then use the provided `env` variable to determine
// which configuration to unmarshal into v.
func UnmarshalEnv(data []byte, v interface{}, env string) error {
	var envs map[string]json.RawMessage

	// Initially unmarshal everything into a map of raw JSON messages.
	if err := json.Unmarshal(data, &envs); err != nil {
		return err
	}

	conf, ok := envs[env]
	if !ok {
		return EnvNotFoundError{env}
	}

	// Unmarshal the data in to the provided interface `v`
	if err := json.Unmarshal(conf, v); err != nil {
		return err
	}

	return nil
}

// Decoder unmarshals config data from an io.Reader
// into a target struct type.
type Decoder struct {
	*json.Decoder
}

// NewDecoder initialises a new Decoder from the provided reader.
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

// DecodeFromFile unmarshals the env portion of the JSON configuration
// in the file found at path into dest.
func DecodeFromFile(pth string, dest interface{}, env string) error {
	// open the file
	f, err := os.Open(pth)
	if err != nil {
		return err
	}

	if err := NewDecoder(f).DecodeEnv(dest, env); err != nil {
		return err
	}
	return f.Close()
}

// DecodeFromFileP is like DecodeFromFile but panics on error.
func DecodeFromFileP(pth string, dest interface{}, env string) {
	if err := DecodeFromFile(pth, dest, env); err != nil {
		panic(err)
	}
}

// EnvNotFoundError is returned when an environment requested
// to be unmarshalled is not found in the provided data.
type EnvNotFoundError struct {
	env string
}

func (e EnvNotFoundError) Error() string {
	return fmt.Sprintf("Cannot find env '%s' in config", e.env)
}
