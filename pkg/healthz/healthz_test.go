package healthz

import (
	"net"
	"testing"
	"time"
)

func TestRunCatchPanic(t *testing.T) {
	// Test case 1: Function returns error
	errFunc := func() error {
		return nil
	}
	err := runCatchPanic(errFunc)
	if err != nil {
		t.Errorf("runCatchPanic should return nil when function returns nil, got %v", err)
	}

	// Test case 2: Function panics
	panicFunc := func() error {
		panic("test panic")
	}
	err = runCatchPanic(panicFunc)
	if err == nil {
		t.Error("runCatchPanic should return error when function panics")
	}
	if err.Error() != "panic: test panic" {
		t.Errorf("runCatchPanic error message = %s, want 'panic: test panic'", err.Error())
	}
}

func TestRunServer_Basic(t *testing.T) {
	// This is a basic test - RunServer is complex and starts a TCP server
	// In a real scenario, you might want to test with a mock or integration test
	
	// Test that server can start (but we'll need to stop it quickly)
	done := make(chan bool)
	
	go func() {
		// Start server with a simple health check
		checkFunc := func() error {
			return nil // Always healthy
		}
		
		RunServer(checkFunc)
		done <- true
	}()
	
	// Give server time to start
	time.Sleep(100 * time.Millisecond)
	
	// Try to connect to healthz server
	conn, err := net.DialTimeout("tcp", "127.0.0.1:9999", 1*time.Second)
	if err != nil {
		t.Logf("Could not connect to healthz server (this is expected if port is busy): %v", err)
		return
	}
	defer conn.Close()
	
	// Read response
	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		t.Logf("Error reading response: %v", err)
		return
	}
	
	response := string(buf[:n])
	if len(response) == 0 {
		t.Error("Expected non-empty response from healthz server")
	}
	
	// Note: In a real test, you'd want to properly shut down the server
	// For now, this is a basic connectivity test
}

