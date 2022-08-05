package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/firewall"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

var (
	ddosPreventionContext       context.Context
	ddosPreventionCalcelContext context.CancelFunc
	RATE_TCP_SYN                = 2000
	RATE_UDP_FLOOD              = 2000
	RATE_ICMP_FLOOD             = 250
)

func StartDDOSPreventor() {
	ddosPreventionContext, ddosPreventionCalcelContext = context.WithCancel(context.Background())
	ddosPreventor(ddosPreventionContext)
}

func StopDDOSPreventor() {
	ddosPreventionCalcelContext()
}

func RestartDDOSPreventor() {
	StopDDOSPreventor()
	StartDDOSPreventor()
}

func ddosPreventor(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return
			default:
				interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()
				ddosHandler(config.Config.SecureInterface, interfaceIp)
			}
		}
	}(ctx)
	go func(ctx context.Context) {
		for {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return
			default:
				interfaceIp, _ := utils.GetEgressInterfaceIpv4Addr()
				ddosHandler(config.Config.EgressInterface, interfaceIp)
			}

		}
	}(ctx)
}

func ddosHandler(interfaceName, intIpv4 string) {
	RATE_LIMIT := 2000 // Per second
	ignore_port := []string{}

	if !config.PROD_MODE {
		ignore_port = append(ignore_port, "22")
		ignore_port = append(ignore_port, "6379")
	}

	if handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter(fmt.Sprintf("not src %s and dst %s", intIpv4, intIpv4)); err != nil { // optional
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		for packet := range packetSource.Packets() {

			srcip := ""
			proto := ""
			op := constants.NET_OP_FLOOD
			dport := ""
			dPortOnlyNum := ""
			icmpType := layers.ICMPv4TypeCode(0)

			if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
				ip, _ := ipLayer.(*layers.IPv4)
				srcip = ip.SrcIP.String()
			}

			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				proto = constants.NET_PROTO_TCP
				tcp, _ := tcpLayer.(*layers.TCP)
				dport = tcp.DstPort.String()
				regex := regexp.MustCompile("[a-zA-Z()]+")
				dportSplit := regex.Split(dport, -1)
				dPortOnlyNum = dportSplit[0]
				dPortOnlyNum = strings.TrimSpace(dPortOnlyNum)
				if tcp.SYN {
					op = constants.NET_OP_SYN
				}
				if tcp.ACK {
					op = constants.NET_OP_ACK
				}
			}

			if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
				proto = constants.NET_PROTO_UDP
				udp, _ := udpLayer.(*layers.UDP)
				dport = udp.DstPort.String()
				regex := regexp.MustCompile("[a-zA-Z()]+")
				dportSplit := regex.Split(dport, -1)
				dPortOnlyNum = dportSplit[0]
				dPortOnlyNum = strings.TrimSpace(dPortOnlyNum)
			}

			if icmpLayer := packet.Layer(layers.LayerTypeICMPv4); icmpLayer != nil {
				proto = constants.NET_PROTO_ICMP
				icmp, _ := icmpLayer.(*layers.ICMPv4)
				icmpType = icmp.TypeCode
			}

			if !utils.ExistingInArray(ignore_port, dPortOnlyNum) {

				key := fmt.Sprintf("%s-ip%s-p%s-%s", proto, srcip, dport, op)
				switch proto {
				case constants.NET_PROTO_TCP:
					if op == constants.NET_OP_SYN {
						RATE_LIMIT = RATE_TCP_SYN
					}
				case constants.NET_PROTO_UDP:
					RATE_LIMIT = RATE_UDP_FLOOD
				case constants.NET_PROTO_ICMP:
					RATE_LIMIT = RATE_ICMP_FLOOD
				}

				if proto == constants.NET_PROTO_ICMP {
					key = fmt.Sprintf("%s-ip%s-c%s-%s", proto, srcip, icmpType, op)
				}

				if utils.CacheFindExistingKey(constants.SCHEMA_DDOS, key) {
					xtimesStr, err := utils.CacheGetString(constants.SCHEMA_DDOS, key)
					if err != nil {
						xtimesStr = "0"
					}

					if xtimesStr != "blocked" {
						xtimes, err := utils.StringToInt64(xtimesStr)
						if err != nil {
							xtimes = 0
						}
						if xtimes < int64(RATE_LIMIT) {

							if proto == constants.NET_PROTO_ICMP {
								firewall.IPT.Delete("filter", "INPUT", "-i", interfaceName, "-p", proto, "-s", srcip, "-j", "DROP")
							} else {
								firewall.IPT.Delete("filter", "INPUT", "-i", interfaceName, "-p", proto, "-s", srcip, "--dport", dPortOnlyNum, "-j", "DROP")
							}

							xtimes++
							xtimesStr = utils.Int64ToString(xtimes)

							utils.CacheSetWithTimeDuration(constants.SCHEMA_DDOS, key, xtimesStr, (time.Second*1)+(time.Millisecond*100))

						} else {
							utils.CacheSetWithTimeDuration(constants.SCHEMA_DDOS, key, "blocked", time.Minute*1)
							if proto == constants.NET_PROTO_ICMP {
								exist, _ := firewall.IPT.Exists("filter", "INPUT", "-i", interfaceName, "-p", proto, "-s", srcip, "-j", "DROP")
								if !exist {
									firewall.IPT.Insert("filter", "INPUT", 1, "-i", interfaceName, "-p", proto, "-s", srcip, "-j", "DROP")
								}
							} else {
								exist, _ := firewall.IPT.Exists("filter", "INPUT", "-i", interfaceName, "-p", proto, "-s", srcip, "--dport", dPortOnlyNum, "-j", "DROP")
								if !exist {
									firewall.IPT.Insert("filter", "INPUT", 1, "-i", interfaceName, "-p", proto, "-s", srcip, "--dport", dPortOnlyNum, "-j", "DROP")
								}
							}
							config.AppLog.Warn().Msgf("Block DoS from ip %s", srcip)
						}
					}
				} else {
					if proto == constants.NET_PROTO_ICMP {
						firewall.IPT.Delete("filter", "INPUT", "-i", interfaceName, "-p", proto, "-s", srcip, "-j", "DROP")
						utils.CacheSetWithTimeDuration(constants.SCHEMA_DDOS, key, "1", (time.Second*1)+(time.Millisecond*100))
					} else {
						firewall.IPT.Delete("filter", "INPUT", "-i", interfaceName, "-p", proto, "-s", srcip, "--dport", dPortOnlyNum, "-j", "DROP")
						utils.CacheSetWithTimeDuration(constants.SCHEMA_DDOS, key, "1", (time.Second*1)+(time.Millisecond*100))
					}
				}
			}

		}
	}
}
