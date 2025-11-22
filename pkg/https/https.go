package https

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
)

// ErrorStatusNotOK is an error type for non-200 status codes
type ErrorStatusNotOK struct {
	Code int
	Msg  string
	Body string
}

// Error returns the error message
func (e *ErrorStatusNotOK) Error() string {
	return fmt.Sprintf("status code: %d %s\n%s", e.Code, e.Msg, e.Body)
}

// Do makes a request to the given URL,
// The options are used to configure the request.
// The options are applied in order, so the last option will override the previous ones.
// To get the response body, use WithJSONRespTo or WithTextRespTo.
// This function returns an error if the request fails or the status code is not 200.
// To check the status code, use the ErrorStatusNotOK type, for example:
//
//	err := https.Do("http://example.com")
//	if err != nil {
//		if e, ok := err.(*https.ErrorStatusNotOK); ok {
//			fmt.Println(e.Code, e.Msg, e.Body)
//		}
//	}
func Do(url string, options ...func(cfg *Options)) (err error) {
	cfg := &Options{}
	for _, option := range options {
		option(cfg)
	}

	if cfg.method == "" {
		cfg.method = GET
	}

	if cfg.headers == nil {
		cfg.headers = M{}
	}

	if cfg.jsonResp != nil {
		cfg.headers["Accept"] = "application/json"
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	if cfg.proxyProvider != nil {
		proxyHost, proxySecret := cfg.proxyProvider.GetProxy()

		url = "https://" + proxyHost + "/" + strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "http://")

		if len(proxySecret) > 0 {
			req.Header.Set("X-Proxy-Secret", proxySecret)
		}
	}

	req.SetRequestURI(url)
	req.Header.SetMethod(string(cfg.method))

	for k, v := range cfg.headers {
		req.Header.Set(k, v)
	}

	if cfg.query != nil {
		args := req.URI().QueryArgs()
		for k, v := range cfg.query {
			args.Set(k, v)
		}
	}

	if cfg.jsonReq != nil {
		req.Header.SetContentType("application/json")
		body, err := sonic.ConfigFastest.Marshal(cfg.jsonReq)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON request: %w", err)
		}
		req.SetBody(body)
	} else if cfg.postForm != nil {
		req.Header.SetContentType("application/x-www-form-urlencoded")
		args := req.PostArgs()
		for k, v := range cfg.postForm {
			args.Set(k, v)
		}
	} else if cfg.multipartForm != nil {
		writer := multipart.NewWriter(req.BodyWriter())
		req.Header.SetContentType(writer.FormDataContentType())
		for k, v := range cfg.multipartForm {
			_ = writer.WriteField(k, v)
		}
		writer.Close()
	} else if cfg.byteReq != nil {
		req.SetBody(cfg.byteReq)
	}

	req.Header.Set("Accept-Encoding", "gzip, br")

	if cfg.timeout > 0 {
		req.SetTimeout(time.Duration(cfg.timeout) * time.Second)
	} else {
		req.SetTimeout(10 * time.Second)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err = fastHttpClient.Do(req, resp); err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	if cfg.headerResp != nil {
		resp.Header.VisitAll(func(k, v []byte) {
			cfg.headerResp[string(k)] = string(v)
		})
	}

	if code := resp.StatusCode(); code > fasthttp.StatusMultipleChoices {
		respBody, _ := resp.BodyUncompressed()

		return &ErrorStatusNotOK{
			Code: code,
			Msg:  string(resp.Header.StatusMessage()),
			Body: string(respBody),
		}
	}

	if cfg.jsonResp != nil || cfg.textResp != nil {
		respBody, err := resp.BodyUncompressed()

		if err != nil {
			return fmt.Errorf("failed to decompress response body: %w", err)
		}

		if cfg.jsonResp != nil {
			if err = sonic.ConfigFastest.Unmarshal(respBody, cfg.jsonResp); err != nil {
				return fmt.Errorf("failed to unmarshal response body: %w", err)
			}
		} else if cfg.textResp != nil {
			*cfg.textResp = string(respBody)
		}
	}

	return nil
}

// DoText makes a request to the given URL and returns the response body as a string
func DoText(url string, options ...func(cfg *Options)) (*string, error) {
	var resp string
	err := Do(url, append(options, WithTextRespTo(&resp))...)
	return &resp, err
}

// DoJSON makes a request to the given URL and returns the response body as a JSON,
// Generic type T is the response struct,
// Example:
//
//	resp, err := https.DoJSON[RespStructType]("http://example.com")
func DoJSON[T any](url string, options ...func(cfg *Options)) (*T, error) {
	var resp T
	err := Do(url, append(options, WithJSONRespTo(&resp))...)
	return &resp, err
}