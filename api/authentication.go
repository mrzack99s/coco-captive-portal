package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/session"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

// Headers godoc
// @Summary Authentication
// @Description Check credential to get access
// @ID authentication
// @Accept   json
// @Tags	Authentication
// @Produce  json
// @Param params body types.CheckCredentialType true "Parameters"
// @Success 200 {object} types.AuthorizedResponseType
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /v1/authentication [post]
func (ctl *controller) getAuthentication(c *gin.Context) {

	clientIp := c.ClientIP()

	var checkCredential types.CheckCredentialType
	if err := c.ShouldBind(&checkCredential); err != nil {
		msg := fmt.Sprintf("%s bind an interface failed", clientIp)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": msg,
		})
		return
	}

	// Get init session by secret
	var initSession types.InitializedType
	err := utils.CacheGet(constants.SESSION_INITIALIZE, checkCredential.Secret, &initSession)
	if err != nil {
		msg := fmt.Sprintf("wrong initial secret from %s", clientIp)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	if initSession.IPAddress != clientIp {
		config.AppLog.Error().Msgf("%s not match via secret %s", clientIp, initSession.Secret)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "ip address not match",
		})
		return
	}

	// Fine logged by ip
	_, err = utils.CacheGetString(constants.MAP_IP_TO_SESSION, clientIp)
	if err == nil {
		msg := fmt.Sprintf("%s signed", clientIp)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	// Fine logged by username
	var arrSid []string
	utils.CacheGet(constants.MAP_ISSUE_TO_SESSION, checkCredential.Username, &arrSid)

	if len(arrSid) >= int(config.Config.MaxConcurrentSession) {
		msg := fmt.Sprintf("user %s reached the limit concurrent session and sign-in via %s", checkCredential.Username, clientIp)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	if config.Config.LDAP != nil {
		err := config.Config.LDAP.Authentication(checkCredential.Username, checkCredential.Password)
		if err != nil {
			msg := fmt.Sprintf("ldap authentication: user %s login failed, via %s", checkCredential.Username, clientIp)
			config.AppLog.Error().Msg(msg)
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": msg,
			})
			return
		}
	} else if config.Config.Radius != nil {
		err := config.Config.Radius.Authentication(checkCredential.Username, checkCredential.Password)
		if err != nil {
			msg := fmt.Sprintf("radius authentication: user %s login failed, via %s", checkCredential.Username, clientIp)
			config.AppLog.Error().Msg(msg)
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": msg,
			})
			return
		}
	} else {
		msg := "not have an authentication"
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
	}

	newSession := types.SessionType{
		Issue:     checkCredential.Username,
		IPAddress: clientIp,
		LastSeen:  time.Now(),
	}

	err = session.NewSession(&newSession)
	if err != nil {
		msg := fmt.Sprintf("%s via %s", err.Error(), clientIp)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	err = utils.CacheDelete(constants.SESSION_INITIALIZE, initSession.Secret)
	if err != nil {
		msg := fmt.Sprintf("%s via %s", err.Error(), clientIp)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	c.JSON(http.StatusOK, types.AuthorizedResponseType{
		Status:      constants.STATUS_OK,
		Issue:       checkCredential.Username,
		RedirectURL: config.Config.RedirectURL,
	})
}
