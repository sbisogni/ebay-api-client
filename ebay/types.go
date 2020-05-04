package ebay

import "net/http"

// HTTPClient is the interface for interacting through the HTTP channel
type HTTPClient interface {
	// Do perforces an HTTP request and return an HTTP response or an error
	Do(req *http.Request) (*http.Response, error)
}
