package queue

import (
	"testing"
)

type TestPayload struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestNewQueue(t *testing.T) {
	queueName := "test-queue"
	queue := NewQueue[TestPayload](queueName)

	if queue == nil {
		t.Error("NewQueue should return a non-nil queue")
	}

	if queue.queueName != queueName {
		t.Errorf("Expected queue name '%s', got '%s'", queueName, queue.queueName)
	}

	if queue.task == nil {
		t.Error("Queue task should not be nil")
	}
}

func TestQueue_WithData(t *testing.T) {
	queue := NewQueue[TestPayload]("test-queue")
	
	testData := &TestPayload{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}

	queue.WithData(testData)

	if queue.payload == nil {
		t.Error("Queue payload should not be nil after WithData")
	}

	if queue.payload.ID != testData.ID {
		t.Errorf("Expected payload ID %d, got %d", testData.ID, queue.payload.ID)
	}

	if queue.payload.Name != testData.Name {
		t.Errorf("Expected payload Name '%s', got '%s'", testData.Name, queue.payload.Name)
	}
}

func TestQueue_Serialize(t *testing.T) {
	queue := NewQueue[TestPayload]("test-queue")
	
	testData := &TestPayload{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}

	queue.WithData(testData)
	serialized := queue.serialize()

	if serialized == "" {
		t.Error("Serialized payload should not be empty")
	}

	// Verify it's valid JSON by checking it contains expected fields
	if len(serialized) < 10 {
		t.Error("Serialized payload seems too short to be valid JSON")
	}
}

func TestNewWorker(t *testing.T) {
	queueName := "test-worker-queue"
	worker := NewWorker[TestPayload](queueName)

	if worker == nil {
		t.Error("NewWorker should return a non-nil worker")
	}

	if worker.queueName != queueName {
		t.Errorf("Expected queue name '%s', got '%s'", queueName, worker.queueName)
	}

	if worker.pool == nil {
		t.Error("Worker pool should not be nil")
	}

	if worker.options == nil {
		t.Error("Worker options should not be nil")
	}
}

func TestNewWorker_WithOptions(t *testing.T) {
	queueName := "test-worker-queue"
	
	// Test with custom options
	worker := NewWorker[TestPayload](queueName, func(options *Options) {
		options.MaxFails = 5
		options.MaxConcurrency = 20
	})

	if worker.options.MaxFails != 5 {
		t.Errorf("Expected MaxFails 5, got %d", worker.options.MaxFails)
	}

	if worker.options.MaxConcurrency != 20 {
		t.Errorf("Expected MaxConcurrency 20, got %d", worker.options.MaxConcurrency)
	}
}

func TestWorker_Deserialize(t *testing.T) {
	worker := NewWorker[TestPayload]("test-queue")
	
	jsonData := `{"id":1,"name":"Test User","email":"test@example.com"}`
	
	err := worker.deserialize(jsonData)
	if err != nil {
		t.Errorf("Deserialize should not return error, got %v", err)
	}

	if worker.payload == nil {
		t.Error("Worker payload should not be nil after deserialize")
	}

	if worker.payload.ID != 1 {
		t.Errorf("Expected payload ID 1, got %d", worker.payload.ID)
	}

	if worker.payload.Name != "Test User" {
		t.Errorf("Expected payload Name 'Test User', got '%s'", worker.payload.Name)
	}
}

func TestWorker_Deserialize_InvalidJSON(t *testing.T) {
	worker := NewWorker[TestPayload]("test-queue")
	
	invalidJSON := `{"id":1,"name":invalid}`
	
	err := worker.deserialize(invalidJSON)
	if err == nil {
		t.Error("Deserialize should return error for invalid JSON")
	}
}

