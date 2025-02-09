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
	"kayori.io/backend/models"
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

	app.Get("/api/jobs/:id/artifacts", func(c *fiber.Ctx) error {
		id := c.Params("id")

		var artifactIDs []struct {
			ID        gocql.UUID
			Timestamp int64
		}
		iter := session.Query(`
			SELECT artifact_id, timestamp FROM collection_job_lookup WHERE job_id = ?`,
			id).Iter()
		var artifact struct {
			ID        gocql.UUID
			Timestamp int64
		}
		for iter.Scan(&artifact.ID, &artifact.Timestamp) {
			artifact.Timestamp *= 1000
			artifactIDs = append(artifactIDs, artifact)
		}
		if err := iter.Close(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var articles []models.NewsArticle
		for _, artifact := range artifactIDs {
			var articleData []byte
			var sourceType, sourceID string
			if err := session.Query(`
				SELECT source_type, source_id FROM collection_lookup WHERE artifact_id = ?`,
				artifact.ID).Scan(&sourceType, &sourceID); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			log.Printf("Source Type: %v", sourceType)
			log.Printf("Source ID: %v", sourceID)
			log.Printf("Timestamp: %v", artifact.Timestamp)
			log.Printf("Artifact ID: %v", artifact.ID)
			if err := session.Query(`
				SELECT data FROM collection_store WHERE source_type = ? AND source_id = ? AND timestamp = ? AND artifact_id = ?`,
				sourceType, sourceID, artifact.Timestamp, artifact.ID).Scan(&articleData); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			var article models.NewsArticle
			if err := json.Unmarshal(articleData, &article); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			articles = append(articles, article)
		}

		return c.JSON(articles)
	})

	app.Post("/api/rss", func(c *fiber.Ctx) error {
		type RssJob struct {
			Title string   `json:"title"`
			Urls  []string `json:"urls"`
		}

		return c.SendString("Hello, world")
	})

	app.Post("/api/news_article", func(c *fiber.Ctx) error {
		var article models.NewsArticle // Use the shared model
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

		// Insert job ID into collection_job_lookup table
		if err := session.Query(`
			INSERT INTO collection_job_lookup (artifact_id, job_id, timestamp)
			VALUES (?, ?, ?)`,
			article.ArtifactID, article.JobId, article.Timestamp).Exec(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Insert artifact_id, source_type, and source_id into collection_lookup table
		if err := session.Query(`
			INSERT INTO collection_lookup (artifact_id, source_type, source_id)
			VALUES (?, 'rss', ?)`,
			article.ArtifactID, article.SourceID).Exec(); err != nil {
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
