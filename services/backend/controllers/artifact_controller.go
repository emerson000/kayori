package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
		if err := artifact.Delete(context.Background(), db, "artifacts"); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Artifact deleted successfully",
		})
	}
}
