package controllers

import (
	"context"
	"errors"
	"fmt"

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

func GetObjectIdFromParam(c *fiber.Ctx, param string, out *bson.ObjectID) error {
	id := c.Params(param)
	if id == "" {
		return errors.New("A " + param + " ID is required")
	}
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("Invalid " + param + " ID")
	}
	*out = objID
	return nil
}

func AddPaginationToFindOptions(c *fiber.Ctx, findOptions *options.FindOptionsBuilder) {
	_, limit, skip := GetPaginationParams(c)
	findOptions.SetSkip(skip).SetLimit(int64(limit))
}

func GetPaginationParams(c *fiber.Ctx) (int, int, int64) {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	skip := (page - 1) * limit
	return page, limit, int64(skip)
}

func SetHasMoreHeader(c *fiber.Ctx, total int) {
	_, limit, _ := GetPaginationParams(c)
	hasMore := total == limit
	c.Set("X-Has-More", fmt.Sprintf("%v", hasMore))
}

func CreateUpdateStatement(model interface{}) (bson.M, error) {
	modelMap := bson.M{}
	modelBytes, err := bson.Marshal(model)
	if err != nil {
		return nil, err
	}
	bson.Unmarshal(modelBytes, &modelMap)
	delete(modelMap, "id")
	delete(modelMap, "_id")
	delete(modelMap, "created_at")
	statement := bson.M{
		"$set": modelMap,
	}
	return statement, nil
}
