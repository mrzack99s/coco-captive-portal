package session

import (
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/firewall"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

func CutOffSession(sessionUUID string) (err error) {

	ss := types.SessionType{}
	err = utils.CacheGet(constants.SESSION, sessionUUID, &ss)
	if err != nil {
		return
	}

	err = firewall.UnallowAccess(&ss)
	if err != nil {
		return
	}

	err = utils.CacheDelete(constants.SESSION, ss.SessionUUID)
	if err != nil {
		return
	}

	err = utils.CacheDelete(constants.MAP_IP_TO_SESSION, ss.IPAddress)
	if err != nil {
		return
	}

	err = utils.CacheDelete(constants.MAP_ISSUE_TO_SESSION, ss.Issue)
	if err != nil {
		return
	}

	return
}
