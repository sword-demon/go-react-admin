// Package cache provides local LRU cache implementation
package cache

import (
	"container/list"
	"sync"
	"time"
)

// LocalCache implements a thread-safe LRU cache with TTL support
// Layer 1 cache: in-memory, 5min TTL, 80%+ hit rate, <1ms latency
type LocalCache struct {
	mu         sync.RWMutex
	cache      map[string]*cacheEntry
	lruList    *list.List
	maxSize    int
	defaultTTL time.Duration

	// Metrics
	hits   uint64
	misses uint64
}

// cacheEntry represents a cache entry with expiration
type cacheEntry struct {
	key       string
	value     interface{}
	expiresAt time.Time
	element   *list.Element // LRU list element
}

// NewLocalCache creates a new local cache
// maxSize: maximum number of entries (recommend: 1000-10000)
// defaultTTL: default expiration time (recommend: 5 minutes)
func NewLocalCache(maxSize int, defaultTTL time.Duration) *LocalCache {
	if maxSize <= 0 {
		maxSize = 1000 // Default to 1000 entries
	}
	if defaultTTL <= 0 {
		defaultTTL = 5 * time.Minute
	}

	return &LocalCache{
		cache:      make(map[string]*cacheEntry, maxSize),
		lruList:    list.New(),
		maxSize:    maxSize,
		defaultTTL: defaultTTL,
	}
}

// Get retrieves a value from cache
// Returns (value, true) if found and not expired
// Returns (nil, false) if not found or expired
func (c *LocalCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.cache[key]
	if !exists {
		c.misses++
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.expiresAt) {
		c.removeEntry(entry)
		c.misses++
		return nil, false
	}

	// Move to front (most recently used)
	c.lruList.MoveToFront(entry.element)
	c.hits++
	return entry.value, true
}

// Set stores a value in cache with default TTL
func (c *LocalCache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores a value in cache with custom TTL
func (c *LocalCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If key exists, update it
	if entry, exists := c.cache[key]; exists {
		entry.value = value
		entry.expiresAt = time.Now().Add(ttl)
		c.lruList.MoveToFront(entry.element)
		return
	}

	// Evict least recently used if cache is full
	if c.lruList.Len() >= c.maxSize {
		c.evictOldest()
	}

	// Create new entry
	entry := &cacheEntry{
		key:       key,
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	entry.element = c.lruList.PushFront(entry)
	c.cache[key] = entry
}

// Delete removes a key from cache
func (c *LocalCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, exists := c.cache[key]; exists {
		c.removeEntry(entry)
	}
}

// DeletePrefix removes all keys with the given prefix
// Useful for invalidating related cache entries
// Example: DeletePrefix("user:permissions:") clears all user permission cache
func (c *LocalCache) DeletePrefix(prefix string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	for key, entry := range c.cache {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			c.removeEntry(entry)
			count++
		}
	}
	return count
}

// Clear removes all entries from cache
func (c *LocalCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*cacheEntry, c.maxSize)
	c.lruList.Init()
	c.hits = 0
	c.misses = 0
}

// Stats returns cache statistics
func (c *LocalCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(c.hits) / float64(total) * 100
	}

	return CacheStats{
		Hits:    c.hits,
		Misses:  c.misses,
		HitRate: hitRate,
		Size:    c.lruList.Len(),
		MaxSize: c.maxSize,
	}
}

// CleanupExpired removes expired entries (call periodically in background)
func (c *LocalCache) CleanupExpired() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	count := 0

	// Iterate from back (least recently used)
	for e := c.lruList.Back(); e != nil; {
		entry := e.Value.(*cacheEntry)
		prev := e.Prev()

		if now.After(entry.expiresAt) {
			c.removeEntry(entry)
			count++
		}

		e = prev
	}

	return count
}

// StartCleanupWorker starts a background goroutine to cleanup expired entries
func (c *LocalCache) StartCleanupWorker(interval time.Duration) chan struct{} {
	stopCh := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				removed := c.CleanupExpired()
				if removed > 0 {
					// Optional: log cleanup stats
					// log.Printf("LocalCache: cleaned up %d expired entries", removed)
				}
			case <-stopCh:
				return
			}
		}
	}()

	return stopCh
}

// removeEntry removes an entry from cache (must be called with lock held)
func (c *LocalCache) removeEntry(entry *cacheEntry) {
	c.lruList.Remove(entry.element)
	delete(c.cache, entry.key)
}

// evictOldest removes the least recently used entry
func (c *LocalCache) evictOldest() {
	oldest := c.lruList.Back()
	if oldest != nil {
		entry := oldest.Value.(*cacheEntry)
		c.removeEntry(entry)
	}
}

// CacheStats represents cache statistics
type CacheStats struct {
	Hits    uint64  `json:"hits"`
	Misses  uint64  `json:"misses"`
	HitRate float64 `json:"hit_rate"` // Percentage
	Size    int     `json:"size"`
	MaxSize int     `json:"max_size"`
}