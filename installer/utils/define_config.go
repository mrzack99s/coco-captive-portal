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
	config.EgressInterface = scan("Egress interface name (WAN, ex. eth0): ")
	config.SecureInterface = scan("Secure interface name (LAN, ex. eth1): ")
	config.SessionIdle = scanUint64("Session idle (Minutes): ")
	config.MaxConcurrentSession = scanUint64("Max concurrent session (Devices): ")
	authCertFileName := scan("Captive portal certificate file location (ex. /home/coco/cert.pem): ")
	authKeyFileName := scan("Captive portal key of certificate file location (ex. /home/coco/key.pem): ")
	operatorCertFileName := scan("Operator certificate file location (ex. /home/coco/cert.pem): ")
	operatorKeyFileName := scan("Operator key of certificate file location (ex. /home/coco/key.pem): ")

	dstAuthCertFileName := constants.APP_DIR + "/certs/authfullchain.pem"
	dstAuthKeyFileName := constants.APP_DIR + "/certs/authprivkey.pem"
	dstOperatorCertFileName := constants.APP_DIR + "/certs/operatorfullchain.pem"
	dstOperatorKeyFileName := constants.APP_DIR + "/certs/operatorprivkey.pem"

	log.Info().Msg("coping certificate file")
	if e := copy(CopyType{
		Src:  authCertFileName,
		Dst:  dstAuthCertFileName,
		Perm: 0644,
	}); e != nil {
		if IGNORE_VERIFY {
			log.Warn().Msgf("copy %s to %s failed", authCertFileName, dstAuthCertFileName)
		} else {
			log.Warn().Msgf("copy %s to %s failed", authCertFileName, dstAuthCertFileName)
			err = e
			return
		}
	}

	if e := copy(CopyType{
		Src:  authKeyFileName,
		Dst:  dstAuthKeyFileName,
		Perm: 0644,
	}); e != nil {
		if IGNORE_VERIFY {
			log.Warn().Msgf("copy %s to %s failed", authKeyFileName, dstAuthKeyFileName)
		} else {
			log.Warn().Msgf("copy %s to %s failed", authKeyFileName, dstAuthKeyFileName)
			err = e
			return
		}
	}

	if e := copy(CopyType{
		Src:  operatorCertFileName,
		Dst:  dstOperatorCertFileName,
		Perm: 0644,
	}); e != nil {
		if IGNORE_VERIFY {
			log.Warn().Msgf("copy %s to %s failed", operatorCertFileName, dstOperatorCertFileName)
		} else {
			log.Warn().Msgf("copy %s to %s failed", operatorCertFileName, dstOperatorCertFileName)
			err = e
			return
		}
	}

	if e := copy(CopyType{
		Src:  operatorKeyFileName,
		Dst:  dstOperatorKeyFileName,
		Perm: 0644,
	}); e != nil {
		if IGNORE_VERIFY {
			log.Warn().Msgf("copy %s to %s failed", operatorKeyFileName, dstOperatorKeyFileName)
		} else {
			log.Warn().Msgf("copy %s to %s failed", operatorKeyFileName, dstOperatorKeyFileName)
			err = e
			return
		}
	}

	if confirmWithMsg("Are you need to custom domain of captive portal and operator?") {
		config.DomainNames.AuthDomainName = scan("Captive portal domain name (ex. coco-cative-portal.local): ")
		config.DomainNames.OperatorDomainName = scan("Operator domain name (ex. coco-cative-portal.local): ")
	}

	mode := ""
	for !(mode == "ldap" || mode == "radius") {
		mode = scan("Mode (ldap|radius): ")
	}

	if mode == "ldap" {
		config.LDAP = &authentication.LDAPEndpointType{}
		config.LDAP.Hostname = scan("[LDAP] Hostname (ex. coco-cative-portal.local, 172.16.0.254): ")
		config.LDAP.Port = scanUint64("[LDAP] Port: ")
		config.LDAP.TLSEnable = confirmWithMsg("Enable TLS?")
		if config.LDAP.TLSEnable {
			certFileName := scan("Certificate file location: ")
			keyFileName := scan("Key of certificate file location: ")

			dstCertFileName := constants.APP_DIR + "/certs/ldapchain.pem"
			dstKeyFileName := constants.APP_DIR + "/certs/ldapprivkey.pem"

			log.Info().Msg("coping certificate file")
			if e := copy(CopyType{
				Src:  certFileName,
				Dst:  dstCertFileName,
				Perm: 0644,
			}); e != nil {
				if IGNORE_VERIFY {
					log.Warn().Msgf("copy %s to %s failed", certFileName, dstCertFileName)
				} else {
					log.Warn().Msgf("copy %s to %s failed", certFileName, dstCertFileName)
					err = e
					return
				}
			}

			if e := copy(CopyType{
				Src:  keyFileName,
				Dst:  dstKeyFileName,
				Perm: 0644,
			}); e != nil {
				if IGNORE_VERIFY {
					log.Warn().Msgf("copy %s to %s failed", keyFileName, dstKeyFileName)
				} else {
					log.Warn().Msgf("copy %s to %s failed", keyFileName, dstKeyFileName)
					err = e
					return
				}
			}

		}

		config.LDAP.SingleDomain = confirmWithMsg("Single domain?")
		if !config.LDAP.SingleDomain {
			config.LDAP.DomainNames = scanArray("Enter a domain name (ex. coco-cative-portal.local): ")
		} else {
			domain := scan("Enter a domain name (ex. coco-cative-portal.local): ")
			config.LDAP.DomainNames = append(config.LDAP.DomainNames, domain)
		}
	} else if mode == "radius" {
		config.Radius = &authentication.RadiusEndpointType{}
		config.Radius.Hostname = scan("[RADIUS] Hostname (ex. coco-cative-portal.local, 172.16.0.254): ")
		config.Radius.Port = scanUint64("[RADIUS] Port: ")
		config.Radius.Secret = scan("[RADIUS] Secret: ")
	}

	if !confirmWithMsg("Authorized access from any network?") {
		config.AuthorizedNetworks = scanArray("Enter a network cidr (ex. 10.0.0.0/24): ")
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
	config.Administrator.Credential = types.CredentialType{}
	fmt.Print("### Administrator Config\n\n")
	config.Administrator.Credential.Username = scan("Administrator username: ")
	fmt.Print("Administrator password: ")
	password, _ := terminal.ReadPassword(0)
	config.Administrator.Credential.Password = utils.Sha512encode(string(password))

	if !confirmWithMsg("Authorized access an operator from any network?") {
		config.Administrator.AuthorizedNetworks = scanArray("Enter a network cidr (ex. 10.0.0.0/24): ")
	}

	clearTerminal()
	fmt.Print("### System Config\n\n")
	config.DDOSPrevention = confirmWithMsg("Enable DDOS prevention?")

	config.AuthorizedNetworks = append(config.AuthorizedNetworks, "0.0.0.0/0")
	config.Administrator.AuthorizedNetworks = append(config.Administrator.AuthorizedNetworks, "0.0.0.0/0")

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
	fmt.Print("### Enter 'end' to exit adding a message\n")
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
