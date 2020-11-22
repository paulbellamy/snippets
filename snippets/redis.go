package snippets

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	redis "github.com/go-redis/redis/v8"
)

func NewRedisStore(u string) (Store, error) {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	password, _ := parsedUrl.User.Password()
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Addr:     parsedUrl.Host,
			Password: password,
			DB:       0, // use default DB
		}),
	}, nil
}

type RedisStore struct {
	client *redis.Client
}

func (s *RedisStore) Load(name string) (*Snippet, error) {
	raw, err := s.client.Get(context.TODO(), name).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	result := &Snippet{}
	if err := json.Unmarshal([]byte(raw), result); err != nil {
		return nil, err
	}
	// Hack to update the expiresAt. Could probably query the redis expiry or
	// something more clever here.
	return s.Store(name, result.Snippet, result.ExpiresIn)
}

func (s *RedisStore) Store(name, value string, expiresIn time.Duration) (*Snippet, error) {
	snip := &Snippet{
		ExpiresAt: time.Now().Add(expiresIn),
		ExpiresIn: expiresIn,
		Name:      name,
		Snippet:   value,
	}
	raw, err := json.Marshal(snip)
	if err != nil {
		return nil, err
	}
	err = s.client.Set(context.TODO(), name, string(raw), expiresIn).Err()
	if err != nil {
		return nil, err
	}
	return snip, nil
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}
