package snippets

import (
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
	return NewMemoryStore()
}
