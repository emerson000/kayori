package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"kayori.io/backend/data/mongorm"
)

type Job struct {
	mongorm.Model `bson:",inline"`
	Title         string          `bson:"title" json:"title"`
	Service       string          `bson:"service" json:"service"`
	Category      string          `bson:"category" json:"category"`
	Status        string          `bson:"status" json:"status"`
	Task          bson.D          `bson:"task" json:"task"`
	Projects      []bson.ObjectID `bson:"projects" json:"projects"`
	Schedule      *struct {
		Schedule bool      `bson:"schedule" json:"schedule"`
		Duration int       `bson:"duration" json:"duration"`
		Interval string    `bson:"interval" json:"interval"`
		LastRun  time.Time `bson:"last_run" json:"last_run"`
	} `bson:"schedule,omitempty" json:"schedule,omitempty"`
}
