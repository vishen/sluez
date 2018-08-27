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

# Pair bluetooth devices, you will need to put you device into pairing mode,
# if --device or --device-name are not specified, sluez will pair with the first
# device it finds
$ sluez pair

# Connect a bluetooth device by part of its device name, this will do a simple
# fuzzy search on known device names (device must be paired). It will find
# devices thate look like "bose", ie: "Bose QC35 II".
$ sluez connect --device-name=bose

# Or connect to your phone
$ sluez connect --device-name=pixel

# Disconnect a bluetooth device by its MAC
$ sluez disconnect --device=AA:BB:CC:11:22:33

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
