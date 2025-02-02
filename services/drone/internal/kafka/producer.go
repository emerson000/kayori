package kafka

import (
	"context"
	"os"

	"github.com/segmentio/kafka-go"
)

// Producer will contain logic for publishing messages (if necessary).
type Producer struct{}

func NewProducer() *Producer {
	return &Producer{}
}

func (p *Producer) SendMessage(msg string) error {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   "drone-output",
	})
	defer w.Close()

	err := w.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(msg),
	})
	if err != nil {
		return err
	}
	return nil
}
