// kafka.go
package services

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"ble_go_server/configs"
)

var KafkaProducer *kafka.Writer

func InitKafkaProducer() {
	KafkaProducer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{configs.KAFKA_BROKER},
		Topic:    configs.KAFKA_TOPIC,
		Balancer: &kafka.LeastBytes{},
	})
}

func SendToKafka(data map[string]interface{}) error {
	msgBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return KafkaProducer.WriteMessages(context.Background(), kafka.Message{Value: msgBytes})
}