# Sleuz (Simple Bluez) CLI for Linux

`sleuz` is a simple bluez dbus client for linux. It provides a command line interface
for interacting with bluez over the dbus protocol on linux, currently only the basic
commands and functionality are implemented.

This is a replacement for `bluetoothctl` on linux. `bluetoothctl` is very handy, but
the fact that it isn't a command line interface makes it hard to use in scripts. `sleuz`
is usable in scripts as all the basic interactions are invokable via the command line.

Currently none of the authentication capabilities are implemented because I don't
have a bluetooth device that requires authentication.

## Installation
### Using Go

```
$ go get -u github.com/vishen/sluez
```

### Releases
```
# TODO: add releases for linux?
```

## Examples

```
# Get a current overview of your adapters and bluetooth devices
$ sluez status
Adapters:
1) name="jonathan-Blade" alias="jonathan-Blade" address="9C:B6:D0:1C:BB:B0" discoverable=true pairable=true powered=true discovering=false
Connected devices:
1) name="Pixel 2" alias="Pixel 2" address="40:4E:36:9F:1E:EC" adapter="/org/bluez/hci0" paired=true connected=false trusted=false blocked=false
2) name="Bose QC35 II" alias="Bose QC35 II" address="2C:41:A1:49:37:CF" adapter="/org/bluez/hci0" paired=true connected=false trusted=false blocked=fals

# Pair bluetooth devices, you will need to put you device into pairing mode,
# if --device or --device-name are not specified, sluez will pair with the first
# device it finds
$ sluez pairi --debug
[sluez] trying to pair bluetooth devices to "hci0"found no devices similar to specified device= or device-name=
waiting for new bluetooth devices, make sure to put device into pairing mode
[sluez] received signal=org.freedesktop.DBus.ObjectManager.InterfacesAdded => (2)[/org/bluez/hci0/dev_2C_41_A1_49_37_CF map[org.freedesktop.DBus.Introspectable:map[] org.bluez.Device1:map[RSSI:@n -30 Modalias:"bluetooth:v009Ep4020d0251" Icon:"audio-card" Alias:"Bose QC35 II" Trusted:false LegacyPairing:false UUIDs:["0000110d-0000-1000-8000-00805f9b34fb", "0000110b-0000-1000-8000-00805f9b34fb", "0000110a-0000-1000-8000-00805f9b34fb", "0000110e-0000-1000-8000-00805f9b34fb", "0000110f-0000-1000-8000-00805f9b34fb", "00001130-0000-1000-8000-00805f9b34fb", "0000112e-0000-1000-8000-00805f9b34fb", "0000111e-0000-1000-8000-00805f9b34fb", "00001108-0000-1000-8000-00805f9b34fb", "00001131-0000-1000-8000-00805f9b34fb", "00000000-deca-fade-deca-deafdecacaff"] Adapter:@o "/org/bluez/hci0" Address:"2C:41:A1:49:37:CF" Blocked:false Connected:false Name:"Bose QC35 II" Paired:false Class:@u 2360344] org.freedesktop.DBus.Properties:map[]]]
[sluez] trying to pair with device mac "2C:41:A1:49:37:CF"
successfully paired "2C:41:A1:49:37:CF" and "hci0"

# If no device is specified, the first device found will be connected. If
# there is more than one device, you will be asked to select the device.
$ ./sluez connect
Choose a bluetooth device from the following:
1) Pixel 2, 40:4E:36:9F:1E:EC
2) Bose QC35 II, 2C:41:A1:49:37:CF
>> 2
successfully connected "2C:41:A1:49:37:CF" and "hci0"

# Connect a bluetooth device by part of its device name, this will do a simple
# fuzzy search on known device names (device must be paired). It will find
# devices thate look like "bose", ie: "Bose QC35 II".
$ sluez connect --device-name=bose

# Or connect to your phone
$ sluez connect --device-name=pixel

# Disconnect a bluetooth device by its MAC
$ sluez disconnect --device=AA:BB:CC:11:22:33
successfully disconnected "2C:41:A1:49:37:CF" and "hci0

# Discover bluetooth devices as the become pairable or when they disconnect.
# This will watch for new events about bluetooth devices.
$ sluez discover
```

## Usage

```
Simple CLI for Bluez dBus on linux

Usage:
  sluez [command]

Available Commands:
  connect     Connect a device to an adapter
  disconnect  Disconnect a device from an adapter
  discover    Discover will watch for devices as the connect or disconnect to an adapter
  help        Help about any command
  pair        Pair a device from to an adapter, requires your device to be in pairing mode
  remove      Remove a device from an adapter
  status      The current status of known adapters and devices

Flags:
  -a, --adapter string       HCI device adapter. Can be found from 'hciconfig -a' (default "hci0")
      --debug                Print debug logs
  -d, --device string        Bluetooth device MAC address
  -n, --device-name string   Bluetooth device name. A fuzzy search is used to determine which device the name matches for. '--device' will take precedence if both are specified
  -h, --help                 help for sluez

Use "sluez [command] --help" for more information about a command.
```

## TODO

- Add command to be able to set device properties
- Add command "auto" that will try to pair the device, then connect the device
