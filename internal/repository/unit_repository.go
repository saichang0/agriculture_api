package repository

import (
	"context"
	"errors"

	"agriculture-api/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrUnitNameTaken = errors.New("unit name is already used")

type UnitRepository struct {
	col *mongo.Collection
}

func NewUnitRepository(db *mongo.Database) *UnitRepository {
	col := db.Collection("units")

	_, _ = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	return &UnitRepository{col: col}
}

func (r *UnitRepository) FindAll(ctx context.Context) ([]*model.UnitDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var units []*model.UnitDoc
	if err := cur.All(ctx, &units); err != nil {
		return nil, err
	}
	return units, nil
}

func (r *UnitRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.UnitDoc, error) {
	var unit model.UnitDoc
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&unit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &unit, nil
}

func (r *UnitRepository) Create(ctx context.Context, u *model.UnitDoc) error {
	u.ID = bson.NewObjectID()

	_, err := r.col.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return ErrUnitNameTaken
	}
	return err
}

func (r *UnitRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) (*model.UnitDoc, error) {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if mongo.IsDuplicateKeyError(err) {
		return nil, ErrUnitNameTaken
	}
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *UnitRepository) Delete(ctx context.Context, id bson.ObjectID) (bool, error) {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
