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

	enableForward := ReplaceWordInFileType{
		OldWord: "#net.ipv4.ip_forward=1",
		NewWord: "net.ipv4.ip_forward=1",
		File:    "/etc/sysctl.conf",
	}
	if e := ReplaceInFile(enableForward); e != nil {
		if IGNORE_VERIFY {
			log.Warn().Msg("replace was failed")
		} else {
			log.Error().Msg("replace was failed")
			err = e
			return
		}
	}

	disableIpv6 := ReplaceWordInFileType{
		OldWord: "#net.ipv6.conf.all.disable_ipv6=1",
		NewWord: "net.ipv6.conf.all.disable_ipv6=1",
		File:    "/etc/sysctl.conf",
	}
	if e := ReplaceInFile(disableIpv6); e != nil {
		if IGNORE_VERIFY {
			log.Warn().Msg("replace was failed")
		} else {
			log.Error().Msg("replace was failed")
			err = e
			return
		}
	}

	networkHardening := []AppendStringToFileType{
		{
			Str:  fmt.Sprintf("net.netfilter.nf_conntrack_max=%d", (si.Memory.Size*1048576)/16384/2),
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_generic_timeout=60",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_icmp_timeout=10",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_tcp_timeout_established=1800",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_tcp_timeout_close_wait = 20",
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
			Str:  "net.netfilter.nf_conntrack_tcp_timeout_time_wait = 60",
			File: "/etc/sysctl.conf",
		},
		{
			Str:  "net.netfilter.nf_conntrack_udp_timeout_stream = 60",
			File: "/etc/sysctl.conf",
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

	cmd := CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "commit enable sysctl.conf",
		Command: *exec.Command("sysctl", "-p"),
	}
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

	return
}
