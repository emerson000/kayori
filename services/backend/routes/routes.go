package routes

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"kayori.io/backend/controllers"
	"kayori.io/backend/models"
)

func RegisterRoutes(app *fiber.App, session *gocql.Session, rdb *redis.Client, kafkaWriter *kafka.Writer, pubsub *redis.PubSub) {
	app.Get("/api/users", func(c *fiber.Ctx) error {
		if err := rdb.Publish(context.Background(), "drone-status", "it worked").Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.SendString("Hello, world")
	})

	app.Post("/api/task", func(c *fiber.Ctx) error {
		var job models.Job
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

		// Insert job into jobs table
		if err := session.Query(`
			INSERT INTO jobs (job_id, title, task, service)
			VALUES (?, ?, ?, ?)`,
			job.ID, job.Title, job.Task, job.Service).Exec(); err != nil {
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

	app.Get("/api/jobs", func(c *fiber.Ctx) error {
		var jobs []models.Job
		iter := session.Query("SELECT job_id, title, task, service FROM jobs").Iter()
		var job models.Job
		for iter.Scan(&job.ID, &job.Title, &job.Task, &job.Service) {
			jobs = append(jobs, job)
		}
		if err := iter.Close(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(jobs)
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

	app.Get("/api/ws", controllers.WebSocketHandler(pubsub))
	app.Get("/api/ws", controllers.WebSocketConnection(pubsub))
}
