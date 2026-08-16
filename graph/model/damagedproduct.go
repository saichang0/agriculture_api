package model

import "go.mongodb.org/mongo-driver/v2/bson"

type DamagedProductDoc struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	ProductID bson.ObjectID `bson:"productId"`
	UserID    bson.ObjectID `bson:"userId"`
	Quantity  float64       `bson:"quantity"`
	CostPrice float64       `bson:"costPrice"`
	Reason    string        `bson:"reason"`
	Note      *string       `bson:"note,omitempty"`
	Date      int64         `bson:"date"`
}
