package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/services"
	"github.com/mrzack99s/coco-captive-portal/session"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

// Headers godoc
// @Summary To kick via username
// @Description To kick via username
// @ID kick-via-username
// @Accept   json
// @Tags	Operator
// @Produce  json
// @security ApiKeyAuth
// @Param params body types.SessionType true "Parameters"
// @Success 200 {string} string "ok"
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/kick-username [put]
func (ctl *operatorController) kickSessionViaUsername(c *gin.Context) {
	var ss types.SessionType
	if err := c.ShouldBind(&ss); err != nil {
		msg := "kick: bind an interface failed"
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": msg,
		})
		return
	}

	// Get session by ip address
	var sessionIds []string
	err := utils.CacheGet(constants.SCHEMA_MAP_ISSUE_TO_SESSION, ss.Issue, &sessionIds)
	if err != nil {
		msg := fmt.Sprintf("kick: not found session of username %s", ss.Issue)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	for _, s := range sessionIds {
		var nSession types.SessionType
		err = utils.CacheGet(constants.SCHEMA_SESSION, s, &nSession)
		if err != nil {
			msg := fmt.Sprintf("kick: session of id %s", s)
			config.AppLog.Error().Msg(msg)
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": msg,
			})
			return
		}

		err = session.CutOffSession(nSession.SessionUUID)
		if err != nil {
			msg := fmt.Sprintf("kick: %s", err.Error())
			config.AppLog.Error().Msg(msg)
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": msg,
			})
			return
		}
	}

	c.String(http.StatusOK, "kicked")

}

// Headers godoc
// @Summary To kick via ip address
// @Description To kick via ip address
// @ID kick-via-ip-address
// @Accept   json
// @Tags	Operator
// @Produce  json
// @security ApiKeyAuth
// @Param params body types.SessionType true "Parameters"
// @Success 200 {string} string "ok"
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/kick-ip-address [put]
func (ctl *operatorController) kickSessionViaIPAddress(c *gin.Context) {
	var ss types.SessionType
	if err := c.ShouldBind(&ss); err != nil {
		msg := "kick: bind an interface failed"
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": msg,
		})
		return
	}

	// Get session by ip address
	sessionId, err := utils.CacheGetString(constants.SCHEMA_MAP_IP_TO_SESSION, ss.IPAddress)
	if err != nil {
		msg := fmt.Sprintf("kick: not found session of ip %s", ss.IPAddress)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	err = utils.CacheGet(constants.SCHEMA_SESSION, sessionId, &ss)
	if err != nil {
		msg := fmt.Sprintf("kick: session of id %s", sessionId)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	err = session.CutOffSession(ss.SessionUUID)
	if err != nil {
		msg := fmt.Sprintf("kick: %s", err.Error())
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	c.String(http.StatusOK, "kicked")

}

// Headers godoc
// @Summary Get all session
// @Description Get all session
// @ID get-all-session
// @Accept   json
// @Tags	 Operator
// @Produce  json
// @security ApiKeyAuth
// @Success 200 {array} types.SessionType
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/get-all-session [get]
func (ctl *operatorController) getAllSession(c *gin.Context) {
	allSession := []types.SessionType{}
	allKey, err := utils.CacheGetAllKey(constants.SCHEMA_SESSION)
	if err != nil {
		msg := fmt.Sprintf("get-all-session: %s", err.Error())
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	for _, k := range allKey {
		var ss types.SessionType
		err := utils.CacheGetWithRawKey(k, &ss)
		if err != nil {
			msg := fmt.Sprintf("get-all-session: %s", err.Error())
			config.AppLog.Error().Msg(msg)
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": msg,
			})
			return
		}
		allSession = append(allSession, ss)
	}

	c.JSON(http.StatusOK, allSession)

}

// Headers godoc
// @Summary Count all session
// @Description Count all session
// @ID count-all-session
// @Accept   json
// @Tags	 Operator
// @Produce  json
// @security ApiKeyAuth
// @Success 200 {string} string "count"
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/count-all-session [get]
func (ctl *operatorController) countAllSession(c *gin.Context) {
	allKey, err := utils.CacheGetAllKey(constants.SCHEMA_SESSION)
	if err != nil {
		msg := fmt.Sprintf("get-all-session: %s", err.Error())
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	c.JSON(http.StatusOK, len(allKey))

}

// Headers godoc
// @Summary Check is administrator
// @Description Check is administrator
// @ID check-is-administrator
// @Accept   json
// @Tags	Operator
// @Produce  json
// @Param params body types.CredentialType true "Parameters"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/check-is-administrator [post]
func (ctl *operatorController) checkIsAdministrator(c *gin.Context) {
	var ss types.CredentialType
	if err := c.ShouldBind(&ss); err != nil {
		msg := "admin: bind an interface failed"
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": msg,
		})
		return
	}

	if ss.Username == config.Config.Administrator.Credential.Username && utils.Sha512encode(ss.Password) == config.Config.Administrator.Credential.Password {
		newToken := utils.SecretGenerator(64)
		err := utils.CacheSetWithTimeDuration("temp", "admtoken", newToken, time.Hour*1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "admin: generate token failed",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"api_token": newToken,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"msg": "not authorized",
	})

}

// Headers godoc
// @Summary Revoke administrator
// @Description Revoke administrator
// @ID revoke-administrator
// @Accept   json
// @Tags	Operator
// @Produce  json
// @security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/revoke-administrator [get]
func (ctl *operatorController) revokeAdministrator(c *gin.Context) {
	var ss types.CredentialType
	if err := c.ShouldBind(&ss); err != nil {
		msg := "admin: bind an interface failed"
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": msg,
		})
		return
	}

	err := utils.CacheDelete("temp", "admtoken")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "admin: token revoke failed")
		c.Abort()
		return
	}

	c.String(http.StatusOK, "revoked")

}

// Headers godoc
// @Summary Network Interfaces Usage
// @Description Network Interfaces Usage
// @ID net-interfaces-bytes-usage
// @Accept   json
// @Tags	Operator
// @Produce  json
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/net-intf-usage [get]
func (ctl *operatorController) getNetInterfacesUsage(c *gin.Context) {

	secureIntfRx, secureIntfTx := services.GetNetInterfaceBytes(config.Config.SecureInterface)
	egressIntfRx, egressIntfTx := services.GetNetInterfaceBytes(config.Config.EgressInterface)

	c.JSON(http.StatusOK, gin.H{
		"secure_interface": gin.H{
			"rx": secureIntfRx,
			"tx": secureIntfTx,
		},
		"egress_interface": gin.H{
			"rx": egressIntfRx,
			"tx": egressIntfTx,
		},
	})

}
