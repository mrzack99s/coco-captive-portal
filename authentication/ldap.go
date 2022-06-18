package authentication

import (
	"crypto/tls"
	"fmt"

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

	if l.Product != "google" {
		pass := false

		err = l.Connect()
		if err != nil {
			return
		}

		for _, bdn := range l.AllowBaseDN {
			err = l.instance.Bind(fmt.Sprintf("cn=%s,%s", username, bdn), password)
			if err == nil {
				pass = true
				break
			}
		}

		if !pass {
			err = fmt.Errorf("%s authentication failed", username)
			return
		}
	} else {

		err = l.Connect()
		if err != nil {
			return
		}
		err = l.instance.Bind(username, password)
		if err != nil {
			err = fmt.Errorf("%s authentication failed", username)
			return
		}

	}

	l.instance.Unbind()

	return
}
