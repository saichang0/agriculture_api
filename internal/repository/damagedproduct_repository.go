package repository

import (
	"context"

	"agriculture-api/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DamagedProductRepository struct {
	col *mongo.Collection
}

func NewDamagedProductRepository(db *mongo.Database) *DamagedProductRepository {
	return &DamagedProductRepository{col: db.Collection("damagedProducts")}
}

func (r *DamagedProductRepository) FindAll(ctx context.Context) ([]*model.DamagedProductDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var docs []*model.DamagedProductDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *DamagedProductRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.DamagedProductDoc, error) {
	var doc model.DamagedProductDoc
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

func (r *DamagedProductRepository) FindByProduct(ctx context.Context, productID bson.ObjectID) ([]*model.DamagedProductDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{"productId": productID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var docs []*model.DamagedProductDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *DamagedProductRepository) Create(ctx context.Context, d *model.DamagedProductDoc) error {
	d.ID = bson.NewObjectID()
	_, err := r.col.InsertOne(ctx, d)
	return err
}

func (r *DamagedProductRepository) Delete(ctx context.Context, id bson.ObjectID) (bool, error) {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
