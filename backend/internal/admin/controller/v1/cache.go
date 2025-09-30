// Package v1 provides HTTP handlers for API v1
package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sword-demon/go-react-admin/internal/admin/biz"
	"github.com/sword-demon/go-react-admin/internal/pkg/cache"
)

// CacheController handles cache monitoring operations
type CacheController struct {
	permissionBiz biz.IPermissionBiz
	localCache    *cache.LocalCache
}

// NewCacheController creates a new cache controller
func NewCacheController(permissionBiz biz.IPermissionBiz, localCache *cache.LocalCache) *CacheController {
	return &CacheController{
		permissionBiz: permissionBiz,
		localCache:    localCache,
	}
}

// GetStats returns cache statistics
// GET /api/v1/cache/stats
func (c *CacheController) GetStats(ctx *gin.Context) {
	stats := c.localCache.Stats()

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"hits":      stats.Hits,
			"misses":    stats.Misses,
			"hit_rate":  stats.HitRate,
			"size":      stats.Size,
			"max_size":  stats.MaxSize,
			"health":    getHealthStatus(stats.HitRate),
		},
	})
}

// ClearCache clears specific cache entries
// DELETE /api/v1/cache
// Query params: user_id, role_id, prefix
func (c *CacheController) ClearCache(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	roleID := ctx.Query("role_id")
	prefix := ctx.Query("prefix")

	var clearedCount int

	if userID != "" {
		// Clear user cache
		var uid uint64
		if _, err := fmt.Sscanf(userID, "%d", &uid); err == nil {
			c.permissionBiz.ClearCache(ctx, uid)
			clearedCount++
		}
	}

	if roleID != "" {
		// Clear role cache
		var rid uint64
		if _, err := fmt.Sscanf(roleID, "%d", &rid); err == nil {
			c.permissionBiz.ClearRoleCache(ctx, rid)
			clearedCount++
		}
	}

	if prefix != "" {
		// Clear by prefix
		count := c.localCache.DeletePrefix(prefix)
		clearedCount += count
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "cache cleared",
		"data": gin.H{
			"cleared_count": clearedCount,
		},
	})
}

// WarmupCache triggers cache warmup
// POST /api/v1/cache/warmup
func (c *CacheController) WarmupCache(ctx *gin.Context) {
	// TODO: Implement warmup trigger
	// For now, return not implemented
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "warmup triggered",
		"data": gin.H{
			"status": "queued",
		},
	})
}

// getHealthStatus determines cache health based on hit rate
func getHealthStatus(hitRate float64) string {
	switch {
	case hitRate >= 80:
		return "excellent"
	case hitRate >= 60:
		return "good"
	case hitRate >= 40:
		return "fair"
	case hitRate >= 20:
		return "poor"
	default:
		return "critical"
	}
}
