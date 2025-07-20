package shortener

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	//"math/rand"
	"strings"
	"time"
)

type Shortener struct {
	store     *Store
	maxRetries int
	codeLength int
}

func NewShortener(store *Store, codeLength int) *Shortener {
	return &Shortener{
		store:     store,
		maxRetries: 5,         // máximo 5 intentos para evitar bucle infinito
		codeLength: codeLength, // longitud fija entre 6 y 8
	}
}

// Genera un código corto único para una URL larga
func (s *Shortener) GenerateShortCode(longURL string) (string, error) {
	for i := 0; i < s.maxRetries; i++ {
		code := s.createCode(longURL, time.Now().UnixNano(), i)

		if _, exists := s.store.Get(code); !exists {
			return code, nil
		}
	}

	return "", errors.New("failed to generate unique short code after several attempts")
}

// Función privada: crea un código a partir del hash de la URL + tiempo + intento
func (s *Shortener) createCode(longURL string, timestamp int64, salt int) string {
	// Combinar URL + timestamp + intento
	base := fmt.Sprintf("%s:%d:%d", longURL, timestamp, salt)

	// Crear SHA1 hash
	hash := sha1.Sum([]byte(base))

	// Codificar a base64
	encoded := base64.URLEncoding.EncodeToString(hash[:])

	// Quitar caracteres no alfanuméricos y limitar longitud
	clean := cleanAlphanumeric(encoded)

	if len(clean) < s.codeLength {
		return clean // fallback
	}

	return clean[:s.codeLength]
}

// Quita caracteres que no sean a-z, A-Z o 0-9
func cleanAlphanumeric(s string) string {
	var builder strings.Builder
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			builder.WriteRune(c)
		}
	}
	return builder.String()
}

