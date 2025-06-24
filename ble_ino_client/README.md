# BLE Client com NimBLE e ArduinoJson

Este projeto implementa um cliente BLE (Bluetooth Low Energy) baseado na biblioteca **NimBLEDevice** em conjunto com a **ArduinoJson**. O objetivo é conectar-se a um servidor BLE específico, localizar um serviço e uma característica, e enviar dados formatados em JSON.

---

## 📋 Funcionalidades

- **Scan BLE**: Busca por dispositivos BLE próximos.
- **Conexão Automática**: Conecta automaticamente ao dispositivo alvo identificado por um nome específico ("RPi-BLE").
- **Envio de Dados**: Após conectar-se:
  - Envia informações JSON contendo:
    - Um token JWT fixo.
    - Leitura simulada de temperatura.
    - Leitura simulada de umidade.
  - Dados são enviados para uma característica BLE configurada.
- **Desconexão**: Desconecta automaticamente após enviar os dados.

---

## 🛠️ Bibliotecas Utilizadas

- **NimBLEDevice (1.4.3)**: Biblioteca para interação BLE eficiente e leve.
- **esp32 (3.2.0)**: Modulo para dispositivos ESP32S3.
- **ArduinoJson (7.4.1)**: Biblioteca para criação e manipulação de documentos JSON.

---

## 🔧 Como Funciona?

1. **Definições Importantes:**

   - **UUIDs**: O cliente BLE procura por um serviço e uma característica específicos definidos no código.
   - **Token JWT**: Um token JWT fixo é enviado como parte dos dados JSON.

2. **Ciclo de Operação**:
   - Escaneia dispositivos BLE.
   - Busca pelo dispositivo cujo nome contém _"RPi-BLE"_.
   - Conecta ao dispositivo alvo, se encontrado.
   - Lê os serviços e acessa uma característica configurada.
   - Envia um objeto JSON para a característica como valor.
   - Desconecta e reinicia o processo.

---

## 📑 Estrutura JSON Enviada

O cliente BLE envia os dados no seguinte formato (exemplo):

```json
{
  "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "temp": 25.4,
  "hum": 45.6
}
```

- **jwt**: Token JWT fixo para autenticação.
- **temp**: Temperatura simulada (em ºC, com números randômicos entre 20.0 e 30.0).
- **hum**: Umidade simulada (em %, com números randômicos entre 40.0 e 70.0).

---

## 🚀 Iniciando o Cliente BLE

1. **Requisitos:**

   - Placa compatível com BLE (ESP32, por exemplo).
   - Ambiente de desenvolvimento Arduino IDE ou equivalente.

2. **Configuração Inicial:**

   - Carregue o código no seu dispositivo.
   - Certifique-se de que o servidor BLE alvo esteja ativo e configurado com o nome correto e UUIDs correspondentes.

3. **Funcionamento:**

   - Após iniciar, a placa começa a escanear dispositivos BLE disponíveis.
   - Quando o servidor alvo for encontrado, o cliente se conecta automaticamente, envia os dados JSON e desconecta.

4. **Logs:**
   - Conecte o monitor serial para visualizar logs como:
     - Dispositivos escaneados.
     - Mensagens de conexão/desconexão.
     - Dados enviados ao servidor BLE.

---

## 📨 Use Cases

- **Sensores IoT**: Enviar dados de sensores de temperatura e umidade para um servidor central.
- **Autenticação**: Utilizar o token JWT para validar dispositivos conectados.
- **Comunicação BLE**: Implementar interação eficiente entre dispositivos BLE.

---

## ⚠️ Avisos

- Certifique-se de que o servidor BLE está configurado com:

  - O nome `"RPi-BLE"`.
  - Um serviço com UUID `4fafc201-1fb5-459e-8fcc-c5c9c331914b`.
  - Uma característica para escrita com UUID `beb5483e-36e1-4688-b7f5-ea07361b26a8`.

- Ajuste esses valores conforme necessário para seu ambiente BLE.

---

## 🔍 Referências

- [NimBLE-Arduino GitHub](https://github.com/h2zero/NimBLE-Arduino)
- [ArduinoJson Documentation](https://arduinojson.org/)

---

## **Licença**

Este projeto é distribuído sob a licença MIT.

---
