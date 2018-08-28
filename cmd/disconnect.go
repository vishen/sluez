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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// disconnectCmd represents the disconnect command
var disconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnect a device from an adapter",
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := newBluez(cmd)
		if err != nil {
			fmt.Printf("unable to get bluez client: %v\n", err)
			return nil
		}
		device, adapter, err := deviceAndAdapter(b, cmd)
		if err != nil {
			return errors.Wrap(err, "unable to determine device and/or adapter")
		}
		debug("disconnecting to adapter=%s device=%s", adapter, device)
		if err := b.Disconnect(adapter, device); err != nil {
			fmt.Printf("unable to disconnect to device %q: %v", device, err)
			return nil
		}
		fmt.Printf("successfully disconnected %q and %q\n", device, adapter)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(disconnectCmd)
}
