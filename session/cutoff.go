package session

import (
	"time"

	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/firewall"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

func CutOffSession(sessionUUID string) (err error) {

	ss := types.SessionType{}
	err = utils.CacheGet(constants.SCHEMA_SESSION, sessionUUID, &ss)
	if err != nil {
		return
	}

	err = firewall.UnallowAccess(&ss)
	if err != nil {
		return
	}

	err = utils.CacheSetWithTimeDuration(constants.SCHEMA_MAP_IP_TO_OUT_SESSION, ss.IPAddress, ss.SessionUUID, time.Hour*1)
	if err != nil {
		return
	}

	err = utils.CacheDelete(constants.SCHEMA_SESSION, ss.SessionUUID)
	if err != nil {
		return
	}

	err = utils.CacheDelete(constants.SCHEMA_MAP_IP_TO_SESSION, ss.IPAddress)
	if err != nil {
		return
	}

	err = utils.CacheDelete(constants.SCHEMA_MAP_ISSUE_TO_SESSION, ss.Issue)
	if err != nil {
		return
	}

	return
}
