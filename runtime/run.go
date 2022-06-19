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

	go operatorRuntime(flag...)
	authRuntime(flag...)

}

func authRuntime(flag ...bool) {
	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()

	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/docs/index.html")
	})

	if config.Config.ExternalPortalURL == "" {
		router.Use(static.Serve("/", static.LocalFile("dist-auth-ui", true)))
		router.NoRoute(func(c *gin.Context) {
			c.File("dist-auth-ui/index.html")
		})
	}

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "api-token"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}

	if flag[0] {

		if config.Config.ExternalPortalURL != "" {
			corsConfig.AllowOrigins = []string{fmt.Sprintf("https://%s", interfaceIp), config.Config.ExternalPortalURL}
		} else {
			corsConfig.AllowOrigins = []string{fmt.Sprintf("https://%s", interfaceIp)}
		}

		corsConfig.AllowAllOrigins = false

		router.Use(cors.New(corsConfig))
		apiEndpoint := router.Group("/api")
		api.NewAuthController(apiEndpoint)
		err := router.RunTLS(fmt.Sprintf("%s:443", interfaceIp), "./certs/authfullchain.pem", "./certs/authprivkey.pem")
		if err != nil {
			config.AppLog.Error().Msg(err.Error())
			return
		}

	} else {
		router.Use(cors.New(corsConfig))
		apiEndpoint := router.Group("/api")
		api.NewAuthController(apiEndpoint)
		err := router.RunTLS(":443", "./certs/authfullchain.pem", "./certs/authprivkey.pem")
		if err != nil {
			config.AppLog.Error().Msg(err.Error())
			return
		}
	}
}

func operatorRuntime(flag ...bool) {
	interfaceIp, _ := utils.GetEgressInterfaceIpv4Addr()

	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/docs/index.html")
	})
	router.Use(static.Serve("/", static.LocalFile("dist-operator-ui", true)))
	router.NoRoute(func(c *gin.Context) {
		c.File("dist-operator-ui/index.html")
	})

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "api-token"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}

	if flag[0] {

		if config.Config.ExternalPortalURL != "" {
			corsConfig.AllowOrigins = []string{fmt.Sprintf("https://%s", interfaceIp), config.Config.ExternalPortalURL}
		} else {
			corsConfig.AllowOrigins = []string{fmt.Sprintf("https://%s", interfaceIp)}
		}

		corsConfig.AllowAllOrigins = false

		router.Use(cors.New(corsConfig))
		apiEndpoint := router.Group("/api")
		api.NewOperatorController(apiEndpoint)
		err := router.RunTLS(fmt.Sprintf("%s:443", interfaceIp), "./certs/operatorfullchain.pem", "./certs/operatorprivkey.pem")
		if err != nil {
			config.AppLog.Error().Msg(err.Error())
			return
		}
	} else {
		fmt.Println("sadasd")
		router.Use(cors.New(corsConfig))
		apiEndpoint := router.Group("/api")
		api.NewOperatorController(apiEndpoint)
		err := router.RunTLS(":4443", "./certs/operatorfullchain.pem", "./certs/operatorprivkey.pem")
		if err != nil {
			config.AppLog.Error().Msg(err.Error())
			return
		}
	}
}
