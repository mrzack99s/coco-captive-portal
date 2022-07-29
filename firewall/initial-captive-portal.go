package firewall

import (
	"fmt"

	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

func InitializeCaptivePortal() (err error) {
	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()
	// Flush chain
	err = IPT.ClearAll()
	if err != nil {
		return
	}
	err = IPT.ClearChain("nat", "PREROUTING")
	if err != nil {
		return
	}

	// Append Rules
	// err = IPT.AppendUnique("filter", "INPUT", "-p", "icmp", "-j", "DROP")
	// if err != nil {
	// 	return
	// }

	if config.PROD_MODE {
		err = IPT.AppendUnique("filter", "INPUT", "-p", "tcp", "--dport", "22", "-d", interfaceIp, "-j", "DROP")
		if err != nil {
			return
		}
	}

	err = IPT.AppendUnique("filter", "INPUT", "-i", config.Config.EgressInterface, "-m", "state", "--state", "ESTABLISHED,RELATED", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.AppendUnique("filter", "FORWARD", "-s", "0.0.0.0/0", "-p", "udp", "-i", config.Config.SecureInterface, "--dport", "53", "-j", "ACCEPT")
	if err != nil {
		return
	}

	if config.Config.ExternalPortalURL != "" {
		_, host, port, _ := utils.ParseURL(config.Config.ExternalPortalURL)
		hostIp, _ := utils.ResolveIp(host)

		err = IPT.AppendUnique("filter", "FORWARD", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "--match", "multiport", "--dports", fmt.Sprintf("%s,8080,8443", port), "-d", hostIp, "-j", "ACCEPT")
		if err != nil {
			return
		}
	} else {
		err = IPT.AppendUnique("filter", "FORWARD", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "--match", "multiport", "--dports", "443,8080,8443", "-d", interfaceIp, "-j", "ACCEPT")
		if err != nil {
			return
		}
	}

	initFQDNBlocklist()
	initFinally()

	bypassNetworks()

	return
}

func initFQDNBlocklist() (err error) {
	for _, pattern := range config.Config.FQDNBlocklist {
		err = AddFQDNBlacklist(pattern)
		if err != nil {
			return
		}
	}
	return
}

func bypassNetworks() (err error) {
	for _, snet := range config.Config.BypassNetworks {
		err = AllowAccessBypass(&types.SessionType{
			IPAddress: snet,
		})
		if err != nil {
			return
		}
	}
	return
}

func initFinally() (err error) {
	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()
	err = IPT.AppendUnique("filter", "FORWARD", "-s", "0.0.0.0/0", "-i", config.Config.SecureInterface, "-j", "DROP")
	if err != nil {
		return
	}

	err = IPT.AppendUnique("nat", "PREROUTING", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "-d", interfaceIp, "-m", "tcp", "--dport", "443", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.AppendUnique("nat", "PREROUTING", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "-d", "1.1.1.1", "--dport", "80", "-j", "DNAT", "--to-destination", interfaceIp+":8080")
	if err != nil {
		return
	}

	err = IPT.AppendUnique("nat", "PREROUTING", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "-d", "1.1.1.1", "--dport", "443", "-j", "DNAT", "--to-destination", interfaceIp+":8443")
	if err != nil {
		return
	}

	for _, s := range config.Config.AllowEndpoints {
		AddAllowEndpoint(&s)
	}

	err = IPT.AppendUnique("nat", "PREROUTING", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "--dport", "80", "-j", "DNAT", "--to-destination", interfaceIp+":8080")
	if err != nil {
		return
	}

	err = IPT.AppendUnique("nat", "PREROUTING", "-s", "0.0.0.0/0", "-p", "tcp", "-i", config.Config.SecureInterface, "--dport", "443", "-j", "DNAT", "--to-destination", interfaceIp+":8443")
	if err != nil {
		return
	}

	return
}
