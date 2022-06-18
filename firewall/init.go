package firewall

import "github.com/coreos/go-iptables/iptables"

func init() {
	var err error
	ipt, err = iptables.New()
	if err != nil {
		return
	}
}
