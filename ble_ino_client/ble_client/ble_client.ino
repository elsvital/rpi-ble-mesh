#include <NimBLEDevice.h>
#include <ArduinoJson.h>
#include <Update.h>
#include "esp_app_desc.h"

#include <AESLib.h>
#include <base64.h>  // ou mbedtls/base64.h, dependendo do seu ESP32 core

AESLib aesLib;


#define SERVICE_UUID            "4fafc201-1fb5-459e-8fcc-c5c9c331914b"
#define CHARACTERISTIC_UUID     "beb5483e-36e1-4688-b7f5-ea07361b26a8"
#define FW_VERSION_UUID         "12345678-1234-5678-1234-56789abcdef1"
#define FW_DATA_UUID            "12345678-1234-5678-1234-56789abcdef2"

const char* jwt_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoibWNfMDAxIiwiZXhwIjoxNzUxMDcyMzIzfQ.CCou1ozAz4axE8zjZYrBruRTDVC3EAVsQ9jhePAydsM";

static NimBLEAdvertisedDevice* targetDevice = nullptr;
unsigned long lastOTACheck = 0;

byte aes_key[] = {
  0x00, 0x01, 0x02, 0x03,
  0x04, 0x05, 0x06, 0x07,
  0x08, 0x09, 0x0A, 0x0B,
  0x0C, 0x0D, 0x0E, 0x0F
};

byte aes_iv[16];  // ✅ tipo correto

String encryptAndEncode(const char* msg) {
  for (int i = 0; i < 16; i++) aes_iv[i] = random(0, 256);

  int msgLen = strlen(msg);
  int padLen = 16 - (msgLen % 16);
  int paddedLen = msgLen + padLen;

  byte padded[paddedLen];
  memcpy(padded, msg, msgLen);
  for (int i = 0; i < padLen; i++) padded[msgLen + i] = padLen;

  byte encrypted[paddedLen];
  // aesLib.encrypt(padded, paddedLen, encrypted, aes_key, 128, (byte*)aes_iv);
  byte iv_copy[16];
  memcpy(iv_copy, aes_iv, 16);
  aesLib.encrypt(padded, paddedLen, encrypted, aes_key, 128, iv_copy);

  byte combined[16 + paddedLen];
  memset(combined, 0, sizeof(combined)); 
  memcpy(combined, aes_iv, 16);
  memcpy(combined + 16, encrypted, paddedLen);

  String encoded = base64::encode(combined, 16 + paddedLen);
  return encoded;

}

class ClientCallbacks : public NimBLEClientCallbacks {
  void onConnect(NimBLEClient* pClient) override {
    Serial.println("✅ Conectado ao servidor BLE!");
  }

  void onDisconnect(NimBLEClient* pClient) override {
    Serial.println("❌ Desconectado.");
  }
};

bool connectAndSend() {
  NimBLEClient* pClient = NimBLEDevice::createClient();
  pClient->setClientCallbacks(new ClientCallbacks(), false);

  Serial.print("🔗 Conectando a ");
  Serial.println(targetDevice->getAddress().toString().c_str());

  if (!pClient->connect(targetDevice)) {
    Serial.println("❌ Falha ao conectar.");
    NimBLEDevice::deleteClient(pClient);
    return false;
  }

  NimBLERemoteService* pService = pClient->getService(SERVICE_UUID);
  if (!pService) {
    Serial.println("❌ Serviço não encontrado.");
    pClient->disconnect();
    NimBLEDevice::deleteClient(pClient);
    return false;
  }

  NimBLERemoteCharacteristic* pChar = pService->getCharacteristic(CHARACTERISTIC_UUID);
  if (!pChar || !pChar->canWrite()) {
    Serial.println("❌ Característica inválida.");
    pClient->disconnect();
    NimBLEDevice::deleteClient(pClient);
    return false;
  }

  StaticJsonDocument<256> doc;
  doc["jwt"] = jwt_token;
  doc["micro_id"] = "mc_001";
  doc["user_id"] = random(200, 300);
  doc["total_reps"] = random(6, 12);
  doc["failed_reps"] = random(0, 3);
  doc["total_series"] = random(0, 3);

  char jsonStr[256];
  serializeJson(doc, jsonStr);

  String encryptedPayload = encryptAndEncode(jsonStr);

  if (pChar->writeValue((uint8_t*)encryptedPayload.c_str(), encryptedPayload.length(), false)) {
    Serial.print("📤 Enviado criptografado: ");
    Serial.println(encryptedPayload);
  } else {
    Serial.println("❌ Falha ao enviar.");
  }

  // if (pChar->writeValue((uint8_t*)jsonStr, strlen(jsonStr), false)) {
  //   Serial.print("📤 Enviado: ");
  //   Serial.println(jsonStr);
  // } else {
  //   Serial.println("❌ Falha ao escrever.");
  // }

  pClient->disconnect();
  NimBLEDevice::deleteClient(pClient);
  return true;
}

