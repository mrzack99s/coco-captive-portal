package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/types"
	_ "github.com/mrzack99s/coco-captive-portal/types"
)

// Headers godoc
// @Summary Captive portal config fundamental
// @Description Captive portal config fundamental
// @ID get-captive-portal-config-fundamental
// @Accept   json
// @Tags	Utils
// @Produce  json
// @Success 200 {object} types.CaptivePortalConfigFundamentalType
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/get-captive-portal-config-fundamental [get]
func (ctl *authController) getCaptivePortalConfigFundamental(c *gin.Context) {
	response := types.CaptivePortalConfigFundamentalType{
		HTML: config.Config.HTML,
	}
	if config.Config.LDAP != nil {
		response.Mode = "ldap"
		response.SingleDomain = config.Config.LDAP.SingleDomain
		response.DomainNames = config.Config.LDAP.DomainNames
	} else {
		response.Mode = "radius"
	}
	c.JSON(http.StatusOK, response)
}
