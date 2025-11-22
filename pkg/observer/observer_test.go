package observer

import (
	"sync"
	"testing"
	"time"
)

// MockObserver is a test implementation of Observer
type MockObserver struct {
	name        string
	receivedData []interface{}
	receivedTopic TopicName
	mu          sync.Mutex
}

func (m *MockObserver) Handle(topic TopicName, data interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.receivedTopic = topic
	m.receivedData = append(m.receivedData, data)
}

func (m *MockObserver) Name() string {
	return m.name
}

func (m *MockObserver) GetReceivedData() []interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.receivedData
}

func TestSubscribe(t *testing.T) {
	// Reset subject for clean test
	subject = &Subject[any]{
		observers: make(map[TopicName][]Observer),
	}

	observer := &MockObserver{name: "test-observer"}
	topic := TopicName("test-topic")

	Subscribe(topic, observer)

	if len(subject.observers[topic]) != 1 {
		t.Errorf("Expected 1 observer for topic, got %d", len(subject.observers[topic]))
	}

	if subject.observers[topic][0].Name() != "test-observer" {
		t.Errorf("Expected observer name 'test-observer', got '%s'", subject.observers[topic][0].Name())
	}
}

func TestNotify(t *testing.T) {
	// Reset subject for clean test
	subject = &Subject[any]{
		observers: make(map[TopicName][]Observer),
	}

	observer1 := &MockObserver{name: "observer-1"}
	observer2 := &MockObserver{name: "observer-2"}
	topic := TopicName("test-topic")

	Subscribe(topic, observer1)
	Subscribe(topic, observer2)

	testData := "test data"
	Notify(topic, testData)

	// Give goroutines time to execute
	// In a real scenario, you might use channels or wait groups
	time.Sleep(100 * time.Millisecond)

	data1 := observer1.GetReceivedData()
	if len(data1) == 0 {
		t.Error("Observer 1 should have received data")
	}
	if len(data1) > 0 && data1[0] != testData {
		t.Errorf("Observer 1 received %v, want %v", data1[0], testData)
	}

	data2 := observer2.GetReceivedData()
	if len(data2) == 0 {
		t.Error("Observer 2 should have received data")
	}
	if len(data2) > 0 && data2[0] != testData {
		t.Errorf("Observer 2 received %v, want %v", data2[0], testData)
	}
}

func TestNotify_NoObservers(t *testing.T) {
	// Reset subject for clean test
	subject = &Subject[any]{
		observers: make(map[TopicName][]Observer),
	}

	// Notify on a topic with no observers (should not panic)
	topic := TopicName("empty-topic")
	Notify(topic, "some data")
	
	// If we reach here without panicking, test passes
}

func TestNotify_PanicRecovery(t *testing.T) {
	// Reset subject for clean test
	subject = &Subject[any]{
		observers: make(map[TopicName][]Observer),
	}

	// Create an observer that panics
	panicObserver := &panicObserverImpl{name: "panic-observer"}

	topic := TopicName("panic-topic")
	Subscribe(topic, panicObserver)

	// This should not panic due to recovery mechanism
	Notify(topic, "test data")
	
	// Give goroutine time to execute and recover
	time.Sleep(100 * time.Millisecond)
	
	// If we reach here, panic was recovered successfully
}

// panicObserverImpl is an observer that panics when handling
type panicObserverImpl struct {
	name string
}

func (p *panicObserverImpl) Handle(topic TopicName, data interface{}) {
	panic("test panic")
}

func (p *panicObserverImpl) Name() string {
	return p.name
}

