package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	marshaller string
	address    string
	username   string
	password   string
	persist    bool
	debug      bool
)

var rootCmd = &cobra.Command{
	Use:   "ringo-cli",
	Short: "ringo cli",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	InitRootCmd()
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version number",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("devel")
	},
}

func InitRootCmd() {
	rootCmd.PersistentFlags().StringVarP(&marshaller, "output", "o", "json", "Output format for server responses (client mode), must be yaml or json")
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "Address to connect to")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "admin", "Username to authenticate as")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password to authenticate with")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug log level")
	rootCmd.PersistentFlags().BoolVar(&persist, "persist", false, "Persist changes such as colour and such, can wear the flash out")

	initConfigCmd()
	initcolourCmd()

	rootCmd.AddCommand(VersionCmd)
	rootCmd.AddCommand(StatusCmd)
	rootCmd.AddCommand(RebootCmd)
	rootCmd.AddCommand(ConfigCmd)
	rootCmd.AddCommand(ColourCmd)

	rootCmd.MarkFlagRequired("address")
}
