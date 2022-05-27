package domain

import (
	"crypto/rand"
	"encoding/base64"
)

// RandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, InternalError{Err: err}
	}

	return b, nil
}

// RandomStringURLSafe returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func RandomStringURLSafe(n int) (string, error) {
	b, err := RandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}
