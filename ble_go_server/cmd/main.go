// main.go (final com introspecção de características)
package main

import (
	"fmt"
	"log"

	"ble_go_server/configs"
	"ble_go_server/services"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

func main() {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatalf("Erro ao conectar ao system bus: %v", err)
	}

	//services.InitKafkaProducer()

	app := &services.Application{}
	service := &services.GATTService{Path: configs.SERVICE_PATH}
	advert := &services.Advertisement{Path: "/org/bluez/example/advertisement0"}

	charMachine := &services.MachineDataCharacteristic{Path: configs.CHARACTERISTIC_PATH_MACHINE_DATA}
	charFirmwareVer := &services.FirmwareVersionCharacteristic{Path: configs.CHARACTERISTIC_PATH_FIRMWARE_VERSION}
	charFirmwareData := &services.FirmwareDataCharacteristic{Path: configs.CHARACTERISTIC_PATH_FIRMWARE_DATA}

	// Exporta objetos
	_ = conn.Export(app, dbus.ObjectPath(configs.APP_PATH), "org.freedesktop.DBus.ObjectManager")
	_ = conn.Export(service, dbus.ObjectPath(configs.SERVICE_PATH), "org.bluez.GattService1")
	_ = conn.Export(charMachine, dbus.ObjectPath(configs.CHARACTERISTIC_PATH_MACHINE_DATA), "org.bluez.GattCharacteristic1")
	_ = conn.Export(charFirmwareVer, dbus.ObjectPath(configs.CHARACTERISTIC_PATH_FIRMWARE_VERSION), "org.bluez.GattCharacteristic1")
	_ = conn.Export(charFirmwareData, dbus.ObjectPath(configs.CHARACTERISTIC_PATH_FIRMWARE_DATA), "org.bluez.GattCharacteristic1")
	_ = conn.Export(advert, advert.Path, "org.bluez.LEAdvertisement1")
	_ = conn.Export(advert, advert.Path, "org.freedesktop.DBus.Properties")

	// Introspecção principal
	_ = conn.Export(introspect.NewIntrospectable(&introspect.Node{
		Name: configs.APP_PATH,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{Name: "org.freedesktop.DBus.ObjectManager"},
			{Name: "org.bluez.GattService1"},
			{Name: "org.bluez.GattCharacteristic1"},
		},
	}), dbus.ObjectPath(configs.APP_PATH), "org.freedesktop.DBus.Introspectable")

	// Introspecção para characteristics
	_ = conn.Export(introspect.NewIntrospectable(&introspect.Node{
		Name: string(configs.CHARACTERISTIC_PATH_MACHINE_DATA),
		Interfaces: []introspect.Interface{
			{
				Name: "org.bluez.GattCharacteristic1",
				Methods: []introspect.Method{
					{
						Name: "WriteValue",
						Args: []introspect.Arg{
							{Name: "value", Direction: "in", Type: "ay"},
							{Name: "options", Direction: "in", Type: "a{sv}"},
						},
					},
				},
				Properties: []introspect.Property{
					{Name: "UUID", Type: "s", Access: "read"},
					{Name: "Service", Type: "o", Access: "read"},
					{Name: "Flags", Type: "as", Access: "read"},
				},
			},
		},
	}), dbus.ObjectPath(configs.CHARACTERISTIC_PATH_MACHINE_DATA), "org.freedesktop.DBus.Introspectable")

	// Introspecção advertisement
	_ = conn.Export(introspect.NewIntrospectable(&introspect.Node{
		Name: string(advert.Path),
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{
				Name: "org.bluez.LEAdvertisement1",
				Methods: []introspect.Method{
					{Name: "Release"},
				},
				Properties: []introspect.Property{
					{Name: "Type", Type: "s", Access: "read"},
					{Name: "ServiceUUIDs", Type: "as", Access: "read"},
					{Name: "LocalName", Type: "s", Access: "read"},
					{Name: "IncludeTxPower", Type: "b", Access: "read"},
				},
			},
		},
	}), advert.Path, "org.freedesktop.DBus.Introspectable")

	// Registra advertisement BLE
	adMgr := conn.Object("org.bluez", "/org/bluez/hci0")
	adCall := adMgr.Call("org.bluez.LEAdvertisingManager1.RegisterAdvertisement", 0,
		advert.Path, map[string]dbus.Variant{})
	if adCall.Err != nil {
		log.Fatalf("Erro ao registrar advertisement: %v", adCall.Err)
	}
	fmt.Println("📡 Advertisement BLE registrado com sucesso.")

	// Registra aplicação GATT
	appObj := conn.Object("org.bluez", "/org/bluez/hci0")
	gattCall := appObj.Call("org.bluez.GattManager1.RegisterApplication", 0,
		dbus.ObjectPath(configs.APP_PATH), map[string]dbus.Variant{})
	if gattCall.Err != nil {
		log.Fatalf("❌ Erro ao registrar aplicação GATT: %v", gattCall.Err)
	}
	fmt.Println("✅ Aplicação BLE com GATT registrada no D-Bus")

	select {} // mantém processo ativo
}
