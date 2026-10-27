// Package cache holds a small in-memory TTL cache.
package cache

import (
	"sync"
	"time"
)

type entry[V any] struct {
	value   V
	expires time.Time
}

// TTL is a concurrency-safe in-memory cache whose entries expire after a fixed duration.
type TTL[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]entry[V]
	ttl   time.Duration
	now   func() time.Time
}

// NewTTL builds a cache. A nil now defaults to time.Now.
func NewTTL[K comparable, V any](ttl time.Duration, now func() time.Time) *TTL[K, V] {
	if now == nil {
		now = time.Now
	}
	return &TTL[K, V]{items: make(map[K]entry[V]), ttl: ttl, now: now}
}

// Get returns the value for key, treating expired entries as missing.
func (c *TTL[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || !c.now().Before(e.expires) {
		var zero V
		return zero, false
	}
	return e.value, true
}

// Set stores v under key and drops every expired entry, so the map cannot grow unbounded.
func (c *TTL[K, V]) Set(key K, v V) {
	now := c.now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, e := range c.items {
		if !now.Before(e.expires) {
			delete(c.items, k)
		}
	}
	c.items[key] = entry[V]{value: v, expires: now.Add(c.ttl)}
}
