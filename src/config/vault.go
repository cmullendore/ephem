package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Vault struct {
}

func LoadVaultConfig() *Vault {
	filePtr, _ := os.Open("vault.json")
	defer filePtr.Close()

	var conf Vault
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
