package installer_utils

import (
	"os"

	"github.com/rs/zerolog/log"
)

func createWorkspace() (err error) {
	if _, err = os.Stat(APP_DIR); os.IsNotExist(err) {
		err = os.Mkdir(APP_DIR, 0744)
		if err != nil {
			log.Error().Msg("workspace create failed")
			return
		} else {
			log.Info().Msg("workspace created")
		}
	} else {
		log.Info().Msg("exist workspace")
	}

	if _, err = os.Stat(APP_DIR + "/certs"); os.IsNotExist(err) {
		err = os.Mkdir(APP_DIR+"/certs", 0744)
		if err != nil {
			log.Error().Msg("certs workspace create failed")
			return
		} else {
			log.Info().Msg("certs workspace created")
		}
	} else {
		log.Info().Msg("exist certs workspace")
	}
	return
}
