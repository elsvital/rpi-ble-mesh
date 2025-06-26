// characteristic_machine.go
package services

import (
	"ble_go_server/configs"
	"fmt"
	"log"

	//jwt "github.com/dgrijalva/jwt-go"
	"github.com/godbus/dbus/v5"
	//"time"
)

type MachineDataCharacteristic struct {
	Path string
}

func (c *MachineDataCharacteristic) WriteValue(value []byte, options map[string]dbus.Variant) *dbus.Error {
	fmt.Println("📥 Valor recebido criptografado (raw):", string(value))

	pyld, err := ExtractDataFromJSON(value)
	if err != nil {
		return dbus.MakeFailedError(fmt.Errorf("Falha ao extrair dados JSON: %v", err))
	}

	// Você já obteve os dados de dentro do JWT, não precisa fazer Unmarshal de novo aqui
	tokenStr := pyld.JWT
	if tokenStr == "" {
		return dbus.MakeFailedError(fmt.Errorf("JWT ausente no payload"))
	}

	dbPath := "../db/raspi_edge.sqlite"
	dbc, err := NewDBController(dbPath)
	if err != nil {
		log.Fatalf("Erro ao abrir o banco: %v", err)
	}
	defer dbc.DB.Close()

	err = dbc.UpsertSessao(pyld.MicroID, tokenStr)
	if err != nil {
		log.Printf("Erro ao gravar session: %v", err)
	}

	err = dbc.SaveEdgeData(pyld)
	if err != nil {
		log.Printf("Erro ao gravar treino: %v", err)
	} else {
		log.Println("Treino gravado com sucesso.")
	}

	return nil
}

func (c *MachineDataCharacteristic) GetProperties() map[string]map[string]dbus.Variant {
	return map[string]map[string]dbus.Variant{
		"org.bluez.GattCharacteristic1": {
			"UUID":    dbus.MakeVariant(configs.CHARACTERISTIC_UUID_MACHINE_DATA),
			"Service": dbus.MakeVariant(dbus.ObjectPath(configs.SERVICE_PATH)),
			"Flags":   dbus.MakeVariant([]string{"write", "write-without-response"}),
		},
	}
}
