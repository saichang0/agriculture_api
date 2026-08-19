package repository

import (
	"context"
	"time"

	"agriculture-api/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RefreshTokenRepository struct {
	col *mongo.Collection
}

func NewRefreshTokenRepository(db *mongo.Database) *RefreshTokenRepository {
	col := db.Collection("refreshTokens")

	_, _ = col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "tokenHash", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "expiresAt", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0)},
	})

	return &RefreshTokenRepository{col: col}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, t *model.RefreshTokenDoc) error {
	t.ID = bson.NewObjectID()
	t.CreatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, t)
	return err
}

func (r *RefreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (*model.RefreshTokenDoc, error) {
	var doc model.RefreshTokenDoc
	err := r.col.FindOne(ctx, bson.M{"tokenHash": tokenHash}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// DeleteByHash removes a single refresh token, used when it is rotated or
// explicitly logged out.
func (r *RefreshTokenRepository) DeleteByHash(ctx context.Context, tokenHash string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"tokenHash": tokenHash})
	return err
}

// DeleteAllForUser revokes every refresh token belonging to a user, e.g. on
// password change or "log out of all devices".
func (r *RefreshTokenRepository) DeleteAllForUser(ctx context.Context, userID bson.ObjectID) error {
	_, err := r.col.DeleteMany(ctx, bson.M{"userId": userID})
	return err
}
