package watcher

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

func CaptivePortalDetector(ctx context.Context, flag ...bool) {
	go func(ctx context.Context, flag ...bool) {

		intIp, err := utils.GetSecureInterfaceIpv4Addr()
		if err != nil {
			panic(err)
		}

		for {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return
			default:
				if flag[0] {
					gin.SetMode(gin.ReleaseMode)
				}

				router := gin.Default()
				router.Use(cors.Default())
				router.Any("/", func(c *gin.Context) {
					if utils.StringContains(config.URL_CAPTIVE_PORTAL_DETECTION, c.Request.Host+c.Request.URL.Path) {
						if config.Config.ExternalPortalURL != "" {
							c.Redirect(http.StatusFound, config.Config.ExternalPortalURL)
							return
						} else {
							redirect(c, intIp)
						}
					} else {
						if config.Config.ExternalPortalURL != "" {
							c.Redirect(http.StatusFound, config.Config.ExternalPortalURL)
							return
						} else {
							redirect(c, intIp)
						}
					}
				})
				router.NoRoute(func(c *gin.Context) {
					if config.Config.ExternalPortalURL != "" {
						c.Redirect(http.StatusFound, config.Config.ExternalPortalURL)
						return
					} else {
						redirect(c, intIp)
					}
				})
				err := router.Run(fmt.Sprintf("%s:8080", intIp))
				if err != nil {
					config.AppLog.Error().Msg("captive-portal-detect-http: " + err.Error())
					return
				}
			}
		}
	}(ctx, flag...)

	go func(ctx context.Context, flag ...bool) {

		intIp, err := utils.GetSecureInterfaceIpv4Addr()
		if err != nil {
			panic(err)
		}

		for {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return
			default:
				if flag[0] {
					gin.SetMode(gin.ReleaseMode)
				}

				router := gin.Default()
				router.Use(cors.Default())
				router.Any("/", func(c *gin.Context) {
					if utils.StringContains(config.URL_CAPTIVE_PORTAL_DETECTION, c.Request.Host+c.Request.URL.Path) {
						if config.Config.ExternalPortalURL != "" {
							c.Redirect(http.StatusFound, config.Config.ExternalPortalURL)
							return
						} else {
							redirect(c, intIp)

						}
					} else {
						if config.Config.ExternalPortalURL != "" {
							c.Redirect(http.StatusFound, config.Config.ExternalPortalURL)
							return
						} else {
							redirect(c, intIp)
						}
					}
				})
				router.NoRoute(func(c *gin.Context) {
					if config.Config.ExternalPortalURL != "" {
						c.Redirect(http.StatusFound, config.Config.ExternalPortalURL)
						return
					} else {
						redirect(c, intIp)
					}
				})

				err := router.RunTLS(fmt.Sprintf("%s:8443", intIp), "./certs/authfullchain.pem", "./certs/authprivkey.pem")
				if err != nil {
					config.AppLog.Error().Msg("captive-portal-detect-https: " + err.Error())
					return
				}
			}
		}
	}(ctx, flag...)
}

func redirect(c *gin.Context, intIp string) {

	if config.Config.DomainNames.AuthDomainName != "" {
		c.Redirect(http.StatusFound, fmt.Sprintf("https://%s", config.Config.DomainNames.AuthDomainName))
		return
	} else {
		c.Redirect(http.StatusFound, fmt.Sprintf("https://%s", intIp))
		return
	}

}
