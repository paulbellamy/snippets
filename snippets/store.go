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
	s := &MemoryStore{
		snippets: map[string]*Snippet{},
		ticker:   time.NewTicker(1 * time.Second),
	}
	go s.garbageCollector()
	return s, nil
}

// Note: Not using a sync.Map here so we get better atomicity around expiry
type MemoryStore struct {
	snippets map[string]*Snippet
	ticker   *time.Ticker
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

// TODO: stop-the-world gc is fine for dev, but might want something better if
// scalability and liveness matter.
func (s *MemoryStore) garbageCollector() {
	for now := range s.ticker.C {
		s.Lock()
		// Two-stage here, to prevent modifying the map while iterating over it.
		var expired []string
		for name, v := range s.snippets {
			if now.After(v.ExpiresAt) {
				expired = append(expired, name)
			}
		}
		for _, name := range expired {
			delete(s.snippets, name)
		}
		s.Unlock()
	}
}

func (s *MemoryStore) Close() error {
	s.Lock()
	defer s.Unlock()
	// TODO: Actually wait for garbage collector to exit here.
	s.ticker.Stop()
	return nil
}
