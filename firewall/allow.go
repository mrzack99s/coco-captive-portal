package firewall

import (
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/types"
)

func AllowAccess(ss *types.SessionType) (err error) {
	err = IPT.Insert("nat", "PREROUTING", last_allow_endpoint_rule_num+1, "-s", ss.IPAddress, "-p", "tcp", "-i", config.Config.SecureInterface, "-m", "tcp", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Insert("filter", "FORWARD", 1, "-s", ss.IPAddress, "-i", config.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	return
}
