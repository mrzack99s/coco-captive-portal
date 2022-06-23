package firewall

import "github.com/coreos/go-iptables/iptables"

var (
	IPT                          *iptables.IPTables
	last_allow_endpoint_rule_num int = 3
)
