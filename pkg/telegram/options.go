package telegram

// WithAppName returns an option function to set the application name.
// This name is included in health check messages for identification.
//
// Example:
//   tg := telegram.NewTelegram("example.com", telegram.WithAppName("My API"))
func WithAppName(name string) func(option *Options) {
	return func(option *Options) {
		option.appName = name
	}
}

// WithChannelID returns an option function to set the Telegram channel ID.
// Messages will be sent to this channel.
//
// How to get channel ID:
//   - For public channels: Use @channel_username
//   - For private channels: Use numeric ID (e.g., "-1001234567890")
//
// Example:
//   tg := telegram.NewTelegram("example.com", telegram.WithChannelID("@mychannel"))
func WithChannelID(id string) func(option *Options) {
	return func(option *Options) {
		option.channelID = id
	}
}

// WithToken returns an option function to set the Telegram bot token.
// Get your bot token from @BotFather on Telegram.
//
// Example:
//   tg := telegram.NewTelegram("example.com", telegram.WithToken("123456:ABC-DEF..."))
func WithToken(token string) func(option *Options) {
	return func(option *Options) {
		option.key = token
	}
}

// WithMetadata returns an option function to set additional metadata.
// This metadata is included in install/uninstall messages for context.
//
// Example:
//   tg := telegram.NewTelegram("example.com",
//       telegram.WithMetadata(map[string]string{"version": "1.0.0"}),
//   )
func WithMetadata(meta interface{}) func(option *Options) {
	return func(option *Options) {
		option.metadata = meta
	}
}
