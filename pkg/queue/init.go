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
	poolWaterMill = "watermill"
	poolWork      = "work"
)

var (
	namespace    = os.Getenv("APP_NAME")
	poolProvider = func() string {
		if os.Getenv("POOL_PROVIDER") != "" && os.Getenv("POOL_PROVIDER") == poolWaterMill {
			return poolWaterMill
		}

		return poolWork
	}

	redisPool     *redis.Pool
	taskInstance  *Task
	maxConcurrent uint = 10
	queueOnce     sync.Once
	taskOnce      sync.Once
)

type Dispatcher[T any] interface {
	Dispatch() error
	WithData(data *T)
	DispatchUnique() error
}

type Listen[T any] interface {
	RunWithContext(f func(ctx context.Context, data *T) error)
	Stop()
}

type Task struct {
	enqueue *work.Enqueuer
}

func instancePool() *redis.Pool {
	queueOnce.Do(func() {
		redisPool = newPoolRedis()
	})

	return redisPool
}

func newPoolRedis() *redis.Pool {
	return &redis.Pool{
		Wait: false,
		Dial: func() (redis.Conn, error) {

			dbNumber, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
				redis.DialDatabase(dbNumber),
				redis.DialPassword(os.Getenv("REDIS_PASSWORD")),
			)
		},
	}
}

// Ping check connection to redis
func Ping() error {
	_, err := instancePool().Dial()
	return err
}

func initQueue() *Task {
	taskOnce.Do(func() {
		taskInstance = &Task{
			enqueue: work.NewEnqueuer(namespace, instancePool()),
		}
	})

	return taskInstance
}