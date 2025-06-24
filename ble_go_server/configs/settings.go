// settings.go
package configs

const (
	// Segurança
	SECRET_KEY = "chave-secreta"

	// Identificadores de serviço BLE
	SERVICE_UUID = "4fafc201-1fb5-459e-8fcc-c5c9c331914b"
	SERVICE_PATH = "/org/bluez/example/service0"
	APP_PATH     = "/org/bluez/example"

	// Característica: Machine Data
	CHARACTERISTIC_UUID_MACHINE_DATA = "beb5483e-36e1-4688-b7f5-ea07361b26a8"
	CHARACTERISTIC_PATH_MACHINE_DATA = "/org/bluez/example/service0/char0"

	// Característica: Firmware Version (escrita)
	CHARACTERISTIC_UUID_FIRMWARE_VERSION = "12345678-1234-5678-1234-56789abcdef1"
	CHARACTERISTIC_PATH_FIRMWARE_VERSION = "/org/bluez/example/service0/char1"

	// Característica: Firmware Data (leitura)
	CHARACTERISTIC_UUID_FIRMWARE_DATA = "12345678-1234-5678-1234-56789abcdef2"
	CHARACTERISTIC_PATH_FIRMWARE_DATA = "/org/bluez/example/service0/char2"

	// Kafka
	KAFKA_BROKER = "192.168.15.160:9092"
	KAFKA_TOPIC  = "sensor.data"
)
