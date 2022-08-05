package main

import (
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/constants"
	"github.com/mrzack99s/coco-captive-portal/runtime"
	"github.com/mrzack99s/coco-captive-portal/utils"
	"github.com/spf13/cobra"
)

// @title COCO Captive Portal
// @version 1
// @description This is a COCO Captive Portal

// @license.name Apache License Version 2.0
// @license.url https://github.com/mrzack99s/coco-captive-portal

// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name api-token
func main() {

	netLogger := config.LoggingConfig{
		ConsoleLoggingEnabled: false,
		FileLoggingEnabled:    true,
		Directory:             constants.LOG_DIR,
		Filename:              "netlog",
		MaxSize:               50,
		MaxAge:                90,
		MaxBackups:            90,
	}
	config.NetLog = netLogger.ConfigureWithoutDisplay()

	appLogger := config.LoggingConfig{
		ConsoleLoggingEnabled: false,
		FileLoggingEnabled:    true,
		Directory:             constants.LOG_DIR,
		Filename:              "applog",
		MaxSize:               50,
		MaxAge:                90,
		MaxBackups:            90,
	}
	config.AppLog = appLogger.Configure()

	loginLogger := config.LoggingConfig{
		ConsoleLoggingEnabled: false,
		FileLoggingEnabled:    true,
		Directory:             constants.LOG_DIR,
		Filename:              "loginlog",
		MaxSize:               50,
		MaxAge:                90,
		MaxBackups:            90,
	}
	config.LoginLog = loginLogger.Configure()

	if !utils.IsRootPrivilege() {
		panic(`this application needs the ability to run commands as root. We are unable to find either "sudo" or "su" available to make this happen.`)
	}

	var cmdRun = &cobra.Command{
		Use:   "run",
		Short: "To run an application. Default run in development mode",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			runtime.AppRunner(config.PROD_MODE)
		},
	}

	var cmdCertificate = &cobra.Command{
		Use:   "gencert",
		Short: "To generate a self-signed certificate.",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			utils.GenerateSelfSignCert()
		},
	}

	var cmdRenewApiToken = &cobra.Command{
		Use:   "renew-api-token",
		Short: "To renew an api token.",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			utils.GenerateApiToken()
		},
	}

	cmdRun.Flags().BoolVarP(&config.PROD_MODE, "production", "r", false, "Run in production mode")

	var rootCmd = &cobra.Command{Use: "coco"}
	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdCertificate)
	rootCmd.AddCommand(cmdRenewApiToken)
	rootCmd.Execute()
}
