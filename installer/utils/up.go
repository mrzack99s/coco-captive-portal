package installer_utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func UpInstaller() {

	resourceVerify(IGNORE_VERIFY)

	if err := createWorkspace(); err != nil {
		log.Error().Msg("create an workspace failed")
		os.Exit(0)
	}

	if IMPORT_FILE_PATH != "" {
		if e := copy(CopyType{
			Src:  IMPORT_FILE_PATH,
			Dst:  fmt.Sprintf("%s/config.yaml", APP_DIR),
			Perm: 0644,
		}); e != nil {
			if IGNORE_VERIFY {
				log.Warn().Msg("import config file failed")
			} else {
				log.Warn().Msg("import config file failed")
				os.Exit(0)
			}
		}
	} else {

		if RE_INSTALL {
			if confirmWithMsg("Create a new configuration file?: ") {
				cmd := CommandType{
					Type:    COMMAND_EXEC_TYPE,
					Name:    "Remove CoCo Captive Portal Binary",
					Command: *exec.Command("rm", "-f", fmt.Sprintf("%s/coco", APP_DIR)),
				}
				cmd.Command.Run()
				if err := defineConfig(); err != nil {
					log.Error().Msg("define configure failed")
					os.Exit(0)
				}
			}
		} else {
			if err := defineConfig(); err != nil {
				log.Error().Msg("define configure failed")
				os.Exit(0)
			}
		}

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
