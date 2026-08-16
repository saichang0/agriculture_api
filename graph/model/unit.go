package model

import "go.mongodb.org/mongo-driver/v2/bson"

type UnitDoc struct {
	ID   bson.ObjectID `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}
