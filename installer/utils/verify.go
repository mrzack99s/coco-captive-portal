package installer_utils

import (
	"fmt"
	"os"

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

}
