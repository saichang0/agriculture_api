package repository

import (
	"context"
	"errors"

	"agriculture-api/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrBarcodeTaken = errors.New("barcode is already used by another product")

type ProductRepository struct {
	col *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	col := db.Collection("products")

	_, _ = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "barcode", Value: 1}},
		Options: options.Index().SetUnique(true).SetSparse(true),
	})

	return &ProductRepository{col: col}
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]*model.ProductDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []*model.ProductDoc
	if err := cur.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.ProductDoc, error) {
	var product model.ProductDoc
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *model.ProductDoc) error {
	p.ID = bson.NewObjectID()
	if p.Status == "" {
		p.Status = "ACTIVE"
	}

	_, err := r.col.InsertOne(ctx, p)
	if mongo.IsDuplicateKeyError(err) {
		return ErrBarcodeTaken
	}
	return err
}

func (r *ProductRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) (*model.ProductDoc, error) {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if mongo.IsDuplicateKeyError(err) {
		return nil, ErrBarcodeTaken
	}
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *ProductRepository) Delete(ctx context.Context, id bson.ObjectID) (bool, error) {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}

// AdjustStock atomically adds delta (positive or negative) to the product's stockQty,
// and optionally sets a new costPrice. Used internally by imports, sales, and
// damagedProducts resolvers — not exposed directly via GraphQL.
func (r *ProductRepository) AdjustStock(ctx context.Context, id bson.ObjectID, delta float64, newCostPrice *float64) error {
	update := bson.M{"$inc": bson.M{"stockQty": delta}}
	if newCostPrice != nil {
		update["$set"] = bson.M{"costPrice": *newCostPrice}
	}
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
