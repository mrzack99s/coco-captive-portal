package installer_utils

import (
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func finally() (err error) {

	cmds := []CommandType{}

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "grant coco-captive-portal to be executalbe",
		Command: *exec.Command("chmod", "+x", fmt.Sprintf("%s/coco", APP_DIR)),
	})

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "reload daemon",
		Command: *exec.Command("systemctl", "daemon-reload"),
	})

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "enable redis service",
		Command: *exec.Command("systemctl", "enable", "--now", "redis-server"),
	})

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "enable coco-captive-portal service",
		Command: *exec.Command("systemctl", "enable", "coco-captive-portal"),
	})

	for _, cmd := range cmds {
		log.Info().Msg(getMessage(cmd, DOING_STATE))
		if e := cmd.Command.Run(); e != nil {
			if IGNORE_VERIFY {
				log.Warn().Msg(getMessage(cmd, FAILED_STATE))
			} else {
				log.Error().Msg(getMessage(cmd, FAILED_STATE))
				err = e
				return
			}
		} else {
			log.Info().Msg(getMessage(cmd, DONE_STATE))
		}
	}

	exec.Command("bash", "-c", "cat /dev/null > /etc/sysctl.conf").Run()

	networkHardening := []AppendStringToFileType{
		{
			Str:  "net.ipv4.ip_forward=1",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.ipv6.conf.all.disable_ipv6=1",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  fmt.Sprintf("net.netfilter.nf_conntrack_max=%d", (si.Memory.Size*1048576)/16384/2),
			File: "/etc/sysctl.conf",
		},
		{
			Str:  fmt.Sprintf("net.netfilter.nf_conntrack_buckets=%d", ((si.Memory.Size*1048576)/16384/2)/4),
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_generic_timeout=300",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_icmp_timeout=15",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_tcp_timeout_established=86400",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_tcp_timeout_close = 10",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_tcp_timeout_close_wait = 30",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_tcp_timeout_fin_wait = 30",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_tcp_timeout_syn_recv = 30",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_tcp_timeout_syn_sent = 60",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_tcp_timeout_time_wait = 30",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_udp_timeout_stream = 30",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "nf_conntrack",
			File: "/etc/modules",
		},
	}
	for _, hh := range networkHardening {
		if e := AppendStringToFile(hh); e != nil {
			if IGNORE_VERIFY {
				log.Warn().Msg("append was failed")
			} else {
				log.Error().Msg("append was failed")
				err = e
				return
			}
		}
	}

	log.Info().Msg("Reboot Server")
	exec.Command("reboot").Run()

	return
}
