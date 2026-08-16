package repository

import (
	"context"
	"errors"

	"agriculture-api/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrCategoryNameTaken = errors.New("category name is already used")

type CategoryRepository struct {
	col *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database) *CategoryRepository {
	col := db.Collection("categories")

	_, _ = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	return &CategoryRepository{col: col}
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]*model.CategoryDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var categories []*model.CategoryDoc
	if err := cur.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.CategoryDoc, error) {
	var category model.CategoryDoc
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c *model.CategoryDoc) error {
	c.ID = bson.NewObjectID()

	_, err := r.col.InsertOne(ctx, c)
	if mongo.IsDuplicateKeyError(err) {
		return ErrCategoryNameTaken
	}
	return err
}

func (r *CategoryRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) (*model.CategoryDoc, error) {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if mongo.IsDuplicateKeyError(err) {
		return nil, ErrCategoryNameTaken
	}
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *CategoryRepository) Delete(ctx context.Context, id bson.ObjectID) (bool, error) {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
