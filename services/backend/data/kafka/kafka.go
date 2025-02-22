package kafka

import (
	"sync"

	"github.com/segmentio/kafka-go"
)

var (
	writer *kafka.Writer
	once   sync.Once
)

func GetKafkaWriter() *kafka.Writer {
	once.Do(func() {
		writer = kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{"kafka:29092"},
			Topic:   "drone-dispatch",
		})
	})
	return writer
}

func CloseKafkaWriter() error {
	if writer == nil {
		return nil
	}
	return writer.Close()
}
