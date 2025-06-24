// advertisement.go
package services

import (
	"fmt"

	"ble_go_server/configs"
	"github.com/godbus/dbus/v5"
)

type Advertisement struct {
	Path dbus.ObjectPath
}

func (a *Advertisement) GetProperties() map[string]map[string]dbus.Variant {
	return map[string]map[string]dbus.Variant{
		"org.bluez.LEAdvertisement1": {
			"Type":           dbus.MakeVariant("peripheral"),
			"ServiceUUIDs":   dbus.MakeVariant([]string{configs.SERVICE_UUID}),
			"LocalName":      dbus.MakeVariant("RPi-BLE"),
			"IncludeTxPower": dbus.MakeVariant(true),
		},
	}
}

func (a *Advertisement) GetAll(interfaceName string) (map[string]dbus.Variant, *dbus.Error) {
	return a.GetProperties()[interfaceName], nil
}

func (a *Advertisement) Release() *dbus.Error {
	fmt.Println("🔕 Anúncio liberado")
	return nil
}
