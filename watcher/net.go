package watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/constants"
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

				if handle, err := pcap.OpenLive(config.Config.SecureInterface, 1600, true, pcap.BlockForever); err != nil {
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

						config.NetLog.Info().Msgf("proto=%s src=%s,%s dst=%s,%s", proto, fqdn, srcip, sport, dstip, dport)

						go func(srcip string) {
							sessionId := ""
							err := utils.CacheGet(constants.SCHEMA_MAP_IP_TO_SESSION, srcip, &sessionId)
							if err == nil {
								session := types.SessionType{}
								err := utils.CacheGet(constants.SCHEMA_SESSION, sessionId, &session)
								if err == nil {
									session.LastSeen = time.Now()
									utils.CacheSet(
										constants.SCHEMA_SESSION, sessionId, session,
									)
								}
							}
						}(srcip)

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
