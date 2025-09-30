// Package cache provides cache warming utilities
package cache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Warmer defines cache warming interface
type Warmer interface {
	// Warm preloads cache with hot data
	Warm(ctx context.Context) error
}

// WarmupConfig defines cache warmup configuration
type WarmupConfig struct {
	// SuperAdminUserIDs are IDs of super admin users to preload
	SuperAdminUserIDs []uint64

	// CommonRoleIDs are IDs of frequently used roles
	CommonRoleIDs []uint64

	// Concurrency controls parallel warmup goroutines
	Concurrency int

	// Timeout for entire warmup process
	Timeout time.Duration

	// EnableLogging enables warmup progress logs
	EnableLogging bool
}

// DefaultWarmupConfig returns default warmup configuration
func DefaultWarmupConfig() *WarmupConfig {
	return &WarmupConfig{
		SuperAdminUserIDs: []uint64{1}, // Default: user ID 1 is super admin
		CommonRoleIDs:     []uint64{1, 2, 3}, // Admin, Manager, User
		Concurrency:       5,
		Timeout:           30 * time.Second,
		EnableLogging:     true,
	}
}

// WarmupStats represents cache warmup statistics
type WarmupStats struct {
	TotalItems    int           `json:"total_items"`
	SuccessCount  int           `json:"success_count"`
	FailureCount  int           `json:"failure_count"`
	Duration      time.Duration `json:"duration_ms"`
	ErrorMessages []string      `json:"errors,omitempty"`
}

// WarmupResult represents the result of a single warmup operation
type warmupResult struct {
	itemType string // "user" or "role"
	id       uint64
	err      error
}

// PermissionWarmer implements cache warming for permissions
type PermissionWarmer struct {
	config *WarmupConfig
	loader PermissionLoader // Delegate to load permissions
}

// PermissionLoader defines interface to load permissions
type PermissionLoader interface {
	// LoadUserPermissions loads permissions for a user
	LoadUserPermissions(ctx context.Context, userID uint64) ([]string, error)

	// LoadRolePermissions loads permissions for a role
	LoadRolePermissions(ctx context.Context, roleID uint64) ([]string, error)
}

// NewPermissionWarmer creates a new permission warmer
func NewPermissionWarmer(config *WarmupConfig, loader PermissionLoader) *PermissionWarmer {
	if config == nil {
		config = DefaultWarmupConfig()
	}
	if config.Concurrency <= 0 {
		config.Concurrency = 5
	}
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}

	return &PermissionWarmer{
		config: config,
		loader: loader,
	}
}

// Warm preloads permissions into cache
// This should be called on application startup
func (w *PermissionWarmer) Warm(ctx context.Context) (*WarmupStats, error) {
	startTime := time.Now()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, w.config.Timeout)
	defer cancel()

	stats := &WarmupStats{
		TotalItems:    len(w.config.SuperAdminUserIDs) + len(w.config.CommonRoleIDs),
		ErrorMessages: make([]string, 0),
	}

	if w.config.EnableLogging {
		log.Printf("🔥 Starting cache warmup: %d users, %d roles",
			len(w.config.SuperAdminUserIDs), len(w.config.CommonRoleIDs))
	}

	// Create worker pool
	resultCh := make(chan warmupResult, stats.TotalItems)
	var wg sync.WaitGroup

	// Worker pool for concurrent warmup
	semaphore := make(chan struct{}, w.config.Concurrency)

	// Warm up user permissions
	for _, userID := range w.config.SuperAdminUserIDs {
		wg.Add(1)
		go func(uid uint64) {
			defer wg.Done()

			semaphore <- struct{}{} // Acquire
			defer func() { <-semaphore }() // Release

			if w.config.EnableLogging {
				log.Printf("  → Warming user permissions: userID=%d", uid)
			}

			_, err := w.loader.LoadUserPermissions(ctx, uid)
			resultCh <- warmupResult{
				itemType: "user",
				id:       uid,
				err:      err,
			}
		}(userID)
	}

	// Warm up role permissions
	for _, roleID := range w.config.CommonRoleIDs {
		wg.Add(1)
		go func(rid uint64) {
			defer wg.Done()

			semaphore <- struct{}{} // Acquire
			defer func() { <-semaphore }() // Release

			if w.config.EnableLogging {
				log.Printf("  → Warming role permissions: roleID=%d", rid)
			}

			_, err := w.loader.LoadRolePermissions(ctx, rid)
			resultCh <- warmupResult{
				itemType: "role",
				id:       rid,
				err:      err,
			}
		}(roleID)
	}

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results
	for result := range resultCh {
		if result.err != nil {
			stats.FailureCount++
			errMsg := fmt.Sprintf("%s:%d - %v", result.itemType, result.id, result.err)
			stats.ErrorMessages = append(stats.ErrorMessages, errMsg)

			if w.config.EnableLogging {
				log.Printf("  ✗ Failed: %s", errMsg)
			}
		} else {
			stats.SuccessCount++

			if w.config.EnableLogging {
				log.Printf("  ✓ Success: %s:%d", result.itemType, result.id)
			}
		}
	}

	stats.Duration = time.Since(startTime)

	if w.config.EnableLogging {
		log.Printf("🎉 Cache warmup completed: %d/%d success, %d failed, took %v",
			stats.SuccessCount, stats.TotalItems, stats.FailureCount, stats.Duration)
	}

	// Return error if all items failed
	if stats.FailureCount == stats.TotalItems && stats.TotalItems > 0 {
		return stats, fmt.Errorf("cache warmup failed: all %d items failed", stats.TotalItems)
	}

	return stats, nil
}

// WarmAsync starts warmup in background (non-blocking)
// Returns a channel that receives the result when warmup completes
func (w *PermissionWarmer) WarmAsync(ctx context.Context) <-chan *WarmupStats {
	resultCh := make(chan *WarmupStats, 1)

	go func() {
		stats, err := w.Warm(ctx)
		if err != nil && w.config.EnableLogging {
			log.Printf("⚠️  Cache warmup error: %v", err)
		}
		resultCh <- stats
		close(resultCh)
	}()

	return resultCh
}
