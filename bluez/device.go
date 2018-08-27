package bluez

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
)

const (
	dbusBluetoothPath        = "org.bluez"
	dbusPropertiesGetAllPath = "org.freedesktop.DBus.Properties.GetAll"
	dbusIntrospectPath       = "org.freedesktop.DBus.Introspectable.Introspect"
	dbusListNamesPath        = "org.freedesktop.DBus.ListNames"
	dbusObjectManagerPath    = "org.freedesktop.DBus.ObjectManager.GetManagedObjects"
)

type Adapter struct {
	Path         string
	Name         string
	Alias        string
	Address      string
	Discoverable bool
	Pairable     bool
	Powered      bool
	Discovering  bool
}

type Device struct {
	Path      string
	Name      string
	Alias     string
	Address   string
	Adapter   string
	Paired    bool
	Connected bool
	Trusted   bool
	Blocked   bool
}

type Bluez struct {
	conn *dbus.Conn

	adapterDefaultIndex int
	deviceDefaultIndex  int

	Adapters []Adapter
	Devices  []Device
}

func NewBluez(conn *dbus.Conn) *Bluez {
	return &Bluez{conn: conn}
}

func (b *Bluez) ConvertToDevices(path string, values map[string]map[string]dbus.Variant) []Device {
	/*
		org.bluez.Device1
			Icon => dbus.Variant{sig:dbus.Signature{str:"s"}, value:"audio-card"}
			LegacyPairing => dbus.Variant{sig:dbus.Signature{str:"b"}, value:false}
			Address => dbus.Variant{sig:dbus.Signature{str:"s"}, value:"2C:41:A1:49:37:CF"}
			Trusted => dbus.Variant{sig:dbus.Signature{str:"b"}, value:false}
			Connected => dbus.Variant{sig:dbus.Signature{str:"b"}, value:true}
			Paired => dbus.Variant{sig:dbus.Signature{str:"b"}, value:true}
			RSSI => dbus.Variant{sig:dbus.Signature{str:"n"}, value:-36}
			Modalias => dbus.Variant{sig:dbus.Signature{str:"s"}, value:"bluetooth:v009Ep4020d0251"}
			Name => dbus.Variant{sig:dbus.Signature{str:"s"}, value:"Bose QC35 II"}
			UUIDs => dbus.Variant{sig:dbus.Signature{str:"as"}, value:[]string{"00000000-deca-fade-deca-deafdecacaff", "00001101-0000-1000-8000-00805f9b34fb", "00001108-0000-1000-8000-00805f9b34fb", "0000110b-0000-1000-8000-00805f9b34fb", "0000110c-0000-1000-8000-00805f9b34fb", "0000110e-0000-1000-8000-00805f9b34fb", "0000111e-0000-1000-8000-00805f9b34fb", "00001200-0000-1000-8000-00805f9b34fb", "81c2e72a-0591-443e-a1ff-05f988593351", "f8d1fbe4-7966-4334-8024-ff96c9330e15"}}
			Adapter => dbus.Variant{sig:dbus.Signature{str:"o"}, value:"/org/bluez/hci0"}
			Blocked => dbus.Variant{sig:dbus.Signature{str:"b"}, value:false}
			Alias => dbus.Variant{sig:dbus.Signature{str:"s"}, value:"Bose QC35 II"}
			Class => dbus.Variant{sig:dbus.Signature{str:"u"}, value:0x240418}

	*/
	devices := []Device{}
	for k, v := range values {
		switch k {
		case "org.bluez.Device1":
			adapter, _ := v["Adapter"].Value().(dbus.ObjectPath)
			devices = append(devices, Device{
				Path:      path,
				Name:      v["Name"].Value().(string),
				Alias:     v["Alias"].Value().(string),
				Address:   v["Address"].Value().(string),
				Adapter:   string(adapter),
				Paired:    v["Paired"].Value().(bool),
				Connected: v["Connected"].Value().(bool),
				Trusted:   v["Trusted"].Value().(bool),
				Blocked:   v["Blocked"].Value().(bool),
			})
		}
	}
	return devices
}

