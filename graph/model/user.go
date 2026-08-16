package model

import "go.mongodb.org/mongo-driver/v2/bson"

type UserDoc struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	Username     string        `bson:"username"`
	PasswordHash string        `bson:"passwordHash"`
	Role         string        `bson:"role"`
	FirstName    string        `bson:"firstName"`
	LastName     string        `bson:"lastName"`
	Phone        string        `bson:"phone"`
	Status       string        `bson:"status"`
}
