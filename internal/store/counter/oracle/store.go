package oracle

import (
	"context"
	"errors"
	"time"
)

// Store is a placeholder for oracle counter persistence.
type Store struct {
	DSN string
}

// New creates an oracle counter store scaffold.
func New(dsn string) *Store {
	return &Store{DSN: dsn}
}

func (s *Store) Increment(_ context.Context, _ string, _ int64, _ time.Duration) (int64, error) {
	return 0, errors.New("oracle counter store scaffold: implement SQL counters")
}

func (s *Store) Get(_ context.Context, _ string) (int64, error) {
	return 0, errors.New("oracle counter store scaffold: implement SQL counters")
}

func (s *Store) Set(_ context.Context, _ string, _ int64, _ time.Duration) error {
	return errors.New("oracle counter store scaffold: implement SQL counters")
}

func (s *Store) CompareAndSwap(_ context.Context, _ string, _, _ int64, _ time.Duration) (bool, error) {
	return false, errors.New("oracle counter store scaffold: implement SQL counters")
}

func (s *Store) Delete(_ context.Context, _ string) error {
	return errors.New("oracle counter store scaffold: implement SQL counters")
}

func (s *Store) AddTimestamp(_ context.Context, _ string, _ time.Time, _ time.Duration) error {
	return errors.New("oracle counter store scaffold: implement SQL counters")
}

func (s *Store) CountAfter(_ context.Context, _ string, _ time.Time) (int64, error) {
	return 0, errors.New("oracle counter store scaffold: implement SQL counters")
}

func (s *Store) TrimBefore(_ context.Context, _ string, _ time.Time) error {
	return errors.New("oracle counter store scaffold: implement SQL counters")
}

func (s *Store) AcquireLease(_ context.Context, _ string, _ int64, _ time.Duration) (bool, int64, error) {
	return false, 0, errors.New("oracle counter store scaffold: implement SQL counters")
}

func (s *Store) ReleaseLease(_ context.Context, _ string) error {
	return errors.New("oracle counter store scaffold: implement SQL counters")
}
