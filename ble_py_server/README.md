# Servidor Bluetooth Low Energy (BLE) em Python

Este projeto implementa um servidor BLE usando o **D-Bus**, **BlueZ** (a pilha oficial Bluetooth para Linux) e **GObject**. Ele registra serviços e características GATT personalizadas e configura anúncios BLE para ser detectado por dispositivos BLE próximos.

## Funcionalidades

1. **Configuração do Adaptador Bluetooth**:

   - Ativa o adaptador Bluetooth (ex.: `hci0`) e o prepara no modo BLE.

2. **Registro de Serviços GATT**:

   - Define serviços GATT com base no UUID configurado.
   - Registra características GATT que podem ser lidas, gravadas ou notificadas.

3. **Anúncio de Disponibilidade**:

   - Configura e registra anúncios BLE, tornando o dispositivo detectável por outros dispositivos BLE.

4. **Manutenção de um Loop de Eventos Assíncrono**:
   - Usa **GLib** para gerenciar conexões BLE e interações com características em tempo de execução.

## Estrutura do Código

- **Adaptador Bluetooth**:
  O adaptador é configurado para habilitar o modo BLE:

  ```python
  adapter_path = "/org/bluez/hci0"
  adapter.Set("org.bluez.Adapter1", "Powered", dbus.Boolean(1))
  ```

- **Aplicativo GATT**:
  Um aplicativo é registrado com serviços e características, permitindo interações com dispositivos BLE.
  Exemplo de criação de um serviço GATT e registro:

  ```python
  app = Application(bus)
  service = Service(bus, 0, TEST_SERVICE_UUID, True)
  characteristic = Characteristic(bus, 0, TEST_CHARACTERISTIC_UUID, ['write'], service)
  service.characteristics.append(characteristic)
  app.add_service(service)
  ```

- **Registro do Serviço GATT**:
  Após configurar serviços e características, o aplicativo é registrado com `GattManager1`:

  ```python
  service_manager.RegisterApplication(app.get_path(), {}, ...)
  ```

- **Anúncios BLE**:
  O objeto `Advertisement` é usado para configurar anúncios:
  ```python
  adv = Advertisement(bus, 0, "peripheral")
  ad_manager.RegisterAdvertisement(adv.get_path(), {}, ...)
  ```

## Como Executar

1. Certifique-se de que o **BlueZ** está instalado e funcionando no sistema:

   ```bash
   systemctl start bluetooth
   ```

2. Execute o programa:

   ```bash
   python3 main.py
   ```

3. Veja as mensagens de status no terminal:

   - `✅ Serviço GATT registrado`: O serviço foi registrado com sucesso.
   - `📡 Anúncio BLE registrado`: O dispositivo está anunciando sua presença.
   - `❌` Erros serão exibidos em casos de falha.

4. Interrompa o servidor com `Ctrl + C`.

## Dependências

- **Python** 3.11.2
- **BlueZ**
- Bibliotecas Python:
  - `dbus`
  - `gi.repository`

Instale as dependências com:

```bash
pip install dbus-python PyGObject
```

## Fluxo do Programa

1. Inicializa o adaptador Bluetooth (`hci0`).
2. Registra um serviço GATT com características.
3. Inicia anúncios BLE para dispositivos próximos.
4. Mantém o loop principal GLib para interagir com os clientes BLE.

## Exemplo de Caso de Uso

Este projeto é útil para dispositivos de IoT, como:

- Sensores de temperatura.
- Monitores de frequência cardíaca.
- Rastreamento de proximidade BLE.

Configure o `TEST_SERVICE_UUID` e `TEST_CHARACTERISTIC_UUID` em `configs/settings.py` para personalizar o serviço e a característica.

## **Licença**

Este projeto é distribuído sob a licença MIT.

---
