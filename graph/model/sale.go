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
	// UnitID/Factor snapshot which selling unit this line was sold in, so historical
	// invoices stay correct even if the product's packaging units are edited/removed
	// later. UnitID is nil when sold in the product's base unit (Factor is always 1 then).
	UnitID *bson.ObjectID `bson:"unitId,omitempty"`
	Factor float64        `bson:"factor"`
}
