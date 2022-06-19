package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/session"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

// Headers godoc
// @Summary Sign out
// @Description Sign out
// @ID sign-out
// @Accept   json
// @Tags	Authentication
// @Produce  json
// @Success 200 {string} string "ok"
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/sign-out [get]
func (ctl *authController) signout(c *gin.Context) {

	clientIp := c.ClientIP()

	// Get session by ip address
	sessionId, err := utils.CacheGetString(constants.MAP_IP_TO_SESSION, clientIp)
	if err != nil {
		msg := fmt.Sprintf("not found session of %s", clientIp)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	var ss types.SessionType
	err = utils.CacheGet(constants.SESSION, sessionId, &ss)
	if err != nil {
		msg := fmt.Sprintf("session of id %s", sessionId)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	err = session.CutOffSession(ss.SessionUUID)
	if err != nil {
		msg := fmt.Sprintf("%s via %s", err.Error(), clientIp)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	c.String(http.StatusOK, "singed-out")
}
