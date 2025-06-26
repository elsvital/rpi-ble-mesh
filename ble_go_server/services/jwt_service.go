// jwt_service.go - Serviço para extração e validação de usuário a partir do token JWT
package services

import (
	"ble_go_server/configs"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	//"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var aesKey = []byte{
	0x00, 0x01, 0x02, 0x03,
	0x04, 0x05, 0x06, 0x07,
	0x08, 0x09, 0x0A, 0x0B,
	0x0C, 0x0D, 0x0E, 0x0F,
}

func decryptPayload(base64Input string) (string, error) {
	fmt.Println("🔐 Iniciando decryptPayload")
	// 1. Decode base64
	cipherData, err := base64.StdEncoding.DecodeString(base64Input)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %v", err)
	}

	if len(cipherData) < aes.BlockSize {
		return "", errors.New("cipher too short")
	}

	if (len(cipherData)-aes.BlockSize)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext (%d bytes) is not a multiple of block size", len(cipherData)-aes.BlockSize)
	}

	// 2. Extract IV and ciphertext
	iv := cipherData[:aes.BlockSize]
	ciphertext := cipherData[aes.BlockSize:]

	fmt.Printf("🔐 IV extraído: %x\n", iv)
	fmt.Printf("🔐 Tamanho do ciphertext: %d\n", len(ciphertext))

	// 3. Create cipher
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("aes.NewCipher failed: %v", err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of block size")
	}

	// 4. Decrypt
	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	copy(decrypted, ciphertext)
	mode.CryptBlocks(decrypted, decrypted)
	fmt.Println("🔐 Decriptação feita, tentando remover padding...")
	fmt.Printf("🔓 Decrypted (hex): %x\n", decrypted)

	// 5. Remove padding (PKCS#7)
	plaintext, err := pkcs7Unpad(decrypted)
	if err != nil {
		return "", fmt.Errorf("unpad failed: %v", err)
	}

	return string(plaintext), nil
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("data is empty")
	}
	padLen := int(data[length-1])
	if padLen > aes.BlockSize || padLen == 0 {
		return nil, errors.New("invalid padding size")
	}
	for _, b := range data[length-padLen:] {
		if int(b) != padLen {
			return nil, errors.New("invalid padding bytes")
		}
	}
	return data[:length-padLen], nil
}

func ExtractDataFromJSON(value []byte) (Treino, error) {
	// Passo 1: descriptografar
	decryptedStr, err := decryptPayload(string(value))
	if err != nil {
		return Treino{}, fmt.Errorf("❌ Erro ao descriptografar: %v", err)
	}

	// Passo 2: decodificar JSON
	var payload map[string]interface{}
	err = json.Unmarshal([]byte(decryptedStr), &payload)
	if err != nil {
		return Treino{}, fmt.Errorf("❌ JSON inválido: %v", err)
	}
	fmt.Println("📥 Valor recebido descriptografado (raw):", payload)

	fmt.Println("Extrair dados...")
	// Extrai dados
	var microID string
	var userID string
	var totalRepsFloat float64
	var failedRepsFloat float64
	var totalSeriesFloat float64

	totalRepsFloat = 0.0
	failedRepsFloat = 0.0
	totalSeriesFloat = 0.0

	fmt.Println("Preenche estrutura inicio")
	microID = payload["micro_id"].(string)
	userID = payload["user_id"].(string)
	total_reps, exists := payload["total_reps"]
	if exists {
		totalRepsFloat = total_reps.(float64)
	}
	failed_reps, exists := payload["failed_reps"]
	if exists {
		failedRepsFloat = failed_reps.(float64)
	}
	total_series, exists := payload["total_series"]
	if exists {
		totalSeriesFloat = total_series.(float64)
	}
	fmt.Println("Preenche estrutura")
	// Preenche estrutura
	jdata := Treino{
		MicroID:     microID,
		UserID:      userID,
		TotalReps:   int(totalRepsFloat),
		FailedReps:  int(failedRepsFloat),
		TotalSeries: int(totalSeriesFloat),
	}
	return jdata, nil
}

