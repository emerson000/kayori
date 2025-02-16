package mongorm

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Model struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}
