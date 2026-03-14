package memory

import (
	"context"
	"errors"
	"hash/fnv"
	"sort"
	"sync"
	"time"
)

type valueItem struct {
	value  int64
	expiry time.Time
}

// MemoryStore is an in process counter store for local or development use.
type MemoryStore struct {
	shards []memoryShard
}

type memoryShard struct {
	mu         sync.Mutex
	values     map[string]valueItem
	timestamps map[string][]time.Time
	tsExpiry   map[string]time.Time
	leases     map[string]valueItem
}

// New creates an empty memory store.
func New() *MemoryStore {
	return NewSharded(64)
}

// NewSharded creates a lock-sharded memory store for higher concurrency.
func NewSharded(shardCount int) *MemoryStore {
	if shardCount <= 0 {
		shardCount = 1
	}
	shards := make([]memoryShard, shardCount)
	for i := range shards {
		shards[i] = memoryShard{
			values:     map[string]valueItem{},
			timestamps: map[string][]time.Time{},
			tsExpiry:   map[string]time.Time{},
			leases:     map[string]valueItem{},
		}
	}
	return &MemoryStore{shards: shards}
}

// Increment adds value to a key and applies ttl when provided.
func (m *MemoryStore) Increment(_ context.Context, key string, value int64, ttl time.Duration) (int64, error) {
	shard := m.shardFor(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	it := getValueLocked(shard, key)
	it.value += value
	if ttl > 0 {
		it.expiry = time.Now().Add(ttl)
	}
	shard.values[key] = it
	return it.value, nil
}

// Get reads a counter value.
func (m *MemoryStore) Get(_ context.Context, key string) (int64, error) {
	shard := m.shardFor(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	it := getValueLocked(shard, key)
	return it.value, nil
}

// Set writes a counter value.
func (m *MemoryStore) Set(_ context.Context, key string, value int64, ttl time.Duration) error {
	shard := m.shardFor(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	it := valueItem{value: value}
	if ttl > 0 {
		it.expiry = time.Now().Add(ttl)
	}
	shard.values[key] = it
	return nil
}

// CompareAndSwap updates value when old value matches current value.
func (m *MemoryStore) CompareAndSwap(_ context.Context, key string, oldVal, newVal int64, ttl time.Duration) (bool, error) {
	shard := m.shardFor(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	it := getValueLocked(shard, key)
	if it.value != oldVal {
		return false, nil
	}
	it.value = newVal
	if ttl > 0 {
		it.expiry = time.Now().Add(ttl)
	}
	shard.values[key] = it
	return true, nil
}

// Delete removes all stored data for one key.
func (m *MemoryStore) Delete(_ context.Context, key string) error {
	shard := m.shardFor(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.values, key)
	delete(shard.timestamps, key)
	delete(shard.tsExpiry, key)
	delete(shard.leases, key)
	return nil
}

// AddTimestamp appends a timestamp entry used by log style algorithms.
func (m *MemoryStore) AddTimestamp(_ context.Context, key string, ts time.Time, ttl time.Duration) error {
	shard := m.shardFor(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	arr := append(shard.timestamps[key], ts)
	sort.Slice(arr, func(i, j int) bool { return arr[i].Before(arr[j]) })
	shard.timestamps[key] = arr
	if ttl > 0 {
		shard.tsExpiry[key] = time.Now().Add(ttl)
	}
	return nil
}

// CountAfter counts timestamps newer than or equal to the provided time.
func (m *MemoryStore) CountAfter(_ context.Context, key string, after time.Time) (int64, error) {
	shard := m.shardFor(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	cleanTSLocked(shard, key)
	var cnt int64
	for _, t := range shard.timestamps[key] {
		if !t.Before(after) {
			cnt++
		}
	}
	return cnt, nil
}

// TrimBefore removes timestamps older than the provided time.
func (m *MemoryStore) TrimBefore(_ context.Context, key string, before time.Time) error {
	shard := m.shardFor(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	arr := shard.timestamps[key]
	j := 0
	for _, t := range arr {
		if !t.Before(before) {
			arr[j] = t
			j++
		}
	}
	shard.timestamps[key] = arr[:j]
	return nil
}

// AcquireLease reserves one concurrency slot.
func (m *MemoryStore) AcquireLease(_ context.Context, key string, limit int64, ttl time.Duration) (bool, int64, error) {
	if limit <= 0 {
		return false, 0, errors.New("limit must be positive")
	}
	shard := m.shardFor(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	it := getLeaseLocked(shard, key)
	if it.value >= limit {
		return false, it.value, nil
	}
	it.value++
	if ttl > 0 {
		it.expiry = time.Now().Add(ttl)
	}
	shard.leases[key] = it
	return true, it.value, nil
}

// ReleaseLease frees one concurrency slot.
func (m *MemoryStore) ReleaseLease(_ context.Context, key string) error {
	shard := m.shardFor(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	it := getLeaseLocked(shard, key)
	if it.value > 0 {
		it.value--
	}
	shard.leases[key] = it
	return nil
}

func getValueLocked(shard *memoryShard, key string) valueItem {
	it := shard.values[key]
	if !it.expiry.IsZero() && time.Now().After(it.expiry) {
		delete(shard.values, key)
		return valueItem{}
	}
	return it
}

func getLeaseLocked(shard *memoryShard, key string) valueItem {
	it := shard.leases[key]
	if !it.expiry.IsZero() && time.Now().After(it.expiry) {
		delete(shard.leases, key)
		return valueItem{}
	}
	return it
}

func cleanTSLocked(shard *memoryShard, key string) {
	expiry := shard.tsExpiry[key]
	if expiry.IsZero() || time.Now().Before(expiry) {
		return
	}
	delete(shard.timestamps, key)
	delete(shard.tsExpiry, key)
}

func (m *MemoryStore) shardFor(key string) *memoryShard {
	if len(m.shards) == 1 {
		return &m.shards[0]
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	idx := int(h.Sum32() % uint32(len(m.shards)))
	return &m.shards[idx]
}
