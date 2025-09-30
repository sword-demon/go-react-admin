package permission

import (
	"context"
	"fmt"
	"strings"

	"github.com/sword-demon/go-react-admin/internal/admin/biz"
	"github.com/sword-demon/go-react-admin/internal/admin/store"
	"github.com/sword-demon/go-react-admin/internal/pkg/cache"
)

type permissionBiz struct {
	store store.IStore
	cache *cache.RedisClient
}

func NewPermissionBiz(store store.IStore, cache *cache.RedisClient) biz.IPermissionBiz {
	return &permissionBiz{
		store: store,
		cache: cache,
	}
}

// GetUserPermissions retrieves all permission patterns for a user
func (b *permissionBiz) GetUserPermissions(ctx context.Context, userID uint64) ([]string, error) {
	// TODO: Check cache first
	// if b.cache != nil {
	// 	  cacheKey := fmt.Sprintf("user:permissions:%d", userID)
	//    cached, err := b.cache.Get(ctx, cacheKey)
	//    if err == nil && cached != "" {
	//        return strings.Split(cached, ","), nil
	//    }
	// }

	// Get from database
	permissions, err := b.store.Permissions().GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// TODO: Cache result
	// if b.cache != nil {
	//     cacheKey := fmt.Sprintf("user:permissions:%d", userID)
	//     b.cache.Set(ctx, cacheKey, strings.Join(permissions, ","), 30*time.Minute)
	// }

	return permissions, nil
}

// CheckPermission checks if user has permission using pattern matching
func (b *permissionBiz) CheckPermission(ctx context.Context, userID uint64, pattern string) (bool, error) {
	// Get user's permission patterns
	permissions, err := b.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	// Check if user has matching permission
	return matchPermission(permissions, pattern), nil
}

// ClearCache clears user permission cache
func (b *permissionBiz) ClearCache(ctx context.Context, userID uint64) error {
	if b.cache == nil {
		return nil
	}

	cacheKey := fmt.Sprintf("user:permissions:%d", userID)
	return b.cache.Del(ctx, cacheKey)
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
// User pattern examples: *:*, user:*, user:read, /api/users/*
// Request pattern examples: user:create, /api/users/123
func matchPattern(userPattern, requestPattern string) bool {
	// Global permission: *:*
	if userPattern == "*:*" {
		return true
	}

	// Exact match
	if userPattern == requestPattern {
		return true
	}

	// Module wildcard: user:* matches user:create, user:read, etc.
	if strings.HasSuffix(userPattern, ":*") {
		module := strings.TrimSuffix(userPattern, ":*")
		if strings.HasPrefix(requestPattern, module+":") {
			return true
		}
	}

	// Path wildcard: /api/users/* matches /api/users/123, /api/users/create, etc.
	if strings.HasSuffix(userPattern, "/*") {
		prefix := strings.TrimSuffix(userPattern, "/*")
		if strings.HasPrefix(requestPattern, prefix+"/") || requestPattern == prefix {
			return true
		}
	}

	return false
}
