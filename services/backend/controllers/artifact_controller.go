package controllers

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"kayori.io/backend/models"
)

func DeleteArtifact(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		artifactId := c.Params("id")
		objId, err := bson.ObjectIDFromHex(artifactId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid artifact ID",
			})
		}
		var artifact models.Artifact
		artifact.ID = objId
		updateStatement := bson.M{
			"$set": bson.M{
				"deleted": true,
			},
		}
		if err := artifact.Update(context.Background(), db, "artifacts", updateStatement); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Artifact deleted successfully",
		})
	}
}

func CreateArtifactIndex(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var indexFields map[string]string
		if err := json.Unmarshal(c.Body(), &indexFields); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON",
			})
		}

		keys := bson.D{}
		for k, v := range indexFields {
			keys = append(keys, bson.E{Key: k, Value: v})
		}
		indexModel := mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetDefaultLanguage("english"),
		}

		collection := db.Collection("artifacts")
		_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Index created successfully",
		})
	}
}

func GetArtifactIndexes(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		indexes, err := db.Collection("artifacts").Indexes().List(context.Background())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var indexList []bson.M
		if err := indexes.All(context.Background(), &indexList); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(indexList)
	}
}

func DeleteArtifactIndex(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		indexName := c.Params("name")
		if indexName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Index name is required",
			})
		}
		if indexName == "_id_" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot delete the default _id_ index",
			})
		}
		collection := db.Collection("artifacts")
		if err := collection.Indexes().DropOne(context.Background(), indexName); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Index deleted successfully",
		})
	}
}
