package telegram

import (
	"fmt"
	"net/url"
	"time"

	// "github.com/tranvannghia021/goshared/pkg/https"
)

// Options is a struct to store the options of the telegram
type Options struct {
	key       string      // key is the telegram bot token
	channelID string      // channelID is the telegram channel ID
	appName   string      // appName is the name of the app
	metadata  interface{} // metadata is a map[string]interface{}
}

// Telegram is a struct to store the telegram
type Telegram struct {
	option *Options // option is the options of the telegram
	domain string   // domain is the domain of the app
}

// NewTelegram creates a new telegram
func NewTelegram(domain string, ops ...func(option *Options)) *Telegram {
	// options is a pointer to Options struct
	options := &Options{}
	for _, op := range ops {
		op(options)
	}

	return &Telegram{
		option: options,
		domain: domain,
	}
}

// SendInstall sends the install message
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

	_ = https.Do(
		t.getUrl(query),
	)
}

// SendUnInstall sends the uninstall message
func (t *Telegram) SendUnInstall() {
	message := "UnInstall"

	message += fmt.Sprintf("[%s] - TIMESTAMP=%s", t.domain, time.Now().Format(time.DateTime))
	if t.option.metadata != nil {
		message += fmt.Sprintf(" - metadata=%v", t.option.metadata)
	}

	query := url.Values{}
	query.Add("chat_id", t.option.channelID)
	query.Add("text", message)

	_ = https.Do(
		t.getUrl(query),
	)
}

// getUrl returns the url
func (t *Telegram) getUrl(query url.Values) string {

	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?%s", t.option.key, query.Encode())
}

// Health sends the health message
func (t *Telegram) Health(host string, code int, err error) {
	message := fmt.Sprintf("[%s] - [%s] - status : %d [error] %v", t.option.appName, host, code, err)
	query := url.Values{}
	query.Add("chat_id", t.option.channelID)
	query.Add("text", message)

	_ = https.Do(
		t.getUrl(query),
	)
}