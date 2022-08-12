package installer_utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mrzack99s/coco-captive-portal/utils"
	"github.com/rs/zerolog/log"
)

func resourceVerify(ignore bool) {

	si.GetSysInfo()

	if !ignore {
		if len(si.Network) < 2 {
			log.Error().Msgf("%s need 2 interfaces (can be use bond and link aggregation)", APP_NAME)
			os.Exit(0)
		}

		if si.CPU.Threads < 4 {
			log.Error().Msgf("%s need minimum 4 vCpu (4 threads)", APP_NAME)
			os.Exit(0)
		}

		if si.Memory.Size < 2048 {
			log.Error().Msgf("%s need minimum 2 GB", APP_NAME)
			os.Exit(0)
		}
	}

	if utils.ExistingKeyInMap(OS_SUPPORT, si.OS.Vendor) {

	} else {
		msg := fmt.Sprintf("%s support only ", APP_NAME)
		i := 1
		lenOKey := len(OS_SUPPORT)
		for os, versions := range OS_SUPPORT {
			msg += fmt.Sprintf("%s [", os)
			lenOVersion := len(versions)
			for j, v := range versions {
				if j == lenOVersion-1 {
					msg += fmt.Sprintf("%s]", v)
				} else {
					msg += fmt.Sprintf("%s, ", v)
				}
			}
			if i != lenOKey {
				msg += ", "
			}
			i++
		}

		log.Error().Msg(msg)
		os.Exit(0)
	}

	if RE_INSTALL {
		cmds := []CommandType{}

		cmds = append(cmds, CommandType{
			Type:    COMMAND_EXEC_TYPE,
			Name:    "Remove CoCo Captive Portal Binary",
			Command: *exec.Command("rm", "-f", fmt.Sprintf("%s/coco", APP_DIR)),
		})

		cmds = append(cmds, CommandType{
			Type:    COMMAND_EXEC_TYPE,
			Name:    "Remove Auth-UI",
			Command: *exec.Command("rm", "-rf", fmt.Sprintf("%s/dist-auth-ui", APP_DIR)),
		})

		cmds = append(cmds, CommandType{
			Type:    COMMAND_EXEC_TYPE,
			Name:    "Remove Operator-UI",
			Command: *exec.Command("rm", "-rf", fmt.Sprintf("%s/dist-operator-ui", APP_DIR)),
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

	} else {
		if _, err := os.Stat(fmt.Sprintf("%s/coco", APP_DIR)); os.IsExist(err) {
			log.Warn().Msg("CoCo Captive Portal installed")
			os.Exit(0)
		}

		if _, err := os.Stat(fmt.Sprintf("%s/dist-auth-ui", APP_DIR)); os.IsExist(err) {
			log.Warn().Msg("CoCo Captive Portal installed")
			os.Exit(0)
		}

		if _, err := os.Stat(fmt.Sprintf("%s/dist-operator-ui", APP_DIR)); os.IsExist(err) {
			log.Warn().Msg("CoCo Captive Portal installed")
			os.Exit(0)
		}
	}

}
