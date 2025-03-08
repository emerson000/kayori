package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"kayori.io/backend/data/mongorm"
)

type Cluster struct {
	mongorm.Model `bson:",inline"`
	Artifacts     []bson.ObjectID `bson:"artifacts" json:"artifacts"`
}
