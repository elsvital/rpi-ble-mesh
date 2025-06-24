from kafka import KafkaProducer
import json
from configs.settings import KAFKA_BOOTSTRAP_SERVERS, KAFKA_TOPIC

producer = KafkaProducer(
    bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS,
    value_serializer=lambda v: json.dumps(v).encode("utf-8")
)

def send_to_kafka(data):
    try:
        producer.send(KAFKA_TOPIC, value=data)
        producer.flush()
        print("📤 Enviado ao Kafka com sucesso")
    except Exception as e:
        print("❌ Erro ao enviar para Kafka:", e)