package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
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

	app.Get("/api/users", func(c *fiber.Ctx) error {

		if err := rdb.Publish(context.Background(), "drone-status", "it worked").Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendString("Hello, world")
	})
	// Kafka writer
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:29092"},
		Topic:   "drone-dispatch",
	})
	defer kafkaWriter.Close()

	app.Post("/api/task", func(c *fiber.Ctx) error {
		type Job struct {
			ID      string          `json:"id"`
			Service string          `json:"service"`
			Task    json.RawMessage `json:"task"`
		}
		var job Job
		if err := c.BodyParser(&job); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		jobData, err := json.Marshal(job)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		err = kafkaWriter.WriteMessages(context.Background(),
			kafka.Message{
				Value: jobData,
			},
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendString("Job published to Kafka")
	})
	app.Post("/api/rss", func(c *fiber.Ctx) error {
		type RssJob struct {
			Title string   `json:"title"`
			Urls  []string `json:"urls"`
		}

		return c.SendString("Hello, world")
	})

	app.Post("/api/news_article", func(c *fiber.Ctx) error {
		type NewsArticle struct {
			ArtifactID  gocql.UUID `json:"artifact_id"`
			Title       string     `json:"title"`
			Description string     `json:"description"`
			URL         string     `json:"url"`
			Published   string     `json:"published"`
			Timestamp   int64      `json:"timestamp"`
			Author      string     `json:"author"`
			Categories  []string   `json:"categories"`
			SourceID    string     `json:"source_id"`
			Checksum    string     `json:"checksum"`
		}

		var article NewsArticle
		if err := c.BodyParser(&article); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Convert the timestamp from int64 to time.Time
		articleTime := time.Unix(article.Timestamp, 0)

		// Generate UUID automatically
		article.ArtifactID = gocql.TimeUUID()

		// Check if checksum already exists
		var existingID gocql.UUID
		if err := session.Query(`
			SELECT artifact_id FROM checksum_index WHERE checksum = ?`,
			article.Checksum).Scan(&existingID); err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "The news article already exists.",
			})
		}

		// Convert article to JSON
		articleData, err := json.Marshal(article)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := session.Query(`
			INSERT INTO collection_store (artifact_id, source_type, source_id, checksum, timestamp, data)
			VALUES (?, 'rss', ?, ?, ?, ?)`,
			article.ArtifactID, article.SourceID, article.Checksum, articleTime, articleData).Exec(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Insert checksum into checksum_index table
		if err := session.Query(`
			INSERT INTO checksum_index (checksum, artifact_id)
			VALUES (?, ?)`,
			article.Checksum, article.ArtifactID).Exec(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Publish the created article to the Redis channel
		if err := rdb.Publish(context.Background(), "drone-status", articleData).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(article)
	})

	var (
		connections = make(map[*websocket.Conn]chan []byte)
		connMutex   sync.Mutex
	)

	app.Get("/api/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/api/ws", websocket.New(func(c *websocket.Conn) {
		writeQueue := make(chan []byte, 10)

		connMutex.Lock()
		connections[c] = writeQueue
		connMutex.Unlock()

		defer func() {
			connMutex.Lock()
			delete(connections, c)
			connMutex.Unlock()
			c.Close()
		}()

		var writeMutex sync.Mutex

		// Goroutine to handle all writes to the websocket
		go func() {
			for msg := range writeQueue {
				writeMutex.Lock()
				if c == nil {
					log.Println("write: nil *Conn")
					writeMutex.Unlock()
					break
				}
				if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Println("write:", err)
					writeMutex.Unlock()
					break
				}
				writeMutex.Unlock()
			}
		}()

		// Listen for messages from Redis and broadcast them to all connections
		go func() {
			ch := pubsub.Channel()
			for msg := range ch {
				connMutex.Lock()
				for conn, queue := range connections {
					select {
					case queue <- []byte(msg.Payload):
					default:
						log.Printf("dropping message for connection: %v", conn)
					}
				}
				connMutex.Unlock()
			}
		}()

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			writeMutex.Lock()
			if c == nil {
				log.Println("write: nil *Conn")
				writeMutex.Unlock()
				break
			}
			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				writeMutex.Unlock()
				break
			}
			writeMutex.Unlock()
		}
	}))
	app.Listen(":3001")
}
