package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	configFile = "./config/config.json"
)

func Load(isDev bool) (*Config, error) {
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		return &Config{}, err
	}

	c := &Config{}
	err = json.Unmarshal(raw, &c)
	if err != nil {
		return &Config{}, err
	}
	c.DevEnv = isDev

	return c, nil
}

//func LoadCustom(file string) (*Config, error) {
//	raw, err := ioutil.ReadFile(file)
//	if err != nil {
//		return &Config{}, err
//	}
//
//	c := &Config{}
//	err = json.Unmarshal(raw, &c)
//	if err != nil {
//		return &Config{}, err
//	}
//
//	return c, nil
//}
