package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Gets the status of the device",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := MustGetClient()
		status, err := client.Status()
		if err != nil {
			logrus.WithError(err).Fatal("could not get device status")
		}
		output(status)
	},
}
