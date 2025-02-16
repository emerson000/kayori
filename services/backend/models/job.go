package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"kayori.io/backend/data/mongorm"
)

type Job struct {
	mongorm.Model `bson:",inline"`
	Title         string `bson:"title" json:"title"`
	Service       string `bson:"service" json:"service"`
	Task          bson.D `bson:"task" json:"task"`
}
