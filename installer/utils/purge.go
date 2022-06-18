package installer_utils

import (
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func Purge() {

	cmds := []CommandType{}

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "stop service",
		Command: *exec.Command("systemctl", "stop", "coco-captive-portal"),
	})

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "disable service",
		Command: *exec.Command("systemctl", "disable", "coco-captive-portal"),
	})

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "remove service",
		Command: *exec.Command("rm", "-rf", "/etc/systemd/system/coco-captive-portal.service"),
	})

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "reload daemon",
		Command: *exec.Command("systemctl", "daemon-reload"),
	})

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "remove workspace",
		Command: *exec.Command("rm", "-rf", APP_DIR),
	})

	cmds = append(cmds, CommandType{
		Type:    COMMAND_PURGE_TYPE,
		Name:    "all package",
		Command: *exec.Command("apt-get", "purge", "--auto-remove", "-y", "libpcap0.8", "redis"),
	})

	for _, cmd := range cmds {
		log.Info().Msg(getMessage(cmd, DOING_STATE))
		if e := cmd.Command.Run(); e != nil {
			if IGNORE_VERIFY {
				log.Warn().Msg(getMessage(cmd, FAILED_STATE))
			} else {
				log.Error().Msg(getMessage(cmd, FAILED_STATE))
				os.Exit(0)
			}
		} else {

			log.Info().Msg(getMessage(cmd, DONE_STATE))
		}
	}

}
