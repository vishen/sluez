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

// pairCmd represents the pair command
var pairCmd = &cobra.Command{
	Use:   "pair",
	Short: "Pair a device to an adapter, requires your device to be in pairing mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		adapter, _ := cmd.Flags().GetString("adapter")
		if adapter == "" {
			return errors.New("--adapter is required")
		}
		device, _ := cmd.Flags().GetString("device")
		deviceName, _ := cmd.Flags().GetString("device-name")
		b, err := newBluez(cmd)
		if err != nil {
			return err
		}

		// "pair" is different from the rest of the commands as the device
		// isn't currently connected to bluez. We need to; start discovery, watch for any
		// new devices that are discovered, and then attempt to pair with that device.

		// Check that the device isn't already paired
		for _, d := range b.Devices {
			if device != "" && d.Address == device || deviceName != "" && similar(deviceName, d.Name) {
				fmt.Printf("device %q is already paired\n", d.Name)
				return nil
			}
		}

		debug("trying to pair bluetooth devices to %q", adapter)
		if err := b.StartDiscovery(adapter); err != nil {
			return errors.Wrap(err, "unable to start discovery")
		}
		fmt.Printf("found no devices similar to specified device=%s or device-name=%s\n", device, deviceName)
		fmt.Printf("waiting for new bluetooth devices, make sure to put device into pairing mode\n")

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
					// If no device mac is set, attempt to pair  to the first device found, otherwise
					// if the device mac is set and is the same as the found device mac, then we try
					// and pair that device.
					if (device == "" && deviceName == "") || (device != "" && d.Address == device) || (deviceName != "" && similar(d.Name, deviceName)) {
						device = d.Address
						debug("trying to pair with device mac %q", device)
						if err := b.Pair(adapter, device); err != nil {
							return errors.Wrapf(err, "unable to pair with device %q", device)
						}
						fmt.Printf("successfully paired %q and %q\n", device, adapter)
						return nil
					}
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pairCmd)
}
