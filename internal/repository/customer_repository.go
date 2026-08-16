package repository

import (
	"context"
	"errors"

	"agriculture-api/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrCustomerPhoneTaken = errors.New("phone number is already registered to another customer")

type CustomerRepository struct {
	col *mongo.Collection
}

func NewCustomerRepository(db *mongo.Database) *CustomerRepository {
	col := db.Collection("customers")

	_, _ = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "phone", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	return &CustomerRepository{col: col}
}

func (r *CustomerRepository) FindAll(ctx context.Context) ([]*model.CustomerDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var customers []*model.CustomerDoc
	if err := cur.All(ctx, &customers); err != nil {
		return nil, err
	}
	return customers, nil
}

func (r *CustomerRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.CustomerDoc, error) {
	var customer model.CustomerDoc
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) Create(ctx context.Context, c *model.CustomerDoc) error {
	c.ID = bson.NewObjectID()
	if c.Status == "" {
		c.Status = "ACTIVE"
	}

	_, err := r.col.InsertOne(ctx, c)
	if mongo.IsDuplicateKeyError(err) {
		return ErrCustomerPhoneTaken
	}
	return err
}

func (r *CustomerRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) (*model.CustomerDoc, error) {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if mongo.IsDuplicateKeyError(err) {
		return nil, ErrCustomerPhoneTaken
	}
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *CustomerRepository) Delete(ctx context.Context, id bson.ObjectID) (bool, error) {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}

// AdjustDebt atomically adds delta (positive or negative) to the customer's debt.
// Used internally by sales and debtPayments resolvers — not exposed directly via GraphQL.
func (r *CustomerRepository) AdjustDebt(ctx context.Context, id bson.ObjectID, delta float64) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"debt": delta}})
	return err
}
