package ebay

import (
	"fmt"
	"net/url"
)

// Client is the eBay API client. It provides the methods to access to the eBay resources
type Client struct {
	// httpClient is the HTTP httpClient used to communicate with the API
	httpClient HTTPClient
	// baseURL is the base URL for API request. baseURL should always be
	// specified with a trailing slash
	baseURL *url.URL
	// maxChunkSize is the max chunk size for donwload the feed files
	maxChunkSize int64
	// reuse a singl struct instead of allocating one for each service
	common service
	// Services used to access to the different part of eBay API
	Feed *FeedService
}

type service struct {
	client *Client
}

// NewSandboxClient creates a new eBay API client pointing to eBay Sandbox environment
// The API methods require OAuth2 authentication, provide an HTTPClient that will
// performe the authentication
func NewSandboxClient(httpClient HTTPClient) (*Client, error) {
	return newClient(sboxDefaultBaseURL, sboxDefaultMaxChunkSize, httpClient)
}

// NewProdClient creates a new eBay API client pointing to eBay Sandbox environment
// The API methods require OAuth2 authentication, provide an HTTPClient that will
// performe the authentication
func NewProdClient(httpClient HTTPClient) (*Client, error) {
	return newClient(prodDefaultBaseURL, prodDefaultMaxChunkSize, httpClient)
}

// newClient is helper function to create a new eBay API client
func newClient(baseURL string, maxChunkSize int64, httpClient HTTPClient) (*Client, error) {
	if httpClient == nil {
		return nil, fmt.Errorf("NewClient(): httpClient is nil")
	}

	baseEndpoint, _ := url.Parse(baseURL)

	c := &Client{httpClient: httpClient, baseURL: baseEndpoint, maxChunkSize: maxChunkSize}

	c.common.client = c
	c.Feed = (*FeedService)(&c.common)

	return c, nil
}
