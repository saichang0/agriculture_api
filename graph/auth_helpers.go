package graph

import (
	"context"
	"time"

	"agriculture-api/graph/model"
	"agriculture-api/internal/auth"
)

// issueTokenPair generates a fresh access token + refresh token for a user,
// persisting the refresh token's hash so it can be looked up and rotated later.
// Not a GraphQL resolver — kept in its own file so `gqlgen generate` never treats
// it as orphaned resolver code and moves it into users.resolvers.go's dead-code
// comment block.
func (r *mutationResolver) issueTokenPair(ctx context.Context, doc *model.UserDoc) (*model.AuthPayload, error) {
	accessToken, err := r.JWT.Generate(doc.ID.Hex(), doc.Role)
	if err != nil {
		return nil, err
	}

	refreshSecret, err := auth.NewRefreshTokenSecret()
	if err != nil {
		return nil, err
	}

	if err := r.RefreshTokenRepo.Create(ctx, &model.RefreshTokenDoc{
		UserID:    doc.ID,
		TokenHash: auth.HashRefreshToken(refreshSecret),
		ExpiresAt: time.Now().Add(r.RefreshTokenTTL),
	}); err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		Token:        accessToken,
		RefreshToken: refreshSecret,
		User:         toGraphUser(doc),
	}, nil
}
