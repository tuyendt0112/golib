package env

import (
	"os"
	"testing"
)

func TestMustGet(t *testing.T) {
	// Test case 1: Environment variable exists
	testKey := "TEST_ENV_VAR"
	testValue := "test_value_123"
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	result := MustGet(testKey)
	if result != testValue {
		t.Errorf("MustGet(%s) = %s, want %s", testKey, result, testValue)
	}
}

func TestMustGet_PanicWhenNotExists(t *testing.T) {
	// Test case 2: Environment variable does not exist (should panic)
	testKey := "NON_EXISTENT_VAR"
	os.Unsetenv(testKey) // Ensure it doesn't exist

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGet should panic when environment variable does not exist")
		}
	}()

	MustGet(testKey)
	t.Error("Should not reach here - MustGet should have panicked")
}

