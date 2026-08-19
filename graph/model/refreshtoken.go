package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RefreshTokenDoc struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    bson.ObjectID `bson:"userId"`
	TokenHash string        `bson:"tokenHash"`
	ExpiresAt time.Time     `bson:"expiresAt"`
	CreatedAt time.Time     `bson:"createdAt"`
}
