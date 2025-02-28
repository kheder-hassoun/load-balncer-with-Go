package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	HealthCheckInterval string   `json:"healthCheckInterval"`
	Servers             []Server `json:"servers"`
	ListenPort          string   `json:"listenPort"`
}

func loadConfig(file string) (Config, error) {
	var config Config

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
