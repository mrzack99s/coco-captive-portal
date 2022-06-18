package session

import (
	"github.com/google/uuid"
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/firewall"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

func NewSession(session *types.SessionType) (err error) {

	sessionUUIDObj := uuid.New()
	sessionUUID := sessionUUIDObj.String()
	session.SessionUUID = sessionUUID

	err = firewall.AllowAccess(session)
	if err != nil {
		return
	}

	err = utils.CacheSet(constants.SESSION, sessionUUID, *session)
	if err != nil {
		return
	}

	err = utils.CacheSet(constants.MAP_IP_TO_SESSION, session.IPAddress, sessionUUID)
	if err != nil {
		return
	}

	missue2session := []string{}
	utils.CacheGet(constants.MAP_ISSUE_TO_SESSION, session.Issue, &missue2session)
	missue2session = append(missue2session, sessionUUID)

	err = utils.CacheSet(constants.MAP_ISSUE_TO_SESSION, session.Issue, missue2session)
	if err != nil {
		return
	}

	return
}
