package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

type controller struct {
	router gin.IRouter
}

func NewController(router gin.IRouter) *controller {
	s := &controller{
		router: router,
	}
	s.register()
	return s
}

func tokenMiddleware(c *gin.Context) {
	tokenString := c.Request.Header.Get("api-token")

	token, err := utils.CacheGetString(constants.CONFIG, "api-token")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "not found token in config, please renew an api token")
		c.Abort()
		return
	}
	if tokenString != token {
		token, err = utils.CacheGetString("temp", "admtoken")
		if err != nil || tokenString != token {
			c.JSON(http.StatusUnauthorized, "not authorized")
			c.Abort()
			return
		}
	}

	c.Next()

}

func (ctl *controller) register() {
	api := ctl.router.Group("/v1")
	api.GET("/initialize", ctl.getInitialize)
	api.GET("/sign-out", ctl.signout)
	api.GET("/signed", ctl.getSigned)
	api.GET("/adm-signed", ctl.getAdminSigned)
	api.GET("/html-properties", ctl.getHTMLProperties)
	api.POST("/authentication", ctl.getAuthentication)
	api.POST("/check-is-administrator", ctl.checkIsAdministrator)

	api.PUT("/kick-username", tokenMiddleware, ctl.kickSessionViaUsername)
	api.PUT("/kick-ip-address", tokenMiddleware, ctl.kickSessionViaIPAddress)
	api.GET("/get-all-session", tokenMiddleware, ctl.getAllSession)
	api.GET("/revoke-administrator", tokenMiddleware, ctl.revokeAdministrator)

}
