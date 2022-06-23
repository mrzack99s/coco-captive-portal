package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

// Headers godoc
// @Summary Get initialize secret to get access
// @Description Get initialize secret to get access
// @ID initialize
// @Accept   json
// @Tags	Authentication
// @Produce  json
// @Success 200 {object} types.InitializedType
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/initialize [get]
func (ctl *authController) getInitialize(c *gin.Context) {
	secret := utils.SecretGenerator(32)

	iniObj := types.InitializedType{
		IPAddress: c.ClientIP(),
		Secret:    secret,
	}

	err := utils.CacheSetWithTimeDuration(constants.SCHEMA_SESSION_INITIALIZE, secret, iniObj, time.Minute*2)
	if err != nil {
		msg := fmt.Sprintf("%s via %s", err.Error(), c.ClientIP())
		config.AppLog.Error().Msg(msg)
	}
	c.JSON(http.StatusOK, iniObj)
}

// Headers godoc
// @Summary Exist initialize secret
// @Description Exist initialize secret
// @ID is-exist-initialize-secret
// @Accept   json
// @Tags	Authentication
// @Produce  json
// @Param params body types.InitializedType true "Parameters"
// @Success 200 {string} string "found"
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/is-exist-initialize-secret [post]
func (ctl *authController) isExistInitializeSecret(c *gin.Context) {
	clientIp := c.ClientIP()

	var initialized types.InitializedType
	if err := c.ShouldBind(&initialized); err != nil {
		msg := fmt.Sprintf("%s bind an interface failed", clientIp)
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": msg,
		})
		return
	}

	// Get init session by secret
	var initSession types.InitializedType
	err := utils.CacheGet(constants.SCHEMA_SESSION_INITIALIZE, initialized.Secret, &initSession)
	if err != nil {
		msg := "not found initial secret"
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": msg,
		})
		return
	}

	c.String(http.StatusOK, "found")
}
