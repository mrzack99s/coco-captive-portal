package installer_utils

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/mrzack99s/coco-captive-portal/authentication"
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

func defineConfig() (err error) {

	clearTerminal()

	config := types.ConfigType{}
	fmt.Print("### Common Config\n\n")
	config.EgressInterface = scan("Egress interface name (WAN): ")
	config.SecureInterface = scan("Secure interface name (LAN): ")
	config.SessionIdle = scanUint64("Session idle (Minutes): ")
	config.MaxConcurrentSession = scanUint64("Max concurrent session (Devices): ")
	config.RedirectURL = scan("After authorized then redirect to url: ")

	mode := ""
	for !(mode == "ldap" || mode == "radius") {
		mode = scan("Mode (ldap|radius): ")
	}

	if mode == "ldap" {
		config.LDAP = &authentication.LDAPEndpointType{}
		config.LDAP.Hostname = scan("[LDAP] Hostname: ")
		config.LDAP.Port = scanUint64("[LDAP] Port: ")
		config.LDAP.TLSEnable = confirmWithMsg("Enable TLS?")
		config.LDAP.SingleDomain = confirmWithMsg("Single domain?")
		fmt.Println(config.LDAP.SingleDomain)
		if !config.LDAP.SingleDomain {
			config.LDAP.DomainNames = scanArray("Enter a domain name: ")
		}
	} else if mode == "radius" {
		config.Radius = &authentication.RadiusEndpointType{}
		config.Radius.Hostname = scan("[RADIUS] Hostname: ")
		config.Radius.Port = scanUint64("[RADIUS] Port: ")
		config.Radius.Secret = scan("[RADIUS] Secret: ")
	}

	clearTerminal()
	config.HTML = types.HTMLType{}
	fmt.Print("### HTML Config\n\n")
	config.HTML.EnTitleName = scan("Title name in english laguage: ")
	config.HTML.EnSubTitle = scan("Sub title name in english laguage: ")
	config.HTML.ThTitleName = scan("Title name in thai laguage: ")
	config.HTML.ThSubTitle = scan("Sub title name in thai laguage: ")
	htmlDefaultLang := ""
	for !(htmlDefaultLang == "en" || htmlDefaultLang == "th") {
		htmlDefaultLang = scan("Default language (en|th): ")
	}
	config.HTML.DefaultLanguage = htmlDefaultLang

	if confirmWithMsg("Are you need to change logo file name?") {
		config.HTML.EnTitleName = scan("Logo file name: ")
	}

	if confirmWithMsg("Are you need to change background file name?") {
		config.HTML.EnTitleName = scan("Background file name: ")
	}

	clearTerminal()
	config.Administrator = types.CredentialType{}
	fmt.Print("### Administrator Config\n\n")
	config.Administrator.Username = scan("Administrator username: ")
	fmt.Print("Administrator password: ")
	password, _ := terminal.ReadPassword(0)
	config.Administrator.Password = utils.Sha512encode(string(password))

	clearTerminal()
	fmt.Print("### System Config\n\n")
	config.DDOSPrevention = confirmWithMsg("Enable DDOS prevention?")

	dByte, err := yaml.Marshal(config)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(constants.APP_DIR+"/config.yaml", dByte, 0644)

	clearTerminal()
	log.Info().Msg("config generated")
	return

}

func clearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func scan(words string) string {
	dataInput := ""
	i := 0
	for dataInput == "" {
		if i > 0 {
			log.Warn().Msg("please enter a value")
		}
		fmt.Print(words)
		reader := bufio.NewReader(os.Stdin)
		dataInput, _ = reader.ReadString('\n')
		dataInput = strings.TrimSpace(dataInput)
		i++
	}
	return dataInput
}

func scanUint64(words string) uint64 {
	dataInput := ""
	var resData uint64
	var err error = errors.New("")
	i := 0
	for dataInput == "" || err != nil {
		if i > 0 {
			log.Warn().Msg("please enter a number value")
		}
		fmt.Print(words)
		reader := bufio.NewReader(os.Stdin)
		dataInput, _ = reader.ReadString('\n')
		dataInput = strings.TrimSpace(dataInput)
		i++
		resData, err = utils.StringToUInt64(dataInput)
	}
	return resData
}

func scanArray(words string) []string {
	dataInput := ""
	dataInputArray := []string{}
	for dataInput != "end" || dataInput == "" {
		dataInput = scan(words)
		if dataInput != "end" {
			dataInputArray = append(dataInputArray, dataInput)
		}
	}
	return dataInputArray
}

func confirmWithMsg(words string) bool {
	dataInput := ""
	for !(dataInput == "Y" || dataInput == "N" || dataInput == "y" || dataInput == "n") {
		fmt.Printf("%s (Y|N): ", words)
		reader := bufio.NewReader(os.Stdin)
		dataInput, _ = reader.ReadString('\n')
		dataInput = strings.TrimSpace(dataInput)
	}
	if dataInput == "Y" || dataInput == "y" {
		return true
	} else {
		return false
	}
}
