package controllers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
		if cluster.ID.IsZero() {
			log.Println("Checking for existing cluster with artifacts")
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
					if len(cluster.Centroid) > 0 {
						existingCluster.Centroid = cluster.Centroid
						updateStatement["$set"].(bson.M)["centroid"] = cluster.Centroid
					}
					if err := existingCluster.Update(context.Background(), db, "clusters", updateStatement); err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"error": err.Error(),
						})
					}
					updateArtifacts(db, existingCluster.Artifacts, existingCluster.ID)
					return c.Status(fiber.StatusOK).JSON(existingCluster)
				}
			}
			// Create a new cluster if no existing artifact ID is found
			log.Println("Creating new cluster")
			if err := cluster.Create(context.Background(), db, "clusters", &cluster); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			updateArtifacts(db, cluster.Artifacts, cluster.ID)
			return c.Status(fiber.StatusCreated).JSON(cluster)
		}
		log.Println("Updating existing cluster")
		updateStatement := bson.M{
			"$addToSet": bson.M{
				"artifacts": bson.M{
					"$each": cluster.Artifacts,
				},
			},
		}
		if err := cluster.Update(context.Background(), db, "clusters", updateStatement); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		updateArtifacts(db, cluster.Artifacts, cluster.ID)
		return c.Status(fiber.StatusOK).JSON(cluster)
	}
}

func updateArtifacts(db *mongo.Database, artifacts []bson.ObjectID, clusterId bson.ObjectID) error {
	for _, id := range artifacts {
		artifactUpdate := bson.M{
			"$set": bson.M{
				"cluster_id": clusterId,
			},
		}
		model := models.Artifact{}
		model.ID = id
		if err := model.Update(context.Background(), db, "artifacts", artifactUpdate); err != nil {
			return err
		}
		log.Printf("Updated artifact with ID %s to cluster ID %s", id.Hex(), clusterId.Hex())
	}
	return nil
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

func GetClusters(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		clusters := make([]models.Cluster, 0)
		collection := db.Collection("clusters")
		findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
		filters := bson.M{}
		cursor, err := collection.Find(context.Background(), filters, findOptions)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer cursor.Close(context.Background())

		if err := cursor.All(context.Background(), &clusters); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(clusters)
	}
}
