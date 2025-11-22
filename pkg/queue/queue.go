package queue

import (
	"encoding/json"

	"github.com/gocraft/work"
)

// Queue represents a job queue with a typed payload.
// It provides methods to dispatch jobs to Redis-backed queues using gocraft/work.
//
// Generic type [T any] allows type-safe job payloads:
//   type UserPayload struct { ID int; Name string }
//   q := queue.NewQueue[UserPayload]("user-queue")
type Queue[T any] struct {
	task      *Task  // Task enqueuer for dispatching jobs
	queueName string // Name of the queue (e.g., "user-created", "email-send")
	payload   *T     // The job payload data (set via WithData)
}

// NewQueue creates a new Queue instance for the given queue name.
//
// Example:
//   q := queue.NewQueue[MyPayload]("my-queue")
//   q.WithData(&MyPayload{ID: 1})
//   q.Dispatch()
func NewQueue[T any](queueName string) *Queue[T] {
	return &Queue[T]{
		task:      initQueue(),
		queueName: queueName,
	}
}

// WithData sets the payload data for the job.
// This must be called before Dispatch() or DispatchUnique().
//
// Example:
//   q.WithData(&MyPayload{ID: 1, Name: "test"})
func (q *Queue[T]) WithData(data *T) {
	q.payload = data
}

// Dispatch adds a job to the queue.
// This method may create duplicate jobs if called multiple times with the same data.
// Use DispatchUnique() if you want to prevent duplicates.
//
// Returns an error if the job could not be enqueued (e.g., Redis connection error).
func (q *Queue[T]) Dispatch() error {
	_, err := q.task.enqueue.Enqueue(q.queueName, work.Q{
		"payload": q.serialize(),
	})
	return err
}

// DispatchUnique adds a unique job to the queue.
// If a job with the same payload already exists, it won't create a duplicate.
// Useful for idempotent operations (e.g., sending welcome email only once).
//
// WHY unique jobs?
//   - Prevents duplicate processing
//   - Useful for idempotent operations
//   - Reduces unnecessary work
//
// Returns an error if the job could not be enqueued.
func (q *Queue[T]) DispatchUnique() error {
	_, err := q.task.enqueue.EnqueueUnique(q.queueName, work.Q{
		"payload": q.serialize(),
	})
	return err
}

// serialize converts the payload to JSON string for storage in Redis.
// The payload is stored as JSON so it can be deserialized by workers.
//
// WHY JSON?
//   - Human-readable format (useful for debugging)
//   - Language-agnostic (can be read by other services)
//   - Standard format supported by Redis
//
// NOTE: Errors during marshaling are ignored (returns empty string).
// In production, you might want to handle this error explicitly.
func (q *Queue[T]) serialize() string {
	b, _ := json.Marshal(q.payload)
	return string(b)
}