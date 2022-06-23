package config

import (
	"io/ioutil"

	"github.com/mrzack99s/coco-captive-portal/constants"
	"gopkg.in/yaml.v3"
)

func UpdateConfig() (err error) {

	dByte, err := yaml.Marshal(Config)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(constants.APP_DIR+"/config.yaml", dByte, 0644)
	return
}
