package firewall

import (
	"fmt"

	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

func AllowAccess(ss *types.SessionType) (err error) {
	err = IPT.Insert("nat", "PREROUTING", last_allow_endpoint_rule_num+1, "-s", ss.IPAddress, "-p", "tcp", "-i", config.Config.SecureInterface, "-m", "tcp", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Insert("filter", "FORWARD", last_fqdn_blacklist_rule_num, "-s", ss.IPAddress, "-i", config.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	return
}

func AllowAccessBypass(ss *types.SessionType) (err error) {

	err = IPT.Insert("nat", "PREROUTING", last_allow_endpoint_rule_num+1, "-s", ss.IPAddress, "-p", "tcp", "-i", config.Config.SecureInterface, "-m", "tcp", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Insert("filter", "FORWARD", last_fqdn_blacklist_rule_num, "-s", ss.IPAddress, "-i", config.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()
	err = IPT.Insert("filter", "INPUT", 3, "-s", ss.IPAddress, "-i", config.Config.SecureInterface, "-p", "tcp", "--match", "multiport", "--dports", "443,8080,8443", "-d", interfaceIp, "-j", "DROP")
	if err != nil {
		return
	}

	return
}

func AddAllowEndpoint(ss *types.EndpointType) (err error) {
	allIp, _ := utils.ResolveAllIp(ss.Hostname)
	for _, hostIp := range allIp {
		err = IPT.Insert("filter", "FORWARD", last_allow_endpoint_rule_num, "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "-m", "tcp", "-d", hostIp, "--dport", fmt.Sprintf("%d", ss.Port), "-j", "ACCEPT")
		if err != nil {
			return
		}

		err = IPT.Insert("nat", "PREROUTING", last_allow_endpoint_rule_num+1, "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "-m", "tcp", "-d", hostIp, "--dport", fmt.Sprintf("%d", ss.Port), "-j", "ACCEPT")
		if err != nil {
			return
		}
		last_allow_endpoint_rule_num++
	}
	return
}
