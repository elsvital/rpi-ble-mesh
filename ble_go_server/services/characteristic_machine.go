// characteristic_machine.go
package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"ble_go_server/configs"
	"ble_go_server/services"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/godbus/dbus/v5"
)

type MachineDataCharacteristic struct {
	Path string
}

func (c *MachineDataCharacteristic) WriteValue(value []byte, options map[string]dbus.Variant) *dbus.Error {
	fmt.Println("📥 Valor recebido (raw):", string(value))

	var payload map[string]interface{}
	err := json.Unmarshal(value, &payload)
	if err != nil {
		return dbus.MakeFailedError(fmt.Errorf("JSON inválido: %v", err))
	}

	tokenStr, ok := payload["jwt"].(string)
	if !ok {
		return dbus.MakeFailedError(fmt.Errorf("JWT ausente"))
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		return dbus.MakeFailedError(fmt.Errorf("JWT inválido: %v", err))
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println("✅ JWT válido:", claims)
		if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
			return dbus.MakeFailedError(fmt.Errorf("JWT expirado: exp=%d", int64(exp)))
		}
		payload["user"] = claims["user"]
	}

	dbPath := "db/edge_vollta_sqlite"
	dbc, err := services.NewDBController(dbPath)
	if err != nil {
		log.Fatalf("Erro ao abrir o banco: %v", err)
	}
	defer dbc.DB.Close()

	err = dbc.SaveEdgeData(payload)
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
