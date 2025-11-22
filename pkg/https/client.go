package https

import "github.com/valyala/fasthttp"

// fastHttpClient is a pre-configured HTTP client with a read buffer size of 8192.
var fastHttpClient = &fasthttp.Client{
	ReadBufferSize: 8192,
}

// UpdateClient updates the fastHttpClient by the given function.
// This function is useful to update the client's configuration.
// For example, to set a timeout:
//
//	https.UpdateClient(func(client *fasthttp.Client) {
//		client.ReadTimeout = time.Second * 5
//	})
func UpdateClient(f func(client *fasthttp.Client)) {
	f(fastHttpClient)
}