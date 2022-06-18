package authentication

import "github.com/go-ldap/ldap/v3"

type RadiusEndpointType struct {
	Hostname string `yaml:"hostname"`
	Port     uint   `yaml:"port"`
	Secret   string `yaml:"secret"`
}

type LDAPEndpointType struct {
	Product      string   `yaml:"product"`
	Hostname     string   `yaml:"hostname"`
	BindUsername string   `yaml:"bind_username"`
	BindPassword string   `yaml:"bind_password"`
	Port         uint     `yaml:"port"`
	TLSEnable    bool     `yaml:"tls_enable"`
	AllowBaseDN  []string `yaml:"allow_base_dn"`
	instance     *ldap.Conn
}
