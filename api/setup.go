package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/config"
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
	token, _ := utils.CacheGetString(constants.SCHEMA_CONFIG, "api-token")
	if tokenString != token {
		token, err := utils.CacheGetString("temp", "admtoken")
		if err != nil || tokenString != token {
			c.JSON(http.StatusUnauthorized, "not authorized")
			c.Abort()
			return
		}
	}
	c.Next()
}
func GetUnAuthirizedNetworkMiddleware(domain string, operator bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/unauthorized" || path == "/api/get-captive-portal-config-fundamental" {
			c.Next()
			return
		}

		clientIp := c.ClientIP()

		if operator {
			for _, cidr := range config.Config.Administrator.AuthorizedNetworks {
				if utils.Ipv4InCidr(cidr, clientIp) {
					c.Next()
					return
				}
			}

		} else {
			for _, cidr := range config.Config.AuthorizedNetworks {
				if utils.Ipv4InCidr(cidr, clientIp) {
					c.Next()
					return
				}
			}

		}

		if strings.Contains(path, "/api") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "your network is not authorized",
			})
			c.Abort()
			return
		}

		c.Redirect(http.StatusFound, fmt.Sprintf("https://%s/unauthorized", domain))
		c.Abort()
	}
}

func (ctl *authController) register() {
	ctl.router.GET("/initialize", ctl.getInitialize)
	ctl.router.GET("/sign-out", ctl.signout)
	ctl.router.GET("/signed", ctl.getSigned)
	ctl.router.POST("/authentication", ctl.getAuthentication)
	ctl.router.POST("/is-exist-initialize-secret", ctl.isExistInitializeSecret)
	ctl.router.GET("/get-captive-portal-config-fundamental", ctl.getCaptivePortalConfigFundamental)
}

func (ctl *operatorController) register() {
	ctl.router.PUT("/kick-username", tokenMiddleware, ctl.kickSessionViaUsername)
	ctl.router.PUT("/kick-ip-address", tokenMiddleware, ctl.kickSessionViaIPAddress)
	ctl.router.PUT("/config", tokenMiddleware, ctl.setConfig)
	ctl.router.PUT("/config-with-restart-system", tokenMiddleware, ctl.setConfigWithRestartSystem)
	ctl.router.GET("/get-all-session", tokenMiddleware, ctl.getAllSession)
	ctl.router.GET("/count-all-session", tokenMiddleware, ctl.countAllSession)
	ctl.router.GET("/revoke-administrator", tokenMiddleware, ctl.revokeAdministrator)
	ctl.router.GET("/config", tokenMiddleware, ctl.getConfig)
	ctl.router.GET("/adm-signed", ctl.getAdminSigned)
	ctl.router.POST("/check-is-administrator", ctl.checkIsAdministrator)
	ctl.router.GET("/net-intf-usage", ctl.getNetInterfacesUsage)
}
