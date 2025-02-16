package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"kayori.io/backend/data/mongorm"
	"kayori.io/backend/routes"
)

func main() {
	app := fiber.New()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	client, err := mongorm.Connect("mongodb://root:kayori@mongo:27017")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	db := client.Database("kayori")

	// Subscribe to the "drone-status" topic
	pubsub := rdb.Subscribe(context.Background(), "drone-status")
	defer pubsub.Close()

	// Kafka writer
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:29092"},
		Topic:   "drone-dispatch",
	})
	defer kafkaWriter.Close()

	// Register routes
	routes.RegisterRoutes(app, db, rdb, kafkaWriter, pubsub)

	app.Listen(":3001")
}
