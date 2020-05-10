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

	ts := newTokenSource(ctx, &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		TokenURL:     tokenURL,
		AuthStyle:    oauth2.AuthStyleInHeader,
	})

	return &http.Client{
		Transport: &oauth2.Transport{
			Base:   contextClient(ctx).Transport,
			Source: ts,
		},
	}, nil
}

// newTokenSource creates a wrapper for the oauth2 client credentials TokenSource which is used to force the token_type value to "bearer"
// The issue is that eBay OAuth2 API returns "Application Access Token" as token_type of the access token instead of "bearer".
// The value is set in outh2.Token.TokenType field which is then used to define how to forge the Authentication Header field in the API request.
// eBay API respondes with "1003 OAuth REQUEST Token type in the Authorization header is invalid"
func newTokenSource(ctx context.Context, conf *clientcredentials.Config) oauth2.TokenSource {
	source := &tokenSource{
		ctx:  ctx,
		conf: conf,
		orig: conf.TokenSource(ctx),
	}
	return source
}

type tokenSource struct {
	ctx  context.Context
	conf *clientcredentials.Config
	orig oauth2.TokenSource
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	tk, err := t.orig.Token()
	if err != nil {
		return tk, err
	}

	// Forcing the TokenType to bearer
	tk.TokenType = "bearer"
	return tk, err
}

// Maintaining compatibility with standard oauth2.internal.transport implementation
func contextClient(ctx context.Context) *http.Client {
	if ctx != nil {
		if hc, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
			return hc
		}
	}

	return http.DefaultClient
}
