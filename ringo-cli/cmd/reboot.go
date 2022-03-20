package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	ringo "github.com/thomas-maurice/ringos/go-client"
)

var RebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboots the device",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := MustGetClient()
		var resp ringo.GenericAnswer
		_, err := client.Post(client.URL("api", "reboot"), map[string]bool{"reboot": true}, &resp)
		if err != nil {
			logrus.WithError(err).Fatal("could not reboot device")
		}
		output(resp)
	},
}
