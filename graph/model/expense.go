package model

import "go.mongodb.org/mongo-driver/v2/bson"

type ExpenseDoc struct {
	ID     bson.ObjectID `bson:"_id,omitempty"`
	UserID bson.ObjectID `bson:"userId"`
	Title  string        `bson:"title"`
	Type   string        `bson:"type"`
	Amount float64       `bson:"amount"`
	Date   int64         `bson:"date"`
}
