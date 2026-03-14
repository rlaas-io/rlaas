package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// Store implements CounterStore using Redis commands and Lua scripts.
type Store struct {
	client *goredis.Client
}

// New creates a Redis backed counter store.
func New(addr, password string, db int) *Store {
	return &Store{client: goredis.NewClient(&goredis.Options{Addr: addr, Password: password, DB: db})}
}

// Increment adds value to a key.
func (s *Store) Increment(ctx context.Context, key string, value int64, ttl time.Duration) (int64, error) {
	val, err := s.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, err
	}
	if ttl > 0 {
		_ = s.client.Expire(ctx, key, ttl).Err()
	}
	return val, nil
}

// Get returns key value or zero when key does not exist.
func (s *Store) Get(ctx context.Context, key string) (int64, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err == goredis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	out, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return out, nil
}

// Set writes key value with ttl.
func (s *Store) Set(ctx context.Context, key string, value int64, ttl time.Duration) error {
	return s.client.Set(ctx, key, value, ttl).Err()
}

// CompareAndSwap performs optimistic atomic update with WATCH.
func (s *Store) CompareAndSwap(ctx context.Context, key string, oldVal, newVal int64, ttl time.Duration) (bool, error) {
	err := s.client.Watch(ctx, func(tx *goredis.Tx) error {
		cur, err := tx.Get(ctx, key).Int64()
		if err == goredis.Nil {
			cur = 0
		} else if err != nil {
			return err
		}
		if cur != oldVal {
			return goredis.TxFailedErr
		}
		_, err = tx.TxPipelined(ctx, func(pipe goredis.Pipeliner) error {
			pipe.Set(ctx, key, newVal, ttl)
			return nil
		})
		return err
	}, key)
	if err == goredis.TxFailedErr {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Delete removes a key.
func (s *Store) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

// AddTimestamp writes one timestamp into a sorted set key.
func (s *Store) AddTimestamp(ctx context.Context, key string, ts time.Time, ttl time.Duration) error {
	if err := s.client.ZAdd(ctx, key, goredis.Z{Score: float64(ts.UnixNano()), Member: ts.UnixNano()}).Err(); err != nil {
		return err
	}
	if ttl > 0 {
		_ = s.client.Expire(ctx, key, ttl).Err()
	}
	return nil
}

// CountAfter returns sorted set member count after the given time.
func (s *Store) CountAfter(ctx context.Context, key string, after time.Time) (int64, error) {
	return s.client.ZCount(ctx, key, fmt.Sprintf("%d", after.UnixNano()), "+inf").Result()
}

// TrimBefore removes sorted set members older than the given time.
func (s *Store) TrimBefore(ctx context.Context, key string, before time.Time) error {
	return s.client.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", before.UnixNano())).Err()
}

// AcquireLease uses Lua to atomically enforce a concurrency limit.
func (s *Store) AcquireLease(ctx context.Context, key string, limit int64, ttl time.Duration) (bool, int64, error) {
	script := goredis.NewScript(`
local current = redis.call('INCR', KEYS[1])
if current == 1 and tonumber(ARGV[2]) > 0 then
  redis.call('PEXPIRE', KEYS[1], ARGV[2])
end
if current > tonumber(ARGV[1]) then
  redis.call('DECR', KEYS[1])
  return {0, current - 1}
end
return {1, current}
`)
	ms := ttl.Milliseconds()
	res, err := script.Run(ctx, s.client, []string{key}, limit, ms).Result()
	if err != nil {
		return false, 0, err
	}
	arr, ok := res.([]interface{})
	if !ok || len(arr) != 2 {
		return false, 0, fmt.Errorf("unexpected lua response")
	}
	okVal := asInt64(arr[0])
	curVal := asInt64(arr[1])
	return okVal == 1, curVal, nil
}

// ReleaseLease decrements active lease count safely.
func (s *Store) ReleaseLease(ctx context.Context, key string) error {
	script := goredis.NewScript(`
local current = redis.call('GET', KEYS[1])
if not current then
  return 0
end
current = tonumber(current)
if current <= 0 then
  redis.call('SET', KEYS[1], 0)
  return 0
end
return redis.call('DECR', KEYS[1])
`)
	_, err := script.Run(ctx, s.client, []string{key}).Result()
	return err
}

// asInt64 converts Lua response values into int64.
func asInt64(v interface{}) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case string:
		i, _ := strconv.ParseInt(t, 10, 64)
		return i
	default:
		return 0
	}
}
