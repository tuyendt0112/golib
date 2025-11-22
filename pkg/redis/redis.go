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

// Options is the redis options
type Options struct {
	Host     string
	Port     string
	Password string
	DB       int
	MaxRetry int
}

func NewClientRedis(ops ...*Options) *redis.Client {
	if len(ops) == 0 {
		db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
		ops[0] = &Options{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			DB:       db,
			MaxRetry: 3,
		}
	}

	redisOnce.Do(func() {
		opts := ops[0]

		instanceRedis = redis.NewClient(&redis.Options{
			Addr:       fmt.Sprintf("%s:%s", opts.Host, opts.Port),
			Password:   opts.Password,
			DB:         opts.DB,
			MaxRetries: opts.MaxRetry,
		})

	})

	return instanceRedis
}
