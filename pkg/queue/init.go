package queue

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

const (
	// poolWaterMill is the identifier for Watermill queue provider
	poolWaterMill = "watermill"
	// poolWork is the identifier for gocraft/work queue provider (default)
	poolWork = "work"
)

var (
	// namespace is the Redis namespace prefix for all queues.
	// Set via APP_NAME environment variable. Used to isolate queues from different applications.
	namespace = os.Getenv("APP_NAME")
	
	// poolProvider determines which queue provider to use.
	// Checks POOL_PROVIDER environment variable, defaults to "work" (gocraft/work).
	poolProvider = func() string {
		if os.Getenv("POOL_PROVIDER") != "" && os.Getenv("POOL_PROVIDER") == poolWaterMill {
			return poolWaterMill
		}
		return poolWork
	}

	// redisPool is the singleton Redis connection pool.
	// WHY singleton?
	//   - Connection pools are expensive to create
	//   - One pool is sufficient for the entire application
	//   - Thread-safe initialization with sync.Once
	redisPool *redis.Pool
	
	// taskInstance is the singleton task enqueuer instance.
	taskInstance *Task
	
	// maxConcurrent is the default maximum number of concurrent workers.
	maxConcurrent uint = 10
	
	// queueOnce ensures Redis pool is created only once (thread-safe).
	queueOnce sync.Once
	
	// taskOnce ensures task enqueuer is created only once (thread-safe).
	taskOnce sync.Once
)

// Dispatcher defines the interface for dispatching jobs to the queue.
// This interface allows different queue implementations while maintaining the same API.
type Dispatcher[T any] interface {
	Dispatch() error              // Dispatch a job (may create duplicates)
	WithData(data *T)            // Set the job payload
	DispatchUnique() error       // Dispatch a unique job (prevents duplicates)
}

// Listen defines the interface for consuming jobs from the queue.
// Workers implement this interface to process jobs.
type Listen[T any] interface {
	RunWithContext(f func(ctx context.Context, data *T) error) // Start processing jobs
	Stop()                                                       // Stop processing jobs
}

// Task wraps the gocraft/work enqueuer for dispatching jobs.
type Task struct {
	enqueue *work.Enqueuer
}

// instancePool returns the singleton Redis connection pool.
// Creates the pool on first call (thread-safe).
//
// WHY singleton?
//   - Connection pools are expensive to create
//   - Multiple pools would waste resources
//   - All queue operations can share one pool
func instancePool() *redis.Pool {
	queueOnce.Do(func() {
		redisPool = newPoolRedis()
	})
	return redisPool
}

// newPoolRedis creates a new Redis connection pool.
// Reads connection details from environment variables:
//   - REDIS_HOST: Redis server host
//   - REDIS_PORT: Redis server port
//   - REDIS_DB: Database number (defaults to 0)
//   - REDIS_PASSWORD: Redis password (optional)
//
// WHY connection pool?
//   - Reuses connections instead of creating new ones for each operation
//   - Improves performance and reduces connection overhead
//   - Manages connection lifecycle automatically
func newPoolRedis() *redis.Pool {
	return &redis.Pool{
		Wait: false, // Don't wait for connection if pool is exhausted (fail fast)
		Dial: func() (redis.Conn, error) {
			// Parse database number from environment
			dbNumber, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
			
			// Create new connection with configuration
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
				redis.DialDatabase(dbNumber),
				redis.DialPassword(os.Getenv("REDIS_PASSWORD")),
			)
		},
	}
}

// Ping checks the connection to Redis by attempting to dial a new connection.
// Useful for health checks or verifying Redis availability.
//
// Returns an error if connection fails, nil if successful.
func Ping() error {
	_, err := instancePool().Dial()
	return err
}

// initQueue initializes and returns the singleton task enqueuer.
// Creates the enqueuer on first call (thread-safe).
//
// WHY singleton?
//   - Enqueuer is lightweight but should be shared
//   - Ensures consistent namespace across all queue operations
//   - Thread-safe initialization
func initQueue() *Task {
	taskOnce.Do(func() {
		taskInstance = &Task{
			enqueue: work.NewEnqueuer(namespace, instancePool()),
		}
	})
	return taskInstance
}