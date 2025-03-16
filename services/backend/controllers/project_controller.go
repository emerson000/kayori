package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"kayori.io/backend/models"
)

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
		updateStatement := bson.M{"$set": project}
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
