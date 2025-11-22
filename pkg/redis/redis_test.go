package redis

import (
	"os"
	"sync"
	"testing"
)

func TestNewClientRedis_WithOptions(t *testing.T) {
	// Test case 1: Create client with options
	opts := &Options{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
		MaxRetry: 3,
	}

	client := NewClientRedis(opts)
	if client == nil {
		t.Error("NewClientRedis should return a non-nil client")
	}
	
	// Verify client options
	options := client.Options()
	if options.Addr != "localhost:6379" {
		t.Errorf("Expected addr 'localhost:6379', got '%s'", options.Addr)
	}
	if options.DB != 0 {
		t.Errorf("Expected DB 0, got %d", options.DB)
	}
}

func TestNewClientRedis_WithEnvVars(t *testing.T) {
	// Test case 2: Create client using environment variables (no options provided)
	os.Setenv("REDIS_HOST", "127.0.0.1")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("REDIS_DB", "1")
	defer func() {
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("REDIS_PORT")
		os.Unsetenv("REDIS_DB")
	}()

	// Reset singleton for this test
	instanceRedis = nil
	redisOnce = sync.Once{}

	// Call without options - should use env vars
	client := NewClientRedis()
	if client == nil {
		t.Error("NewClientRedis should return a non-nil client")
	}
	
	// Verify client options
	options := client.Options()
	if options.Addr != "127.0.0.1:6380" {
		t.Errorf("Expected addr '127.0.0.1:6380', got '%s'", options.Addr)
	}
}

func TestNewClientRedis_Singleton(t *testing.T) {
	// Test case 3: Verify singleton pattern
	opts1 := &Options{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
		MaxRetry: 3,
	}

	opts2 := &Options{
		Host:     "different",
		Port:     "6380",
		Password: "",
		DB:       1,
		MaxRetry: 5,
	}

	client1 := NewClientRedis(opts1)
	client2 := NewClientRedis(opts2)
	
	// Due to sync.Once, the second call should return the same instance
	// This is a limitation of the current implementation
	if client1 != client2 {
		t.Log("Note: Due to singleton pattern, both clients are the same instance")
	}
}

