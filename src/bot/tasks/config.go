package tasks

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Webhook string `json:"WEBHOOK"`
	Delay   int    `json:"DELAY"`
	Timeout int    `json:"TIMEOUT"`
}

func ReadConfig() (Config, error) {

	fileContent, err := ioutil.ReadFile("config.json")

	if err != nil {
		return Config{}, err
	}

	var config Config

	if err := json.Unmarshal(fileContent, &config); err != nil {
		return config, err
	}

	return config, nil

}
