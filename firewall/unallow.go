package firewall

import (
	"github.com/mrzack99s/coco-captive-portal/types"
)

func UnallowAccess(ss *types.SessionType) (err error) {
	err = ipt.Delete("nat", "PREROUTING", "-s", ss.IPAddress, "-p", "tcp", "-m", "tcp", "--dport", "80", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = ipt.Delete("nat", "PREROUTING", "-s", ss.IPAddress, "-p", "tcp", "-m", "tcp", "--dport", "443", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = ipt.Delete("filter", "FORWARD", "-s", ss.IPAddress, "-j", "ACCEPT")
	if err != nil {
		return
	}

	return
}
