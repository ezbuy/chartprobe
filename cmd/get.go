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
	"encoding/json"
	"fmt"
	"log"

	"github.com/ezbuy/chartprobe/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var isGetAll bool

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Charts from museum",
	Run: func(_ *cobra.Command, _ []string) {
		c := client.NewClient()
		cs, err := c.GetAll(context.Background())
		if err != nil {
			log.Fatalf("cmd/get: %q", err)
		}

		data, err := json.MarshalIndent(cs, "", "	")
		if err != nil {
			log.Fatalf("cmd/get: %q", err)
		}
		fmt.Println(string(data))
	},
}

func init() {
	MustSetEnv()
	RootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVarP(&isGetAll, "all", "a", viper.GetBool("all"), "get all charts")
}
