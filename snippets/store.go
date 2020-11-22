package snippets

import (
	"fmt"
	"net/url"
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

func NewStore(u string) (Store, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	switch parsed.Scheme {
	case "memory":
		return NewMemoryStore()
	case "redis":
		return NewRedisStore(u)
	default:
		return nil, fmt.Errorf("unknown db scheme: %s", parsed.Scheme)
	}
}
