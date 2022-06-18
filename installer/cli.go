package main

import (
	installer_utils "github.com/mrzack99s/coco-captive-portal/installer/utils"
	"github.com/spf13/cobra"
)

func main() {

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
