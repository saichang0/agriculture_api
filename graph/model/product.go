package model

import "go.mongodb.org/mongo-driver/v2/bson"

type ProductDoc struct {
	ID              bson.ObjectID      `bson:"_id,omitempty"`
	Barcode         *string            `bson:"barcode,omitempty"`
	Name            string             `bson:"name"`
	ImageURL        *string            `bson:"imageUrl,omitempty"`
	CategoryID      bson.ObjectID      `bson:"categoryId"`
	UnitID          bson.ObjectID      `bson:"unitId"`
	CostPrice       float64            `bson:"costPrice"`
	RetailPrice     float64            `bson:"retailPrice"`
	WholesalePrice  float64            `bson:"wholesalePrice"`
	WholesaleMinQty int32              `bson:"wholesaleMinQty"`
	StockQty        float64            `bson:"stockQty"`
	MinStockAlert   float64            `bson:"minStockAlert"`
	Status          string             `bson:"status"`
	PackagingUnits  []PackagingUnitDoc `bson:"packagingUnits,omitempty"`
}

// PackagingUnitDoc is an additional sellable unit for a product, alongside its
// base UnitID/stock tracking. Factor is how many base units make up 1 of this
// unit (e.g. a "ແກັດ" case might have Factor=4 if the base unit is "ຕຸກ" bottles).
// Prices here are independent of the base unit's prices — set directly by the
// shop owner (a case is often cheaper per-bottle than buying loose), never derived.
type PackagingUnitDoc struct {
	UnitID          bson.ObjectID `bson:"unitId"`
	Factor          float64       `bson:"factor"`
	CostPrice       float64       `bson:"costPrice"`
	RetailPrice     float64       `bson:"retailPrice"`
	WholesalePrice  float64       `bson:"wholesalePrice"`
	WholesaleMinQty int32         `bson:"wholesaleMinQty"`
}
