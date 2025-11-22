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

	// "github.com/tranvannghia021/goshared/pkg/translation"
)

// Translator implements Google translation service
type Translator struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

// NewTranslator creates a new Google translator instance
func NewTranslator(client *http.Client, apiKey string) *Translator {
	if client == nil {
		client = http.DefaultClient
	}

	return &Translator{
		client:  client,
		baseURL: "https://translate-pa.googleapis.com/v1/translate",
		apiKey:  apiKey,
	}
}

// TranslationResponse represents the response from Google Translate API
type TranslationResponse struct {
	Translation string `json:"translation"`
	Sentences   []struct {
		Trans string `json:"trans"`
		Orig  string `json:"orig"`
	} `json:"sentences"`
	SourceLanguage string `json:"sourceLanguage"`
}

// TranslateText translates the given text using the provided options
func (t *Translator) TranslateText(ctx context.Context, options *translation.TranslatorOptions) (string, error) {
	if options == nil {
		return "", fmt.Errorf("options cannot be nil")
	}

	if options.Text == "" {
		return "", nil
	}

	// Set default values for retry logic
	maxRetries := options.MaxRetries
	if maxRetries < 0 {
		maxRetries = 0
	}

	retryDelay := options.RetryDelay
	if retryDelay <= 0 {
		retryDelay = time.Second
	}

	var lastErr error

	// Try translation with retry logic
	for attempt := 0; attempt <= maxRetries; attempt++ {
		result, err := t.translateTextOnce(ctx, options.Text, options.SourceLang, options.TargetLang)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// If this is not the last attempt, wait before retrying
		if attempt < maxRetries {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(retryDelay):
				// Continue to next attempt
			}
		}
	}

	return "", fmt.Errorf("translation failed after %d attempts: %w", maxRetries+1, lastErr)
}

// translateTextOnce performs a single translation attempt
func (t *Translator) translateTextOnce(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	// Build the request URL with query parameters
	params := url.Values{}
	params.Add("params.client", "gtx")
	params.Add("query.source_language", sourceLang)
	params.Add("query.target_language", targetLang)
	params.Add("query.text", text)
	params.Add("key", t.apiKey)
	params.Add("data_types", "TRANSLATION")
	params.Add("data_types", "SENTENCE_SPLITS")
	params.Add("data_types", "BILINGUAL_DICTIONARY_FULL")

	// Create request
	reqURL := fmt.Sprintf("%s?%s", t.baseURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("translation request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	// Parse response
	var result TranslationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Return the translated text
	return result.Translation, nil
}
