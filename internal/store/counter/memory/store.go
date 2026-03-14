package memory

import (
	"context"
	"errors"
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
	mu         sync.Mutex
	values     map[string]valueItem
	timestamps map[string][]time.Time
	leases     map[string]valueItem
}

// New creates an empty memory store.
func New() *MemoryStore {
	return &MemoryStore{
		values:     map[string]valueItem{},
		timestamps: map[string][]time.Time{},
		leases:     map[string]valueItem{},
	}
}

// Increment adds value to a key and applies ttl when provided.
func (m *MemoryStore) Increment(_ context.Context, key string, value int64, ttl time.Duration) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	it := m.getValueLocked(key)
	it.value += value
	if ttl > 0 {
		it.expiry = time.Now().Add(ttl)
	}
	m.values[key] = it
	return it.value, nil
}

// Get reads a counter value.
func (m *MemoryStore) Get(_ context.Context, key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	it := m.getValueLocked(key)
	return it.value, nil
}

// Set writes a counter value.
func (m *MemoryStore) Set(_ context.Context, key string, value int64, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	it := valueItem{value: value}
	if ttl > 0 {
		it.expiry = time.Now().Add(ttl)
	}
	m.values[key] = it
	return nil
}

// CompareAndSwap updates value when old value matches current value.
func (m *MemoryStore) CompareAndSwap(_ context.Context, key string, oldVal, newVal int64, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	it := m.getValueLocked(key)
	if it.value != oldVal {
		return false, nil
	}
	it.value = newVal
	if ttl > 0 {
		it.expiry = time.Now().Add(ttl)
	}
	m.values[key] = it
	return true, nil
}

// Delete removes all stored data for one key.
func (m *MemoryStore) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.values, key)
	delete(m.timestamps, key)
	delete(m.leases, key)
	return nil
}

// AddTimestamp appends a timestamp entry used by log style algorithms.
func (m *MemoryStore) AddTimestamp(_ context.Context, key string, ts time.Time, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	arr := append(m.timestamps[key], ts)
	sort.Slice(arr, func(i, j int) bool { return arr[i].Before(arr[j]) })
	m.timestamps[key] = arr
	if ttl > 0 {
		m.values[key+":ts_expiry"] = valueItem{expiry: time.Now().Add(ttl)}
	}
	return nil
}

// CountAfter counts timestamps newer than or equal to the provided time.
func (m *MemoryStore) CountAfter(_ context.Context, key string, after time.Time) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cleanTSLocked(key)
	var cnt int64
	for _, t := range m.timestamps[key] {
		if !t.Before(after) {
			cnt++
		}
	}
	return cnt, nil
}

// TrimBefore removes timestamps older than the provided time.
func (m *MemoryStore) TrimBefore(_ context.Context, key string, before time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	arr := m.timestamps[key]
	j := 0
	for _, t := range arr {
		if !t.Before(before) {
			arr[j] = t
			j++
		}
	}
	m.timestamps[key] = arr[:j]
	return nil
}

// AcquireLease reserves one concurrency slot.
func (m *MemoryStore) AcquireLease(_ context.Context, key string, limit int64, ttl time.Duration) (bool, int64, error) {
	if limit <= 0 {
		return false, 0, errors.New("limit must be positive")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	it := m.getLeaseLocked(key)
	if it.value >= limit {
		return false, it.value, nil
	}
	it.value++
	if ttl > 0 {
		it.expiry = time.Now().Add(ttl)
	}
	m.leases[key] = it
	return true, it.value, nil
}

// ReleaseLease frees one concurrency slot.
func (m *MemoryStore) ReleaseLease(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	it := m.getLeaseLocked(key)
	if it.value > 0 {
		it.value--
	}
	m.leases[key] = it
	return nil
}

func (m *MemoryStore) getValueLocked(key string) valueItem {
	it := m.values[key]
	if !it.expiry.IsZero() && time.Now().After(it.expiry) {
		delete(m.values, key)
		return valueItem{}
	}
	return it
}

func (m *MemoryStore) getLeaseLocked(key string) valueItem {
	it := m.leases[key]
	if !it.expiry.IsZero() && time.Now().After(it.expiry) {
		delete(m.leases, key)
		return valueItem{}
	}
	return it
}

func (m *MemoryStore) cleanTSLocked(key string) {
	exp := m.values[key+":ts_expiry"]
	if exp.expiry.IsZero() || time.Now().Before(exp.expiry) {
		return
	}
	delete(m.timestamps, key)
	delete(m.values, key+":ts_expiry")
}
