package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"kayori.io/backend/controllers"
)

func RegisterRoutes(app *fiber.App, db *mongo.Database, rdb *redis.Client, kafkaWriter *kafka.Writer, pubsub *redis.PubSub) {
	app.Post("/api/jobs", controllers.CreateJob(db, kafkaWriter))
	app.Get("/api/jobs", controllers.GetJobs(db))
	app.Get("/api/jobs/:id", controllers.GetJobByID(db))
	app.Put("/api/jobs/:id", controllers.UpdateJob(db))
	app.Get("/api/jobs/:id/artifacts", controllers.GetJobArtifacts(db))

	app.Post("/api/entities/news_articles", controllers.CreateNewsArticle(db, rdb))
	app.Get("/api/entities/news_articles", controllers.GetNewsArticles(db))

	app.Delete("/api/artifacts/:id", controllers.DeleteArtifact(db))

	app.Post("/api/artifacts/indexes", controllers.CreateArtifactIndex(db))
	app.Get("/api/artifacts/indexes", controllers.GetArtifactIndexes(db))
	app.Delete("/api/artifacts/indexes/:name", controllers.DeleteArtifactIndex(db))

	app.Post("/api/clusters", controllers.CreateCluster(db))
	app.Get("/api/clusters/:id", controllers.GetCluster(db))
	app.Post("/api/clusters/indexes", controllers.CreateClusterIndex(db))

	app.Get("/api/ws", controllers.WebSocketHandler(pubsub))
	app.Get("/api/ws", controllers.WebSocketConnection(pubsub))
}
