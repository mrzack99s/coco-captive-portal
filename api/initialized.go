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
// @Router /v1/initialize [get]
func (ctl *controller) getInitialize(c *gin.Context) {
	secret := utils.SecretGenerator(32)

	iniObj := types.InitializedType{
		IPAddress: c.ClientIP(),
		Secret:    secret,
	}

	err := utils.CacheSetWithTimeDuration(constants.SESSION_INITIALIZE, secret, iniObj, time.Minute*2)
	if err != nil {
		msg := fmt.Sprintf("%s via %s", err.Error(), c.ClientIP())
		config.AppLog.Error().Msg(msg)
	}
	c.JSON(http.StatusOK, iniObj)
}
