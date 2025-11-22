package queue

import (
	"encoding/json"

	"github.com/gocraft/work"
)

// Queue is a struct that contains the task, queue name, and payload.
type Queue[T any] struct {
	task      *Task
	queueName string
	payload   *T
}

// NewQueue is a function that returns a new Queue instance.
func NewQueue[T any](queueName string) *Queue[T] {
	return &Queue[T]{
		task:      initQueue(),
		queueName: queueName,
	}
}

// WithData is a function that sets the payload of the Queue.
func (q *Queue[T]) WithData(data *T) {
	q.payload = data
}

// Dispatch is a function that dispatches the Queue.
func (q *Queue[T]) Dispatch() error {

	_, err := q.task.enqueue.Enqueue(q.queueName, work.Q{
		"payload": q.serialize(),
	})

	return err
}

// DispatchUnique is a function that dispatches the Queue uniquely.
func (q *Queue[T]) DispatchUnique() error {
	_, err := q.task.enqueue.EnqueueUnique(q.queueName, work.Q{
		"payload": q.serialize(),
	})
	return err
}

// serialize is a function that serializes the payload of the Queue.
func (q *Queue[T]) serialize() string {
	b, _ := json.Marshal(q.payload)

	return string(b)
}