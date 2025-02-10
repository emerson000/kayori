package main

import (
	"context"
	"log"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"kayori.io/backend/routes"
)

func main() {
	app := fiber.New()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Connect to Cassandra
	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "kayori_drone_output"
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("unable to connect to cassandra: %v", err)
	}
	defer session.Close()

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
	routes.RegisterRoutes(app, session, rdb, kafkaWriter, pubsub)

	app.Listen(":3001")
}
