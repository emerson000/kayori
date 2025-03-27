package controllers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"kayori.io/backend/data/mongorm"
	"kayori.io/backend/models"
)

// entityTypeRegistry maps entity type strings to their corresponding model types
var entityTypeRegistry = map[string]interface{}{
	"news_article": &models.NewsArticle{},
	"artifact":     &models.Artifact{},
}

// getEntitySlice creates a new slice of the appropriate type for the given entity
func getEntitySlice(entityType string) (interface{}, error) {
	modelType, exists := entityTypeRegistry[entityType]
	if !exists {
		return nil, fmt.Errorf("unknown entity type: %s", entityType)
	}

	sliceType := reflect.SliceOf(reflect.TypeOf(modelType).Elem())
	return reflect.New(sliceType).Interface(), nil
}

func CreateProject(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var project models.Project
		if err := c.BodyParser(&project); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := project.Create(context.Background(), db, "projects", &project); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(project)
	}
}

func GetProject(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Project ID is required",
			})
		}
		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid project ID",
			})
		}
		var project models.Project
		if err := project.Read(context.Background(), db, "projects", bson.M{"_id": objectID}, &project); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(project)
	}
}

func UpdateProject(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Project ID is required",
			})
		}
		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid project ID",
			})
		}
		var project models.Project
		if err := c.BodyParser(&project); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		project.ID = objectID

		updateStatement, err := CreateUpdateStatement(project)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := project.Update(context.Background(), db, "projects", updateStatement); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(project)
	}
}

func GetProjects(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var projects []models.Project
		err := (&models.Project{}).ReadAll(context.Background(), db, "projects", bson.M{}, &projects, nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(projects)
	}
}

func GetProjectArtifacts(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Project ID is required",
			})
		}
		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid project ID",
			})
		}
		entityType := c.Query("type")
		if entityType == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Entity type is required",
			})
		}

		// Get the appropriate model type from registry
		_, exists := entityTypeRegistry[entityType]
		if !exists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Unknown entity type: %s", entityType),
			})
		}

		// Create a slice of the appropriate type
		results, err := getEntitySlice(entityType)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var jobs []models.Job
		if err := (&models.Job{}).ReadAll(context.Background(), db, "jobs", bson.M{"projects": objectID}, &jobs, nil); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		jobIds := make([]bson.ObjectID, len(jobs))
		for i, job := range jobs {
			jobObjectId, err := bson.ObjectIDFromHex(job.ID.Hex())
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			jobIds[i] = jobObjectId
		}

		filter := bson.M{
			"job_id": bson.M{
				"$in": jobIds,
			},
			"entity_type": entityType,
		}

		findOptions := options.Find().SetSort(bson.M{"timestamp": -1})
		AddPaginationToFindOptions(c, findOptions)

		model := &mongorm.Model{}
		if err := model.ReadAll(context.Background(), db, "artifacts", filter, results, findOptions); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		resultsSlice := reflect.ValueOf(results).Elem().Interface()
		SetHasMoreHeader(c, reflect.ValueOf(resultsSlice).Len())
		return c.Status(fiber.StatusOK).JSON(resultsSlice)
	}
}

func DeleteProject(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Project ID is required",
			})
		}
		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid project ID",
			})
		}
		project := &models.Project{}
		project.ID = objectID
		if err := project.Delete(context.Background(), db, "projects"); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Project deleted successfully",
		})
	}
}
