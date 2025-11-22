package telegram

// WithAppName is a function that sets the app name
func WithAppName(name string) func(option *Options) {

	return func(option *Options) {
		option.appName = name
	}
}

// WithChannelID is a function that sets the channel ID
func WithChannelID(id string) func(option *Options) {

	return func(option *Options) {
		option.channelID = id
	}
}

// WithToken is a function that sets the token
func WithToken(token string) func(option *Options) {

	return func(option *Options) {
		option.key = token
	}
}

// WithMetadata is a function that sets the metadata
func WithMetadata(meta interface{}) func(option *Options) {

	return func(option *Options) {
		option.metadata = meta
	}
}
