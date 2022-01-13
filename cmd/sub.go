package cmd

import (
	"github.com/lureiny/v2raymg/config"
	"github.com/lureiny/v2raymg/sub"
	"github.com/spf13/cobra"
)

// subCmd represents the sub command
var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "Get some user's sub uri",
	Run:   getSubURI,
}

func init() {
	subCmd.Flags().StringVarP(&email, "email", "e", "", "Email of user.")
	subCmd.MarkFlagRequired("email")
	subCmd.Flags().StringVarP(&configFile, "config", "c", "/usr/local/etc/v2ray/config.json", "The config file of v2ray.")
}

func getSubURI(cmd *cobra.Command, args []string) {
	uri, err := sub.GetUserSubUri(host, email, uint32(port), configFile)
	if err != nil {
		config.Error.Fatal(err)
	}
	config.Info.Println(uri)
}
