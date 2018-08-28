// Copyright © 2018 Jonathan Pentecost <pentecostjonathan@gmail.com>
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
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "The current status of known adapters and devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := newBluez(cmd)
		if err != nil {
			fmt.Printf("unable to get bluez client: %v\n", err)
			return nil
		}
		fmt.Println("Adapters:")
		for i, a := range b.Adapters {
			// TODO(vishen): add these to methods
			fmt.Printf("%d) name=%q alias=%q address=%q discoverable=%t pairable=%t powered=%t discovering=%t\n", i+1, a.Name, a.Alias, a.Address, a.Discoverable, a.Pairable, a.Powered, a.Discovering)
		}
		fmt.Println("Connected devices:")
		for i, d := range b.Devices {
			// TODO(vishen): add these to methods
			fmt.Printf("%d) name=%q alias=%q address=%q adapter=%q paired=%t connected=%t trusted=%t blocked=%t\n", i+1, d.Name, d.Alias, d.Address, d.Adapter, d.Paired, d.Connected, d.Trusted, d.Blocked)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
