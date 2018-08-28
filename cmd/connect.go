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
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect a paired bluetooth device to an adapter",
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
		debug("connecting to adapter=%s device=%s", adapter, device)
		if err := b.Connect(adapter, device); err != nil {
			fmt.Printf("unable to connect to device %q: %v\n", device, err)
			return nil
		}
		fmt.Printf("successfully connected %q and %q\n", device, adapter)

		// NOTE: Need to manually set the card profile for pulseaudio, this _should_
		// happen already, but for some reason it doesn't always happen. This tends
		// to happen when the computer has been idle for a while.
		// pactl set-card-profile bluez_card.2C_41_A1_49_37_CF a2dp_sink
		// TODO(vishen): pulseaudio can be controlled via dbus, look at using that
		// instead
		exec.Command("pactl", "set-card-profile", fmt.Sprintf("bluez_card.%s", strings.Replace(device, ":", "_", -1)), "a2dp_sink").Run()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
