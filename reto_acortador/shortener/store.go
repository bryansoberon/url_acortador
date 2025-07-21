package shortener

import (
	"errors"
	"sync"
)

type Store struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

// Guarda una nueva relación shortCode → longURL
func (s *Store) Save(shortCode, longURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[shortCode]; exists {
		return errors.New("codigo corto ya existe")
	}

	s.data[shortCode] = longURL
	return nil
}

func (s *Store) Get(shortCode string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.data[shortCode]
	return url, exists
}
