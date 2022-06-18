package firewall

import (
	"github.com/mrzack99s/coco-captive-portal/types"
)

func AllowAccess(ss *types.SessionType) (err error) {
	err = ipt.Insert("nat", "PREROUTING", 1, "-s", ss.IPAddress, "-p", "tcp", "-m", "tcp", "--dport", "80", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = ipt.Insert("nat", "PREROUTING", 1, "-s", ss.IPAddress, "-p", "tcp", "-m", "tcp", "--dport", "443", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = ipt.Insert("filter", "FORWARD", 1, "-s", ss.IPAddress, "-j", "ACCEPT")
	if err != nil {
		return
	}

	return
}
