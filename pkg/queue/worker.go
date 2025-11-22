package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gocraft/work"
)

const (
	// defaultMaxFails is the default maximum number of retry attempts before giving up.
	defaultMaxFails = 3
	// defaultMaxConcurrency is the default maximum number of concurrent jobs.
	defaultMaxConcurrency = 10
)

// Worker processes jobs from a queue.
// It consumes jobs dispatched by Queue and executes them using a worker pool.
//
// Generic type [T any] matches the payload type used in Queue:
//   worker := queue.NewWorker[MyPayload]("my-queue")
//   worker.RunWithContext(func(ctx context.Context, data *MyPayload) error {
//       // Process the job
//       return nil
//   })
type Worker[T any] struct {
	queueName string          // Name of the queue to consume from
	pool      *work.WorkerPool // Worker pool that processes jobs
	payload   *T              // Current job payload (deserialized from queue)
	options   *Options        // Worker configuration options
}

// NewWorker creates a new Worker instance for processing jobs from the specified queue.
//
// Options can be provided using option functions:
//   worker := queue.NewWorker[MyPayload]("my-queue",
//       queue.WithMaxFails(5),
//       queue.WithMaxConcurrency(10),
//   )
func NewWorker[T any](queueName string, ops ...func(options *Options)) *Worker[T] {
	// Start with default options
	options := &Options{
		MaxFails:       defaultMaxFails,
		MaxConcurrency: defaultMaxConcurrency,
	}

	// Apply provided options
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

// RunWithContext starts processing jobs from the queue.
// The provided function is called for each job with the job's context and payload.
//
// This method is blocking - it will process jobs until Stop() is called.
// Typically called in a goroutine:
//   go worker.RunWithContext(processJob)
//
// The context passed to the handler function:
//   - Can be cancelled if MaxTimeout is set
//   - Should be checked for cancellation: if ctx.Done() is closed, stop processing
//
// Example:
//   worker.RunWithContext(func(ctx context.Context, data *MyPayload) error {
//       // Check if context is cancelled
//       if ctx.Err() != nil {
//           return ctx.Err()
//       }
//       // Process the job
//       return processData(data)
//   })
func (w *Worker[T]) RunWithContext(f func(ctx context.Context, data *T) error) {
	// Register the job handler with the worker pool
	w.pool.JobWithOptions(w.queueName, w.getOptions(), func(job *work.Job) error {
		// Get context (with timeout if configured)
		ctxWorker, cancel := w.getContext()
		defer cancel() // Always cancel to free resources

		// Create new payload instance for this job
		w.payload = new(T)
		
		// Deserialize payload from job arguments
		if err := w.deserialize(job.ArgString("payload")); err != nil {
			return err // Return error to trigger retry logic
		}

		// Call the user-provided handler
		return f(ctxWorker, w.payload)
	})

	// Start the worker pool (this blocks until Stop() is called)
	w.pool.Start()
}

// Stop gracefully stops the worker pool.
// Stops accepting new jobs and waits for current jobs to finish.
// Should be called during application shutdown.
func (w *Worker[T]) Stop() {
	w.pool.Stop()
}

// deserialize converts the JSON string payload back to the typed struct.
// This is the reverse of Queue.serialize().
//
// Returns an error if JSON is invalid or doesn't match the expected type.
func (w *Worker[T]) deserialize(data string) error {
	return json.Unmarshal([]byte(data), &w.payload)
}

// getOptions converts internal Options to gocraft/work JobOptions.
// Only includes options that are set (non-zero values).
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

// getContext returns a context for job processing.
// If MaxTimeout is set, returns a context with timeout.
// Otherwise, returns a background context with a no-op cancel function.
//
// WHY timeout?
//   - Prevents jobs from running indefinitely
//   - Allows graceful cancellation of long-running jobs
//   - Protects against resource leaks
func (w *Worker[T]) getContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()

	// Set timeout if configured
	if w.options.MaxTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(w.options.MaxTimeout)*time.Second)
		return ctx, cancel
	}

	// No timeout - return background context with no-op cancel
	return ctx, func() {}
}
