package model

import "go.mongodb.org/mongo-driver/v2/bson"

type CustomerDoc struct {
	ID      bson.ObjectID `bson:"_id,omitempty"`
	Name    string        `bson:"name"`
	Phone   string        `bson:"phone"`
	Address *string       `bson:"address,omitempty"`
	Debt    float64       `bson:"debt"`
	Status  string        `bson:"status"`
}
