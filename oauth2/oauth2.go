// Package oauth2 provides functionalities to authenticate against eBay oAuth2
package oauth2

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (

	// DefaultAPIVersion is the default oAuth2 API version
	DefaultAPIVersion string = "v1"

	defaultSandboxBaseURL string = ""
	defaultProdBaseURL    string = "https://api.ebay.com/identity/"

	// Token API
	pathOAuth2Token string = "oauth2/token"

	// APIClientID is the os variable name storing the eBay oAuth2 client id
	APIClientID string = "EBAY_API_CLIENT_ID"
	// APIClientSecret is the os variable name storing the eBay oAuth2 client secret
	APIClientSecret string = "EBAY_API_CLIENT_SECRET"
)

const (
	// ScopeBuyFeedAPI is the scope identifier required to access to the buy/feed api
	ScopeBuyFeedAPI string = "https://api.ebay.com/oauth/api_scope/buy.item.feed"
)

// SandBoxEndpoint is eBay oAuth2 Sandbox endpoint URLs
var SandBoxEndpoint = oauth2.Endpoint{
	AuthURL:  "https://api.sandbox.ebay.com/identity/v1/oauth2/auth",
	TokenURL: "https://api.sandbox.ebay.com/identity/v1/oauth2/token",
}

// ProdEndpoint is eBay oAuth2 Prod endpoint URLs
var ProdEndpoint = oauth2.Endpoint{
	AuthURL:  "https://api.ebay.com/identity/v1/oauth2/auth",
	TokenURL: "https://api.ebay.com/identity/v1/oauth2/token",
}

// NewSandboxClientCredentialsClient creates a new HTTP client using oAuth2 client credential token flow and pointing to eBay sandbox environment
func NewSandboxClientCredentialsClient(ctx context.Context, scopes []string) (*http.Client, error) {
	return NewClientCredentialsClient(ctx, SandBoxEndpoint.TokenURL, scopes)
}

// NewProdClientCredentialsClient creates a new HTTP client using oAuth2 client credential token flow and pointing to eBay sandbox environment
func NewProdClientCredentialsClient(ctx context.Context, scopes []string) (*http.Client, error) {
	return NewClientCredentialsClient(ctx, SandBoxEndpoint.TokenURL, scopes)
}

// NewClientCredentialsClient creates a new HTTP client using the oAuth2 client credential token flow.
// The functions requires that os environment variables axists as given in the const APIClientID and APIClientSecret
// Check here for more details https://godoc.org/golang.org/x/oauth2/clientcredentials
func NewClientCredentialsClient(ctx context.Context, tokenURL string, scopes []string) (*http.Client, error) {

	clientID := os.Getenv(APIClientID)
	if clientID == "" {
		return nil, fmt.Errorf("Environment variable %v is not set", APIClientID)
	}

	clientSecret := os.Getenv(APIClientSecret)
	if clientSecret == "" {
		return nil, fmt.Errorf("Environment variable %v is not set", APIClientSecret)
	}
	conf := clientcredentials.Config{ClientID: clientID, ClientSecret: clientSecret, Scopes: scopes, TokenURL: tokenURL}
	return conf.Client(ctx), nil
}