// TODO(vishen): Better name?
func (b *Bluez) PopulateCache() error {
	results, err := b.ManagedObjects()
	if err != nil {
		return err
	}
	devices := []Device{}
	adapters := []Adapter{}
	for k, v := range results {
		devices = append(devices, b.ConvertToDevices(string(k), v)...)
		for k1, v1 := range v {
			switch k1 {
			case "org.bluez.Adapter1":
				/*
					/org/bluez/hci0
						org.bluez.Adapter1
								Discoverable => dbus.Variant{sig:dbus.Signature{str:"b"}, value:true}
								UUIDs => dbus.Variant{sig:dbus.Signature{str:"as"}, value:[]string{"00001112-0000-1000-8000-00805f9b34fb", "00001801-0000-1000-8000-00805f9b34fb", "0000110e-0000-1000-8000-00805f9b34fb", "00001800-0000-1000-8000-00805f9b34fb", "00001200-0000-1000-8000-00805f9b34fb", "0000110c-0000-1000-8000-00805f9b34fb", "0000110b-0000-1000-8000-00805f9b34fb", "0000110a-0000-1000-8000-00805f9b34fb"}}
								Modalias => dbus.Variant{sig:dbus.Signature{str:"s"}, value:"usb:v1D6Bp0246d0525"}
								Pairable => dbus.Variant{sig:dbus.Signature{str:"b"}, value:true}
								DiscoverableTimeout => dbus.Variant{sig:dbus.Signature{str:"u"}, value:0x0}
								PairableTimeout => dbus.Variant{sig:dbus.Signature{str:"u"}, value:0x0}
								Powered => dbus.Variant{sig:dbus.Signature{str:"b"}, value:true}
								Class => dbus.Variant{sig:dbus.Signature{str:"u"}, value:0xc010c}
								Discovering => dbus.Variant{sig:dbus.Signature{str:"b"}, value:true}
								Address => dbus.Variant{sig:dbus.Signature{str:"s"}, value:"9C:B6:D0:1C:BB:B0"}
								Name => dbus.Variant{sig:dbus.Signature{str:"s"}, value:"jonathan-Blade"}
								Alias => dbus.Variant{sig:dbus.Signature{str:"s"}, value:"jonathan-Blade"}

				*/
				// TODO(vishen): do the same convert to adapaters, as done for devices
				adapters = append(adapters, Adapter{
					Path:         string(k),
					Name:         v1["Name"].Value().(string),
					Alias:        v1["Alias"].Value().(string),
					Address:      v1["Address"].Value().(string),
					Discoverable: v1["Discoverable"].Value().(bool),
					Pairable:     v1["Pairable"].Value().(bool),
					Powered:      v1["Powered"].Value().(bool),
					Discovering:  v1["Discovering"].Value().(bool),
				})
			}
		}
	}

	b.Adapters = adapters
	b.Devices = devices

	return nil
}

func (b *Bluez) ManagedObjects() (map[dbus.ObjectPath]map[string]map[string]dbus.Variant, error) {
	result := make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)
	if err := b.conn.Object(dbusBluetoothPath, "/").Call(dbusObjectManagerPath, 0).Store(&result); err != nil {
		return result, err
	}
	return result, nil
}

func (b *Bluez) CallAdapter(adapterName, method string, flags dbus.Flags, args ...interface{}) *dbus.Call {
	path := "/org/bluez/" + adapterName
	return b.conn.Object(dbusBluetoothPath, dbus.ObjectPath(path)).Call("org.bluez.Adapter1."+method, flags, args...)
}

func (b *Bluez) StartDiscovery(adapterName string) error {
	if err := b.CallAdapter(adapterName, "StartDiscovery", 0).Store(); err != nil {
		return err
	}
	return nil
}

func (b *Bluez) RemoveDevice(adapterName, deviceMac string) error {
	devicePath := b.devicePath(adapterName, deviceMac)
	if err := b.CallAdapter(adapterName, "RemoveDevice", 0, devicePath).Store(); err != nil {
		return err
	}
	return nil
}

func (b *Bluez) WatchSignal() chan *dbus.Signal {
	signalMatch := "type='signal',interface='org.freedesktop.DBus.ObjectManager',path='/'"
	b.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, signalMatch)
	ch := make(chan *dbus.Signal, 1)
	b.conn.Signal(ch)
	return ch
}

func (b *Bluez) devicePath(adapterName, deviceMac string) dbus.ObjectPath {
	path := fmt.Sprintf(
		"/org/bluez/%s/dev_%s",
		adapterName,
		strings.Replace(deviceMac, ":", "_", -1),
	)
	return dbus.ObjectPath(path)
}

// https://git.kernel.org/pub/scm/bluetooth/bluez.git/tree/doc/device-api.txt
func (b *Bluez) CallDevice(adapterName, deviceMac, method string, flags dbus.Flags, args ...interface{}) *dbus.Call {
	path := b.devicePath(adapterName, deviceMac)
	return b.conn.Object(dbusBluetoothPath, path).Call("org.bluez.Device1."+method, flags, args...)
}

func (b *Bluez) Pair(adapterName, deviceMac string) error {
	return b.CallDevice(adapterName, deviceMac, "Pair", 0).Store()
}

func (b *Bluez) Connect(adapterName, deviceMac string) error {
	return b.CallDevice(adapterName, deviceMac, "Connect", 0).Store()
}

func (b *Bluez) Disconnect(adapterName, deviceMac string) error {
	return b.CallDevice(adapterName, deviceMac, "Disconnect", 0).Store()
}

func (b *Bluez) GetDeviceProperties(adapterName, deviceMac string) (map[string]dbus.Variant, error) {
	result := make(map[string]dbus.Variant)
	path := b.devicePath(adapterName, deviceMac)
	// TODO(vishen): factor this with the CallDevice functionality
	if err := b.conn.Object(dbusBluetoothPath, path).Call("org.freedesktop.DBus.Properties.GetAll", 0, "org.bluez.Device1").Store(&result); err != nil {
		return result, err
	}
	return result, nil
}

func (b *Bluez) SetDeviceProperty(adapterName, deviceMac string, key string, value interface{}) error {
	path := b.devicePath(adapterName, deviceMac)
	// TODO(vishen): factor this with the CallDevice functionality
	return b.conn.Object(dbusBluetoothPath, path).Call("org.freedesktop.DBus.Properties.Set", 0, "org.bluez.Device1", key, dbus.MakeVariant(value)).Store()
}
