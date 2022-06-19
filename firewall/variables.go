package firewall

import "github.com/coreos/go-iptables/iptables"

var (
	ipt                          *iptables.IPTables
	last_allow_endpoint_rule_num int = 3
)
