package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/mrzack99s/coco-captive-portal/constants"
	"gopkg.in/yaml.v3"
)

type AppCredentialsType struct {
	APIToken string `yaml:"api_token"`
}

var (
	AppCredentials AppCredentialsType = AppCredentialsType{}
)

func GenerateApiToken() (err error) {
	appCredential := AppCredentialsType{
		APIToken: SecretGenerator(64),
	}

	err = CacheSet(constants.SCHEMA_CONFIG, "api-token", appCredential.APIToken)
	if err != nil {
		return
	}

	fmt.Println("api-token: " + appCredential.APIToken)

	bData, e := yaml.Marshal(appCredential)
	if e != nil {
		err = e
		return
	}

	err = ioutil.WriteFile("./app_credentials.yaml", bData, 0444)
	if err != nil {
		return
	}
	return
}
