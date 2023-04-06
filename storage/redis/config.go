package redis

import (
	"crypto/tls"
	"runtime"
	"time"
)

// Config defines the config for storage.
type Config struct {
	// Either a single address or a seed list of host:port addresses
	// of cluster/sentinel nodes.
	Addrs []string `yaml:"addrs"`

	// Server username
	//
	// Optional. Default is ""
	Username string `yaml:"username"`

	// Server password
	//
	// Optional. Default is ""
	Password string `yaml:"password"`

	// Database to be selected after connecting to the server.
	//
	// Optional. Default is 0
	Database int `yaml:"database"`

	// Reset clears any existing keys in existing Collection
	//
	// Optional. Default is false
	Reset bool `yaml:"reset"`

	// TLS Config to use. When set TLS will be negotiated.
	TLSConfig *tls.Config

	// Maximum number of socket connections.
	//
	// Optional. Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int `yaml:"poolSize"`

	////////////////////////////////////
	// Adaptor related config options //
	////////////////////////////////////
	SentinelPassword string `yaml:"sentinelPassword"`

	MaxRetries      int           `yaml:"maxRetries"`
	MinRetryBackoff time.Duration `yaml:"minRetryBackoff"`
	MaxRetryBackoff time.Duration `yaml:"maxRetryBackoff"`

	DialTimeout  time.Duration `yaml:"dialTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`

	// PoolFIFO uses FIFO mode for each node connection pool GET/PUT (default LIFO).
	PoolFIFO bool

	MinIdleConns       int           `yaml:"minIdleConns"`
	MaxConnAge         time.Duration `yaml:"maxConnAge"`
	PoolTimeout        time.Duration `yaml:"poolTimeout"`
	IdleTimeout        time.Duration `yaml:"idleTimeout"`
	IdleCheckFrequency time.Duration `yaml:"idleCheckFrequency"`

	// Only cluster clients.

	MaxRedirects   int  `yaml:"maxRedirects"`
	ReadOnly       bool `yaml:"readOnly"`
	RouteByLatency bool `yaml:"routeByLatency"`
	RouteRandomly  bool `yaml:"routeRandomly"`

	// The sentinel master name.
	// Only failover clients.
	MasterName string `yaml:"masterName"`

	// https://pkg.go.dev/github.com/go-redis/redis/v9#Options
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Addrs:              []string{"127.0.0.1:6379"},
	Username:           "",
	Password:           "",
	Database:           0,
	Reset:              false,
	TLSConfig:          nil,
	PoolSize:           10 * runtime.GOMAXPROCS(0),
	SentinelPassword:   "",
	MaxRetries:         0,
	MinRetryBackoff:    0,
	MaxRetryBackoff:    0,
	DialTimeout:        0,
	ReadTimeout:        0,
	WriteTimeout:       0,
	PoolFIFO:           true,
	MinIdleConns:       0,
	MaxConnAge:         0,
	PoolTimeout:        0,
	IdleTimeout:        0,
	IdleCheckFrequency: 0,
	MaxRedirects:       0,
	ReadOnly:           false,
	RouteByLatency:     false,
	RouteRandomly:      false,
	MasterName:         "",
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if len(cfg.Addrs) < 1 {
		cfg.Addrs = ConfigDefault.Addrs
	}

	return cfg
}
