package authentication

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func (l *LDAPEndpointType) Connect() (err error) {
	if l.TLSEnable {
		cer, e := tls.LoadX509KeyPair("./certs/ldapchain.pem", "./certs/ldapprivkey.pem")
		if e != nil {
			err = e
			return
		}
		config := &tls.Config{Certificates: []tls.Certificate{cer}}

		l.instance, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", l.Hostname, l.Port), config)
		if err != nil {
			return
		}

	} else {
		l.instance, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", l.Hostname, l.Port))
		if err != nil {
			return
		}
		defer l.instance.Close()
	}

	return
}

func (l *LDAPEndpointType) Authentication(username, password string) (err error) {

	l.Connect()
	pass := false

	if l.SingleDomain {
		splitString := strings.Split(username, "@")
		if splitString[1] == l.DomainNames[0] {
			pass = true
		} else {
			err = fmt.Errorf("username %s is not authorized by the domain name", username)
			return
		}

	} else {
		splitString := strings.Split(username, "@")
		for _, domain := range l.DomainNames {
			if splitString[1] == domain {
				pass = true
				break
			}
		}
		if !pass {
			err = fmt.Errorf("username %s is not authorized by the domain name", username)
			return
		}
	}

	err = l.instance.Bind(username, password)
	if err != nil {
		err = fmt.Errorf("%s credentials are invalid", username)
		return
	}

	l.instance.Unbind()

	return
}
