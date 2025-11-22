package healthz

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// runCatchPanic wraps and executes a health-check function safely.
// WHY:
//   - Một health check có thể panic (ví dụ: nil pointer khi check DB)
//   - Nếu panic thoát ra ngoài → crash luôn healthz server → orchestrator (K8s, Docker) hiểu sai.
//   - Hàm này đảm bảo panic được catch lại và trả về error thay vì làm server chết.
func runCatchPanic(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
			slog.Error("Healthz panic", "err", err)
		}
	}()

	err = f()
	return
}

// RunServer starts a simple TCP-based health server on port 9999.
// WHY TCP instead of HTTP server?
//   - Healthz thường cần server cực nhẹ, không cần router, middleware.
//   - TCP listener cho performance cao nhất, ít overhead.
//   - Dễ dùng với load balancer hoặc readiness/liveness probe của K8s.
func RunServer(checkFuncs ...func() error) {
	// Healthz server listens on 0.0.0.0:9999 → lắng nghe tất cả interfaces
	addr := net.TCPAddr{IP: net.IPv4zero, Port: 9999}
	// Try to open listener. If port is already used → log error.
	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		slog.Error("Healthz failed to listen", slog.Any("err", err))
	}

	slog.Info("Healthz server started", slog.Any("port", addr.Port))
	// Flag để biết khi server bị kill (SIGINT, SIGTERM)
	// WHY:
	//   - Trong Docker hoặc K8s, container shutdown sẽ gửi SIGTERM
	//   - Cần close listener đúng cách để tránh goroutine leak.
	var killed = false
	// Goroutine bắt tín hiệu dừng server.
	go func() {
		// Lắng nghe tín hiệu hệ thống: Ctrl+C (SIGINT), kill (SIGTERM)
		s := make(chan os.Signal, 1)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
		<-s // Chờ tín hiệu
		killed = true
		listener.Close() // Giải phóng TCP listener
	}()
	// Reusable byte buffers để tránh phải format header nhiều lần.
	// WHY: performance tối ưu.
	var resOKBuf = []byte("HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n")
	var resErrBufPrefix = []byte("HTTP/1.1 503 Service Unavailable\r\nContent-Length: ")
	// Main loop: accept incoming TCP connections
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			if killed {
				// Nếu server đã bị kill → thoát vòng lặp
				return
			}
			slog.Error("Healthz failed to accept", "err", err)
			break
		}
		// Execute all provided health-check functions
		for _, f := range checkFuncs {
			// Catch panic + trả error thay vì crash server
			if err := runCatchPanic(f); err != nil {
				errStr := err.Error()
				fmt.Println("Healthz error: ", errStr)
				// Tính content length để gửi HTTP header đúng chuẩn
				contentLen := len(errStr)
				// WHY deadline:
				//   - Tránh treo khi client không đọc data
				//   - Healthz phải timeout nhanh (<100ms)
				conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
				// Trả HTTP 503 + error message
				conn.Write(resErrBufPrefix)
				conn.Write([]byte(strconv.Itoa(contentLen) + "\r\n\r\n" + errStr))
				conn.Close()
				// continue: không gửi 200 OK nữa
				continue
			}
		}
		// Nếu tất cả checkFunc đều OK → trả HTTP 200
		conn.Write(resOKBuf)
		conn.Close()
	}
}
