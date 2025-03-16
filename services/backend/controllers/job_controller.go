package controllers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"kayori.io/backend/models"
)

func CreateJob(db *mongo.Database, kafkaWriter *kafka.Writer) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		if job.Status == "" {
			job.Status = "pending"
		}

		jobData, err := json.Marshal(job)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = kafkaWriter.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(job.Service),
				Value: jobData,
			},
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(job)
	}
}

func GetJobs(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var projectId bson.ObjectID
		if err := GetObjectIdFromParam(c, "project_id", &projectId); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		jobs := make([]models.Job, 0)
		category := c.Query("category")
		collection := db.Collection("jobs")
		findOptions := options.Find().SetSort(bson.D{{Key: "title", Value: 1}})
		AddPaginationToFindOptions(c, findOptions)
		filters := bson.M{
			"projects": projectId,
		}
		if category != "" {
			filters["category"] = category
		}
		cursor, err := collection.Find(context.Background(), filters, findOptions)
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
		SetHasMoreHeader(c, len(jobs))
		return c.JSON(jobs)
	}
}

func GetJobByID(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var projectObjectId bson.ObjectID
		if err := GetObjectIdFromParam(c, "project_id", &projectObjectId); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		id := c.Params("id")
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid job ID",
			})
		}
		collection := db.Collection("jobs")
		var result bson.M
		if err := collection.FindOne(context.TODO(), bson.M{"_id": objID, "projects": projectObjectId}).Decode(&result); err != nil {
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
	}
}

func UpdateJob(db *mongo.Database, requireProject bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var projectObjectId bson.ObjectID
		if requireProject {
			if err := GetObjectIdFromParam(c, "project_id", &projectObjectId); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		}
		var job models.Job
		if err := c.BodyParser(&job); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		id := c.Params("id")
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid job ID",
			})
		}
		job.ID = objID
		updateStatement := bson.M{
			"$set": bson.M{},
		}
		if requireProject {
			updateStatement["$addToSet"].(bson.M)["projects"] = projectObjectId
		}
		if job.Title != "" {
			updateStatement["$set"].(bson.M)["title"] = job.Title
		}
		if job.Status != "" {
			updateStatement["$set"].(bson.M)["status"] = job.Status
		}
		if err := job.Update(context.Background(), db, "jobs", updateStatement); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		err = job.Read(context.Background(), db, "jobs", bson.M{"_id": job.ID}, &job)
		if err != nil {
			return err
		}
		return c.JSON(job)
	}
}

func GetJobArtifacts(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var projectObjectId bson.ObjectID
		if err := GetObjectIdFromParam(c, "project_id", &projectObjectId); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		id := c.Params("id")
		var artifacts = make([]bson.D, 0)
		collection := db.Collection("artifacts")
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid job ID",
			})
		}
		filters := bson.M{
			"job_id":  objID,
			"deleted": bson.M{"$ne": true},
		}
		findOptions := options.Find().SetSort(bson.M{
			"created_at": -1,
		})
		AddPaginationToFindOptions(c, findOptions)
		cursor, err := collection.Find(context.Background(), filters, findOptions)
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

		SetHasMoreHeader(c, len(artifacts))

		for i := range artifacts {
			for j := range artifacts[i] {
				if artifacts[i][j].Key == "_id" {
					artifacts[i][j].Key = "id"
				}
			}
		}

		return c.JSON(artifacts)
	}
}
