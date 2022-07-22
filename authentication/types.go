package authentication

import "github.com/go-ldap/ldap/v3"

type RadiusEndpointType struct {
	Hostname string `yaml:"hostname" json:"hostname"`
	Port     uint64 `yaml:"port" json:"port"`
	Secret   string `yaml:"secret" json:"secret"`
}

type LDAPEndpointType struct {
	Hostname     string   `yaml:"hostname" json:"hostname"`
	Port         uint64   `yaml:"port" json:"port"`
	TLSEnable    bool     `yaml:"tls_enable" json:"tls_enable"`
	SingleDomain bool     `yaml:"single_domain" json:"single_domain"`
	DomainNames  []string `yaml:"domain_names" json:"domain_names"`
	instance     *ldap.Conn
}
