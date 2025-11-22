package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golib/pkg/translate"
)

// Translator implements the translation.Translator interface using Google Translate API.
// This is a concrete implementation that can be used for translating text.
//
// WHY Google Translate?
//   - Supports 100+ languages
//   - High accuracy for common languages
//   - Reliable API with good uptime
//   - Free tier available for testing
type Translator struct {
	client  *http.Client // HTTP client for making API requests
	baseURL string       // Google Translate API endpoint
	apiKey  string       // API key for authentication
}

// NewTranslator creates a new Google translator instance.
//
// Parameters:
//   - client: HTTP client to use for requests. If nil, uses http.DefaultClient.
//     Useful for setting timeouts, custom transport, or testing with mock clients.
//   - apiKey: Google Translate API key. Get one from Google Cloud Console.
//
// Example:
//
//	translator := google.NewTranslator(http.DefaultClient, "your-api-key")
func NewTranslator(client *http.Client, apiKey string) *Translator {
	// Use default client if none provided
	if client == nil {
		client = http.DefaultClient
	}

	return &Translator{
		client:  client,
		baseURL: "https://translate-pa.googleapis.com/v1/translate",
		apiKey:  apiKey,
	}
}

// TranslationResponse represents the JSON response structure from Google Translate API.
type TranslationResponse struct {
	Translation string `json:"translation"` // The translated text
	Sentences   []struct {
		Trans string `json:"trans"` // Translated sentence
		Orig  string `json:"orig"`  // Original sentence
	} `json:"sentences"` // Sentence-by-sentence breakdown
	SourceLanguage string `json:"sourceLanguage"` // Detected source language
}

// TranslateText translates text using Google Translate API with automatic retry logic.
//
// This method implements the translation.Translator interface.
//
// HOW retry works:
//  1. Attempts translation up to MaxRetries+1 times
//  2. Waits RetryDelay between attempts
//  3. Respects context cancellation (can be cancelled mid-retry)
//  4. Returns error if all attempts fail
//
// WHY retry?
//   - Network issues can cause temporary failures
//   - API rate limits may cause temporary rejections
//   - Improves reliability without manual intervention
//
// Returns the translated text or an error if translation fails after all retries.
func (t *Translator) TranslateText(ctx context.Context, options *translate.TranslatorOptions) (string, error) {
	// Validate options
	if options == nil {
		return "", fmt.Errorf("options cannot be nil")
	}

	// Empty text doesn't need translation
	if options.Text == "" {
		return "", nil
	}

	// Normalize retry settings
	maxRetries := options.MaxRetries
	if maxRetries < 0 {
		maxRetries = 0 // No negative retries
	}

	retryDelay := options.RetryDelay
	if retryDelay <= 0 {
		retryDelay = time.Second // Default 1 second delay
	}

	var lastErr error

	// Retry loop: attempt translation up to MaxRetries+1 times
	// WHY MaxRetries+1? Because attempt 0 is the first attempt, not a retry
	for attempt := 0; attempt <= maxRetries; attempt++ {
		result, err := t.translateTextOnce(ctx, options.Text, options.SourceLang, options.TargetLang)
		if err == nil {
			return result, nil // Success!
		}

		lastErr = err

		// Wait before retrying (unless this was the last attempt)
		if attempt < maxRetries {
			select {
			case <-ctx.Done():
				// Context was cancelled - stop retrying
				return "", ctx.Err()
			case <-time.After(retryDelay):
				// Wait for retry delay, then continue to next attempt
			}
		}
	}

	// All attempts failed
	return "", fmt.Errorf("translation failed after %d attempts: %w", maxRetries+1, lastErr)
}

// translateTextOnce performs a single translation attempt to Google Translate API.
// This is the core translation logic without retry handling.
//
// HOW it works:
//  1. Builds API request URL with query parameters
//  2. Sends HTTP GET request to Google Translate API
//  3. Parses JSON response
//  4. Returns translated text
//
// WHY separate method?
//   - Keeps retry logic separate from API call logic
//   - Makes testing easier (can test API call without retry)
//   - Cleaner code organization
func (t *Translator) translateTextOnce(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	// Build query parameters for Google Translate API
	params := url.Values{}
	params.Add("params.client", "gtx")                    // Client identifier
	params.Add("query.source_language", sourceLang)       // Source language code
	params.Add("query.target_language", targetLang)       // Target language code
	params.Add("query.text", text)                        // Text to translate
	params.Add("key", t.apiKey)                           // API key for authentication
	params.Add("data_types", "TRANSLATION")               // Request translation
	params.Add("data_types", "SENTENCE_SPLITS")           // Request sentence breakdown
	params.Add("data_types", "BILINGUAL_DICTIONARY_FULL") // Request dictionary data

	// Construct full API URL
	reqURL := fmt.Sprintf("%s?%s", t.baseURL, params.Encode())

	// Create HTTP request with context (for cancellation/timeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Execute HTTP request
	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close() // Always close response body

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		// Read error response body for debugging
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("translation request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	// Parse JSON response
	var result TranslationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Return the translated text
	return result.Translation, nil
}
