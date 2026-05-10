package slidingcounter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rlaas-io/rlaas/internal/store"
	"github.com/rlaas-io/rlaas/pkg/model"
)

type swcGetErrStore struct{ store.CounterStore }

func (swcGetErrStore) Get(context.Context, string) (int64, error) { return 0, errors.New("get failed") }

type swcCASErrStore struct{ store.CounterStore }

func (swcCASErrStore) Get(_ context.Context, _ string) (int64, error) { return 0, nil }
func (swcCASErrStore) CompareAndSwap(_ context.Context, _ string, _, _ int64, _ time.Duration) (bool, error) {
	return false, errors.New("cas error")
}

func TestSlidingCounter_GetErrorPath(t *testing.T) {
	e := New(swcGetErrStore{})
	_, err := e.Evaluate(context.Background(), model.Policy{Algorithm: model.AlgorithmConfig{Limit: 1, Window: "1m"}}, model.RequestContext{}, "k")
	require.Error(t, err)
}

func TestSlidingCounter_CASError(t *testing.T) {
	e := New(swcCASErrStore{})
	now := time.Unix(60, 0)
	e.Now = func() time.Time { return now }
	_, err := e.Evaluate(context.Background(), model.Policy{Algorithm: model.AlgorithmConfig{Limit: 10, Window: "1m"}}, model.RequestContext{}, "k")
	require.Error(t, err)
}
