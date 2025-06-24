import dbus
import dbus.service
from configs.settings import TEST_SERVICE_UUID, TEST_CHARACTERISTIC_UUID
from services.kafka import send_to_kafka
import jwt
import json

DBUS_OM_IFACE = 'org.freedesktop.DBus.ObjectManager'

class Application(dbus.service.Object):
    def __init__(self, bus):
        self.path = '/custom/app'
        self.services = []
        self.bus = bus
        dbus.service.Object.__init__(self, bus, self.path)

    def get_path(self):
        return dbus.ObjectPath(self.path)

    def add_service(self, service):
        self.services.append(service)

    def get_dbus_objects(self):
        result = {}
        for service in self.services:
            result[service.get_path()] = service.get_properties()
            for char in service.characteristics:
                result[char.get_path()] = char.get_properties()
        return result

    @dbus.service.method(DBUS_OM_IFACE, out_signature='a{oa{sa{sv}}}')
    def GetManagedObjects(self):
        return self.get_dbus_objects()


class Service(dbus.service.Object):
    PATH_BASE = '/custom/service'

    def __init__(self, bus, index, uuid, primary):
        self.path = self.PATH_BASE + str(index)
        self.bus = bus
        self.uuid = uuid
        self.primary = primary
        self.characteristics = []
        dbus.service.Object.__init__(self, bus, self.path)

    def get_properties(self):
        return {
            'org.bluez.GattService1': {
                'UUID': self.uuid,
                'Primary': self.primary,
                'Characteristics': dbus.Array(
                    [char.get_path() for char in self.characteristics], signature='o'
                )
            }
        }

    def get_path(self):
        return dbus.ObjectPath(self.path)


class Characteristic(dbus.service.Object):
    def __init__(self, bus, index, uuid, flags, service):
        self.path = service.path + '/char' + str(index)
        self.bus = bus
        self.uuid = uuid
        self.flags = flags
        self.service = service
        dbus.service.Object.__init__(self, bus, self.path)

    def get_properties(self):
        return {
            'org.bluez.GattCharacteristic1': {
                'Service': self.service.get_path(),
                'UUID': self.uuid,
                'Flags': self.flags
            }
        }

    def get_path(self):
        return dbus.ObjectPath(self.path)

    @dbus.service.method('org.bluez.GattCharacteristic1', in_signature='aya{sv}')
    def WriteValue(self, value, options):
        try:
            decoded = bytes(value).decode('utf-8')
            print("📥 Valor recebido (raw):", decoded)

            data = json.loads(decoded)
            from configs.settings import SECRET_KEY
            payload = jwt.decode(data['jwt'], SECRET_KEY, algorithms=["HS256"])
            print("✅ JWT válido:", payload)

            data['user'] = payload['user']
            send_to_kafka(data)

        except Exception as e:
            print("❌ Erro no processamento do valor:", e)


class Advertisement(dbus.service.Object):
    PATH = '/custom/advertisement'

    def __init__(self, bus, index, advertising_type):
        self.path = self.PATH + str(index)
        self.bus = bus
        self.ad_type = advertising_type
        dbus.service.Object.__init__(self, bus, self.path)

    def get_properties(self):
        return {
            'org.bluez.LEAdvertisement1': {
                'Type': self.ad_type,
                'ServiceUUIDs': [TEST_SERVICE_UUID],
                'LocalName': 'RPi-BLE',
                'Includes': ['tx-power']
            }
        }

    def get_path(self):
        return dbus.ObjectPath(self.path)

    @dbus.service.method('org.freedesktop.DBus.Properties',
                         in_signature='ss', out_signature='v')
    def Get(self, interface, prop):
        return self.get_properties()[interface][prop]

    @dbus.service.method('org.freedesktop.DBus.Properties',
                         in_signature='s', out_signature='a{sv}')
    def GetAll(self, interface):
        return self.get_properties()[interface]

    @dbus.service.method('org.bluez.LEAdvertisement1')
    def Release(self):
        print("🔌 Anúncio liberado")