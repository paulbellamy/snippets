package snippets

import (
	"sync"
	"time"
)

type Store interface {
	Load(name string) (*Snippet, error)
	Store(name, value string, ttl time.Duration) (*Snippet, error)
	Close() error
}

type Snippet struct {
	ExpiresAt time.Time     `json:"expires_at"`
	ExpiresIn time.Duration `json:"expires_in"`
	Name      string        `json:"name"`
	Snippet   string        `json:"snippet"`
}

func NewStore() (Store, error) {
	return &MemoryStore{
		snippets: map[string]*Snippet{},
	}, nil
}

// Note: Not using a sync.Map here so we get better atomicity around expiry
type MemoryStore struct {
	snippets map[string]*Snippet
	sync.RWMutex
}

func (s *MemoryStore) Load(name string) (*Snippet, error) {
	s.RLock()
	defer s.RUnlock()
	el, ok := s.snippets[name]
	if !ok || time.Now().After(el.ExpiresAt) {
		// nil, nil will mean not found.
		return nil, nil
	}
	return el, nil
}

func (s *MemoryStore) Store(name, value string, expiresIn time.Duration) (*Snippet, error) {
	s.Lock()
	defer s.Unlock()
	s.snippets[name] = &Snippet{
		ExpiresAt: time.Now().Add(expiresIn),
		ExpiresIn: expiresIn,
		Name:      name,
		Snippet:   value,
	}
	return s.snippets[name], nil
}

func (s *MemoryStore) Close() error {
	return nil
}
