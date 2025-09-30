package permission

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sword-demon/go-react-admin/internal/admin/biz"
	"github.com/sword-demon/go-react-admin/internal/admin/store"
	"github.com/sword-demon/go-react-admin/internal/pkg/cache"
	"github.com/sword-demon/go-react-admin/internal/pkg/errors"
)

type permissionBiz struct {
	store      store.IStore
	localCache *cache.LocalCache
	redis      *cache.RedisClient
}

// NewPermissionBiz creates a new permission biz with three-tier cache
func NewPermissionBiz(store store.IStore, localCache *cache.LocalCache, redis *cache.RedisClient) biz.IPermissionBiz {
	return &permissionBiz{
		store:      store,
		localCache: localCache,
		redis:      redis,
	}
}

// GetUserPermissions retrieves all permission patterns for a user (with three-tier cache)
// Cache strategy:
// - L1 (Local): 5min TTL, <1ms latency
// - L2 (Redis): 30min TTL, <10ms latency
// - L3 (MySQL): persistent, 10-50ms latency
func (b *permissionBiz) GetUserPermissions(ctx context.Context, userID uint64) ([]string, error) {
	cacheKey := cache.PermissionCacheKey(userID)

	// Layer 1: Try local cache
	if val, ok := b.localCache.Get(cacheKey); ok {
		if permissions, ok := val.([]string); ok {
			return permissions, nil
		}
	}

	// Layer 2: Try Redis
	if b.redis != nil {
		data, err := b.redis.Get(ctx, cacheKey)
		if err == nil {
			var permissions []string
			if err := json.Unmarshal([]byte(data), &permissions); err == nil {
				// Backfill to local cache
				b.localCache.SetWithTTL(cacheKey, permissions, 5*time.Minute)
				return permissions, nil
			}
		}
		// Ignore redis.Nil (cache miss), propagate other errors
		if err != nil && err != redis.Nil {
			// Log error but continue to database
			// log.Printf("Redis error: %v", err)
		}
	}

	// Layer 3: Load from database
	permissions, err := b.store.Permissions().GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternalServer, "failed to get user permissions", err)
	}

	// Backfill to L2 (Redis) and L1 (Local)
	if b.redis != nil {
		data, _ := json.Marshal(permissions)
		_ = b.redis.Set(ctx, cacheKey, data, 30*time.Minute) // Ignore error
	}
	b.localCache.SetWithTTL(cacheKey, permissions, 5*time.Minute)

	return permissions, nil
}

// CheckPermission checks if user has permission using pattern matching (with cache)
func (b *permissionBiz) CheckPermission(ctx context.Context, userID uint64, pattern string) (bool, error) {
	// Get user's permission patterns (cached)
	permissions, err := b.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	// Check if user has matching permission
	return matchPermission(permissions, pattern), nil
}

// ClearCache clears user permission cache (all three tiers)
// Called when:
// - User roles changed
// - Role permissions changed
// - User deleted
func (b *permissionBiz) ClearCache(ctx context.Context, userID uint64) error {
	cacheKey := cache.PermissionCacheKey(userID)

	// Clear local cache
	b.localCache.Delete(cacheKey)

	// Clear Redis cache
	if b.redis != nil {
		if err := b.redis.Del(ctx, cacheKey); err != nil {
			return errors.Wrap(errors.ErrInternalServer, "failed to clear Redis cache", err)
		}
	}

	return nil
}

// ClearUserCache clears all cache related to a user
// This includes permissions, roles, profile, etc.
func (b *permissionBiz) ClearUserCache(ctx context.Context, userID uint64) error {
	prefix := cache.UserCacheKeyPrefix(userID)

	// Clear local cache
	b.localCache.DeletePrefix(prefix)

	// Clear Redis cache
	if b.redis != nil {
		if _, err := b.redis.DeleteByPrefix(ctx, prefix); err != nil {
			return errors.Wrap(errors.ErrInternalServer, "failed to clear user cache", err)
		}
	}

	return nil
}

// ClearRoleCache clears all cache related to a role
// Called when:
// - Role permissions changed
// - Role deleted
func (b *permissionBiz) ClearRoleCache(ctx context.Context, roleID uint64) error {
	// 1. Clear role permission cache
	roleCacheKey := cache.RolePermissionCacheKey(roleID)
	b.localCache.Delete(roleCacheKey)
	if b.redis != nil {
		_ = b.redis.Del(ctx, roleCacheKey)
	}

	// 2. Clear all users with this role
	// TODO: This requires getting all users with this role from database
	// For now, we can use a wildcard delete for role-related cache
	prefix := cache.RoleCacheKeyPrefix(roleID)
	b.localCache.DeletePrefix(prefix)
	if b.redis != nil {
		_, _ = b.redis.DeleteByPrefix(ctx, prefix)
	}

	return nil
}

// GetCacheStats returns cache statistics
func (b *permissionBiz) GetCacheStats() cache.CacheStats {
	return b.localCache.Stats()
}

// matchPermission checks if pattern matches any user permission
func matchPermission(userPermissions []string, requestPattern string) bool {
	for _, userPerm := range userPermissions {
		if matchPattern(userPerm, requestPattern) {
			return true
		}
	}
	return false
}

// matchPattern checks if a permission pattern matches the request pattern
// Priority (high → low):
// 1. Exact match: /api/users:GET
// 2. Path wildcard: /api/users/*
// 3. Module permission: user:read, user:write
// 4. Module wildcard: user:*
// 5. Global wildcard: *:*
//
// User pattern examples: *:*, user:*, user:read, /api/users/*
// Request pattern examples: user:create, /api/users/123
func matchPattern(userPattern, requestPattern string) bool {
	// Global permission: *:* (priority 5)
	if userPattern == "*:*" {
		return true
	}

	// Exact match (priority 1)
	if userPattern == requestPattern {
		return true
	}

	// Module wildcard: user:* matches user:create, user:read, etc. (priority 4)
	if strings.HasSuffix(userPattern, ":*") {
		module := strings.TrimSuffix(userPattern, ":*")
		if strings.HasPrefix(requestPattern, module+":") {
			return true
		}
	}

	// Path wildcard: /api/users/* matches /api/users/123, /api/users/create, etc. (priority 2)
	if strings.HasSuffix(userPattern, "/*") {
		prefix := strings.TrimSuffix(userPattern, "/*")
		if strings.HasPrefix(requestPattern, prefix+"/") || requestPattern == prefix {
			return true
		}
	}

	return false
}
