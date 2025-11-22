package redis

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	instanceRedis *redis.Client
	redisOnce     sync.Once
)

// Options contains configuration for the Redis client.
type Options struct {
	Host     string // Redis server host (e.g., "localhost")
	Port     string // Redis server port (e.g., "6379")
	Password string // Redis password (empty if no auth)
	DB       int    // Redis database number (0-15, default: 0)
	MaxRetry int    // Maximum number of retries for failed commands
}

// NewClientRedis creates or returns the singleton Redis client.
//
// Usage:
//   // Option 1: Use environment variables
//   // Set: REDIS_HOST, REDIS_PORT, REDIS_DB
//   client := redis.NewClientRedis()
//
//   // Option 2: Provide options explicitly
//   opts := &redis.Options{
//       Host:     "localhost",
//       Port:     "6379",
//       Password: "",
//       DB:       0,
//       MaxRetry: 3,
//   }
//   client := redis.NewClientRedis(opts)
//
// WHY singleton pattern?
//   - Efficient: One connection pool shared across the application
//   - Resource-friendly: Avoids creating multiple connections unnecessarily
//   - Thread-safe: sync.Once ensures safe concurrent access
//
// NOTE: The first call determines the configuration. Subsequent calls
// with different options will be ignored (due to sync.Once).
func NewClientRedis(ops ...*Options) *redis.Client {
	var opts *Options
	
	// If no options provided, read from environment variables
	if len(ops) == 0 {
		// Parse DB number from environment (defaults to 0 if invalid)
		db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
		opts = &Options{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			DB:       db,
			MaxRetry: 3, // Default retry count
		}
	} else {
		// Use provided options
		opts = ops[0]
	}

	// Create client only once (thread-safe)
	redisOnce.Do(func() {
		instanceRedis = redis.NewClient(&redis.Options{
			Addr:       fmt.Sprintf("%s:%s", opts.Host, opts.Port),
			Password:   opts.Password,
			DB:         opts.DB,
			MaxRetries: opts.MaxRetry,
		})
	})

	return instanceRedis
}
