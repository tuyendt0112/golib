package telegram

import (
	"fmt"
	"golib/pkg/https"
	"net/url"
	"time"
	// NOTE: https package is commented out. You need to implement or uncomment it.
)

// Options contains configuration for the Telegram bot.
type Options struct {
	key       string      // Telegram bot token (from @BotFather)
	channelID string      // Telegram channel ID where messages are sent
	appName   string      // Application name (for identification in messages)
	metadata  interface{} // Additional metadata to include in messages
}

// Telegram provides methods to send notifications via Telegram Bot API.
// Useful for monitoring, alerts, and application lifecycle events.
//
// WHY use Telegram for notifications?
//   - Real-time: Instant delivery to mobile/desktop
//   - Free: No cost for basic bot usage
//   - Reliable: Telegram's infrastructure is robust
//   - Easy setup: Simple bot creation via @BotFather
type Telegram struct {
	option *Options // Bot configuration options
	domain string   // Application domain (for context in messages)
}

// NewTelegram creates a new Telegram notification client.
//
// Usage:
//
//	tg := telegram.NewTelegram("example.com",
//	    telegram.WithToken("bot-token"),
//	    telegram.WithChannelID("@mychannel"),
//	    telegram.WithAppName("My App"),
//	)
//
// The domain parameter is used to identify which application/service sent the message.
func NewTelegram(domain string, ops ...func(option *Options)) *Telegram {
	// Initialize options with defaults
	options := &Options{}

	// Apply all provided option functions
	for _, op := range ops {
		op(options)
	}

	return &Telegram{
		option: options,
		domain: domain,
	}
}

// SendInstall sends a notification when the application is installed.
// Useful for tracking when new installations occur.
//
// Parameters:
//   - newUser: true for first-time install, false for re-install
//
// Message format: "Install[domain] - TIMESTAMP=... - metadata=..."
func (t *Telegram) SendInstall(newUser bool) {
	message := "Install"
	if !newUser {
		message = "Re-install"
	}

	message += fmt.Sprintf("[%s] - TIMESTAMP=%s", t.domain, time.Now().Format(time.DateTime))
	if t.option.metadata != nil {
		message += fmt.Sprintf(" - metadata=%v", t.option.metadata)
	}

	query := url.Values{}
	query.Add("chat_id", t.option.channelID)
	query.Add("text", message)

	// NOTE: https.Do is commented out - needs to be implemented
	_ = https.Do(
		t.getUrl(query),
	)
}

// SendUnInstall sends a notification when the application is uninstalled.
// Useful for tracking when users remove the application.
//
// Message format: "UnInstall[domain] - TIMESTAMP=... - metadata=..."
func (t *Telegram) SendUnInstall() {
	message := "UnInstall"
	message += fmt.Sprintf("[%s] - TIMESTAMP=%s", t.domain, time.Now().Format(time.DateTime))
	if t.option.metadata != nil {
		message += fmt.Sprintf(" - metadata=%v", t.option.metadata)
	}

	query := url.Values{}
	query.Add("chat_id", t.option.channelID)
	query.Add("text", message)

	// NOTE: https.Do is commented out - needs to be implemented
	_ = https.Do(
		t.getUrl(query),
	)
}

// getUrl constructs the Telegram Bot API URL for sending messages.
// Uses the sendMessage endpoint with query parameters.
//
// API endpoint: https://api.telegram.org/bot{token}/sendMessage
func (t *Telegram) getUrl(query url.Values) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?%s", t.option.key, query.Encode())
}

// Health sends a health check notification.
// Useful for monitoring application health and alerting on failures.
//
// Parameters:
//   - host: The host/service being checked (e.g., "api.example.com")
//   - code: HTTP status code or health status code
//   - err: Error if health check failed, nil if healthy
//
// Message format: "[appName] - [host] - status : {code} [error] {error}"
//
// Example:
//
//	tg.Health("api.example.com", 200, nil)  // Healthy
//	tg.Health("db.example.com", 500, err)   // Unhealthy
func (t *Telegram) Health(host string, code int, err error) {
	message := fmt.Sprintf("[%s] - [%s] - status : %d [error] %v", t.option.appName, host, code, err)
	query := url.Values{}
	query.Add("chat_id", t.option.channelID)
	query.Add("text", message)

	// NOTE: https.Do is commented out - needs to be implemented
	_ = https.Do(
		t.getUrl(query),
	)
}
