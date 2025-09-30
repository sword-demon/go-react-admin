// Package config provides configuration loading utilities
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/sword-demon/go-react-admin/internal/pkg/cache"
	"github.com/sword-demon/go-react-admin/internal/pkg/db"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"` // debug, release, test
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"name"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"` // in minutes
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	Expiration int    `yaml:"expiration"` // in hours
}

// Load loads configuration from YAML file
func Load(configPath string) (*Config, error) {
	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	cfg.setDefaults()

	return &cfg, nil
}

// setDefaults sets default values for configuration
func (c *Config) setDefaults() {
	// Server defaults
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.Mode == "" {
		c.Server.Mode = "debug"
	}

	// Database connection pool defaults
	if c.Database.MaxIdleConns == 0 {
		c.Database.MaxIdleConns = 10
	}
	if c.Database.MaxOpenConns == 0 {
		c.Database.MaxOpenConns = 100
	}
	if c.Database.ConnMaxLifetime == 0 {
		c.Database.ConnMaxLifetime = 60 // 60 minutes
	}

	// Redis connection pool defaults
	if c.Redis.PoolSize == 0 {
		c.Redis.PoolSize = 10
	}
	if c.Redis.MinIdleConns == 0 {
		c.Redis.MinIdleConns = 5
	}

	// JWT defaults
	if c.JWT.Secret == "" {
		c.JWT.Secret = "go-react-admin-secret-key-change-in-production"
	}
	if c.JWT.Expiration == 0 {
		c.JWT.Expiration = 24 // 24 hours
	}
}

// ToDBConfig converts DatabaseConfig to db.Config
func (c *DatabaseConfig) ToDBConfig() *db.Config {
	return &db.Config{
		Host:            c.Host,
		Port:            c.Port,
		Username:        c.Username,
		Password:        c.Password,
		Database:        c.Database,
		MaxIdleConns:    c.MaxIdleConns,
		MaxOpenConns:    c.MaxOpenConns,
		ConnMaxLifetime: time.Duration(c.ConnMaxLifetime) * time.Minute,
	}
}

// ToRedisConfig converts RedisConfig to cache.Config
func (c *RedisConfig) ToRedisConfig() *cache.Config {
	return &cache.Config{
		Host:         c.Host,
		Port:         c.Port,
		Password:     c.Password,
		DB:           c.DB,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
	}
}

// Default returns default configuration
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            3306,
			Username:        "root",
			Password:        "",
			Database:        "go_react_admin",
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: 60,
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
			DB:   0,
		},
		JWT: JWTConfig{
			Secret:     "go-react-admin-secret-key-change-in-production",
			Expiration: 24,
		},
	}
}
