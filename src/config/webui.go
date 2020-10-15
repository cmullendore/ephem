package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type WebUI struct {
	Listener     string
	ListenerCert Cert
	APIURL       string
	APICert      Cert
	PathLength   int
}

func LoadWebUIConfig() *WebUI {
	filePtr, _ := os.Open("webui.json")
	defer filePtr.Close()

	var conf WebUI
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
