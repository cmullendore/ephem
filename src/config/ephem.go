package config

import (
	"encoding/json"
	"io/ioutil"
)

type Ephem struct {
	Listener         string
	ListenerCert     Cert
	CertCAPath       string
	CleanupFrequency string
	Aes256Key        string
	MaxReads         int
	MaxAgeSeconds    int
	Aes256KeySource  string
	Database         Database
}

type Database struct {
	Driver   string
	Server   string
	Name     string
	Username string
	Password string
}

func LoadEphemConfig() *Ephem {
	var conf Ephem
	b, err := ioutil.ReadFile("ephem.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, &conf)
	if err != nil {
		panic(err)
	}

	return &conf
}
