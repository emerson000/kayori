package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"kayori.io/backend/models"
)

func CreateCluster(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cluster models.Cluster
		if err := c.BodyParser(&cluster); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Check if any artifact ID already exists
		for _, artifactID := range cluster.Artifacts {
			var existingCluster models.Cluster
			if err := existingCluster.Read(context.Background(), db, "clusters", bson.M{"artifacts": artifactID}, &existingCluster); err == nil {
				// Append new artifact IDs to the existing cluster
				artifactMap := make(map[bson.ObjectID]bool)
				for _, id := range existingCluster.Artifacts {
					artifactMap[id] = true
				}
				for _, id := range cluster.Artifacts {
					artifactMap[id] = true
				}
				uniqueArtifacts := make([]bson.ObjectID, 0, len(artifactMap))
				for id := range artifactMap {
					uniqueArtifacts = append(uniqueArtifacts, id)
				}
				existingCluster.Artifacts = uniqueArtifacts
				updateStatement := bson.M{
					"$set": bson.M{
						"artifacts": existingCluster.Artifacts,
					},
				}
				if err := existingCluster.Update(context.Background(), db, "clusters", updateStatement); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
				return c.Status(fiber.StatusOK).JSON(existingCluster)
			}
		}

		// Create a new cluster if no existing artifact ID is found
		if err := cluster.Create(context.Background(), db, "clusters", &cluster); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(cluster)
	}
}

func GetCluster(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		objectId, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}
		var cluster models.Cluster
		if err := cluster.Read(context.Background(), db, "clusters", bson.M{"_id": objectId}, &cluster); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(cluster)
	}
}

func CreateClusterIndex(db *mongo.Database) fiber.Handler {
	return BuildCreateIndexRoute(db, "clusters")
}
