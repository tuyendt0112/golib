package log

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

// SetLogHandler configures slog's default logger for the entire application.
//
// - Sử dụng stderr làm output (phù hợp cho log server / container).
// - Tự động chọn log level theo ENV (dev: DEBUG, prod: INFO).
// - Dùng tint để có log nhiều màu, dễ đọc hơn trong development.
// - Ẩn timestamp ở root-level khi không cần (thường đã có timestamp từ log aggregator).
//
// Gọi hàm này một lần ở early stage (ví dụ main.init() hoặc ngay đầu hàm main)
// trước khi dùng slog.Info / slog.Error ở bất kỳ đâu.
func SetLogHandler() {
	// Ghi log ra stderr (thông lệ của các app chạy trong container / systemd).
	w := os.Stderr
	// Xác định môi trường production dựa trên biến ENV.
	// Có thể set ENV=prod / production / master khi deploy.
	isProduction := os.Getenv("ENV") == "prod" || os.Getenv("ENV") == "production" || os.Getenv("ENV") == "master"
	// Mặc định:
	// - Production: chỉ log từ INFO trở lên (INFO, WARN, ERROR).
	// - Non-production (dev, staging, v.v.): log cả DEBUG trở lên.
	slogLevel := slog.LevelInfo
	if !isProduction {
		slogLevel = slog.LevelDebug
	}
	// Thiết lập logger mặc định cho slog.
	// tint.Handler:
	// - format gọn, đẹp, có màu trong terminal.
	// - AddSource: true => tự động log file:line gọi slog (hữu ích khi debug).
	// - ReplaceAttr: tuỳ chỉnh các field log (ở đây ẩn timestamp root).
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:     slogLevel,
			AddSource: true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Mặc định slog sẽ thêm Time ở root (không có group).
				// Nếu hệ thống log/observability (ELK, Loki, v.v.) đã có timestamp sẵn,
				// ta có thể bỏ bớt để log không bị trùng thông tin.
				if a.Key == slog.TimeKey && len(groups) == 0 {
					return slog.Attr{}
				}
				return a
			},
		}),
	))
}
