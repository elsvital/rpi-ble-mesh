# Servidor BLE para Dispositivos ESP32 com Suporte a OTA

Este servidor BLE, desenvolvido em Go, oferece uma interface baseada em GATT (Generic Attribute Profile) para comunicação com dispositivos ESP32. Ele permite tanto o envio de dados operacionais quanto a verificação e atualização remota de firmware via OTA (Over-The-Air), utilizando características BLE dedicadas.

## 🔧 Funcionalidades

### ✅ 1. Registro de Dados da Máquina
- **Characteristic:** `MachineDataCharacteristic`
- **UUID:** `beb5483e-36e1-4688-b7f5-ea07361b26a8`
- **Função:** Recebe dados em JSON via BLE contendo informações de uso e um token JWT. Valida o token, extrai os dados e armazena no banco SQLite.

### 🔄 2. Verificação de Versão de Firmware (OTA)
- **Characteristic:** `FirmwareVersionCharacteristic`
- **UUID:** `12345678-1234-5678-1234-56789abcdef1`
- **Função:** O ESP32 envia sua versão atual de firmware. O servidor consulta o banco e, se houver nova versão disponível, prepara o binário para envio via BLE.

### 📤 3. Transferência de Binário OTA via BLE
- **Characteristic:** `FirmwareDataCharacteristic`
- **UUID:** `12345678-1234-5678-1234-56789abcdef2`
- **Função:** Fornece o binário OTA em blocos de 512 bytes para o ESP32 ler sequencialmente.

## 💾 Banco de Dados
O servidor utiliza SQLite para armazenar:
- Dados operacionais (`central_micro`)
- Versões disponíveis para OTA (`versao_online`)
- Registro dos dados recebidos de cada ESP32

## 📡 Anúncio BLE
- O servidor se anuncia como `RPi-BLE`, com UUID do serviço: `4fafc201-1fb5-459e-8fcc-c5c9c331914b`.

## 🧱 Estrutura Modular
O projeto está organizado de forma modular:
- `characteristic_machine.go`: dados da máquina
- `characteristic_firmware.go`: controle de versão e binário OTA
- `service.go`, `application.go`: definição de serviço e integração com D-Bus
- `advertisement.go`: anúncio BLE
- `main.go`: registro da aplicação no BlueZ e ciclo principal

## ⚙️ Tecnologias
- Go + BlueZ (via D-Bus)
- SQLite (persistência local)
- Kafka (opcional – integração comentada)
- JWT (validação de segurança)

## 📝 Licença
Este projeto é distribuído sob a licença MIT.

---