package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Sign(key string, msg string) (string, error) {
	hash := hmac.New(sha256.New, []byte(key))
	if _, err := hash.Write([]byte(msg)); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
