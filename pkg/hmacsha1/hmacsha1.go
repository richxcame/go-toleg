package hmacsha1

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

// Generate decodes given key to base64 and creates new hmac hash from it
// Returns generated hash by encoding base64
func Generate(key, message string) string {
	secretKey, _ := base64.StdEncoding.DecodeString(key)
	h := hmac.New(sha1.New, secretKey)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
