package config

import (
	"encoding/json"
	"io/ioutil"

	"hobee-be/pkg/herrors"
)

const (
	configFile = "./config/config.json"
)

func Load() (*Config, error) {
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		return &Config{}, err
	}

	c := &Config{}
	err = json.Unmarshal(raw, &c)
	if err != nil {
		return &Config{}, err
	}

	if c.Port == "" {
		return &Config{}, herrors.New("Loading config failed")
	}

	return c, nil
}
