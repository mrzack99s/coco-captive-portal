package types

import "github.com/mrzack99s/coco-captive-portal/authentication"

type ConfigType struct {
	ExternalPortalURL    string                             `yaml:"external_portal_url"`
	EgressInterface      string                             `yaml:"egress_interface"`
	SecureInterface      string                             `yaml:"secure_interface"`
	SessionIdle          uint                               `yaml:"session_idle"`
	MaxConcurrentSession uint                               `yaml:"max_concurrent_session"`
	BypassNetworks       []string                           `yaml:"bypass_networks"`
	AllowEndpoints       []EndpointType                     `yaml:"allow_endpoints"`
	RedirectURL          string                             `yaml:"redirect_url"`
	Radius               *authentication.RadiusEndpointType `yaml:"radius"`
	LDAP                 *authentication.LDAPEndpointType   `yaml:"ldap"`
	HTML                 HTMLType                           `yaml:"html"`
	Administrator        CredentialType                     `yaml:"administrator"`
}

type EndpointType struct {
	Hostname string `yaml:"hostname"`
	Port     uint   `yaml:"port"`
}

type CredentialType struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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
