package models

import (
	"kayori.io/backend/data/mongorm"
)

type Artifact struct {
	mongorm.Model `bson:",inline"`
}
