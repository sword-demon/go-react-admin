// Package cache provides three-tier cache implementation
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ThreeTierCache implements three-tier cache: Local (L1) → Redis (L2) → Database (L3)
// Design goals:
// - L1: In-memory LRU, 5min TTL, 80%+ hit rate, <1ms latency
// - L2: Redis, 30min TTL, 95%+ hit rate, <10ms latency
// - L3: MySQL, persistent, 100% hit rate, 10-50ms latency
type ThreeTierCache struct {
	local *LocalCache
	redis *RedisClient

	// TTL configuration
	localTTL time.Duration
	redisTTL time.Duration
}

// NewThreeTierCache creates a new three-tier cache
func NewThreeTierCache(local *LocalCache, redis *RedisClient) *ThreeTierCache {
	return &ThreeTierCache{
		local:    local,
		redis:    redis,
		localTTL: 5 * time.Minute,
		redisTTL: 30 * time.Minute,
	}
}

// GetString retrieves a string value from three-tier cache
// Returns (value, cacheLevel, error)
// cacheLevel: "L1"=local, "L2"=redis, "L3"=database
func (c *ThreeTierCache) GetString(ctx context.Context, key string, dbLoader func() (string, error)) (string, string, error) {
	// Layer 1: Try local cache
	if val, ok := c.local.Get(key); ok {
		return val.(string), "L1", nil
	}

	// Layer 2: Try Redis
	if c.redis != nil {
		val, err := c.redis.Get(ctx, key)
		if err == nil {
			// Backfill to local cache
			c.local.Set(key, val)
			return val, "L2", nil
		}
		// Ignore redis.Nil (key not found), propagate other errors
		if err != redis.Nil {
			return "", "L2", fmt.Errorf("redis error: %w", err)
		}
	}

	// Layer 3: Load from database
	val, err := dbLoader()
	if err != nil {
		return "", "L3", err
	}

	// Backfill to L2 and L1
	if c.redis != nil {
		_ = c.redis.Set(ctx, key, val, c.redisTTL) // Ignore error
	}
	c.local.Set(key, val)

	return val, "L3", nil
}

// GetJSON retrieves a JSON object from three-tier cache
// target must be a pointer
func (c *ThreeTierCache) GetJSON(ctx context.Context, key string, target interface{}, dbLoader func() (interface{}, error)) (string, error) {
	// Layer 1: Try local cache
	if val, ok := c.local.Get(key); ok {
		// Deep copy to avoid mutation
		if ptr, ok := val.(interface{}); ok {
			*target.(*interface{}) = ptr
			return "L1", nil
		}
	}

	// Layer 2: Try Redis
	if c.redis != nil {
		err := c.redis.GetJSON(ctx, key, target)
		if err == nil {
			// Backfill to local cache
			c.local.Set(key, target)
			return "L2", nil
		}
		if err != redis.Nil {
			return "L2", fmt.Errorf("redis error: %w", err)
		}
	}

	// Layer 3: Load from database
	val, err := dbLoader()
	if err != nil {
		return "L3", err
	}

	// Backfill to L2 and L1
	if c.redis != nil {
		_ = c.redis.SetJSON(ctx, key, val, c.redisTTL)
	}
	c.local.Set(key, val)

	// Copy to target
	*target.(*interface{}) = val

	return "L3", nil
}

// Set stores a value in all cache tiers
func (c *ThreeTierCache) Set(ctx context.Context, key string, value interface{}) error {
	// Store in local cache
	c.local.Set(key, value)

	// Store in Redis
	if c.redis != nil {
		return c.redis.SetJSON(ctx, key, value, c.redisTTL)
	}

	return nil
}

// Delete removes a key from all cache tiers
func (c *ThreeTierCache) Delete(ctx context.Context, key string) error {
	// Delete from local cache
	c.local.Delete(key)

	// Delete from Redis
	if c.redis != nil {
		return c.redis.Del(ctx, key)
	}

	return nil
}

// DeletePrefix removes all keys with prefix from all cache tiers
// Returns (localCount, redisCount, error)
func (c *ThreeTierCache) DeletePrefix(ctx context.Context, prefix string) (int, int, error) {
	localCount := c.local.DeletePrefix(prefix)

	redisCount := 0
	var err error
	if c.redis != nil {
		redisCount, err = c.redis.DeleteByPrefix(ctx, prefix)
	}

	return localCount, redisCount, err
}

// Stats returns cache statistics for all tiers
func (c *ThreeTierCache) Stats() ThreeTierStats {
	stats := ThreeTierStats{
		Local: c.local.Stats(),
	}

	// TODO: Add Redis stats if needed
	// stats.Redis = RedisStats{...}

	return stats
}

// ThreeTierStats represents statistics for all cache tiers
type ThreeTierStats struct {
	Local CacheStats `json:"local"`
	// Redis RedisStats  `json:"redis"` // TODO: Add Redis stats
}

// PermissionCacheKey generates cache key for user permissions
func PermissionCacheKey(userID uint64) string {
	return fmt.Sprintf("user:permissions:%d", userID)
}

// RolePermissionCacheKey generates cache key for role permissions
func RolePermissionCacheKey(roleID uint64) string {
	return fmt.Sprintf("role:permissions:%d", roleID)
}

// UserCacheKeyPrefix returns prefix for all user-related cache keys
func UserCacheKeyPrefix(userID uint64) string {
	return fmt.Sprintf("user:%d:", userID)
}

// RoleCacheKeyPrefix returns prefix for all role-related cache keys
func RoleCacheKeyPrefix(roleID uint64) string {
	return fmt.Sprintf("role:%d:", roleID)
}