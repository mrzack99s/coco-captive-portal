package main

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"

	installer_utils "github.com/mrzack99s/coco-captive-portal/installer/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	var cmdUp = &cobra.Command{
		Use:   "up",
		Short: "To install the " + installer_utils.APP_NAME,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			installer_utils.UpInstaller()
		},
	}

	var cmdPurge = &cobra.Command{
		Use:   "purge",
		Short: "To uninstall the " + installer_utils.APP_NAME,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			installer_utils.Purge()
		},
	}

	cmdUp.Flags().BoolVar(&installer_utils.IGNORE_VERIFY, "ignore", false, "Ignore some resource verify")

	var rootCmd = &cobra.Command{Use: "coco-installer"}
	rootCmd.AddCommand(cmdUp)
	rootCmd.AddCommand(cmdPurge)
	rootCmd.Execute()
}
