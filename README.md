# 🛰️ RPi-BLE: Gateway BLE com ESP32-S3

Este projeto é um sistema distribuído para coleta de dados via **BLE (Bluetooth Low Energy)** usando **ESP32-S3** como cliente e um **Raspberry Pi** como servidor BLE, com envio dos dados para uma infraestrutura de mensagens baseada em **Apache Kafka** (opcional). O sistema possui módulos escritos em Go, Python, C++ (Arduino) e integração com interface web em Django para armazenamento de dados de borda.

---

## 🧱 Estrutura do Projeto

```text
rpi-ble/
├── ble_go_server/          # Servidor BLE em Go + Kafka
│   ├── cmd/                # Ponto de entrada do servidor
│   ├── configs/            # Configurações como chaves JWT
│   ├── firmware/           # Firmware OTA armazenado para atualização remota
│   ├── services/           # Serviços: GATT, DB, Kafka, JWT
│   ├── install-and-run.sh  # Script de instalação e execução
│
├── ble_ino_client/         # Firmware para ESP32-S3 (Arduino)
│   ├── ble_client/         # Código principal do cliente BLE
│
├── ble_py_server/          # Versão alternativa do servidor BLE em Python
│
├── db/
│   └── raspi_edge.sqlite   # Banco de dados local (SQLite)
│
├── raspi_web_mgmt/         # Interface Web para gerenciar sessões e treinos
│   ├── core/               # Aplicação principal em Django com o modelo de dados 
│
├── arch.png                # Diagrama da arquitetura
├── go.mod                  # Dependências Go
└── README.md
```

---
## 🧠 Conceitos e Tecnologias

### 📡 BLE (Bluetooth Low Energy)
- **Cliente (ESP32-S3)**: envia JSON criptografado com dados da sessão (ex.: identificação, usuário, repetições etc).
- **Servidor (Raspberry Pi com Go)**: atua como GATT Server, recebendo dados por `Characteristic` escrita pelo cliente.

### 🔐 Criptografia AES + JWT
- O JSON é criptografado com **AES-CBC de 128 bits**, com IV aleatório, e depois codificado em base64.
- O campo `jwt` é um **JSON Web Token** assinado no ESP32 com chave secreta e validado no servidor.
- Exemplo de payload criptografado:
  ```json
  {
    "jwt": "<token>",
    "user_id": 245,
    "total_reps": 10,
    "failed_reps": 1,
    "total_series": 2,
    "micro_id": "mc_001"
  }
  ```

### 🐹 Go
- Utilizado no backend BLE Server (`ble_go_server/`).
- Principais bibliotecas:
  - [`github.com/godbus/dbus`](https://github.com/godbus/dbus): para manipular o BlueZ via D-Bus.
  - `crypto/aes`, `encoding/base64`: para descriptografar o payload do ESP32.
  - `github.com/golang-jwt/jwt`: para validar tokens JWT.

### 🧪 SQLite
- Banco local para armazenar sessões de treino e atualizações OTA e sessões do clientes (micro-controladores).
- Arquivo `raspi_edge.sqlite`.

### ☁️ Apache Kafka
- Usado para envio dos dados do BLE para um sistema de backend analítico.
- Produtor Kafka está no arquivo `kafka.go`.

### 🧠 Python (opcional)
- Servidor BLE alternativo usando `dbus-python`.
- Útil para testes e desenvolvimento rápido (`ble_py_server/`).

### ⚙️ ESP32 com Arduino
- Utiliza **NimBLE-Arduino** (leve, ideal para BLE Client).
- A criptografia é feita com a biblioteca `AESLib`.
- Firmware está em `ble_ino_client/ble_client/ble_client.ino`.

### 🌐 Django Web
- Painel web para visualização e gestão dos dados coletados.
- Estrutura padrão do Django dentro de `raspi_web_mgmt/`.

---

## 🚀 Como Executar

### 1. Requisitos

#### Raspberry Pi:
- Go ≥ 1.20
- BlueZ (versão compatível com GATT + D-Bus)
- SQLite
- Apache Kafka (em rede ou local)
- Python 3 (opcional para servidor BLE em Python)

#### ESP32:
- Placa: `ESP32-S3 Dev Module`
- Bibliotecas:
  - `AESLib`
  - `ArduinoJson`
  - `NimBLE-Arduino`
  - `base64`

### 2. Compilar e Executar Servidor Go

```bash
cd ble_go_server
chmod + x install-and-run.sh
./install-and-run.sh
```

### 3. Atualizar Firmware OTA (opcional)
Coloque o novo `.bin` em `ble_go_server/firmware/esp32-s3/` com o nome `ble_client.ino.bin`.

### 4. Interface Web

```bash
cd raspi_web_mgmt
python -m  venv ble_env
source ble_env/bin/activate
pip install -r requirements.txt
python manage.py runserver
```

---

## 📷 Diagrama da Arquitetura

![Arquitetura do Servidor BLE](arch.png)

---

## 📄 Licença

MIT License

Copyright (c) 2025 Elinilson Vital

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.



