package redis

import (
	"context"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func newBrokenStore() *Store {
	return &Store{client: goredis.NewClient(&goredis.Options{Addr: "127.0.0.1:0", DialTimeout: 50 * time.Millisecond, ReadTimeout: 50 * time.Millisecond, WriteTimeout: 50 * time.Millisecond})}
}

func TestRedisStoreErrorPaths(t *testing.T) {
	s := newBrokenStore()
	ctx := context.Background()
	if _, err := s.Increment(ctx, "k", 1, time.Second); err == nil {
		t.Fatalf("expected increment error")
	}
	if _, err := s.Get(ctx, "k"); err == nil {
		t.Fatalf("expected get error")
	}
	if err := s.Set(ctx, "k", 1, time.Second); err == nil {
		t.Fatalf("expected set error")
	}
	if _, err := s.CompareAndSwap(ctx, "k", 1, 2, time.Second); err == nil {
		t.Fatalf("expected cas error")
	}
	if err := s.Delete(ctx, "k"); err == nil {
		t.Fatalf("expected delete error")
	}
	if err := s.AddTimestamp(ctx, "k", time.Now(), time.Second); err == nil {
		t.Fatalf("expected add timestamp error")
	}
	if _, err := s.CountAfter(ctx, "k", time.Now()); err == nil {
		t.Fatalf("expected count after error")
	}
	if err := s.TrimBefore(ctx, "k", time.Now()); err == nil {
		t.Fatalf("expected trim before error")
	}
	if _, _, err := s.AcquireLease(ctx, "k", 1, time.Second); err == nil {
		t.Fatalf("expected acquire lease error")
	}
	if err := s.ReleaseLease(ctx, "k"); err == nil {
		t.Fatalf("expected release lease error")
	}
}
