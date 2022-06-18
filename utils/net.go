package utils

import (
	"fmt"
	"net"
	"regexp"

	"github.com/mrzack99s/coco-captive-portal/config"
)

func GetSecureInterfaceIpv4Addr() (addr string, err error) {
	var (
		ief      *net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
	)

	if ief, err = net.InterfaceByName(config.Config.SecureInterface); err != nil { // get interface
		return
	}
	if addrs, err = ief.Addrs(); err != nil { // get addresses
		return
	}
	for _, addr := range addrs { // get ipv4 address
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil {
		return "", fmt.Errorf(fmt.Sprintf("interface %s don't have an ipv4 address\n", config.Config.SecureInterface))
	}
	return ipv4Addr.String(), nil
}

func ResolveIp(hostname string) (address string, err error) {
	addr, err := net.LookupIP(hostname)
	if err != nil {
		return
	} else {
		address = addr[0].String()
	}
	return
}

func ResolveAllIp(hostname string) (address []string, err error) {
	addr, err := net.LookupIP(hostname)
	if err != nil {
		return
	} else {
		for _, a := range addr {
			if IsIpv4(a.String()) {
				address = append(address, a.String())
			}

		}
	}
	return
}

func IsIpv4(ip string) bool {
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	return re.MatchString(ip)
}
