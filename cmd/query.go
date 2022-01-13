package cmd

import (
	"fmt"
	"strings"

	"github.com/lureiny/v2raymg/config"
	"github.com/lureiny/v2raymg/stats"
	"github.com/spf13/cobra"
)

const (
	k = 1024
	m = 1024 * 1024
	g = 1024 * 1024 * 1024
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query user's stats",
	Run:   queryStats,
}

func init() {
	queryCmd.Flags().StringVarP(&unit, "unit", "u", "m", "Unit of stats. K/k or M/m or G/g")
	queryCmd.Flags().StringVarP(&email, "email", "e", "", "user name/email which to query")
}

func queryStats(cmd *cobra.Command, args []string) {
	unitBase := m
	unitSign := "M"

	switch strings.ToLower(unit) {
	case "k":
		unitBase = k
		unitSign = "K"
	case "g":
		unitBase = g
		unitSign = "G"
	}

	var statsResult *map[string]*stats.MyStat
	var err error
	switch email {
	case "":
		statsResult, err = stats.QueryAllStats(host, port)
		if err != nil {
			config.Error.Fatal(err)
		}
	default:
		statsResult, err = stats.QueryUserStat(host, port, email)
		if err != nil {
			config.Error.Fatal(err)
		}
		if len(*statsResult) == 0 {
			config.Info.Printf("No user: %s", email)
		}
	}

	if len(*statsResult) > 0 {
		fmt.Printf("%20s%21s%21s\n", "User", "Downlink", "Uplink")
		for key, value := range *statsResult {
			if value.Type != "user" {
				continue
			}
			fmt.Printf("%20[2]s%20[3]d%[1]s%20[4]d%[1]s\n", unitSign, key, value.Downlink/int64(unitBase), value.Uplink/int64(unitBase))
		}
	}
}
