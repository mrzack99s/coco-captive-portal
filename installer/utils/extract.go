package installer_utils

import (
	"os/exec"

	"github.com/rs/zerolog/log"
)

func extract() (err error) {
	cmds := []CommandType{}

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "extract coco-auth-ui",
		Command: *exec.Command("tar", "-zxf", "/tmp/coco-dist-auth-ui.tar.gz", "--directory", APP_DIR),
	})

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "extract coco-operator-ui",
		Command: *exec.Command("tar", "-zxf", "/tmp/coco-dist-operator-ui.tar.gz", "--directory", APP_DIR),
	})

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "remove temp coco-auth-ui",
		Command: *exec.Command("rm", "-rf", "/tmp/coco-dist-auth-ui.tar.gz"),
	})

	cmds = append(cmds, CommandType{
		Type:    COMMAND_EXEC_TYPE,
		Name:    "remove temp coco-operator-ui",
		Command: *exec.Command("rm", "-rf", "/tmp/coco-dist-operator-ui.tar.gz"),
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
