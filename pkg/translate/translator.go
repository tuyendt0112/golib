package translation

import (
	"context"
	"time"
)

// TranslatorOptions is the options for the translator
type TranslatorOptions struct {
	// Text is the text to translate
	Text string
	// SourceLang is the source language
	SourceLang string
	// TargetLang is the target language
	TargetLang string
	// MaxRetries is the maximum number of retries if the translation fails
	// Default is 0, which means no retries
	MaxRetries int
	// RetryDelay is the delay between retries
	// Default is 1 second
	RetryDelay time.Duration
}

// Translator defines methods for text translation
type Translator interface {
	// TranslateText translates text from source language to target language
	TranslateText(ctx context.Context, options *TranslatorOptions) (string, error)
}