func ExtractDataFromJSONWithTOKEN(value []byte) (Treino, error) {
	// Passo 1: descriptografar
	decryptedStr, err := decryptPayload(string(value))
	if err != nil {
		return Treino{}, fmt.Errorf("❌ Erro ao descriptografar: %v", err)
	}

	// Passo 2: decodificar JSON
	var payload map[string]interface{}
	err = json.Unmarshal([]byte(decryptedStr), &payload)
	if err != nil {
		return Treino{}, fmt.Errorf("❌ JSON inválido: %v", err)
	}
	fmt.Println("📥 Valor recebido descriptografado (raw):", payload)

	tokenStr, ok := payload["jwt"].(string)
	if !ok {
		return Treino{}, errors.New("JWT ausente no payload")
	}
	fmt.Println("<UNK> Valor recebido token:", tokenStr)
	// Valida JWT
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		return Treino{}, fmt.Errorf("JWT inválido: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Treino{}, errors.New("Falha ao extrair claims do token")
	}

	// Valida expiração
	if exp, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			return Treino{}, fmt.Errorf("JWT expirado em %v", expTime)
		}
	}

	// Extrai dados
	microID, ok1 := claims["user"].(string)
	var userID string
	var totalRepsFloat float64
	var failedRepsFloat float64
	var totalSeriesFloat float64

	userIDFloat := payload["user_id"].(float64)
	userID = strconv.Itoa(int(userIDFloat))

	totalRepsFloat = payload["total_reps"].(float64)
	failedRepsFloat = payload["failed_reps"].(float64)
	totalSeriesFloat = payload["total_series"].(float64)

	if !ok1 {
		return Treino{}, errors.New("❌ Campos obrigatórios ausentes ou com tipo inválido no token")
	}

	// Preenche estrutura
	jdata := Treino{
		MicroID:     microID,
		UserID:      userID,
		TotalReps:   int(totalRepsFloat),
		FailedReps:  int(failedRepsFloat),
		TotalSeries: int(totalSeriesFloat),
		JWT:         tokenStr,
	}
	fmt.Println("<UNK> Valor micro:", microID)
	return jdata, nil
}

func ExtractDataFromJSON_old(value []byte) (*Treino, error) {

	//var payload map[string]interface{}
	//err := json.Unmarshal(value, &payload)
	//if err != nil {
	//	return nil, fmt.Errorf("JSON inválido: %v", err)
	//}

	// Passo 1: descriptografar
	decryptedStr, err := decryptPayload(string(value))
	if err != nil {
		return nil, fmt.Errorf("❌ Erro ao descriptografar: %v", err)
	}

	// Passo 2: decodificar JSON
	var payload map[string]interface{}
	err = json.Unmarshal([]byte(decryptedStr), &payload)
	if err != nil {
		return nil, fmt.Errorf("❌ JSON inválido: %v", err)
	}

	tokenStr, ok := payload["jwt"].(string)
	if !ok {
		return nil, errors.New("JWT ausente no payload")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("JWT inválido: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Falha ao extrair claims do token")
	}

	if exp, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			return nil, fmt.Errorf("JWT expirado em %v", expTime)
		}
	}

	MicroID, ok := claims["user"].(string)
	userID, ok := claims["user_id"].(string)
	total_reps, ok := claims["total_reps"].(int)
	failed_reps, ok := claims["failed_reps"].(int)
	total_series, ok := claims["total_series"].(int)

	var jdata Treino
	jdata.MicroID = MicroID
	jdata.UserID = userID
	jdata.TotalReps = total_reps
	jdata.FailedReps = failed_reps
	jdata.TotalSeries = total_series

	if !ok {
		return nil, errors.New("Campo 'user' ausente no token")
	}

	return &jdata, nil
}
