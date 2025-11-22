package translation

import (
	"context"
	"time"
)

// TranslatorOptions contains configuration for a translation request.
type TranslatorOptions struct {
	// Text is the text to be translated.
	Text string
	
	// SourceLang is the source language code (e.g., "en", "vi", "ja").
	// Use ISO 639-1 language codes or Google Translate language codes.
	SourceLang string
	
	// TargetLang is the target language code (e.g., "en", "vi", "ja").
	// Use ISO 639-1 language codes or Google Translate language codes.
	TargetLang string
	
	// MaxRetries is the maximum number of retry attempts if translation fails.
	// Default is 0 (no retries). Set to 3-5 for production use.
	// WHY retry?
	//   - Network issues can cause temporary failures
	//   - API rate limits may cause temporary rejections
	//   - Improves reliability without manual intervention
	MaxRetries int
	
	// RetryDelay is the delay between retry attempts.
	// Default is 1 second. Increase for rate-limited APIs.
	// WHY delay?
	//   - Prevents hammering the API on failures
	//   - Gives transient issues time to resolve
	//   - Respects API rate limits
	RetryDelay time.Duration
}

// Translator defines the interface for text translation services.
// This interface allows different translation providers (Google, DeepL, etc.)
// to be used interchangeably.
//
// WHY interface?
//   - Allows switching translation providers without changing calling code
//   - Enables testing with mock implementations
//   - Supports multiple providers in the same application
type Translator interface {
	// TranslateText translates text from source language to target language.
	//
	// The context can be used to:
	//   - Cancel long-running translations
	//   - Set timeouts
	//   - Pass request-scoped values
	//
	// Returns the translated text and an error if translation fails.
	// If MaxRetries > 0, will automatically retry on failure.
	TranslateText(ctx context.Context, options *TranslatorOptions) (string, error)
}