package tools

import (
	"crypto/sha256"
	"encoding/hex"
)

type HashOption struct {
	Key []byte
}

// Hash make hash use sha256
func Hash(str string, options ...WithHashOption) (string, error) {
	var option HashOption
	for _, o := range options {
		o(&option)
	}
	cryptor := sha256.New()
	_, err := cryptor.Write([]byte(str))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(cryptor.Sum(option.Key)), nil
}

// WithKey set the key to sha256
func WithKey(str []byte) WithHashOption {
	return func(option *HashOption) {
		option.Key = str
	}
}

type WithHashOption func(option *HashOption)
