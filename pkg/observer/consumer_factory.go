package observer

// ConsumerFactory is the interface that wraps the basic CreateConsumer and Topics method.
type ConsumerFactory interface {
	CreateConsumer() Observer
	Topics() []TopicName
}

// ConsumerRegistry is the struct that wraps the factories map.
type ConsumerRegistry struct {
	factories map[TopicName][]ConsumerFactory
}

// NewConsumerRegistry is a function to create the instance of ConsumerRegistry.
func NewConsumerRegistry() *ConsumerRegistry {
	return &ConsumerRegistry{
		factories: make(map[TopicName][]ConsumerFactory),
	}
}

// Register is a method to register the factory to the topic.
func (r *ConsumerRegistry) Register(factory ConsumerFactory) {
	for _, topic := range factory.Topics() {
		r.factories[topic] = append(r.factories[topic], factory)
	}
}

// Initialize is a method to initialize the consumer registry.
func (r *ConsumerRegistry) Initialize() {
	for topic, factories := range r.factories {
		for _, factory := range factories {
			Subscribe(
				topic,
				factory.CreateConsumer(),
			)
		}
	}
}
