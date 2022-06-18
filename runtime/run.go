package runtime

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/api"
	"github.com/mrzack99s/coco-captive-portal/config"
	_ "github.com/mrzack99s/coco-captive-portal/docs"
	"github.com/mrzack99s/coco-captive-portal/firewall"
	"github.com/mrzack99s/coco-captive-portal/utils"
	"github.com/mrzack99s/coco-captive-portal/watcher"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AppRunner(flag ...bool) {

	err := firewall.InitializeCaptivePortal()
	if err != nil {
		panic(err)
	}

	watcher.NetWatcher(context.Background())
	watcher.NetIdleChecking(context.Background())
	watcher.CaptivePortalDetector(context.Background(), flag...)

	if config.Config.LDAP != nil {
		if err := config.Config.LDAP.Connect(); err != nil {
			config.AppLog.Error().Msg("Cannot connect to LDAP Server")
			os.Exit(0)
		}
	}

	if flag[0] {
		gin.SetMode(gin.ReleaseMode)
	}

	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()

	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/docs/index.html")
	})
	router.Use(static.Serve("/", static.LocalFile("dist-ui", true)))
	router.NoRoute(func(c *gin.Context) {
		c.File("dist-ui/index.html")
	})

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "api-token"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}

	if flag[0] {

		if config.Config.ExternalPortalURL != "" {

			_, host, port, err := utils.ParseURL(config.Config.ExternalPortalURL)
			if err != nil {
				panic(err)
			}

			corsConfig.AllowOrigins = []string{config.Config.ExternalPortalURL}
			corsConfig.AllowAllOrigins = false

			router.Use(cors.New(corsConfig))
			apiEndpoint := router.Group("/api")
			api.NewController(apiEndpoint)
			router.RunTLS(fmt.Sprintf("%s:%s", host, port), "./certs/fullchain.pem", "./certs/privkey.pem")

		} else {
			corsConfig.AllowOrigins = []string{fmt.Sprintf("https://%s:1800", interfaceIp)}
			corsConfig.AllowAllOrigins = false

			router.Use(cors.New(corsConfig))
			apiEndpoint := router.Group("/api")
			api.NewController(apiEndpoint)
			router.RunTLS(fmt.Sprintf("%s:1800", interfaceIp), "./certs/fullchain.pem", "./certs/privkey.pem")
		}
	} else {

		router.Use(cors.New(corsConfig))
		apiEndpoint := router.Group("/api")
		api.NewController(apiEndpoint)
		router.RunTLS(":1800", "./certs/fullchain.pem", "./certs/privkey.pem")
	}

}
