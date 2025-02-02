package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gocql/gocql"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
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
		return c.SendString("Hello, world")
	})

	app.Post("/api/news_article", func(c *fiber.Ctx) error {
		type NewsArticle struct {
			ID            gocql.UUID `json:"id"`
			Title         string     `json:"title"`
			Description   string     `json:"description"`
			URL           string     `json:"url"`
			Published     string     `json:"published"`
			PublishedUnix string     `json:"published_unix"`
			Author        string     `json:"author"`
			Categories    []string   `json:"categories"`
			Source        string     `json:"source"`
			Checksum      string     `json:"checksum"`
		}

		var article NewsArticle
		if err := c.BodyParser(&article); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Generate UUID automatically
		article.ID = gocql.TimeUUID()

		// Check if checksum already exists
		var existingID gocql.UUID
		if err := session.Query(`
			SELECT object_id FROM checksum_index WHERE checksum = ?`,
			article.Checksum).Scan(&existingID); err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "The news article already exists.",
			})
		}

		if err := session.Query(`
			INSERT INTO news_article (id, title, description, url, published, published_unix, author, categories, source, checksum)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			article.ID, article.Title, article.Description, article.URL, article.Published, article.PublishedUnix, article.Author, article.Categories, article.Source, article.Checksum).Exec(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Insert into checksum_index table
		if err := session.Query(`
			INSERT INTO checksum_index (checksum, object_id)
			VALUES (?, ?)`,
			article.Checksum, article.ID).Exec(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Convert article to JSON
		articleJSON, err := json.Marshal(article)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Publish the created article to the Redis channel
		if err := rdb.Publish(context.Background(), "drone-status", articleJSON).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(article)
	})

	app.Get("/api/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/api/ws", websocket.New(func(c *websocket.Conn) {
		var writeMutex sync.Mutex

		// Listen for messages from Redis and send them to the websocket
		go func() {
			ch := pubsub.Channel()
			for msg := range ch {
				writeMutex.Lock()
				if err := c.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
					log.Println("write:", err)
					writeMutex.Unlock()
					break
				}
				writeMutex.Unlock()
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
