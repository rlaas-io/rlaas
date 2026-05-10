package slidinglog

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rlaas-io/rlaas/internal/store"
	"github.com/rlaas-io/rlaas/internal/store/counter/memory"
	"github.com/rlaas-io/rlaas/pkg/model"
)

type swlCountErrStore struct{ store.CounterStore }

func (swlCountErrStore) TrimBefore(context.Context, string, time.Time) error { return nil }
func (swlCountErrStore) CountAfter(context.Context, string, time.Time) (int64, error) {
	return 0, errors.New("count failed")
}

type swlTrimErrStore struct{ store.CounterStore }

func (swlTrimErrStore) TrimBefore(context.Context, string, time.Time) error {
	return errors.New("trim failed")
}

func TestSlidingLog_CountErrorPath(t *testing.T) {
	e := New(swlCountErrStore{})
	_, err := e.Evaluate(context.Background(), model.Policy{Algorithm: model.AlgorithmConfig{Limit: 1, Window: "1m"}}, model.RequestContext{}, "k")
	require.Error(t, err)
}

// swlFallbackDenyStore: non-atomic store that returns count=3 to trigger the
// fallback-path deny branch (count+cost > limit).
type swlFallbackDenyStore struct{ store.CounterStore }

func (swlFallbackDenyStore) TrimBefore(_ context.Context, _ string, _ time.Time) error { return nil }
func (swlFallbackDenyStore) CountAfter(_ context.Context, _ string, _ time.Time) (int64, error) {
	return 3, nil
}

func TestSlidingLog_TrimErrorPath(t *testing.T) {
	e := New(swlTrimErrStore{})
	_, err := e.Evaluate(context.Background(), model.Policy{Algorithm: model.AlgorithmConfig{Limit: 1, Window: "1m"}}, model.RequestContext{}, "k")
	require.Error(t, err)
}

func TestSlidingLog_FallbackDenyPath(t *testing.T) {
	// Uses a non-atomic store (no CheckAndAddTimestamps) so the fallback path
	// is exercised. CountAfter returns 3, limit=3, cost=1 → 3+1>3 → deny.
	e := New(swlFallbackDenyStore{})
	p := model.Policy{Algorithm: model.AlgorithmConfig{Limit: 3, Window: "1m"}, Action: model.ActionDeny}
	d, err := e.Evaluate(context.Background(), p, model.RequestContext{}, "k")
	require.NoError(t, err)
	require.False(t, d.Allowed, "fallback deny: count+cost > limit")
	assert.Positive(t, d.RetryAfter, "RetryAfter should be positive")
}

func TestSlidingLog_ComputeRetryAfterZeroCount(t *testing.T) {
	// When the log is empty (count=0) but cost exceeds the limit,
	// computeRetryAfter takes the count<=0 branch and returns the full window.
	e := New(memory.New())
	now := time.Unix(1000, 0)
	e.Now = func() time.Time { return now }
	p := model.Policy{
		Algorithm: model.AlgorithmConfig{Limit: 1, Window: "1m"},
		Action:    model.ActionDeny,
	}
	d, err := e.Evaluate(context.Background(), p, model.RequestContext{Quantity: 2}, "k")
	require.NoError(t, err)
	assert.False(t, d.Allowed, "cost=2 > limit=1 should deny even on empty log")
	assert.Equal(t, time.Minute, d.RetryAfter, "zero-count path must return full window")
}
