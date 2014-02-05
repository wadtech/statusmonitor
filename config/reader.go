package config

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/wadtech/statusmonitor/service"
)

type Config struct {
	Filepath string
	Config   Data
}

type Data struct {
	Port     string            `json:port`
	Workers  int               `json:workers`
	Delay    int               `json:worker_delay`
	Services []service.Service `json:services`
}

func NewConfig(filepath string) (r *Config) {
	config, err := readConfig(filepath)
	if err != nil {
		//dead in the water, we couldn't find or unmarshal the config file so crash!
		log.Panic("configuration failed to initialise with message: ", err)
	}

	return &Config{filepath, config}
}

func readConfig(filepath string) (config Data, err error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
