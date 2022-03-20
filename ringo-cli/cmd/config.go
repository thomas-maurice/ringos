package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	ringo "github.com/thomas-maurice/ringos/go-client"
)

var (
	configHostname        string
	configLEDs            int
	configPassword        string
	configDisablePassword bool
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Device configuration",
	Long:  ``,
}

var ConfigShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays device configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := MustGetClient()
		config, err := client.GetConfig()
		if err != nil {
			logrus.WithError(err).Fatal("could not fetch config")
		}
		output(config)
	},
}

var ConfigUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		newConfig := ringo.Config{}
		if cmd.Flag("hostname").Changed {
			newConfig.Hostname = ringo.String(configHostname)
		}
		if cmd.Flag("leds").Changed {
			newConfig.LEDs = ringo.Int(configLEDs)
		}
		if cmd.Flag("new-password").Changed {
			newConfig.Password = ringo.String(configPassword)
		}
		if cmd.Flag("disable-password").Changed {
			newConfig.DisablePassword = ringo.Bool(configDisablePassword)
		}

		client := MustGetClient()
		resp, err := client.UpdateConfig(newConfig)
		if err != nil {
			logrus.WithError(err).Fatal("could not update config")
		}
		output(resp)
	},
}

func initConfigCmd() {
	ConfigCmd.AddCommand(ConfigShowCmd)
	ConfigCmd.AddCommand(ConfigUpdateCmd)

	ConfigUpdateCmd.PersistentFlags().StringVar(&configHostname, "hostname", "", "Changes the hostname of the device")
	ConfigUpdateCmd.PersistentFlags().IntVarP(&configLEDs, "leds", "l", 0, "Changes the number of leds of the device")
	ConfigUpdateCmd.PersistentFlags().StringVar(&configPassword, "new-password", "", "Changes the password of the device")
	ConfigUpdateCmd.PersistentFlags().BoolVar(&configDisablePassword, "disable-password", false, "Disables password on the device")
}
