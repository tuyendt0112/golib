package observer

// ConsumerFactory is an interface for creating observers dynamically.
// This pattern is useful when you want to register multiple observers
// for multiple topics in a structured way (e.g., from configuration).
//
// WHY use factory pattern?
//   - Allows lazy initialization: observers created only when needed
//   - Supports dependency injection: factories can receive dependencies
//   - Enables batch registration: one factory can create observers for multiple topics
type ConsumerFactory interface {
	// CreateConsumer creates a new observer instance.
	// Each call should return a new instance (not a singleton).
	CreateConsumer() Observer
	
	// Topics returns all topics this factory should register observers for.
	// A factory can handle multiple topics (e.g., "user.*" events).
	Topics() []TopicName
}

// ConsumerRegistry manages a collection of consumer factories.
// It provides a centralized way to register and initialize all observers
// at application startup.
//
// WHY use registry pattern?
//   - Centralized configuration: all observers registered in one place
//   - Batch initialization: set up all observers at once
//   - Separation of concerns: registration logic separate from business logic
type ConsumerRegistry struct {
	// factories maps each topic to a list of factories that can create observers for that topic.
	// Multiple factories can handle the same topic.
	factories map[TopicName][]ConsumerFactory
}

// NewConsumerRegistry creates a new ConsumerRegistry instance.
func NewConsumerRegistry() *ConsumerRegistry {
	return &ConsumerRegistry{
		factories: make(map[TopicName][]ConsumerFactory),
	}
}

// Register adds a factory to the registry.
// The factory will be used to create observers for all topics it returns.
//
// Example:
//   registry := NewConsumerRegistry()
//   registry.Register(&EmailFactory{})
//   registry.Register(&SMSFactory{})
func (r *ConsumerRegistry) Register(factory ConsumerFactory) {
	// Register the factory for each topic it handles
	for _, topic := range factory.Topics() {
		r.factories[topic] = append(r.factories[topic], factory)
	}
}

// Initialize creates and subscribes all observers from registered factories.
// This should be called once at application startup, after all factories are registered.
//
// HOW it works:
//   1. Iterates through all registered factories
//   2. Creates an observer instance from each factory
//   3. Subscribes the observer to its corresponding topic(s)
//
// WHY separate Register and Initialize?
//   - Allows registration phase (collecting factories)
//   - Then initialization phase (creating and subscribing observers)
//   - Useful for dependency injection and testing
func (r *ConsumerRegistry) Initialize() {
	// For each topic, create observers from all factories
	for topic, factories := range r.factories {
		for _, factory := range factories {
			// Create a new observer instance and subscribe it
			Subscribe(
				topic,
				factory.CreateConsumer(),
			)
		}
	}
}
