package timezone

import "os"

// init automatically configures the application's timezone when the package is imported.
//
// Usage:
//   import _ "golb/pkg/timezone"  // Blank import triggers init()
//
//   func main() {
//       // Timezone is already set to UTC or APP_TIMEZONE value
//       now := time.Now() // Uses configured timezone
//   }
//
// WHY use init()?
//   - Automatic setup: timezone configured before any code runs
//   - Zero configuration: just import the package
//   - Global setting: affects all time operations in the application
//
// WHY set TZ environment variable?
//   - Go's time package respects the TZ environment variable
//   - Setting TZ affects time.Now(), time.Parse(), etc.
//   - Works across all packages without passing timezone explicitly
//
// Environment Variable:
//   - APP_TIMEZONE: Timezone name (e.g., "Asia/Ho_Chi_Minh", "America/New_York")
//   - Default: "UTC" if APP_TIMEZONE is not set
//
// NOTE: This must be imported with blank identifier (_) to trigger init().
func init() {
	// Read timezone from environment variable
	timezone := os.Getenv("APP_TIMEZONE")

	// Default to UTC if not specified
	if timezone == "" {
		timezone = "UTC"
	}

	// Set the TZ environment variable for Go's time package
	// Ignore error (Setenv rarely fails, and there's nothing we can do if it does)
	_ = os.Setenv("TZ", timezone)
}