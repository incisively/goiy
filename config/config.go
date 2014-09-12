package config

import (
	"encoding/json"
	"fmt"
)

func Unmarshal(data []byte, v interface{}, env string) (err error) {
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
