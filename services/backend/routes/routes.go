package routes

import (
	// "context"
	// "encoding/json"
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"kayori.io/backend/controllers"
	"kayori.io/backend/models"
)

func RegisterRoutes(app *fiber.App, db *mongo.Database, rdb *redis.Client, kafkaWriter *kafka.Writer, pubsub *redis.PubSub) {
	app.Post("/api/jobs", func(c *fiber.Ctx) error {
		var job models.Job
		if err := c.BodyParser(&job); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := job.Create(context.Background(), db, "jobs", &job); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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
		return c.JSON(job)
	})

	app.Get("/api/jobs", func(c *fiber.Ctx) error {
		jobs := make([]models.Job, 0)
		collection := db.Collection("jobs")
		cursor, err := collection.Find(context.Background(), bson.D{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer cursor.Close(context.Background())

		if err := cursor.All(context.Background(), &jobs); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(jobs)
	})

	app.Get("/api/jobs/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid job ID",
			})
		}
		collection := db.Collection("jobs")
		var result bson.M
		if err := collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		log.Printf("result: %+v\n", result["_id"])
		var job models.Job
		log.Printf("job ID: %+v\n", id)
		if err := job.Read(context.Background(), db, "jobs", bson.M{"_id": objID}, &job); err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Job not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(job)
	})

	app.Get("/api/jobs/:id/artifacts", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var artifacts = make([]bson.D, 0)
		collection := db.Collection("artifacts")
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid job ID",
			})
		}
		cursor, err := collection.Find(context.Background(), bson.M{"job_id": objID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := cursor.All(context.Background(), &artifacts); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(artifacts)
	})

	app.Post("/api/entities/news_articles", func(c *fiber.Ctx) error {
		var article models.NewsArticle
		if err := c.BodyParser(&article); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		article.EntityType = "news_article"
		var existingArticle models.NewsArticle
		if err := existingArticle.Read(context.Background(), db, "artifacts", bson.M{"checksum": article.Checksum}, &existingArticle); err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Article already exists",
			})
		}

		if err := article.Create(context.Background(), db, "artifacts", &article); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		articleData, err := json.Marshal(article)
		if err != nil {
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
		articles := make([]models.NewsArticle, 0)
		if err := (&models.NewsArticle{}).ReadAll(context.Background(), db, "artifacts", bson.M{"entity_type": "news_article"}, &articles); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(articles)
	})

	app.Get("/api/ws", controllers.WebSocketHandler(pubsub))
	app.Get("/api/ws", controllers.WebSocketConnection(pubsub))
}
