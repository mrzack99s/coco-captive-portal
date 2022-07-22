package types

import "github.com/mrzack99s/coco-captive-portal/authentication"

type ConfigType struct {
	ExternalPortalURL    string                             `yaml:"external_portal_url" json:"external_portal_url"`
	EgressInterface      string                             `yaml:"egress_interface" json:"egress_interface"`
	SecureInterface      string                             `yaml:"secure_interface" json:"secure_interface"`
	SessionIdle          uint64                             `yaml:"session_idle" json:"session_idle"`
	MaxConcurrentSession uint64                             `yaml:"max_concurrent_session" json:"max_concurrent_session"`
	BypassNetworks       []string                           `yaml:"bypass_networks" json:"bypass_networks"`
	AllowEndpoints       []EndpointType                     `yaml:"allow_endpoints" json:"allow_endpoints"`
	RedirectURL          string                             `yaml:"redirect_url" json:"redirect_url"`
	Radius               *authentication.RadiusEndpointType `yaml:"radius" json:"radius"`
	LDAP                 *authentication.LDAPEndpointType   `yaml:"ldap" json:"ldap"`
	HTML                 HTMLType                           `yaml:"html" json:"html"`
	Administrator        CredentialType                     `yaml:"administrator" json:"administrator"`
	DomainNames          struct {
		OperatorDomainName string `yaml:"operator_domain_name" json:"operator_domain_name"`
		AuthDomainName     string `yaml:"auth_domain_name" json:"auth_domain_name"`
	} `yaml:"domain_names" json:"domain_names"`
	DDOSPrevention bool     `yaml:"ddos_prevention" json:"ddos_prevention"`
	FQDNBlocklist  []string `yaml:"fqdn_blocklist" json:"fqdn_blocklist"`
}

type EndpointType struct {
	Hostname string `yaml:"hostname" json:"hostname"`
	Port     uint64 `yaml:"port" json:"port"`
}

type CredentialType struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type HTMLType struct {
	DefaultLanguage    string `yaml:"default_language" json:"default_language"`
	EnTitleName        string `yaml:"en_title_name" json:"en_title_name"`
	EnSubTitle         string `yaml:"en_sub_title" json:"en_sub_title"`
	ThTitleName        string `yaml:"th_title_name" json:"th_title_name"`
	ThSubTitle         string `yaml:"th_sub_title" json:"th_sub_title"`
	LogoFileName       string `yaml:"logo_file_name" json:"logo_file_name"`
	BackgroundFileName string `yaml:"background_file_name" json:"background_file_name"`
}

type ExtendConfigType struct {
	ConfigType
	Status struct {
		EgressIPAddress string `yaml:"egress_ip_address" json:"egress_ip_address"`
		SecureIPAddress string `yaml:"secure_ip_address" json:"secure_ip_address"`
	} `yaml:"status" json:"status"`
}
