package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func ParseConfig() {
	yfile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yfile, &Config)
	if err != nil {
		panic(err)
	}

	if Config.HTML.DefaultLanguage == "" {
		Config.HTML.DefaultLanguage = "en"
	}

}
