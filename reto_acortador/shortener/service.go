package shortener

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Shortener struct {
	store      *Store
	maxRetries int
	codeLength int
}

func NewShortener(store *Store, codeLength int) *Shortener {
	return &Shortener{
		store:      store,
		maxRetries: 5,
		codeLength: codeLength,
	}
}

func (s *Shortener) GenerateShortCode(longURL string) (string, error) {
	for i := 0; i < s.maxRetries; i++ {
		code := s.createCode(longURL, time.Now().UnixNano(), i)

		if _, exists := s.store.Get(code); !exists {
			return code, nil
		}
	}
	return "", errors.New("failed to generate unique short code after several attempts")
}

func (s *Shortener) createCode(longURL string, timestamp int64, salt int) string {
	base := fmt.Sprintf("%s:%d:%d", longURL, timestamp, salt)
	hash := sha1.Sum([]byte(base))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	clean := cleanAlphanumeric(encoded)

	if len(clean) < s.codeLength {
		return clean
	}
	return clean[:s.codeLength]
}

func cleanAlphanumeric(s string) string {
	var builder strings.Builder
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			builder.WriteRune(c)
		}
	}
	return builder.String()
}
