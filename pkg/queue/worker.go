package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gocraft/work"
)

const (
	// defaultMaxFails is the default maximum number of fails.
	defaultMaxFails = 3
	// defaultMaxConcurrency is the default maximum concurrency.
	defaultMaxConcurrency = 10
)

// Worker is a struct that contains the queue name, worker pool, payload, and options.
type Worker[T any] struct {
	queueName string
	pool      *work.WorkerPool
	payload   *T
	options   *Options
}

// NewWorker is a function that returns a new Worker instance.
func NewWorker[T any](queueName string, ops ...func(options *Options)) *Worker[T] {
	options := &Options{
		MaxFails:       defaultMaxFails,
		MaxConcurrency: defaultMaxConcurrency,
	}

	for _, op := range ops {
		op(options)
	}

	return &Worker[T]{
		queueName: queueName,
		pool:      work.NewWorkerPool(context.Background(), maxConcurrent, namespace, instancePool()),
		payload:   new(T),
		options:   options,
	}
}

// RunWithContext is a function that runs the worker with context and data.
func (w *Worker[T]) RunWithContext(f func(ctx context.Context, data *T) error) {
	// Register the job with the worker pool
	w.pool.JobWithOptions(w.queueName, w.getOptions(), func(job *work.Job) error {
		ctxWorker, cancel := w.getContext()
		defer cancel()

		w.payload = new(T)
		if err := w.deserialize(job.ArgString("payload")); err != nil {
			return err
		}

		return f(ctxWorker, w.payload)
	})

	w.pool.Start()
}

// Stop is a function that stops the worker.
func (w *Worker[T]) Stop() {
	w.pool.Stop()
}

// deserialize is a function that deserializes the data.
func (w *Worker[T]) deserialize(data string) error {
	return json.Unmarshal([]byte(data), &w.payload)
}

// getOptions is a function that returns the job options.
func (w *Worker[T]) getOptions() work.JobOptions {
	ops := work.JobOptions{}

	if w.options.Priority > 0 {
		ops.Priority = w.options.Priority
	}

	if w.options.MaxFails > 0 {
		ops.MaxFails = w.options.MaxFails
	}

	if w.options.SkipDead {
		ops.SkipDead = w.options.SkipDead
	}

	if w.options.MaxConcurrency > 0 {
		ops.MaxConcurrency = w.options.MaxConcurrency
	}

	return ops
}

func (w *Worker[T]) getContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()

	// Set a timeout for the context if MaxTimeout is set
	if w.options.MaxTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(w.options.MaxTimeout)*time.Second)
		return ctx, cancel
	}

	return ctx, func() {}
}
