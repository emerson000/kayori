package models

import (
	"kayori.io/backend/data/mongorm"
)

type Project struct {
	mongorm.Model `bson:",inline"`
	Title         string `bson:"title" json:"title"`
	Number        string `bson:"number" json:"number"`
	Status        string `bson:"status" json:"status"`
}
