package repository

import (
	"context"
	"errors"

	"agriculture-api/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrUsernameTaken = errors.New("username is already registered")
	ErrPhoneTaken    = errors.New("phone number is already registered")
)

type UserRepository struct {
	col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	col := db.Collection("users")

	_, _ = col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "phone", Value: 1}}, Options: options.Index().SetUnique(true)},
	})

	return &UserRepository{col: col}
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*model.UserDoc, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []*model.UserDoc
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.UserDoc, error) {
	var user model.UserDoc
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.UserDoc, error) {
	var user model.UserDoc
	err := r.col.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*model.UserDoc, error) {
	var user model.UserDoc
	err := r.col.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, u *model.UserDoc) error {
	if existing, err := r.FindByUsername(ctx, u.Username); err != nil {
		return err
	} else if existing != nil {
		return ErrUsernameTaken
	}

	if existing, err := r.FindByPhone(ctx, u.Phone); err != nil {
		return err
	} else if existing != nil {
		return ErrPhoneTaken
	}

	u.ID = bson.NewObjectID()
	if u.Status == "" {
		u.Status = "ACTIVE"
	}

	_, err := r.col.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return ErrPhoneTaken
	}
	return err
}

func (r *UserRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) (*model.UserDoc, error) {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *UserRepository) Delete(ctx context.Context, id bson.ObjectID) (bool, error) {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
