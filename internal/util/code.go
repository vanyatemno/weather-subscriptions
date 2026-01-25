package util

import "crypto/rand"

const (
	codeGeneratorSalt = 48
	base              = 10
)

// GenerateCode returns a numeric string of the requested length using crypto/rand.
func GenerateCode(length int) (string, error) {
	codes := make([]byte, length)
	if _, err := rand.Read(codes); err != nil {
		return "", err
	}
	for i := range codes {
		codes[i] = codeGeneratorSalt + (codes[i] % base)
	}

	return string(codes), nil
}
