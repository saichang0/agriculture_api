package model

import "go.mongodb.org/mongo-driver/v2/bson"

type DebtPaymentDoc struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	SaleID      bson.ObjectID `bson:"saleId"`
	CustomerID  bson.ObjectID `bson:"customerId"`
	UserID      bson.ObjectID `bson:"userId"`
	AmountPaid  float64       `bson:"amountPaid"`
	PaymentDate int64         `bson:"paymentDate"`
	Note        *string       `bson:"note,omitempty"`
}
