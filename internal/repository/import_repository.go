package repository

import (
	"context"

	"agriculture-api/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ImportRepository struct {
	col *mongo.Collection
}

func NewImportRepository(db *mongo.Database) *ImportRepository {
	return &ImportRepository{col: db.Collection("imports")}
}

func (r *ImportRepository) FindAll(ctx context.Context) ([]*model.ImportDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var imports []*model.ImportDoc
	if err := cur.All(ctx, &imports); err != nil {
		return nil, err
	}
	return imports, nil
}

func (r *ImportRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.ImportDoc, error) {
	var imp model.ImportDoc
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&imp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &imp, nil
}

func (r *ImportRepository) FindByProduct(ctx context.Context, productID bson.ObjectID) ([]*model.ImportDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{"productId": productID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var imports []*model.ImportDoc
	if err := cur.All(ctx, &imports); err != nil {
		return nil, err
	}
	return imports, nil
}

func (r *ImportRepository) Create(ctx context.Context, imp *model.ImportDoc) error {
	imp.ID = bson.NewObjectID()
	_, err := r.col.InsertOne(ctx, imp)
	return err
}

func (r *ImportRepository) Delete(ctx context.Context, id bson.ObjectID) (bool, error) {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
