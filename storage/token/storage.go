package token

import "sync"

type Storage struct {
	accessToken  string
	refreshToken string
	mu           sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) SetTokens(access, refresh string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accessToken = access
	s.refreshToken = refresh
}

func (s *Storage) GetToken() (string, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.accessToken, s.refreshToken
}

func (s *Storage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accessToken = ""
	s.refreshToken = ""
}
