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

func CaptivePortalDetector(ctx context.Context) {
	go func(ctx context.Context) {

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
				router := gin.Default()
				router.Use(cors.Default())
				router.Any("/", func(c *gin.Context) {
					if utils.StringContains(config.URL_CAPTIVE_PORTAL_DETECTION, c.Request.Host+c.Request.URL.Path) {
						if config.Config.ExternalPortalURL != "" {
							c.Redirect(http.StatusFound, config.Config.ExternalPortalURL)
							return
						} else {
							c.Redirect(http.StatusFound, fmt.Sprintf("https://%s:1800", intIp))
							return
						}
					} else {
						if config.Config.ExternalPortalURL != "" {
							c.Redirect(http.StatusFound, config.Config.ExternalPortalURL)
							return
						} else {
							c.Redirect(http.StatusFound, fmt.Sprintf("https://%s:1800", intIp))
							return
						}
					}
				})
				router.NoRoute(func(c *gin.Context) {
					if config.Config.ExternalPortalURL != "" {
						c.Redirect(http.StatusFound, config.Config.ExternalPortalURL)
						return
					} else {
						c.Redirect(http.StatusFound, fmt.Sprintf("https://%s:1800", intIp))
						return
					}
				})
				router.Run(":8080")
			}
		}
	}(ctx)

	go func(ctx context.Context) {

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
				router := gin.Default()
				router.Use(cors.Default())
				router.Any("/", func(c *gin.Context) {
					if utils.StringContains(config.URL_CAPTIVE_PORTAL_DETECTION, c.Request.Host+c.Request.URL.Path) {
						if config.Config.ExternalPortalURL != "" {
							c.Redirect(http.StatusFound, config.Config.ExternalPortalURL)
							return
						} else {
							c.Redirect(http.StatusFound, fmt.Sprintf("https://%s:1800", intIp))
							return
						}
					} else {
						if config.Config.ExternalPortalURL != "" {
							c.Redirect(http.StatusFound, config.Config.ExternalPortalURL)
							return
						} else {
							c.Redirect(http.StatusFound, fmt.Sprintf("https://%s:1800", intIp))
							return
						}
					}
				})
				router.NoRoute(func(c *gin.Context) {
					if config.Config.ExternalPortalURL != "" {
						c.Redirect(http.StatusFound, config.Config.ExternalPortalURL)
						return
					} else {
						c.Redirect(http.StatusFound, fmt.Sprintf("https://%s:1800", intIp))
						return
					}
				})
				router.RunTLS(":8443", "./certs/fullchain.pem", "./certs/privkey.pem")
			}
		}
	}(ctx)
}
