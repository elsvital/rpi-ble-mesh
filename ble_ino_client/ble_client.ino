#include <NimBLEDevice.h>
#include <ArduinoJson.h>

#define SERVICE_UUID        "4fafc201-1fb5-459e-8fcc-c5c9c331914b"
#define CHARACTERISTIC_UUID "beb5483e-36e1-4688-b7f5-ea07361b26a8"


// Token JWT fixo gerado com a mesama chave que sera recuperado no servidor BLE
const char* jwt_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiZXNwMzItdGVzdGUiLCJleHAiOjE3NDY5Nzc5OTZ9.8wK7fwuYi0Iw--XiyT5e-ttgwkTKmemyagtvkmHq3Zs";

static NimBLEAdvertisedDevice* targetDevice = nullptr;

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

  delay(500);

  Serial.print("🔗 Conectando a ");
  Serial.println(targetDevice->getAddress().toString().c_str());

  if (!pClient->connect(targetDevice)) {
    Serial.println("❌ Falha ao conectar.");
    NimBLEDevice::deleteClient(pClient);
    return false;
  }

  // Debug: listar serviços encontrados
  std::vector<NimBLERemoteService*>* services = pClient->getServices(true);
  Serial.println("📜 Serviços disponíveis:");
  for (auto* svc : *services) {
    Serial.println(svc->getUUID().toString().c_str());
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
  doc["temp"] = random(200, 300) / 10.0;
  doc["hum"] = random(400, 700) / 10.0;

  char jsonStr[256];
  serializeJson(doc, jsonStr);

  if (pChar->writeValue((uint8_t*)jsonStr, strlen(jsonStr), false)) {
    Serial.print("📤 Enviado: ");
    Serial.println(jsonStr);
  } else {
    Serial.println("❌ Falha ao escrever.");
  }

  pClient->disconnect();
  NimBLEDevice::deleteClient(pClient);
  return true;
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
    connectAndSend();
    delay(5000);
  } else {
    Serial.println("🔎 Buscando servidor BLE...");
    NimBLEDevice::getScan()->start(5, true);
  }
}