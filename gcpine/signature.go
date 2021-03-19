package gcpine

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// ValidateSignature - perform signature verification
func ValidateSignature(secret, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	hash := hmac.New(sha256.New, []byte(secret))
	if _, err = hash.Write(body); err != nil {
		return false
	}

	return hmac.Equal(decoded, hash.Sum(nil))
}
