package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/razedwell/traxex/services/auth/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisSessionStore struct {
	rdb *redis.Client
}

func sessionKey(refresh string) string {
	return fmt.Sprintf("session:%s", refresh)
}

func NewRedisSessionStore(rdb *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{rdb}
}

func (s *RedisSessionStore) Save(ctx context.Context, sesh domain.Session, ttl time.Duration) error {
	seshVal, err := json.Marshal(sesh)

	if err != nil {
		return fmt.Errorf("redis session json marshal err: %w", err)
	}

	seshKey := sessionKey(sesh.RefreshToken)

	if err := s.rdb.Set(ctx, seshKey, seshVal, ttl).Err(); err != nil {
		return fmt.Errorf("redis set session kv err: %w", err)
	}

	return nil
}

func (s *RedisSessionStore) Get(ctx context.Context, refresh string) (domain.Session, error) {
	seshKey := sessionKey(refresh)
	seshVal, err := s.rdb.Get(ctx, seshKey).Bytes()

	if errors.Is(err, redis.Nil) {
		return domain.Session{}, domain.ErrTokenInvalid()
	}

	if err != nil {
		return domain.Session{}, fmt.Errorf("redis get session err: %w", err)
	}

	var sesh domain.Session
	err = json.Unmarshal(seshVal, &sesh)
	if err != nil {
		return domain.Session{}, fmt.Errorf("redis session unmarshall err: %w", err)
	}

	return sesh, nil
}

func (s *RedisSessionStore) Delete(ctx context.Context, refresh string) error {
	seshKey := sessionKey(refresh)
	res := s.rdb.Del(ctx, seshKey)
	if err := res.Err(); err != nil {
		return fmt.Errorf("redis session delete err: %w", err)
	}

	if res.Val() == 0 {
		return domain.ErrSessionNotFound()
	}

	return nil
}
