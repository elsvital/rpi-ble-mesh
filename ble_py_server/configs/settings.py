# Chave secreta usada para decodificar o JWT
SECRET_KEY = "chave-secreta"  # deve ser a mesma usada no cliente

# UUIDs do serviço e característica BLE
TEST_SERVICE_UUID = '4fafc201-1fb5-459e-8fcc-c5c9c331914b'
TEST_CHARACTERISTIC_UUID = 'beb5483e-36e1-4688-b7f5-ea07361b26a8'

# Configuração do Kafka
KAFKA_BOOTSTRAP_SERVERS = "192.168.15.160:9092"
KAFKA_TOPIC = "iot-data"