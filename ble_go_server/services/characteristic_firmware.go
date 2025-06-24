// characteristic_firmware.go
package services

import (
	"fmt"
	"log"
	"os"
	"sync"

	"ble_go_server/configs"
	"ble_go_server/services"

	"github.com/godbus/dbus/v5"
)

var firmwareChunks = make(map[string][][]byte)
var firmwareCursor = make(map[string]int)
var firmwareLock sync.Mutex

func prepareFirmware(deviceID string, bin []byte) {
	const chunkSize = 512
	var chunks [][]byte
	for i := 0; i < len(bin); i += chunkSize {
		end := i + chunkSize
		if end > len(bin) {
			end = len(bin)
		}
		chunks = append(chunks, bin[i:end])
	}
	firmwareLock.Lock()
	firmwareChunks[deviceID] = chunks
	firmwareCursor[deviceID] = 0
	firmwareLock.Unlock()
}

func nextChunk(deviceID string) []byte {
	firmwareLock.Lock()
	defer firmwareLock.Unlock()
	chunks, ok := firmwareChunks[deviceID]
	if !ok || firmwareCursor[deviceID] >= len(chunks) {
		delete(firmwareChunks, deviceID)
		firmwareCursor[deviceID] = 0
		return nil
	}
	chunk := chunks[firmwareCursor[deviceID]]
	firmwareCursor[deviceID]++
	return chunk
}

type FirmwareVersionCharacteristic struct {
	Path string
}

func (c *FirmwareVersionCharacteristic) WriteValue(value []byte, options map[string]dbus.Variant) *dbus.Error {
	deviceID := extractDeviceID(options)
	clientVersion := string(value)
	log.Printf("📥 %s enviou versão %s", deviceID, clientVersion)

	dbPath := "db/edge_vollta_sqlite"
	dbc, err := services.NewDBController(dbPath)
	if err != nil {
		return dbus.MakeFailedError(fmt.Errorf("falha ao abrir banco: %v", err))
	}
	defer dbc.DB.Close()

	var currentVersion string
	err = dbc.DB.QueryRow("SELECT versao_atual FROM central_micro WHERE nome_dispositivo = ?", deviceID).Scan(&currentVersion)
	if err != nil {
		return dbus.MakeFailedError(fmt.Errorf("dispositivo não encontrado: %v", err))
	}

	if currentVersion == clientVersion {
		log.Println("✅ Firmware já está atualizado.")
		return nil
	}

	var caminho string
	err = dbc.DB.QueryRow(`
		SELECT caminho_binario FROM versao_online
		WHERE tipo_controlador = 'ESP32' AND versao > ?
		ORDER BY versao DESC LIMIT 1`, clientVersion).Scan(&caminho)
	if err != nil {
		return dbus.MakeFailedError(fmt.Errorf("sem nova versão: %v", err))
	}

	bin, err := os.ReadFile(caminho)
	if err != nil {
		return dbus.MakeFailedError(fmt.Errorf("erro ao ler binário: %v", err))
	}

	prepareFirmware(deviceID, bin)
	log.Printf("🚀 OTA preparada para %s (%d bytes)", deviceID, len(bin))
	return nil
}

func (c *FirmwareVersionCharacteristic) GetProperties() map[string]map[string]dbus.Variant {
	return map[string]map[string]dbus.Variant{
		"org.bluez.GattCharacteristic1": {
			"UUID":    dbus.MakeVariant(configs.CHARACTERISTIC_UUID_FIRMWARE_VERSION),
			"Service": dbus.MakeVariant(dbus.ObjectPath(configs.SERVICE_PATH)),
			"Flags":   dbus.MakeVariant([]string{"write", "write-without-response"}),
		},
	}
}

type FirmwareDataCharacteristic struct {
	Path string
}

func (c *FirmwareDataCharacteristic) ReadValue(options map[string]dbus.Variant) ([]byte, *dbus.Error) {
	deviceID := extractDeviceID(options)
	chunk := nextChunk(deviceID)
	if chunk == nil {
		return []byte{}, nil
	}
	return chunk, nil
}

func (c *FirmwareDataCharacteristic) GetProperties() map[string]map[string]dbus.Variant {
	return map[string]map[string]dbus.Variant{
		"org.bluez.GattCharacteristic1": {
			"UUID":    dbus.MakeVariant(configs.CHARACTERISTIC_UUID_FIRMWARE_DATA),
			"Service": dbus.MakeVariant(dbus.ObjectPath(configs.SERVICE_PATH)),
			"Flags":   dbus.MakeVariant([]string{"read"}),
		},
	}
}

func extractDeviceID(options map[string]dbus.Variant) string {
	if val, ok := options["device"]; ok {
		if path, ok := val.Value().(dbus.ObjectPath); ok {
			return string(path)
		}
	}
	return "unknown-device"
}
