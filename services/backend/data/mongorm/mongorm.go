package mongorm

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to MongoDB")
	return client, nil
}

func (m *Model) Create(ctx context.Context, db *mongo.Database, collectionName string, model interface{}) error {
	collection := db.Collection(collectionName)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	res, err := collection.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	m.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (m *Model) Read(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}, result interface{}) error {
	collection := db.Collection(collectionName)

	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) ReadAll(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}, results interface{}, opts *options.FindOptionsBuilder) error {
	collection := db.Collection(collectionName)

	var cursor *mongo.Cursor
	var err error
	if opts != nil {
		cursor, err = collection.Find(ctx, filter, opts)
	} else {
		cursor, err = collection.Find(ctx, filter)
	}
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, results); err != nil {
		return err
	}

	return nil
}

func (m *Model) Update(ctx context.Context, db *mongo.Database, collectionName string, update interface{}) error {
	collection := db.Collection(collectionName)

	m.UpdatedAt = time.Now()
	updateMap, ok := update.(bson.M)
	if !ok {
		return fmt.Errorf("update is not of type bson.M")
	}
	setMap, ok := updateMap["$set"].(bson.M)
	if !ok {
		setMap = bson.M{}
		updateMap["$set"] = setMap
	}
	delete(setMap, "id")
	delete(setMap, "_id")
	delete(setMap, "created_at")
	setMap["updated_at"] = m.UpdatedAt
	_, err := collection.UpdateOne(ctx, bson.M{"_id": m.ID}, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) Delete(ctx context.Context, db *mongo.Database, collectionName string) error {
	collection := db.Collection(collectionName)
	filter := bson.M{"_id": m.ID}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
