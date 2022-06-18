package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

// Headers godoc
// @Summary Get signed by ip address
// @Description Get signed by ip address
// @ID signed
// @Accept   json
// @Tags	Authentication
// @Produce  json
// @Success 200 {object} types.SessionType
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /v1/signed [get]
func (ctl *controller) getSigned(c *gin.Context) {

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

	c.JSON(http.StatusOK, ss)
}

// Headers godoc
// @Summary Get admin signed
// @Description Get admin signed
// @ID adm-signed
// @Accept   json
// @Tags	Authentication
// @Produce  json
// @Success 200 {string} string "ok"
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /v1/adm-signed [get]
func (ctl *controller) getAdminSigned(c *gin.Context) {

	// Get session by ip address
	_, err := utils.CacheGetString("temp", "admtoken")
	if err != nil {
		msg := "admin: not found token"
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	c.String(http.StatusOK, "found")
}
