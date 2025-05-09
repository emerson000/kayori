package kafka

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct{}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Start(processMessage func([]byte)) error {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          []string{broker},
		GroupID:          "drone-workers",
		Topic:            "drone-dispatch",
		CommitInterval:   time.Second * 10,
		StartOffset:      kafka.FirstOffset,
		ReadBatchTimeout: 200 * time.Millisecond,
		RebalanceTimeout: 5 * time.Second,
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			return fmt.Errorf("could not read message: %v", err)
		}
		processMessage(m.Value)
	}
}
