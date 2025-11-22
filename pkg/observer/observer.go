package observer

import "log/slog"

// TopicName is the type of topic name
type TopicName string

// Observer is the interface that wraps the basic Handle method.
type Observer interface {
	Handle(topic TopicName, data interface{})
	Name() string
}

// Subject is the struct that wraps the observers map.
type Subject[T any] struct {
	observers map[TopicName][]Observer
}

// subject is the instance of Subject.
var subject = &Subject[any]{
	observers: make(map[TopicName][]Observer),
}

// Subscribe is a function to subscribe the observer to the topic.
func Subscribe(topic TopicName, observer Observer) {
	subject.observers[topic] = append(subject.observers[topic], observer)
}

// Notify is a function to notify the observer with the topic and data.
func Notify(topic TopicName, data interface{}) {
	if observers, found := subject.observers[topic]; found {
		for _, observer := range observers {

			go func(observer Observer) {
				defer func() {
					if r := recover(); r != nil {
						// log error // retry when panic (todo)
						slog.Error("[consumer] panic", "topic", topic, "name", observer.Name(), "error", r)
					}
				}()
				observer.Handle(topic, data)
			}(observer)
		}
	}
}