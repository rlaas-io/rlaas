package memory

import (
	"context"
	"testing"
	"time"
)

func TestMemoryStoreCounterOps(t *testing.T) {
	s := New()
	if v, _ := s.Increment(context.Background(), "k", 2, time.Second); v != 2 {
		t.Fatalf("unexpected increment")
	}
	if v, _ := s.Get(context.Background(), "k"); v != 2 {
		t.Fatalf("unexpected get")
	}
	if ok, _ := s.CompareAndSwap(context.Background(), "k", 1, 5, 0); ok {
		t.Fatalf("cas should fail")
	}
	if ok, _ := s.CompareAndSwap(context.Background(), "k", 2, 5, 0); !ok {
		t.Fatalf("cas should pass")
	}
	_ = s.Set(context.Background(), "k2", 1, 0)
	_ = s.Delete(context.Background(), "k2")
}

func TestMemoryStoreTimestampOps(t *testing.T) {
	s := New()
	now := time.Now()
	_ = s.AddTimestamp(context.Background(), "ts", now.Add(-time.Second), time.Second)
	_ = s.AddTimestamp(context.Background(), "ts", now, time.Second)
	if c, _ := s.CountAfter(context.Background(), "ts", now.Add(-500*time.Millisecond)); c != 1 {
		t.Fatalf("unexpected count")
	}
	_ = s.TrimBefore(context.Background(), "ts", now.Add(-500*time.Millisecond))
	if c, _ := s.CountAfter(context.Background(), "ts", now.Add(-2*time.Second)); c != 1 {
		t.Fatalf("unexpected count after trim")
	}
}

func TestMemoryStoreLeaseOps(t *testing.T) {
	s := New()
	if _, _, err := s.AcquireLease(context.Background(), "l", 0, time.Second); err == nil {
		t.Fatalf("expected invalid limit error")
	}
	ok, cur, _ := s.AcquireLease(context.Background(), "l", 1, time.Second)
	if !ok || cur != 1 {
		t.Fatalf("expected first lease")
	}
	ok, _, _ = s.AcquireLease(context.Background(), "l", 1, time.Second)
	if ok {
		t.Fatalf("expected lease deny")
	}
	_ = s.ReleaseLease(context.Background(), "l")
}

func TestMemoryStoreExpiryHelpers(t *testing.T) {
	s := New()
	_ = s.Set(context.Background(), "exp", 1, 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if v, _ := s.Get(context.Background(), "exp"); v != 0 {
		t.Fatalf("expired value should reset")
	}
	_, _, _ = s.AcquireLease(context.Background(), "lexp", 1, 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	_ = s.ReleaseLease(context.Background(), "lexp")

	_ = s.AddTimestamp(context.Background(), "ts-exp", time.Now(), 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if c, _ := s.CountAfter(context.Background(), "ts-exp", time.Now().Add(-time.Hour)); c != 0 {
		t.Fatalf("expired timestamps should be cleaned")
	}
}
