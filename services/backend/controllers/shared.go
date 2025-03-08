package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func BuildCreateIndexRoute(db *mongo.Database, collectionName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var indexFields map[string]interface{}
		if err := c.BodyParser(&indexFields); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON",
			})
		}

		keys := bson.D{}
		for k, v := range indexFields {
			if floatValue, ok := v.(float64); ok {
				v = int(floatValue)
			}
			keys = append(keys, bson.E{Key: k, Value: v})
		}
		indexModel := mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetDefaultLanguage("english"),
		}
		collection := db.Collection(collectionName)
		if _, err := collection.Indexes().CreateOne(context.Background(), indexModel); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Index created successfully",
		})
	}
}
