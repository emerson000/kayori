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
	app.Post("/api/jobs", func(c *fiber.Ctx) error {
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
		return c.JSON(fiber.Map{
			"message": "Job published to Kafka",
		})
	})

	app.Get("/api/jobs", func(c *fiber.Ctx) error {
		jobs := make([]models.Job, 0)
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
			SELECT artifact_id, timestamp FROM artifact_job_lookup WHERE job_id = ?`,
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
		articles := make([]models.NewsArticle, 0)
		for _, artifact := range artifactIDs {
			var articleData []byte
			var serviceType, serviceID string
			if err := session.Query(`
				SELECT service_type, service_id FROM artifact_service_lookup WHERE artifact_id = ?`,
				artifact.ID).Scan(&serviceType, &serviceID); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			if err := session.Query(`
				SELECT data FROM artifacts WHERE service_type = ? AND service_id = ? AND timestamp = ? AND artifact_id = ?`,
				serviceType, serviceID, artifact.Timestamp, artifact.ID).Scan(&articleData); err != nil {
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
		var article models.NewsArticle
		if err := c.BodyParser(&article); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		articleTime := time.Unix(article.Timestamp, 0)
		article.ArtifactID = gocql.TimeUUID()
		var existingID gocql.UUID
		if err := session.Query(`
			SELECT artifact_id FROM artifact_checksum_lookup WHERE checksum = ?`,
			article.Checksum).Scan(&existingID); err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "The news article already exists.",
			})
		}
		articleData, err := json.Marshal(article)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := session.Query(`
			INSERT INTO artifacts (artifact_id, service_type, service_id, checksum, timestamp, data)
			VALUES (?, 'rss', ?, ?, ?, ?)`,
			article.ArtifactID, article.ServiceID, article.Checksum, articleTime, articleData).Exec(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := session.Query(`
			INSERT INTO artifact_checksum_lookup (checksum, artifact_id)
			VALUES (?, ?)`,
			article.Checksum, article.ArtifactID).Exec(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := session.Query(`
			INSERT INTO artifact_job_lookup (artifact_id, job_id, timestamp)
			VALUES (?, ?, ?)`,
			article.ArtifactID, article.JobId, article.Timestamp).Exec(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := session.Query(`
			INSERT INTO artifact_service_lookup (artifact_id, service_type, service_id)
			VALUES (?, 'rss', ?)`,
			article.ArtifactID, article.ServiceID).Exec(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := rdb.Publish(context.Background(), "drone-status", articleData).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(article)
	})

	app.Get("/api/entities/news_articles", func(c *fiber.Ctx) error {
		placeholder := []models.NewsArticle{
			{
				ArtifactID:  gocql.UUIDFromTime(time.Now()),
				Title:       "7.6-magnitude earthquake shakes Caribbean islands",
				Description: "<img src='https://i.cbc.ca/1.7454593.1739068555!/fileImage/httpImage/image.jpg_gen/derivatives/16x9_620/earthquake-map-cayman-islands.jpg' alt='A digital map shows countries and islands in the Caribbean.' width='620' height='349' title=''/><p>A magnitude-7.6 earthquake shook the Caribbean Sea south of the Cayman Islands Saturday, according to the U.S. Geological Survey. Several islands and countries urged people near the coastline to move inland but authorities in most places later lifted the tsunami alerts.</p>",
				URL:         "https://www.cbc.ca/news/world/earthquake-caribbean-islands-1.7454589?cmp=rss",
				Published:   "2025-02-08T21:52:10Z",
				Timestamp:   1739051530000,
				Author:      "",
				Categories:  []string{"News/World"},
				ServiceID:   "https://www.cbc.ca/webfeed/rss/rss-world",
				Checksum:    "c79c1bc7ea60a4cadaa4f8b99441eefe0da0d6272e609bb9bea91daf5f321dce",
				JobId:       "0194f26f-ec41-7007-94c7-397ae413c91b",
			},
		}
		return c.JSON(placeholder)
	})

	app.Get("/api/ws", controllers.WebSocketHandler(pubsub))
	app.Get("/api/ws", controllers.WebSocketConnection(pubsub))
}
