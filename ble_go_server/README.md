# Servidor Bluetooth Low Energy (BLE) em Go

Este projeto é uma aplicação em Go que registra um **GATT (Generic Attribute Profile)** e um **advertisement BLE** no barramento do sistema (**System D-Bus**), utilizando o BlueZ, a pilha de Bluetooth do Linux.

## Descrição do Projeto

O objetivo do projeto é criar uma aplicação BLE (Bluetooth Low Energy), que define serviços e características GATT, e registrar uma propaganda Bluetooth (advertisement BLE) que pode interagir com dispositivos externos via Bluetooth.

O programa realiza os seguintes passos principais:

1. Conecta ao **System D-Bus**, usado para comunicação entre processos no sistema operacional Linux.
2. Inicializa objetos relacionados ao BLE, incluindo:
   - Aplicação GATT.
   - Características e serviços GATT.
   - Advertisement BLE.
3. Exporta os objetos e suas interfaces para o D-Bus, permitindo introspecção e comunicação.
4. Registra o advertisement BLE no `LEAdvertisingManager1`, responsável por configurar propaganda Bluetooth.
5. Registra a aplicação GATT no `GattManager1`, permitindo que serviços e características sejam consultados e interagidos por dispositivos externos.

### **Principais Componentes**

1. **Application GATT**

   - A aplicação BLE que gerencia serviços e características GATT para dispositivos externos.
   - Utiliza a interface `org.freedesktop.DBus.ObjectManager`.

2. **Service GATT**

   - Define serviços BLE, agrupando características relacionadas.
   - Registrado utilizando a interface `org.bluez.GattService1`.

3. **Characteristic GATT**

   - Representa atributos BLE individuais que podem ser lidos, escritos ou enviados via notificações.
   - Exemplo de característica: Nível de bateria, temperatura, etc.
   - Registrado com a interface `org.bluez.GattCharacteristic1`.

4. **Advertisement BLE**

   - O advertisement é responsável pela "propaganda" do dispositivo, anunciando informações e serviços BLE disponíveis.
   - Configurado e registrado usando as interfaces `org.bluez.LEAdvertisement1` e `LEAdvertisingManager1`.

5. **Interfaces de Introspecção**
   - Permite outros processos (no D-Bus) descobrirem os métodos e as propriedades disponíveis para os objetos exportados.

---

## **Como Funciona**

### Passo a Passo do Funcionamento

1. **Conexão com o System D-Bus**
   O programa estabelece uma conexão com o barramento do sistema usando `dbus.SystemBus()`. Caso ocorra um erro, o processo é encerrado.

2. **Inicialização de Objetos**
   São criados objetos para representar a aplicação, serviços GATT, características GATT e o advertisement BLE.

3. **Exportação de Objetos**
   Os objetos criados são registrados no barramento D-Bus com os caminhos apropriados. Isso torna suas interfaces e métodos disponíveis para outros dispositivos ou aplicativos.

4. **Registro do Advertisement BLE**
   A propaganda do dispositivo Bluetooth é registrada no `LEAdvertisingManager1`, permitindo que dispositivos Bluetooth externos detectem o dispositivo.

5. **Registro da Aplicação GATT**
   A aplicação BLE é registrada no `GattManager1`, expondo serviços e características para dispositivos externos.

6. **Manutenção do Processo Ativo**
   O processo principal utiliza `select {}` para permanecer ativo e continuar respondendo a solicitações enquanto a aplicação está em execução.

---

## **Pré-requisitos**

- Go 1.20+ (ou versão compatível).
- Biblioteca **`github.com/godbus/dbus/v5`** para comunicação com o D-Bus.
- BlueZ configurado e rodando no sistema operacional Linux.
- Permissões de acesso ao D-Bus e ao adaptador Bluetooth.

---

## **Como Executar**

1. Certifique-se de que o adaptador Bluetooth está ativado e que o _BlueZ_ está rodando no sistema.

```bash
hciconfig hci0
```

✅ Saída esperada:

```
hci0:	Type: Primary  Bus: USB
	BD Address: XX:XX:XX:XX:XX:XX
	UP RUNNING
```

Se o adaptador estiver `DOWN`, ative com:

```bash
sudo hciconfig hci0 up
```

---

1.1. Verificar o serviço BlueZ (bluetoothd)

```bash
sudo systemctl status bluetooth
```

✅ Saída esperada:

```
● bluetooth.service - Bluetooth service
   Loaded: loaded (/lib/systemd/system/bluetooth.service; enabled)
   Active: active (running)
```

Se o serviço não estiver ativo:

```bash
sudo systemctl enable bluetooth
sudo systemctl start bluetooth
```

---

1.2. Confirmar se o BlueZ expõe as interfaces BLE no D-Bus

```bash
busctl introspect org.bluez /org/bluez/hci0
```

Você deve ver as seguintes interfaces:

- `org.bluez.Adapter1`
- `org.bluez.GattManager1`
- `org.bluez.LEAdvertisingManager1`

---

1.3. Inspecionar adaptadores e status via bluetoothctl

```bash
bluetoothctl
[bluetooth]# list
[bluetooth]# show
```

Esses comandos mostram o status do adaptador, nome, poder de transmissão e status do controlador.

2. Instale as dependências, compile e execute:
   ```bash
   chmod +x install-dependencies.sh
   ./install-dependencies
   ```
3. Se a execução for bem-sucedida, você verá as mensagens:
   - `📡 Advertisement BLE registrado com sucesso.`
   - `✅ Aplicação BLE com GATT registrada no D-Bus`.

---

## **Pontos Importantes**

### Introspecção

O programa configura as interfaces de introspecção (`org.freedesktop.DBus.Introspectable`) para permitir que outros objetos possam descobrir métodos, propriedades e informações disponíveis no D-Bus.

### Características GATT

As características são configuradas com:

- **UUID**: Identificador único.
- **Propriedades e Métodos**:
  - `WriteValue`: Método que permite a escrita de dados na característica.
  - `Flags`: Propriedades que definem o comportamento da característica (ex.: `read`, `write`, `notify`).

---

## **Mensagens de Erro**

- Se houver erro na conexão com o **System D-Bus**, uma mensagem como esta será exibida:
  ```
  Erro ao conectar ao system bus: <ERRO>
  ```
- Caso haja falha ao registrar o advertisement BLE:
  ```
  Erro ao registrar advertisement: <ERRO>
  ```
- Se o registro da aplicação GATT falhar:
  ```
  ❌ Erro ao registrar aplicação GATT: <ERRO>
  ```

---

## **Licença**

Este projeto é distribuído sob a licença MIT.

---
