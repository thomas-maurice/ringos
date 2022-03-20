package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	ringo "github.com/thomas-maurice/ringos/go-client"
)

var (
	colourNewColour     string
	colourNewMode       string
	colourNewBrightness int
)

var ColourCmd = &cobra.Command{
	Use:   "colour",
	Short: "Device colour config",
	Long:  ``,
}

var ColourShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays device colour configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := MustGetClient()
		colour, err := client.Colour()
		if err != nil {
			logrus.WithError(err).Fatal("could not fetch colour state")
		}
		output(colour)
	},
}

var ColourUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the colour configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		colour := ringo.ColourRequest{}
		if cmd.Flag("mode").Changed {
			colour.Mode = ringo.String(colourNewMode)
		}
		if cmd.Flag("colour").Changed {
			colour.Colour = ringo.String(colourNewColour)
		}
		if cmd.Flag("brightness").Changed {
			colour.Brightness = ringo.Int64(int64(colourNewBrightness))
		}
		colour.Persist = &persist

		client := MustGetClient()
		resp, err := client.SetColour(&colour)
		if err != nil {
			logrus.WithError(err).Fatal("could not update config")
		}
		output(resp)
	},
}

func initcolourCmd() {
	ColourCmd.AddCommand(ColourShowCmd)
	ColourCmd.AddCommand(ColourUpdateCmd)

	ColourUpdateCmd.PersistentFlags().StringVarP(&colourNewMode, "mode", "m", "", "Changes the mode of the device")
	ColourUpdateCmd.PersistentFlags().StringVarP(&colourNewColour, "colour", "c", "", "Changes the colour of the device")
	ColourUpdateCmd.PersistentFlags().IntVarP(&colourNewBrightness, "brightness", "b", 0, "Changes the brightness of the device")
}
