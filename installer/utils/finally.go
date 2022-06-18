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
		Name:    "enable service",
		Command: *exec.Command("systemctl", "enable", "--now", "coco-captive-portal"),
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

	return
}
