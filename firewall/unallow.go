package firewall

import (
	"fmt"

	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

func UnallowAccess(ss *types.SessionType) (err error) {
	err = IPT.Delete("nat", "PREROUTING", "-s", ss.IPAddress, "-p", "tcp", "-i", config.Config.SecureInterface, "-m", "tcp", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Delete("filter", "FORWARD", "-s", ss.IPAddress, "-i", config.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	return
}

func UnallowAccessBypass(ss *types.SessionType) (err error) {
	err = IPT.Delete("nat", "PREROUTING", "-s", ss.IPAddress, "-p", "tcp", "-i", config.Config.SecureInterface, "-m", "tcp", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Delete("filter", "FORWARD", "-s", ss.IPAddress, "-i", config.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()
	err = IPT.Delete("filter", "INPUT", "-s", ss.IPAddress, "-i", config.Config.SecureInterface, "-p", "tcp", "--match", "multiport", "--dports", "443,8080,8443", "-d", interfaceIp, "-j", "DROP")
	if err != nil {
		return
	}

	return
}

func DelAllowEndpoint(ss *types.EndpointType) {
	allIp, _ := utils.ResolveAllIp(ss.Hostname)
	for _, hostIp := range allIp {
		err := IPT.Delete("nat", "PREROUTING", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "-m", "tcp", "-d", hostIp, "--dport", fmt.Sprintf("%d", ss.Port), "-j", "ACCEPT")
		if err == nil {
			last_allow_endpoint_rule_num--
		}
	}
}
