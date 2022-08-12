package installer_utils

import "github.com/zcalusic/sysinfo"

var (
	APP_DIR          = "/var/coco-captive-portal"
	IGNORE_VERIFY    bool
	RE_INSTALL       bool
	IMPORT_FILE_PATH string
	OS_SUPPORT       = map[string][]string{
		"ubuntu": {"18.04", "20.04"},
		"debian": {"10", "11"},
	}
	si sysinfo.SysInfo
)
