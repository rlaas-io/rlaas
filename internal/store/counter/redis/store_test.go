package redis

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
)

func TestRedisStoreOps(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis failed: %v", err)
	}
	defer mr.Close()

	s := New(mr.Addr(), "", 0)
	ctx := context.Background()

	v, err := s.Increment(ctx, "k", 1, time.Second)
	if err != nil || v != 1 {
		t.Fatalf("increment failed")
	}
	g, err := s.Get(ctx, "k")
	if err != nil || g != 1 {
		t.Fatalf("get failed")
	}
	if err := s.Set(ctx, "k", 3, time.Second); err != nil {
		t.Fatalf("set failed")
	}
	ok, err := s.CompareAndSwap(ctx, "k", 3, 4, time.Second)
	if err != nil || !ok {
		t.Fatalf("cas should pass")
	}
	ok, err = s.CompareAndSwap(ctx, "k", 3, 5, time.Second)
	if err != nil || ok {
		t.Fatalf("cas should fail")
	}
	if err := s.Delete(ctx, "k"); err != nil {
		t.Fatalf("delete failed")
	}
	if g0, err := s.Get(ctx, "missing"); err != nil || g0 != 0 {
		t.Fatalf("missing key should return zero")
	}
	if err := s.client.Set(ctx, "badint", "abc", time.Second).Err(); err != nil {
		t.Fatalf("setup parse failure failed")
	}
	if _, err := s.Get(ctx, "badint"); err == nil {
		t.Fatalf("expected parse error")
	}

	now := time.Now()
	if err := s.AddTimestamp(ctx, "ts", now.Add(-time.Second), time.Second); err != nil {
		t.Fatalf("add ts failed")
	}
	if err := s.AddTimestamp(ctx, "ts", now, time.Second); err != nil {
		t.Fatalf("add ts failed")
	}
	if c, err := s.CountAfter(ctx, "ts", now.Add(-500*time.Millisecond)); err != nil || c != 1 {
		t.Fatalf("count after failed")
	}
	if err := s.TrimBefore(ctx, "ts", now.Add(-500*time.Millisecond)); err != nil {
		t.Fatalf("trim failed")
	}

	ok, cur, err := s.AcquireLease(ctx, "lease", 1, time.Second)
	if err != nil || !ok || cur != 1 {
		t.Fatalf("lease should pass")
	}
	ok, _, err = s.AcquireLease(ctx, "lease", 1, time.Second)
	if err != nil || ok {
		t.Fatalf("lease should fail")
	}
	if err := s.ReleaseLease(ctx, "lease"); err != nil {
		t.Fatalf("release failed")
	}

	if asInt64(int64(2)) != 2 || asInt64(int(3)) != 3 || asInt64("4") != 4 || asInt64(struct{}{}) != 0 {
		t.Fatalf("asInt64 conversion failed")
	}
}
