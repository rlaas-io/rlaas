package metrics

import "sync/atomic"

// Collector stores lightweight in process counters for core limiter events.
type Collector struct {
	DecisionsTotal     atomic.Int64
	DecisionsAllowed   atomic.Int64
	DecisionsDenied    atomic.Int64
	DecisionsShadow    atomic.Int64
	PolicyCacheHit     atomic.Int64
	PolicyCacheMiss    atomic.Int64
	BackendFailOpen    atomic.Int64
	BackendFailClosed  atomic.Int64
	CounterStoreErrors atomic.Int64
}

// New returns an empty metrics collector.
func New() *Collector {
	return &Collector{}
}
