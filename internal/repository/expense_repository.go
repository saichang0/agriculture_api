package repository

import (
	"context"

	"agriculture-api/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ExpenseRepository struct {
	col *mongo.Collection
}

func NewExpenseRepository(db *mongo.Database) *ExpenseRepository {
	return &ExpenseRepository{col: db.Collection("expenses")}
}

func (r *ExpenseRepository) FindAll(ctx context.Context) ([]*model.ExpenseDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var expenses []*model.ExpenseDoc
	if err := cur.All(ctx, &expenses); err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.ExpenseDoc, error) {
	var expense model.ExpenseDoc
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&expense)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepository) Create(ctx context.Context, e *model.ExpenseDoc) error {
	e.ID = bson.NewObjectID()
	_, err := r.col.InsertOne(ctx, e)
	return err
}

func (r *ExpenseRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) (*model.ExpenseDoc, error) {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *ExpenseRepository) Delete(ctx context.Context, id bson.ObjectID) (bool, error) {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
