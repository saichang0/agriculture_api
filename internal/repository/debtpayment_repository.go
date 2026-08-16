package repository

import (
	"context"

	"agriculture-api/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DebtPaymentRepository struct {
	col *mongo.Collection
}

func NewDebtPaymentRepository(db *mongo.Database) *DebtPaymentRepository {
	return &DebtPaymentRepository{col: db.Collection("debtPayments")}
}

func (r *DebtPaymentRepository) FindAll(ctx context.Context) ([]*model.DebtPaymentDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var payments []*model.DebtPaymentDoc
	if err := cur.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *DebtPaymentRepository) FindBySale(ctx context.Context, saleID bson.ObjectID) ([]*model.DebtPaymentDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{"saleId": saleID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var payments []*model.DebtPaymentDoc
	if err := cur.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *DebtPaymentRepository) FindByCustomer(ctx context.Context, customerID bson.ObjectID) ([]*model.DebtPaymentDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{"customerId": customerID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var payments []*model.DebtPaymentDoc
	if err := cur.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *DebtPaymentRepository) Create(ctx context.Context, p *model.DebtPaymentDoc) error {
	p.ID = bson.NewObjectID()
	_, err := r.col.InsertOne(ctx, p)
	return err
}