const char* getFirmwareVersion() {
  return esp_app_get_description()->version;
}


void checkFirmwareUpdate() {
  NimBLEClient* pClient = NimBLEDevice::createClient();
  pClient->setClientCallbacks(new ClientCallbacks(), false);

  if (!pClient->connect(targetDevice)) {
    Serial.println("❌ Falha ao conectar (OTA).");
    NimBLEDevice::deleteClient(pClient);
    return;
  }

  NimBLERemoteService* service = pClient->getService(SERVICE_UUID);
  if (!service) {
    Serial.println("❌ Serviço OTA não encontrado.");
    pClient->disconnect();
    NimBLEDevice::deleteClient(pClient);
    return;
  }

  // auto verChar = service->getCharacteristic(FW_VERSION_UUID);
  // if (verChar && verChar->canWrite()) {
  //   const char* fwVersion = getFirmwareVersion();
  //   Serial.print("📨 Enviando versão atual: ");
  //   Serial.println(fwVersion);
  //   verChar->writeValue(fwVersion);
  //   delay(200);
  // }



  auto verChar = service->getCharacteristic(FW_VERSION_UUID);

  StaticJsonDocument<256> doc;
  const char* fwVersion = getFirmwareVersion();
  doc["jwt"] = jwt_token;
  doc["micro_id"] = "mc_001";
  doc["user_id"] = fwVersion;
  
  char jsonStr[256];
  serializeJson(doc, jsonStr);
  
  String encryptedPayload = encryptAndEncode(jsonStr);
  
  if (verChar->writeValue((uint8_t*)encryptedPayload.c_str(), encryptedPayload.length(), false)) {
    Serial.print("📤 Enviado solicitação de atualização: ");
    Serial.println(encryptedPayload);
  } else {
    Serial.println("❌ Falha ao enviar solicitação de atualização OTA.");
  }

  // if (verChar->writeValue((uint8_t*)jsonStr, strlen(jsonStr), false)) {
  //   Serial.print("📤 Enviado solicitação de atualização");
  // } else {
  //   Serial.println("❌ Falha ao enviar solicitação de atualização OTA.");
  // }


  auto binChar = service->getCharacteristic(FW_DATA_UUID);
  if (binChar && binChar->canRead()) {
    Serial.println("📦 Iniciando leitura de firmware...");
    Update.begin(UPDATE_SIZE_UNKNOWN);
    int total = 0;
    while (true) {
      std::string chunk = binChar->readValue();
      if (chunk.empty()) break;
      std::vector<uint8_t> buffer(chunk.begin(), chunk.end());
      Update.write(buffer.data(), buffer.size());
      total += buffer.size();
    }
    if (Update.end()) {
      Serial.printf("✅ Firmware atualizado (%d bytes). Reiniciando...\n", total);
      ESP.restart();
    } else {
      Serial.printf("❌ Erro OTA: %s\n", Update.errorString());
    }
  }

  pClient->disconnect();
  NimBLEDevice::deleteClient(pClient);
}

bool shouldCheckOTA() {
  unsigned long now = millis();
  // if (now - lastOTACheck > 86400000UL) { // simulação de 24h
  if (now - lastOTACheck > 600000UL) { // 10 minutos = 600.000 ms
    lastOTACheck = now;
    return true;
  }
  return false;
}

bool shouldCheckUserData() {
  unsigned long now = millis();
  if (now - lastOTACheck > 60000UL) { // 10 minutos = 600.000 ms
    lastOTACheck = now;
    return true;
  }
  return false;
}


class AdvertisedDeviceCallbacks : public NimBLEAdvertisedDeviceCallbacks {
public:
  void onResult(NimBLEAdvertisedDevice* advertisedDevice) override {
    Serial.print("📡 Visto: ");
    Serial.println(advertisedDevice->toString().c_str());

    if (advertisedDevice->haveName() &&
        advertisedDevice->getName().find("RPi-BLE") != std::string::npos) {
      Serial.println("🎯 Dispositivo alvo encontrado!");
      targetDevice = advertisedDevice;
      NimBLEDevice::getScan()->stop();
    }
  }
};

void setup() {
  Serial.begin(115200);
  delay(500);

  Serial.println("🚀 Iniciando cliente BLE...");
  NimBLEDevice::init("");

  NimBLEScan* pScan = NimBLEDevice::getScan();
  pScan->setAdvertisedDeviceCallbacks(new AdvertisedDeviceCallbacks(), false);
  pScan->setInterval(100);
  pScan->setWindow(99);
  pScan->setActiveScan(true);
  pScan->start(5, false);
}

void loop() {
  if (targetDevice) {
    if (shouldCheckUserData()){
      connectAndSend();
    }
    delay(5000);
    if (shouldCheckOTA()) {
      checkFirmwareUpdate();
    }
  } else {
    Serial.println("🔎 Buscando servidor BLE...");
    NimBLEDevice::getScan()->start(5, true);
  }
}