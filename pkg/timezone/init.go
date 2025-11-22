package timezone

import "os"

// init initializes the timezone.
func init() {
	timezone := os.Getenv("APP_TIMEZONE")

	if timezone == "" {
		timezone = "UTC"
	}

	_ = os.Setenv("TZ", timezone)
}