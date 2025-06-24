// main.go (refatorado)
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

	services.InitKafkaProducer()

	app := &services.Application{}
	service := &services.GATTService{Path: configs.SERVICE_PATH}
	advert := &services.Advertisement{Path: "/org/bluez/example/advertisement0"}

	// Exporta objetos
	_ = conn.Export(app, dbus.ObjectPath(configs.APP_PATH), "org.freedesktop.DBus.ObjectManager")
	_ = conn.Export(service, dbus.ObjectPath(configs.SERVICE_PATH), "org.bluez.GattService1")
	_ = conn.Export(advert, advert.Path, "org.bluez.LEAdvertisement1")
	_ = conn.Export(advert, advert.Path, "org.freedesktop.DBus.Properties")
	_ = conn.Export(introspect.NewIntrospectable(&introspect.Node{
		Name: string(advert.Path),
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{Name: "org.bluez.LEAdvertisement1"},
			{Name: "org.freedesktop.DBus.Properties"},
		},
	}), advert.Path, "org.freedesktop.DBus.Introspectable")

	// Introspecção para APP_PATH
	_ = conn.Export(introspect.NewIntrospectable(&introspect.Node{
		Name: configs.APP_PATH,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{Name: "org.freedesktop.DBus.ObjectManager"},
			{Name: "org.bluez.GattService1"},
		},
	}), dbus.ObjectPath(configs.APP_PATH), "org.freedesktop.DBus.Introspectable")

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
