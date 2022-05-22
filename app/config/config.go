package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Port    int
	LogZKey string
	*Mailgun
	*Db
	SecretKey string
}

type Db struct {
	Host     string
	DbName   string
	UserName string
	Port     string
	Password string
}

type Mailgun struct {
	PrivateAPIKey string
	Domain        string
}

func InitCfg() *Config {
	t := Config{}

	filename := "config.yaml"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &t
}

const settingsFilename = "config.yaml"

func Write(p Config) {
	rawDataOut, err := yaml.Marshal(&p)
	if err != nil {
		log.Println("YAML marshaling failed:", err)
	}

	err = ioutil.WriteFile(settingsFilename, rawDataOut, 0644)
	if err != nil {
		log.Println("Cannot write updated settings file:", err)
	}
}
