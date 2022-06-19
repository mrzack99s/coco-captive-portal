package firewall

import (
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/types"
)

func UnallowAccess(ss *types.SessionType) (err error) {
	err = ipt.Delete("nat", "PREROUTING", "-s", ss.IPAddress, "-p", "tcp", "-i", config.Config.SecureInterface, "-m", "tcp", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = ipt.Delete("filter", "FORWARD", "-s", ss.IPAddress, "-i", config.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	return
}
