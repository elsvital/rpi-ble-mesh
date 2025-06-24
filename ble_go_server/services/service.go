// service.go
package services

import (
	"ble_go_server/configs"
	"github.com/godbus/dbus/v5"
)

type GATTService struct {
	Path string
}

func (s *GATTService) GetProperties() map[string]map[string]dbus.Variant {
	return map[string]map[string]dbus.Variant{
		"org.bluez.GattService1": {
			"UUID":    dbus.MakeVariant(configs.SERVICE_UUID),
			"Primary": dbus.MakeVariant(true),
			"Characteristics": dbus.MakeVariant([]dbus.ObjectPath{
				dbus.ObjectPath(configs.CHARACTERISTIC_PATH_MACHINE_DATA),
				dbus.ObjectPath(configs.CHARACTERISTIC_PATH_FIRMWARE_VERSION),
				dbus.ObjectPath(configs.CHARACTERISTIC_PATH_FIRMWARE_DATA),
			}),
		},
	}
}
