package watcher

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/firewall"
	"github.com/mrzack99s/coco-captive-portal/session"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

func NetWatcher(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return
			default:
				intIpv4, err := utils.GetSecureInterfaceIpv4Addr()
				if err != nil {
					panic(err)
				}

				if handle, err := pcap.OpenLive(config.Config.SecureInterface, 65536, true, pcap.BlockForever); err != nil {
					panic(err)
				} else if err := handle.SetBPFFilter(fmt.Sprintf("not src %s and not dst %s", intIpv4, intIpv4)); err != nil { // optional
					panic(err)
				} else {
					packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

					for packet := range packetSource.Packets() {

						srcip := ""
						dstip := ""
						fqdn := ""
						proto := ""
						sport := ""
						dport := ""

						if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
							ip, _ := ipLayer.(*layers.IPv4)
							srcip = ip.SrcIP.String()
							dstip = ip.DstIP.String()
						}

						if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
							proto = "TCP"
							tcp, _ := tcpLayer.(*layers.TCP)
							sport = tcp.SrcPort.String()
							dport = tcp.DstPort.String()
						}

						if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
							proto = "UDP"
							udp, _ := udpLayer.(*layers.UDP)
							sport = udp.SrcPort.String()
							dport = udp.DstPort.String()
						}

						if udpLayer := packet.Layer(layers.LayerTypeDNS); udpLayer != nil {
							dns, _ := udpLayer.(*layers.DNS)
							for _, dnsQuestion := range dns.Questions {
								fqdn = string(dnsQuestion.Name)
							}
						}

						config.NetLog.Info().Msgf("proto=%s fqdn=%s src=%s,%s dst=%s,%s", proto, fqdn, srcip, sport, dstip, dport)

					}
				}
			}

		}
	}(ctx)
}

func NetIdleChecking(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return
			default:
				allKey, err := utils.CacheGetAllKey("session")
				if err != nil {
					config.AppLog.Error().Msg(err.Error())
					continue
				}

				for _, key := range allKey {

					ss := types.SessionType{}
					err := utils.CacheGetWithRawKey(key, &ss)
					if err != nil {
						config.AppLog.Error().Msg(err.Error())
						continue
					}

					now := time.Now()
					diffTime := now.Sub(ss.LastSeen)

					if diffTime.Minutes() > float64(config.Config.SessionIdle) {
						err = session.CutOffSession(ss.SessionUUID)
						if err != nil {
							config.AppLog.Error().Msg(err.Error())
							continue
						}
					} else {
						ss.LastSeen = time.Now()
						if utils.CacheFindExistingRawKey(key) {
							utils.CacheSet("session", ss.SessionUUID, ss)
						}
					}
				}
			}
		}
	}(ctx)
}

func DDOSPreventor(ctx context.Context) {
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
	RATE_LIMIT := 1024 // Per second
	ignore_port := []string{
		"22",
	}

	if !config.PROD_MODE {
		ignore_port = append(ignore_port, "6379")
	}

	if handle, err := pcap.OpenLive(interfaceName, 65536, true, pcap.BlockForever); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter(fmt.Sprintf("not src %s and dst %s", intIpv4, intIpv4)); err != nil { // optional
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		for packet := range packetSource.Packets() {

			srcip := ""
			proto := ""
			dport := ""
			dPortOnlyNum := ""

			if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
				ip, _ := ipLayer.(*layers.IPv4)
				srcip = ip.SrcIP.String()
			}

			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				proto = "TCP"
				tcp, _ := tcpLayer.(*layers.TCP)
				dport = tcp.DstPort.String()
				regex := regexp.MustCompile("[a-zA-Z()]+")
				dportSplit := regex.Split(dport, -1)
				dPortOnlyNum = dportSplit[0]
			}

			if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
				proto = "UDP"
				udp, _ := udpLayer.(*layers.UDP)
				dport = udp.DstPort.String()
				regex := regexp.MustCompile("[a-zA-Z()]+")
				dportSplit := regex.Split(dport, -1)
				dPortOnlyNum = dportSplit[0]
			}

			if !utils.ExistingInArray(ignore_port, dPortOnlyNum) {
				if utils.CacheFindExistingKey(constants.SCHEMA_DDOS, fmt.Sprintf("%s-ip%s-p%s", proto, srcip, dport)) {
					xtimesStr, err := utils.CacheGetString(constants.SCHEMA_DDOS, fmt.Sprintf("%s-ip%s-p%s", proto, srcip, dport))
					if err != nil {
						xtimesStr = "0"
					}

					if xtimesStr != "blocked" {
						xtimes, err := utils.StringToInt64(xtimesStr)
						if err != nil {
							xtimes = 0
						}
						if xtimes < int64(RATE_LIMIT) {

							firewall.IPT.Delete("filter", "INPUT", "-i", interfaceName, "-p", proto, "-s", srcip, "--dport", dPortOnlyNum, "-j", "DROP")

							xtimes++
							xtimesStr = utils.Int64ToString(xtimes)

							utils.CacheSetWithTimeDuration(constants.SCHEMA_DDOS, fmt.Sprintf("%s-ip%s-p%s", proto, srcip, dport), xtimesStr, time.Second*1)

						} else {
							utils.CacheSetWithTimeDuration(constants.SCHEMA_DDOS, fmt.Sprintf("%s-ip%s-p%s", proto, srcip, dport), "blocked", time.Minute*1)
							exist, _ := firewall.IPT.Exists("filter", "INPUT", "-i", interfaceName, "-p", proto, "-s", srcip, "--dport", dPortOnlyNum, "-j", "DROP")
							if !exist {
								firewall.IPT.Insert("filter", "INPUT", 1, "-i", interfaceName, "-p", proto, "-s", srcip, "--dport", dPortOnlyNum, "-j", "DROP")
							}
						}
					}
				} else {

					firewall.IPT.Delete("filter", "INPUT", "-i", interfaceName, "-p", proto, "-s", srcip, "--dport", dPortOnlyNum, "-j", "DROP")
					utils.CacheSetWithTimeDuration(constants.SCHEMA_DDOS, fmt.Sprintf("%s-ip%s-p%s", proto, srcip, dport), "1", time.Second*30)

				}
			}

		}
	}
}
