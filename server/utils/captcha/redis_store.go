package captcha

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const redisKeyPrefix = "captcha:"

var bg = context.Background()

// redisStore 基于 Redis 的验证码存储
type redisStore struct {
	cli *redis.Client
	ttl time.Duration
}

func newRedisStore(cli *redis.Client) *redisStore {
	return &redisStore{cli: cli, ttl: 3 * time.Minute}
}

func (s *redisStore) Set(id, val string) error {
	return s.cli.Set(bg, redisKeyPrefix+id, val, s.ttl).Err()
}

func (s *redisStore) Get(id string, clear bool) string {
	key := redisKeyPrefix + id
	val, err := s.cli.Get(bg, key).Result()
	if err != nil {
		return ""
	}
	if clear {
		s.cli.Del(bg, key)
	}
	return val
}

func (s *redisStore) Verify(id, answer string, clear bool) bool {
	want := s.Get(id, clear)
	return want != "" && want == answer
}
