package observer

import "log/slog"

// TopicName represents the name of a topic/event in the observer pattern.
// Topics are used to categorize different types of events (e.g., "user.created", "order.paid").
type TopicName string

// Observer is the interface that defines how observers handle events.
// Any struct implementing this interface can subscribe to topics and receive notifications.
//
// WHY this interface?
//   - Decouples the publisher from specific observer implementations
//   - Allows multiple different observers to handle the same event
//   - Makes the system extensible without modifying existing code
type Observer interface {
	// Handle processes an event for a given topic with the provided data.
	// This method is called asynchronously in a goroutine, so it should be thread-safe.
	Handle(topic TopicName, data interface{})
	
	// Name returns a unique identifier for this observer (useful for logging/debugging).
	Name() string
}

// Subject maintains a registry of observers grouped by topic.
// It follows the Observer Pattern: subjects notify observers when events occur.
//
// WHY generic type [T any]?
//   - Allows type-safe subjects in the future if needed
//   - Currently uses [any] for maximum flexibility
type Subject[T any] struct {
	// observers maps each topic to a list of observers that are interested in that topic.
	// Multiple observers can subscribe to the same topic.
	observers map[TopicName][]Observer
}

// subject is the global singleton instance of Subject.
// WHY singleton?
//   - Provides a single point of access for the entire application
//   - Simplifies usage: no need to pass subject instance around
//   - All parts of the application can use the same observer registry
var subject = &Subject[any]{
	observers: make(map[TopicName][]Observer),
}

// Subscribe registers an observer to receive notifications for a specific topic.
//
// Example:
//   observer.Subscribe("user.created", &EmailObserver{})
//   observer.Subscribe("user.created", &SMSObserver{})
//
// WHY append instead of replace?
//   - Allows multiple observers for the same topic
//   - Each observer can handle the event independently
//   - Supports one-to-many notification pattern
func Subscribe(topic TopicName, observer Observer) {
	subject.observers[topic] = append(subject.observers[topic], observer)
}

// Notify sends an event to all observers subscribed to the given topic.
//
// HOW it works:
//   1. Looks up all observers for the topic
//   2. Spawns a goroutine for each observer (non-blocking)
//   3. Each observer handles the event independently
//   4. Panics in observers are caught and logged (doesn't crash the app)
//
// WHY async (goroutines)?
//   - Non-blocking: publisher doesn't wait for observers to finish
//   - Parallel processing: multiple observers run concurrently
//   - Prevents slow observers from blocking the main flow
//
// WHY panic recovery?
//   - One observer's panic shouldn't crash the entire application
//   - Allows other observers to continue processing
//   - Logs the error for debugging while keeping the app running
//
// Example:
//   observer.Notify("user.created", userData)
func Notify(topic TopicName, data interface{}) {
	// Check if there are any observers for this topic
	if observers, found := subject.observers[topic]; found {
		// Notify each observer in a separate goroutine
		for _, observer := range observers {
			// WHY capture observer in closure parameter?
			//   - Prevents race condition: each goroutine gets its own copy
			//   - Without this, all goroutines might use the last observer in the loop
			go func(observer Observer) {
				// Recover from panics to prevent app crash
				defer func() {
					if r := recover(); r != nil {
						// Log the panic with context for debugging
						// TODO: Add retry mechanism for failed observers
						slog.Error("[consumer] panic", "topic", topic, "name", observer.Name(), "error", r)
					}
				}()
				// Call the observer's handler
				observer.Handle(topic, data)
			}(observer)
		}
	}
	// If no observers found, silently return (no error, just no-op)
}