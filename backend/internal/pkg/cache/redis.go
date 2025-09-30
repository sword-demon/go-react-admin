// Package cache provides Redis cache utilities
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps redis.Client with additional utilities
type RedisClient struct {
	*redis.Client
}

// Config represents Redis configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	// Connection pool settings
	PoolSize     int
	MinIdleConns int
}

// InitRedis initializes Redis connection
func InitRedis(cfg *Config) (*RedisClient, error) {
	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("✅ Redis connected successfully: %s:%d (DB: %d)", cfg.Host, cfg.Port, cfg.DB)
	return &RedisClient{Client: rdb}, nil
}

// Close closes the Redis connection
func (c *RedisClient) Close() error {
	if err := c.Client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis: %w", err)
	}
	log.Println("✅ Redis connection closed")
	return nil
}

// DefaultConfig returns default Redis configuration
func DefaultConfig() *Config {
	return &Config{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
	}
}

// Set sets a key-value pair with expiration
func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

// Get gets a value by key
func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

// Del deletes keys
func (c *RedisClient) Del(ctx context.Context, keys ...string) error {
	return c.Client.Del(ctx, keys...).Err()
}

// Exists checks if key exists
func (c *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.Client.Exists(ctx, keys...).Result()
}

// SetNX sets a key only if it doesn't exist (for distributed lock)
func (c *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.Client.SetNX(ctx, key, value, expiration).Result()
}

// Expire sets expiration for a key
func (c *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.Client.Expire(ctx, key, expiration).Err()
}

// Keys returns keys matching pattern
func (c *RedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.Client.Keys(ctx, pattern).Result()
}

// FlushDB clears current database (use carefully!)
func (c *RedisClient) FlushDB(ctx context.Context) error {
	return c.Client.FlushDB(ctx).Err()
}

// SetJSON stores a Go object as JSON string
func (c *RedisClient) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return c.Set(ctx, key, data, expiration)
}

// GetJSON retrieves a JSON string and unmarshals into target
// target must be a pointer
func (c *RedisClient) GetJSON(ctx context.Context, key string, target interface{}) error {
	data, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), target)
}

// DeleteByPrefix deletes all keys matching the prefix pattern
// Returns the number of keys deleted
func (c *RedisClient) DeleteByPrefix(ctx context.Context, prefix string) (int, error) {
	keys, err := c.Keys(ctx, prefix+"*")
	if err != nil {
		return 0, err
	}
	if len(keys) == 0 {
		return 0, nil
	}
	if err := c.Del(ctx, keys...); err != nil {
		return 0, err
	}
	return len(keys), nil
}
