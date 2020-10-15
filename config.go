package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type MainConfig struct {
	Mode          string
	SecretsEngine string
}

func LoadMainConfig() *MainConfig {
	filePtr, _ := os.Open("main.json")
	defer filePtr.Close()

	var conf MainConfig
	b, err := ioutil.ReadAll(filePtr)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, &conf)
	if err != nil {
		panic(err)
	}

	return &conf
}
