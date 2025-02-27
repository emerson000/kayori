package controllers

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"kayori.io/backend/models"
)

func CreateNewsArticle(db *mongo.Database, rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

func GetNewsArticles(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		skip := (page - 1) * 10
		articles := make([]models.NewsArticle, 0)
		findOptions := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
		findOptions = findOptions.SetLimit(10).SetSkip(int64(skip))
		filters := bson.M{"entity_type": "news_article", "deleted": false}
		if err := (&models.NewsArticle{}).ReadAll(context.Background(), db, "artifacts", filters, &articles, findOptions); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(articles)
	}
}
