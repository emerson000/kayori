package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Model struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

type NewsArticle struct {
	Model       `bson:",inline"`
	EntityType  string        `bson:"entity_type" json:"entity_type"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	URL         string        `bson:"url" json:"url"`
	Published   string        `bson:"published" json:"published"`
	Timestamp   time.Time     `bson:"timestamp" json:"timestamp"`
	Author      string        `bson:"author" json:"author"`
	Categories  []string      `bson:"categories" json:"categories"`
	Service     string        `bson:"service" json:"service"`
	ServiceID   string        `bson:"service_id" json:"service_id"`
	Checksum    string        `bson:"checksum" json:"checksum"`
	JobId       bson.ObjectID `bson:"job_id" json:"job_id"`
}
