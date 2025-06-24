import dbus
import dbus.mainloop.glib
from gi.repository import GLib

from services.bluetooth import Application, Service, Characteristic, Advertisement
from configs.settings import TEST_SERVICE_UUID, TEST_CHARACTERISTIC_UUID

def main():
    dbus.mainloop.glib.DBusGMainLoop(set_as_default=True)
    bus = dbus.SystemBus()

    adapter_path = "/org/bluez/hci0"
    adapter = dbus.Interface(bus.get_object("org.bluez", adapter_path), "org.freedesktop.DBus.Properties")
    adapter.Set("org.bluez.Adapter1", "Powered", dbus.Boolean(1))

    app = Application(bus)
    service = Service(bus, 0, TEST_SERVICE_UUID, True)
    characteristic = Characteristic(bus, 0, TEST_CHARACTERISTIC_UUID, ['write'], service)

    service.characteristics.append(characteristic)
    app.add_service(service)

    # Registrar GATT
    service_manager = dbus.Interface(bus.get_object("org.bluez", adapter_path), "org.bluez.GattManager1")
    service_manager.RegisterApplication(app.get_path(), {},
                                        reply_handler=lambda: print("✅ Serviço GATT registrado"),
                                        error_handler=lambda e: print("❌ Falha ao registrar serviço GATT:", e))

    # Registrar anúncio BLE
    adv = Advertisement(bus, 0, "peripheral")
    ad_manager = dbus.Interface(bus.get_object("org.bluez", adapter_path), "org.bluez.LEAdvertisingManager1")
    ad_manager.RegisterAdvertisement(adv.get_path(), {},
                                     reply_handler=lambda: print("📡 Anúncio BLE registrado"),
                                     error_handler=lambda e: print("❌ Falha no anúncio BLE:", e))

    try:
        print("🚀 Servidor BLE iniciado")
        GLib.MainLoop().run()
    except KeyboardInterrupt:
        print("🛑 Encerrando servidor BLE...")
        ad_manager.UnregisterAdvertisement(adv.get_path())
        adv.remove_from_connection()
        service_manager.UnregisterApplication(app.get_path())

if __name__ == "__main__":
    main()