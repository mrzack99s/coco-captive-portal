package firewall

import (
	"fmt"

	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

func InitializeCaptivePortal() (err error) {

	// Flush chain
	err = ipt.ClearAll()
	if err != nil {
		return
	}
	err = ipt.ClearChain("nat", "PREROUTING")
	if err != nil {
		return
	}

	// Append Rules
	err = ipt.AppendUnique("filter", "INPUT", "-i", config.Config.EgressInterface, "-m", "state", "--state", "ESTABLISHED,RELATED", "-j", "ACCEPT")
	if err != nil {
		return
	}

	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()

	err = ipt.AppendUnique("filter", "FORWARD", "-s", "0.0.0.0/0", "-p", "udp", "-i", config.Config.SecureInterface, "--dport", "53", "-j", "ACCEPT")
	if err != nil {
		return
	}

	if config.Config.ExternalPortalURL != "" {
		_, host, port, _ := utils.ParseURL(config.Config.ExternalPortalURL)
		hostIp, _ := utils.ResolveIp(host)

		err = ipt.AppendUnique("filter", "FORWARD", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "--dport", port, "-d", hostIp, "-j", "ACCEPT")
		if err != nil {
			return
		}
	} else {
		err = ipt.AppendUnique("filter", "FORWARD", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "--dport", "443", "-d", interfaceIp, "-j", "ACCEPT")
		if err != nil {
			return
		}
	}

	err = ipt.AppendUnique("filter", "FORWARD", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "--dport", "8080", "-d", interfaceIp, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = ipt.AppendUnique("filter", "FORWARD", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "--dport", "8443", "-d", interfaceIp, "-j", "ACCEPT")
	if err != nil {
		return
	}

	allowEndpoints()
	bypassNetworks()
	initFinally()

	return
}

func bypassNetworks() (err error) {
	for _, snet := range config.Config.BypassNetworks {
		err = AllowAccess(&types.SessionType{
			IPAddress: snet,
		})
		if err != nil {
			return
		}
	}
	return
}

func allowEndpoints() (err error) {
	for _, s := range config.Config.AllowEndpoints {
		allIp, _ := utils.ResolveAllIp(s.Hostname)
		for _, hostIp := range allIp {
			err = ipt.AppendUnique("filter", "FORWARD", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "--dport", fmt.Sprintf("%d", s.Port), "-d", hostIp, "-j", "ACCEPT")
			if err != nil {
				return
			}
		}

	}
	return
}

func initFinally() (err error) {
	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()
	err = ipt.AppendUnique("filter", "FORWARD", "-s", "0.0.0.0/0", "-i", config.Config.SecureInterface, "-j", "DROP")
	if err != nil {
		return
	}

	err = ipt.AppendUnique("nat", "PREROUTING", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "-d", interfaceIp, "-m", "tcp", "--dport", "443", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = ipt.AppendUnique("nat", "PREROUTING", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "-d", "1.1.1.1", "--dport", "80", "-j", "DNAT", "--to-destination", interfaceIp+":8080")
	if err != nil {
		return
	}

	err = ipt.AppendUnique("nat", "PREROUTING", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "-d", "1.1.1.1", "--dport", "443", "-j", "DNAT", "--to-destination", interfaceIp+":8443")
	if err != nil {
		return
	}

	for _, s := range config.Config.AllowEndpoints {
		allIp, _ := utils.ResolveAllIp(s.Hostname)
		for _, hostIp := range allIp {
			err = ipt.AppendUnique("nat", "PREROUTING", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "-m", "tcp", "-d", hostIp, "--dport", fmt.Sprintf("%d", s.Port), "-j", "ACCEPT")
			if err != nil {
				return
			}
			last_allow_endpoint_rule_num++
		}

	}

	err = ipt.AppendUnique("nat", "PREROUTING", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "--dport", "80", "-j", "DNAT", "--to-destination", interfaceIp+":8080")
	if err != nil {
		return
	}

	err = ipt.AppendUnique("nat", "PREROUTING", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "--dport", "443", "-j", "DNAT", "--to-destination", interfaceIp+":8443")
	if err != nil {
		return
	}

	return
}
