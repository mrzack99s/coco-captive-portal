package installer_utils

import (
	"os/exec"

	"github.com/rs/zerolog/log"
)

func installPackages() (err error) {
	packages := []CommandType{}
	log.Info().Msg("# install packages")
	switch si.OS.Vendor {
	case "ubuntu":
		packages = append(packages, CommandType{
			Type:    COMMAND_INSTALL_TYPE,
			Name:    "libpcap0.8",
			Command: *exec.Command("apt-get", "install", "-y", "libpcap0.8"),
		})

		packages = append(packages, CommandType{
			Type:    COMMAND_EXEC_TYPE,
			Name:    "download redis key",
			Command: *exec.Command("curl", "-L", "https://packages.redis.io/gpg", "-o", "/tmp/redis-gpg"),
		})

		packages = append(packages, CommandType{
			Type:    COMMAND_IMPORT_KEY_TYPE,
			Name:    "import redis key",
			Command: *exec.Command("gpg", "--dearmor", "-o", "/usr/share/keyrings/redis-archive-keyring.gpg", "/tmp/redis-gpg"),
		})

		packages = append(packages, CommandType{
			Type:    COMMAND_IMPORT_KEY_TYPE,
			Name:    "import apt list",
			Command: *exec.Command("echo", "\"deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main\"", ">", "/etc/apt/sources.list.d/redis.list"),
		})

		packages = append(packages, CommandType{
			Type:    COMMAND_UPDATE_TYPE,
			Name:    "repo",
			Command: *exec.Command("apt-get", "update"),
		})

		packages = append(packages, CommandType{
			Type:    COMMAND_INSTALL_TYPE,
			Name:    "redis",
			Command: *exec.Command("apt-get", "install", "-y", "redis"),
		})

		err = exeCmds(packages)
		if err != nil {
			return
		}

	}
	return
}

func exeCmds(cmds []CommandType) (err error) {
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
