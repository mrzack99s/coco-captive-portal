package installer_utils

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

func UpInstaller() {

	resourceVerify(IGNORE_VERIFY)

	if err := createWorkspace(); err != nil {
		log.Error().Msg("create an workspace failed")
		os.Exit(0)
	}

	if err := defineConfig(); err != nil {
		log.Error().Msg("define configure failed")
		os.Exit(0)
	}

	fmt.Print("\n### Installation\n\n")

	if err := installPackages(); err != nil {
		log.Error().Msg("package install failed")
		os.Exit(0)
	}

	if err := downloadPackages(); err != nil {
		log.Error().Msg("download failed")
		os.Exit(0)
	}

	if err := extract(); err != nil {
		log.Error().Msg("extract coco-dist-ui failed")
		os.Exit(0)
	}

	if err := finally(); err != nil {
		log.Error().Msg("extract coco-dist-ui failed")
		os.Exit(0)
	}

}
