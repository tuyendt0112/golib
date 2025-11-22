package https

import "encoding/base64"

// M is a type alias for a map with string keys and values.
type M map[string]string

// Method represents an HTTP method.
type Method string

// Constants for the different HTTP methods.
const (
	GET    Method = "GET"    // GET is used to request data from a specified resource.
	POST   Method = "POST"   // POST is used to send data to a server to create/update a resource.
	PUT    Method = "PUT"    // PUT is used to update a current resource.
	DELETE Method = "DELETE" // DELETE is used to delete the specified resource.
	PATCH  Method = "PATCH"  // PATCH is used to apply partial modifications to a resource.
)

// Options represents the options for an HTTP request.
type Options struct {
	method        Method            // The HTTP method to use for the request.
	query         M                 // The query parameters to include in the request.
	headers       M                 // The headers to include in the request.
	postForm      M                 // The form data to include in the request body.
	multipartForm M                 // The multipart form data to include in the request body.
	byteReq       []byte            // The byte data to include in the request body.
	jsonReq       any               // The JSON data to include in the request body.
	jsonResp      any               // Reference to a variable where the JSON response body will be stored.
	textResp      *string           // Reference to a variable where the text response body will be stored.
	headerResp    map[string]string // Reference to a variable where the response headers will be stored.
	timeout       int               // The request timeout in seconds.
	proxyProvider GoProxyProvider   // The Go proxy provider to use for the request.
}

// WithMethod sets the request method (GET, POST, PUT, DELETE, PATCH)
func WithMethod(method Method) func(cfg *Options) {
	return func(cfg *Options) {
		cfg.method = method
	}
}

// WithQueries sets the query parameters in the URL (e.g. ?key=value)
func WithQueries(query M) func(cfg *Options) {
	return func(cfg *Options) {
		if cfg.query == nil {
			cfg.query = query
		} else {
			for k, v := range query {
				cfg.query[k] = v
			}
		}
	}
}

// WithQuery sets a single query parameter in the URL
func WithQuery(key, value string) func(cfg *Options) {
	return func(cfg *Options) {
		if cfg.query == nil {
			cfg.query = M{}
		}
		cfg.query[key] = value
	}
}

// WithHeaders sets the request headers (e.g. Content-Type, Accept, Authorization)
func WithHeaders(headers M) func(cfg *Options) {
	return func(cfg *Options) {
		if cfg.headers == nil {
			cfg.headers = headers
		} else {
			for k, v := range headers {
				cfg.headers[k] = v
			}
		}
	}
}

// WithHeader sets a single request header
func WithHeader(key, value string) func(cfg *Options) {
	return func(cfg *Options) {
		if cfg.headers == nil {
			cfg.headers = M{}
		}
		cfg.headers[key] = value
	}
}

// WithGraphQLReq sets the request body as a GraphQL query,
// sets the Content-Type header to application/json and sets the method to POST.
// The query is a string, and the variables is a json object.
// If there are no variables, you can omit the second argument.
func WithGraphQLReq(query string, variables ...any) func(cfg *Options) {
	return func(cfg *Options) {
		cfg.method = POST
		jsonReq := map[string]any{
			"query": query,
		}
		if len(variables) > 0 {
			jsonReq["variables"] = variables[0]
		}
		cfg.jsonReq = jsonReq
	}
}

// WithJSONReq sets the request body as JSON,
// and sets the Content-Type header to application/json.
// Don't forget set the method by using WithMethod
func WithJSONReq(jsonReq any) func(cfg *Options) {
	return func(cfg *Options) {
		cfg.jsonReq = jsonReq
	}
}

// WithFormReq sets the request body as form data,
// and sets the Content-Type header to application/x-www-form-urlencoded.
// Don't forget set the method by using WithMethod
func WithFormReq(postFormReq M) func(cfg *Options) {
	return func(cfg *Options) {
		cfg.postForm = postFormReq
	}
}

// WithStrReq sets the request body as a string
func WithStrReq(str string) func(cfg *Options) {
	return func(cfg *Options) {
		cfg.byteReq = []byte(str)
	}
}

// WithByteReq sets the request body as a byte slice
func WithByteReq(b []byte) func(cfg *Options) {
	return func(cfg *Options) {
		cfg.byteReq = b
	}
}

// WithMultipartFormReq sets the request body as multipart form data,
// and sets the Content-Type header to multipart/form-data.
// Don't forget set the method by using WithMethod
func WithMultipartFormReq(form M) func(cfg *Options) {
	return func(cfg *Options) {
		cfg.multipartForm = form
	}
}

// WithBasicAuth sets the request header Authorization to Basic,
// with the given username and password
func WithBasicAuth(username, password string) func(cfg *Options) {
	return func(cfg *Options) {
		if cfg.headers == nil {
			cfg.headers = M{}
		}
		cfg.headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	}
}

// WithJSONRespTo sets the response body as JSON,
// and sets the Accept header to application/json.
// The jsonResp must be a pointer to the response struct.
// Example:
//
//	resp := &struct{ Name string }{}
//	https.Do("http://example.com", https.WithJSONRespTo(resp))
func WithJSONRespTo(jsonResp any) func(cfg *Options) {
	return func(cfg *Options) {
		cfg.jsonResp = jsonResp
	}
}

// WithTextRespTo sets the response body as text,
// The textResp must be a pointer to the response string.
// Example:
//
//	var resp string
//	https.Do("http://example.com", https.WithTextRespTo(&resp))
func WithTextRespTo(textResp *string) func(cfg *Options) {
	return func(cfg *Options) {
		cfg.textResp = textResp
	}
}

// WithHeaderRespTo sets the response header to a map
func WithHeaderRespTo(headers M) func(cfg *Options) {
	return func(cfg *Options) {
		cfg.headerResp = headers
	}
}

// WithTimeout sets the request timeout in seconds
func WithTimeout(seconds int) func(cfg *Options) {
	return func(cfg *Options) {
		cfg.timeout = seconds
	}
}