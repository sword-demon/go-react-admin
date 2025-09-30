// Package cache_test provides benchmark tests for cache system
package cache_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/sword-demon/go-react-admin/internal/pkg/cache"
)

// BenchmarkLocalCacheGet tests local cache read performance
func BenchmarkLocalCacheGet(b *testing.B) {
	localCache := cache.NewLocalCache(10000, 5*time.Minute)

	// Preload cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("user:permissions:%d", i)
		localCache.Set(key, []string{"user:read", "user:write"})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("user:permissions:%d", i%1000)
			_, _ = localCache.Get(key)
			i++
		}
	})
}

// BenchmarkLocalCacheSet tests local cache write performance
func BenchmarkLocalCacheSet(b *testing.B) {
	localCache := cache.NewLocalCache(10000, 5*time.Minute)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("user:permissions:%d", i)
			localCache.Set(key, []string{"user:read", "user:write"})
			i++
		}
	})
}

// BenchmarkLocalCacheConcurrent tests concurrent read/write
func BenchmarkLocalCacheConcurrent(b *testing.B) {
	localCache := cache.NewLocalCache(10000, 5*time.Minute)

	// Preload cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("user:permissions:%d", i)
		localCache.Set(key, []string{"user:read", "user:write"})
	}

	b.ResetTimer()

	var wg sync.WaitGroup
	readers := 8
	writers := 2

	// Start readers
	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < b.N/readers; i++ {
				key := fmt.Sprintf("user:permissions:%d", i%1000)
				_, _ = localCache.Get(key)
			}
		}()
	}

	// Start writers
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < b.N/writers; i++ {
				key := fmt.Sprintf("user:permissions:%d", i%1000)
				localCache.Set(key, []string{"user:read", "user:write"})
			}
		}()
	}

	wg.Wait()
}

// TestCacheHitRateSimulation simulates real-world usage pattern
func TestCacheHitRateSimulation(t *testing.T) {
	localCache := cache.NewLocalCache(1000, 5*time.Minute)

	// Simulate 10000 requests with Zipf distribution (80-20 rule)
	// 80% of requests access 20% of users
	hotUsers := 200   // Top 20% of 1000 users
	totalUsers := 1000
	totalRequests := 10000

	for i := 0; i < totalRequests; i++ {
		var userID int

		// 80% probability to access hot users
		if i%10 < 8 {
			userID = i % hotUsers
		} else {
			userID = hotUsers + (i % (totalUsers - hotUsers))
		}

		key := fmt.Sprintf("user:permissions:%d", userID)

		// Try to get from cache
		if _, ok := localCache.Get(key); !ok {
			// Cache miss - simulate database load
			permissions := []string{"user:read", "user:write"}
			localCache.Set(key, permissions)
		}
	}

	stats := localCache.Stats()
	t.Logf("Cache Hit Rate: %.2f%%", stats.HitRate)
	t.Logf("Hits: %d, Misses: %d", stats.Hits, stats.Misses)
	t.Logf("Cache Size: %d/%d", stats.Size, stats.MaxSize)

	// Assert hit rate is above 70% (expected for 80-20 distribution)
	if stats.HitRate < 70.0 {
		t.Errorf("Hit rate too low: %.2f%%, expected >= 70%%", stats.HitRate)
	}
}

// BenchmarkLRUEviction tests LRU eviction performance
func BenchmarkLRUEviction(b *testing.B) {
	cacheSize := 100
	localCache := cache.NewLocalCache(cacheSize, 5*time.Minute)

	// Fill cache to capacity
	for i := 0; i < cacheSize; i++ {
		key := fmt.Sprintf("user:permissions:%d", i)
		localCache.Set(key, []string{"user:read"})
	}

	b.ResetTimer()

	// Keep inserting new items to trigger eviction
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("user:permissions:%d", cacheSize+i)
		localCache.Set(key, []string{"user:read"})
	}
}

// TestCachePrefixDelete tests batch deletion performance
func TestCachePrefixDelete(t *testing.T) {
	localCache := cache.NewLocalCache(10000, 5*time.Minute)

	// Create 1000 users, each with 5 cache entries
	userCount := 1000
	entriesPerUser := 5

	for userID := 1; userID <= userCount; userID++ {
		for entryID := 0; entryID < entriesPerUser; entryID++ {
			key := fmt.Sprintf("user:%d:cache:%d", userID, entryID)
			localCache.Set(key, "dummy_data")
		}
	}

	initialSize := localCache.Stats().Size
	t.Logf("Initial cache size: %d", initialSize)

	// Delete all cache for user 123
	start := time.Now()
	deleted := localCache.DeletePrefix("user:123:")
	duration := time.Since(start)

	t.Logf("Deleted %d entries in %v", deleted, duration)

	if deleted != entriesPerUser {
		t.Errorf("Expected to delete %d entries, but deleted %d", entriesPerUser, deleted)
	}

	finalSize := localCache.Stats().Size
	if finalSize != initialSize-deleted {
		t.Errorf("Cache size mismatch: expected %d, got %d", initialSize-deleted, finalSize)
	}
}

// BenchmarkWarmupConcurrency tests cache warmup with different concurrency levels
func BenchmarkWarmupConcurrency(b *testing.B) {
	concurrencyLevels := []int{1, 5, 10, 20}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			localCache := cache.NewLocalCache(10000, 5*time.Minute)

			// Mock loader
			loader := &mockPermissionLoader{
				localCache: localCache,
			}

			config := &cache.WarmupConfig{
				SuperAdminUserIDs: generateUserIDs(100),
				CommonRoleIDs:     []uint64{1, 2, 3, 4, 5},
				Concurrency:       concurrency,
				Timeout:           30 * time.Second,
				EnableLogging:     false,
			}

			warmer := cache.NewPermissionWarmer(config, loader)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = warmer.Warm(context.Background())
			}
		})
	}
}

// mockPermissionLoader mocks permission loading for tests
type mockPermissionLoader struct {
	localCache *cache.LocalCache
}

func (m *mockPermissionLoader) LoadUserPermissions(ctx context.Context, userID uint64) ([]string, error) {
	// Simulate database query delay
	time.Sleep(5 * time.Millisecond)

	key := fmt.Sprintf("user:permissions:%d", userID)
	permissions := []string{"user:read", "user:write", "user:delete"}
	m.localCache.Set(key, permissions)

	return permissions, nil
}

func (m *mockPermissionLoader) LoadRolePermissions(ctx context.Context, roleID uint64) ([]string, error) {
	// Simulate database query delay
	time.Sleep(5 * time.Millisecond)

	key := fmt.Sprintf("role:permissions:%d", roleID)
	permissions := []string{"role:read", "role:write"}
	m.localCache.Set(key, permissions)

	return permissions, nil
}

// generateUserIDs generates a slice of user IDs
func generateUserIDs(count int) []uint64 {
	ids := make([]uint64, count)
	for i := 0; i < count; i++ {
		ids[i] = uint64(i + 1)
	}
	return ids
}
