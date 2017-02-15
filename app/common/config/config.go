package config

import (
	"encoding/json"
	"io/ioutil"
)

var c Configuration

func ReadConfig(f string) (Configuration, error) {

	data, err := ioutil.ReadFile(f)

	config := Configuration{}

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &c)

	if err != nil {
		return config, err
	}

	return c, nil
}
