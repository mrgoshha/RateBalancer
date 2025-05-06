package hash

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Hasher interface {
	Hash(password string) (string, error)
}

type SHA1Hasher struct{}

func NewSHA1Hasher() *SHA1Hasher {
	return &SHA1Hasher{}
}

func (h *SHA1Hasher) Hash(input string) (string, error) {
	hasher := sha1.New()

	if _, err := hasher.Write([]byte(input)); err != nil {
		return "", fmt.Errorf("failed to hash input: %w", err)
	}

	digest := hasher.Sum(nil)
	return hex.EncodeToString(digest), nil
}
