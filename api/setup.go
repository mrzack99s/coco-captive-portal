package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

type authController struct {
	router gin.IRouter
}

type operatorController struct {
	router gin.IRouter
}

func NewAuthController(router gin.IRouter) *authController {
	s := &authController{
		router: router,
	}
	s.register()
	return s
}

func NewOperatorController(router gin.IRouter) *operatorController {
	s := &operatorController{
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

func (ctl *authController) register() {
	ctl.router.GET("/initialize", ctl.getInitialize)
	ctl.router.GET("/sign-out", ctl.signout)
	ctl.router.GET("/signed", ctl.getSigned)
	ctl.router.GET("/html-properties", ctl.getHTMLProperties)
	ctl.router.POST("/authentication", ctl.getAuthentication)
	ctl.router.POST("/is-exist-initialize-secret", ctl.isExistInitializeSecret)
}

func (ctl *operatorController) register() {

	ctl.router.PUT("/kick-username", tokenMiddleware, ctl.kickSessionViaUsername)
	ctl.router.PUT("/kick-ip-address", tokenMiddleware, ctl.kickSessionViaIPAddress)
	ctl.router.GET("/get-all-session", tokenMiddleware, ctl.getAllSession)
	ctl.router.GET("/revoke-administrator", tokenMiddleware, ctl.revokeAdministrator)
	ctl.router.GET("/adm-signed", ctl.getAdminSigned)
	ctl.router.POST("/check-is-administrator", ctl.checkIsAdministrator)
}
