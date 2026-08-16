package model

import "go.mongodb.org/mongo-driver/v2/bson"

type ProductDoc struct {
	ID              bson.ObjectID `bson:"_id,omitempty"`
	Barcode         *string       `bson:"barcode,omitempty"`
	Name            string        `bson:"name"`
	ImageURL        *string       `bson:"imageUrl,omitempty"`
	CategoryID      bson.ObjectID `bson:"categoryId"`
	UnitID          bson.ObjectID `bson:"unitId"`
	CostPrice       float64       `bson:"costPrice"`
	RetailPrice     float64       `bson:"retailPrice"`
	WholesalePrice  float64       `bson:"wholesalePrice"`
	WholesaleMinQty int32         `bson:"wholesaleMinQty"`
	StockQty        float64       `bson:"stockQty"`
	MinStockAlert   float64       `bson:"minStockAlert"`
	Status          string        `bson:"status"`
}
