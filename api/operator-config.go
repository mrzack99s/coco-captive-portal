package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/types"
)

// Headers godoc
// @Summary Get config
// @Description Get config
// @ID get-config
// @Accept   json
// @Tags	Operator
// @Produce  json
// @security ApiKeyAuth
// @Success 200 {object} types.ConfigType
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/get-config [get]
func (ctl *operatorController) getConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.Config)
}

// Headers godoc
// @Summary Get config
// @Description Get config
// @ID get-config
// @Accept   json
// @Tags	Operator
// @Produce  json
// @security ApiKeyAuth
// @Param params body types.ConfigType true "Parameters"
// @Success 200 {string} string "updated"
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/set-config [get]
func (ctl *operatorController) setConfig(c *gin.Context) {
	var configs types.ConfigType
	if err := c.ShouldBind(&configs); err != nil {
		msg := "bind a config interface failed"
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": msg,
		})
		return
	}

	config.Config = configs
	config.UpdateConfig()

	c.String(http.StatusOK, "updated")
}
