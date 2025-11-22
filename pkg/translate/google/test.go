package google

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	// "github.com/tranvannghia021/goshared/pkg/translation"
	// "github.com/tranvannghia021/goshared/pkg/translation/google"
)

func ExampleTranslator_TranslateText() {
	// In a real application, get API key from environment or configuration
	apiKey := os.Getenv("GOOGLE_TRANSLATE_API_KEY")
	if apiKey == "" {
		// For example purposes only, not recommended for production
		apiKey = "YOUR_API_KEY"
	}

	// Create a new translator with default HTTP client
	translator := google.NewTranslator(http.DefaultClient, apiKey)

	// Translate "hello" from English to Filipino with retry options
	ctx := context.Background()
	options := &translation.TranslatorOptions{
		Text:       "hello",
		SourceLang: "en",
		TargetLang: "pam",
		MaxRetries: 3,
		RetryDelay: 2 * time.Second,
	}

	translation, err := translator.TranslateText(ctx, options)
	if err != nil {
		fmt.Printf("Translation error: %v\n", err)
		return
	}

	fmt.Printf("Translation: %s\n", translation)
	// Output: Translation: Komusta
}
