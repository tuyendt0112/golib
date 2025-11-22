package https

import (
	"fmt"
	"os"
	"time"
)

// WithShopifyAccessToken sets the request header X-Shopify-Access-Token
func WithShopifyAccessToken(accessToken string) func(cfg *Options) {
	return func(cfg *Options) {
		if cfg.headers == nil {
			cfg.headers = M{}
		}
		cfg.headers["X-Shopify-Access-Token"] = accessToken
	}
}

// shopifyAPIVersion is the Shopify API version,
// do not use this variable directly, use getShopifyAPIVersion instead.
var shopifyAPIVersion = ""

// shopifyAPIVersion is the Shopify API version.
// It can be set by the environment variable SHOPIFY_API_VERSION.
// Check https://shopify.dev/docs/api/release-notes for the latest version.
// If SHOPIFY_API_VERSION is not set, it will use the version 3 months ago.
var getShopifyAPIVersion = func() string {
	if shopifyAPIVersion == "" {
		if shopifyAPIVersion = os.Getenv("SHOPIFY_API_VERSION"); shopifyAPIVersion == "" {
			d := time.Now().AddDate(0, -2, 0)
			month := ((d.Month()-1)/3)*3 + 1 // 1, 4, 7, 10
			shopifyAPIVersion = fmt.Sprintf("%d-%02d", d.Year(), month)

			fmt.Println("WARNING: SHOPIFY_API_VERSION is not set, using", shopifyAPIVersion)
		}
	}

	return shopifyAPIVersion
}

// MakeShopifyGraphqlURL returns the Shopify GraphQL URL
func MakeShopifyGraphqlURL(myShopifyDomain string) string {
	return fmt.Sprintf("https://%s/admin/api/%s/graphql.json", myShopifyDomain, getShopifyAPIVersion())
}

// MakeShopifyRestURL returns the Shopify REST URL,
// resources is a list of resources (e.g. orders, products, customers).
// Example:
//
//	MakeShopifyRestURL("abc.myshopify.com", "orders", 123, "transactions")
//	MakeShopifyRestURL("abc.myshopify.com", "products")
//
// For query parameters, use WithQuery in the options,
// Example: to get products with ids 632910392 and 921728736 (/admin/api/2024-01/products.json?ids=632910392,921728736)
//
//	https.Do(https.MakeShopifyRestURL("abc.myshopify.com", "products"),
//		https.WithQuery(https.M{"ids": "632910392,921728736"})
func MakeShopifyRestURL(myShopifyDomain string, resources ...any) string {
	resourcesStr := ""
	for _, r := range resources {
		resourcesStr += fmt.Sprintf("/%v", r)
	}

	return fmt.Sprintf("https://%s/admin/api/%s%s.json", myShopifyDomain, getShopifyAPIVersion(), resourcesStr)
}