// Copyright © 2021 ezbuy & LITB TEAM
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"log"
	"time"

	"github.com/ezbuy/chartprobe/internal/client"
	"github.com/spf13/cobra"
)

var isDeleteAll bool
var version string
var prefix string
var period string
var perChartMaxVersion uint32

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete charts from chartmuseum",
	Run: func(_ *cobra.Command, args []string) {

		var charts []*client.Chart
		for _, arg := range args {
			chart, err := client.ParseChart(arg)
			if err != nil {
				log.Printf("WARNING: %q", err)
				continue
			}
			charts = append(charts, chart)
		}

		c := client.NewClient()
		var opts []client.DeleteOption
		if isDeleteAll {
			opts = append(opts, client.WithPurgeOption())
		}
		if prefix != "" {
			opts = append(opts, client.WithPrefix(prefix))
		}
		if period != "" {
			if period[0] != '-' {
				log.Fatalf("period must start with '-'")
			}
			du, err := time.ParseDuration(period)
			if err != nil {
				log.Fatalf("delete: parse period: %q", err)
			}
			opts = append(opts, client.WithPeriod(du))
		}
		if perChartMaxVersion > 0 {
			opts = append(opts, client.WithPerChartVersion(perChartMaxVersion))
		}
		dc, err := c.Del(context.Background(), charts, opts...)
		if err != nil {
			log.Fatalf("delete: %q", err)
		}
		log.Printf("chartprobe: deleted %d charts", dc)
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)

	deleteCmd.PersistentFlags().BoolVarP(&isDeleteAll, "all", "a", false, "Delete All Charts")
	deleteCmd.PersistentFlags().StringVar(&version, "chart_version", "", "Specified Chart Version")
	deleteCmd.PersistentFlags().StringVar(&prefix, "prefix", "", "Chart prefix")
	deleteCmd.PersistentFlags().StringVar(&period, "period", "", "period defines the purge period based on the chart create date, and the value should follow Go's time.ParseDuration(https://golang.org/pkg/time/#ParseDuration)")
	deleteCmd.PersistentFlags().Uint32Var(&perChartMaxVersion, "per-chart-versions", 0, "per-chart-versions defines the max version existed per chart")
}
