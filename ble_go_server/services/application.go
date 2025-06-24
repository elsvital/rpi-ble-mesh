// application.go
package services

import (
	"ble_go_server/configs"
	"github.com/godbus/dbus/v5"
)

type Application struct{}

func (a *Application) GetManagedObjects() (map[dbus.ObjectPath]map[string]map[string]dbus.Variant, *dbus.Error) {
	objects := make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)

	service := &GATTService{Path: configs.SERVICE_PATH}
	charMachine := &MachineDataCharacteristic{Path: configs.CHARACTERISTIC_PATH_MACHINE_DATA}
	charFirmwareVer := &FirmwareVersionCharacteristic{Path: configs.CHARACTERISTIC_PATH_FIRMWARE_VERSION}
	charFirmwareData := &FirmwareDataCharacteristic{Path: configs.CHARACTERISTIC_PATH_FIRMWARE_DATA}

	objects[dbus.ObjectPath(configs.SERVICE_PATH)] = service.GetProperties()
	objects[dbus.ObjectPath(configs.CHARACTERISTIC_PATH_MACHINE_DATA)] = charMachine.GetProperties()
	objects[dbus.ObjectPath(configs.CHARACTERISTIC_PATH_FIRMWARE_VERSION)] = charFirmwareVer.GetProperties()
	objects[dbus.ObjectPath(configs.CHARACTERISTIC_PATH_FIRMWARE_DATA)] = charFirmwareData.GetProperties()

	return objects, nil
}
