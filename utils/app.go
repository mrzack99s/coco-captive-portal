package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mrzack99s/coco-captive-portal/constants"
	"gopkg.in/yaml.v3"
)

type AppCredentialsType struct {
	APIToken string `yaml:"api_token"`
}

var (
	AppCredentials AppCredentialsType = AppCredentialsType{}
)

func SetupAppCredential() (err error) {
	_, err = os.Stat("./app_credentials.yaml")
	if os.IsNotExist(err) {
		err = GenerateApiToken()
		if err != nil {
			return
		}

	}
	return
}

func GenerateApiToken() (err error) {
	appCredential := AppCredentialsType{
		APIToken: SecretGenerator(64),
	}

	err = CacheSet(constants.CONFIG, "api-token", appCredential.APIToken)
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
