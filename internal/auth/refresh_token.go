package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// NewRefreshTokenSecret generates a high-entropy random string to hand out as a
// refresh token. Only its hash is ever persisted — GenerateRefreshToken's caller
// is responsible for storing HashRefreshToken(secret), never the secret itself.
func NewRefreshTokenSecret() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// HashRefreshToken hashes a refresh token secret for storage/lookup. Refresh
// tokens are already high-entropy random values (unlike passwords), so a fast
// deterministic hash is appropriate here instead of bcrypt.
func HashRefreshToken(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}
