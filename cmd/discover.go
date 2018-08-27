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

	"github.com/godbus/dbus"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover will watch for devices as the connect or disconnect to an adapter",
	RunE: func(cmd *cobra.Command, args []string) error {
		adapter, _ := cmd.Flags().GetString("adapter")
		if adapter == "" {
			return errors.New("--adapter is required")
		}
		b, err := newBluez(cmd)
		if err != nil {
			return err
		}
		if err := b.StartDiscovery(adapter); err != nil {
			return errors.Wrap(err, "unable to start discovery")
		}

		fmt.Println("Adapters:")
		for i, a := range b.Adapters {
			fmt.Printf("%d) name=%q alias=%q address=%q discoverable=%t pairable=%t powered=%t discovering=%t\n", i+1, a.Name, a.Alias, a.Address, a.Discoverable, a.Pairable, a.Powered, a.Discovering)
		}
		fmt.Println("Connected devices:")
		for i, d := range b.Devices {
			fmt.Printf("%d) name=%q alias=%q address=%q, adapter=%q paired=%t connected=%t trusted=%t blocked=%t\n", i+1, d.Name, d.Alias, d.Address, d.Adapter, d.Paired, d.Connected, d.Trusted, d.Blocked)
		}

		fmt.Printf("watching for new bluetooth events, make sure to put device into pairing mode\n")
		signalChan := b.WatchSignal()
		for signal := range signalChan {
			debug("received signal=%s => (%d)%v\n", signal.Name, len(signal.Body), signal.Body)
			if signal.Name == "org.freedesktop.DBus.ObjectManager.InterfacesAdded" {
				if len(signal.Body) != 2 {
					continue
				}
				devicePath, ok := signal.Body[0].(dbus.ObjectPath)
				if !ok {
					debug("unable to cast '%#v' to dbus.ObjectPath", signal.Body[0])
					continue
				}
				deviceMap, ok := signal.Body[1].(map[string]map[string]dbus.Variant)
				if !ok {
					debug("unable to cast '%#v' to device map[string]dbus.Variant", signal.Body[1])
					continue
				}
				devices := b.ConvertToDevices(string(devicePath), deviceMap)
				for _, d := range devices {
					fmt.Printf("name=%q alias=%q address=%q, adapter=%q paired=%t connected=%t trusted=%t blocked=%t\n", d.Name, d.Alias, d.Address, d.Adapter, d.Paired, d.Connected, d.Trusted, d.Blocked)
				}
			}
		}
		return nil

	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}
