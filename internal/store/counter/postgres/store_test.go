package postgres

import (
	"context"
	"testing"
	"time"
)

func TestPostgresCounterStoreScaffoldErrors(t *testing.T) {
	s := New("dsn")
	if _, err := s.Increment(context.Background(), "k", 1, time.Second); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := s.Get(context.Background(), "k"); err == nil {
		t.Fatalf("expected error")
	}
	if err := s.Set(context.Background(), "k", 1, time.Second); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := s.CompareAndSwap(context.Background(), "k", 1, 2, time.Second); err == nil {
		t.Fatalf("expected error")
	}
	if err := s.Delete(context.Background(), "k"); err == nil {
		t.Fatalf("expected error")
	}
	if err := s.AddTimestamp(context.Background(), "k", time.Now(), time.Second); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := s.CountAfter(context.Background(), "k", time.Now()); err == nil {
		t.Fatalf("expected error")
	}
	if err := s.TrimBefore(context.Background(), "k", time.Now()); err == nil {
		t.Fatalf("expected error")
	}
	if _, _, err := s.AcquireLease(context.Background(), "k", 1, time.Second); err == nil {
		t.Fatalf("expected error")
	}
	if err := s.ReleaseLease(context.Background(), "k"); err == nil {
		t.Fatalf("expected error")
	}
}
