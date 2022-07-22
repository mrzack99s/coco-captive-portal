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
// @Router /api/signed [get]
func (ctl *authController) getSigned(c *gin.Context) {

	clientIp := c.ClientIP()

	// Get session by ip address
	sessionId, err := utils.CacheGetString(constants.SCHEMA_MAP_IP_TO_SESSION, clientIp)
	if err != nil {
		msg := fmt.Sprintf("not found session of %s", clientIp)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	var ss types.SessionType
	err = utils.CacheGet(constants.SCHEMA_SESSION, sessionId, &ss)
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
// @Router /api/adm-signed [get]
func (ctl *operatorController) getAdminSigned(c *gin.Context) {

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

	tokenString := c.Request.Header.Get("api-token")
	token, _ := utils.CacheGetString("temp", "admtoken")
	if tokenString != token {
		msg := "admin: token not correct"
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	c.String(http.StatusOK, "found")
}
