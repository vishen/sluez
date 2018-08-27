package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/godbus/dbus"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/vishen/sluez/bluez"
)

var (
	debugging = false
)

func newBluez(cmd *cobra.Command) (*bluez.Bluez, error) {
	debugging, _ = cmd.Flags().GetBool("debug")
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create dbus system bus:")
	}
	b := bluez.NewBluez(conn)
	if err := b.PopulateCache(); err != nil {
		return nil, errors.Wrapf(err, "unable to populate cache")
	}
	return b, nil
}

func deviceAndAdapter(b *bluez.Bluez, cmd *cobra.Command) (device string, adapter string, err error) {
	adapter, _ = cmd.Flags().GetString("adapter")
	if adapter == "" {
		return "", "", errors.New("--adapter is required")
	}
	device, _ = cmd.Flags().GetString("device")
	deviceName, _ := cmd.Flags().GetString("device-name")

	// If no device is specified we will try to grab one from the
	// cached/known devices.
	if device == "" || deviceName != "" {
		debug("no bluetooth mac specified in flags")
		switch len(b.Devices) {
		case 0:
			debug("no bluetooth devices found")
			return "", "", errors.New("no bluetooth devices found, please specify a --device or --device-name")
		default:
			// If a device name was specified, we should check all the connected devices
			// and if one of them has a similar name to the one specified, use that.
			if deviceName != "" {
				for _, d := range b.Devices {
					if similar(deviceName, d.Name) {
						device = d.Address
						// debug("device name matches %q, using %q", d.Name, device)
						break
					}
				}
			}

			// If we have found a device from the above searching.
			if device != "" {
				break
			}

			// Ask the user to choose a bluetooth device from the connected devices.
			for {
				fmt.Printf("Choose a bluetooth device from the following:\n")
				for i, d := range b.Devices {
					fmt.Printf("%d) %s, %s\n", i+1, d.Name, d.Address)
				}
				fmt.Printf(">> ")
				reader := bufio.NewReader(os.Stdin)
				text, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("unable to read from stdin: %s\n", err)
					continue
				}
				text = strings.TrimSpace(text)
				i, err := strconv.Atoi(text)
				if err != nil || i < 1 || i > len(b.Devices) {
					fmt.Printf("'%s' is an invalid choice, please select the number for the device you want to connect\n", text)
					continue
				}
				device = b.Devices[i-1].Address
				break
			}

		}
	}
	return device, adapter, nil

}

func similar(match, similarTo string) bool {
	r := strings.NewReplacer(" ", "", "_", "", "-", "", "/", "")
	match = strings.ToLower(r.Replace(match))
	similarTo = strings.ToLower(r.Replace(similarTo))

	if strings.HasPrefix(similarTo, match) {
		return true
	}
	if strings.HasSuffix(similarTo, match) {
		return true
	}
	if strings.Contains(similarTo, match) {
		return true
	}
	return false
}

func debug(message string, args ...interface{}) {
	if debugging {
		fmt.Printf("[sluez] %s", fmt.Sprintf(message, args...))
	}
}
