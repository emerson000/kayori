package controllers

import (
	"context"
	"encoding/json"
	"fmt"

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
		search := c.Query("search")
		columns := c.Query("columns")
		limit := c.QueryInt("limit", 10)
		skip := (page - 1) * limit
		articles := make([]models.NewsArticle, 0)
		findOptions := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
		findOptions = findOptions.SetLimit(int64(limit)).SetSkip(int64(skip))
		filters := bson.M{"entity_type": "news_article", "deleted": bson.M{"$ne": true}}
		if search != "" {
			filters["$text"] = bson.M{"$search": search}
		}
		if columns != "" {
			var columnsArray []string
			if err := json.Unmarshal([]byte(columns), &columnsArray); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid JSON array in columns parameter",
				})
			}
			projection := bson.M{}
			for _, column := range columnsArray {
				projection[column] = 1
			}
			findOptions = findOptions.SetProjection(projection)
			var results []bson.M
			cursor, err := db.Collection("artifacts").Find(context.Background(), filters, findOptions)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			if err := cursor.All(context.Background(), &results); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			hasMore := len(results) == limit
			c.Set("X-Has-More", fmt.Sprintf("%v", hasMore))
			return c.JSON(results)
		}
		if err := (&models.NewsArticle{}).ReadAll(context.Background(), db, "artifacts", filters, &articles, findOptions); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		hasMore := len(articles) == limit
		c.Set("X-Has-More", fmt.Sprintf("%v", hasMore))
		return c.JSON(articles)
	}
}
