package firewall

import "github.com/mrzack99s/coco-captive-portal/config"

func AddFQDNBlacklist(pattern string) (err error) {
	err = IPT.Insert("filter", "FORWARD", last_fqdn_blacklist_rule_num, "-i", config.Config.SecureInterface,
		"-m", "string", "--string", pattern, "--algo", "kmp", "-j", "DROP")
	if err != nil {
		return
	}
	last_fqdn_blacklist_rule_num++

	return
}

func DelFQDNBlacklist(pattern string) (err error) {
	err = IPT.Delete("filter", "FORWARD", "-i", config.Config.SecureInterface,
		"-m", "string", "--string", pattern, "--algo", "kmp", "-j", "DROP")
	if err != nil {
		return
	}
	last_fqdn_blacklist_rule_num--

	return
}
