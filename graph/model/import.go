package model

import "go.mongodb.org/mongo-driver/v2/bson"

type ImportDoc struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	ProductID bson.ObjectID `bson:"productId"`
	Quantity  float64       `bson:"quantity"`
	CostPrice float64       `bson:"costPrice"`
	UserID    bson.ObjectID `bson:"userId"`
	Note      *string       `bson:"note,omitempty"`
	Date      int64         `bson:"date"`
}
