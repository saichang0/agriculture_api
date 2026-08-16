package model

import "go.mongodb.org/mongo-driver/v2/bson"

type SaleDoc struct {
	ID            bson.ObjectID  `bson:"_id,omitempty"`
	Code          string         `bson:"code"`
	CustomerID    *bson.ObjectID `bson:"customerId,omitempty"`
	UserID        bson.ObjectID  `bson:"userId"`
	SaleDate      int64          `bson:"saleDate"`
	Total         float64        `bson:"total"`
	Paid          float64        `bson:"paid"`
	Debt          float64        `bson:"debt"`
	PaymentStatus string         `bson:"paymentStatus"`
	DueDate       *string        `bson:"dueDate,omitempty"`
	PaymentMethod *string        `bson:"paymentMethod,omitempty"`
}

type SaleItemDoc struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	SaleID    bson.ObjectID `bson:"saleId"`
	ProductID bson.ObjectID `bson:"productId"`
	Quantity  float64       `bson:"quantity"`
	CostPrice float64       `bson:"costPrice"`
	UnitPrice float64       `bson:"unitPrice"`
	PriceType string        `bson:"priceType"`
	Subtotal  float64       `bson:"subtotal"`
}
