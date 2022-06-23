package config

import (
	"io/ioutil"

	"github.com/mrzack99s/coco-captive-portal/constants"
	"gopkg.in/yaml.v3"
)

func ParseConfig() {
	yfile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	if PROD_MODE {
		yfile, err = ioutil.ReadFile(constants.APP_DIR + "/config.yaml")
		if err != nil {
			panic(err)
		}
	}

	err = yaml.Unmarshal(yfile, &Config)
	if err != nil {
		panic(err)
	}

	if Config.HTML.DefaultLanguage == "" {
		Config.HTML.DefaultLanguage = "en"
	}

}
