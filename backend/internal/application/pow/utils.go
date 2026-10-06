package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func verifySolution(challenge, nonce string, difficulty int) bool {
	data := challenge + nonce
	hash := sha256.Sum256([]byte(data))
	hashHex := hex.EncodeToString(hash[:])
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hashHex, prefix)
}