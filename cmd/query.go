/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
}

func queryStats(cmd *cobra.Command, args []string) {
	statsResult, err := stats.QueryAllStats(host, port)

	if err != nil {
		config.Error.Fatal(err)
	}

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

	fmt.Printf("%20s%21s%21s\n", "User", "Downlink", "Uplink")
	for key, value := range *statsResult {
		if value.Type != "user" {
			continue
		}
		fmt.Printf("%20[2]s%20[3]d%[1]s%20[4]d%[1]s\n", unitSign, key, value.Downlink/int64(unitBase), value.Uplink/int64(unitBase))
	}
}
