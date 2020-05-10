package oauth2

import (
	"os"
	"testing"

	"golang.org/x/oauth2"
	"gotest.tools/assert"
)

func shutdown() {
	os.Unsetenv(APIClientID)
	os.Unsetenv(APIClientSecret)
}

func TestIsNewClientCredentialsClientCreated(t *testing.T) {

	os.Setenv(APIClientID, "EBAY_CLIENT_ID")
	os.Setenv(APIClientSecret, "EBAY_CLIENT_SECRET")

	_, err := NewClientCredentialsClient(oauth2.NoContext, SandBoxEndpoint.TokenURL, []string{ScopeBuyFeedAPI})
	assert.NilError(t, err)

	shutdown()
}

func TestIsNewClientCredentialsClientReturningErrorIfNoAPIClientID(t *testing.T) {
	os.Setenv(APIClientSecret, "EBAY_CLIENT_SECRET")

	_, err := NewClientCredentialsClient(oauth2.NoContext, SandBoxEndpoint.TokenURL, []string{ScopeBuyFeedAPI})
	assert.Error(t, err, "Environment variable EBAY_API_CLIENT_ID is not set")

	shutdown()
}

func TestIsNewClientCredentialsClientReturningErrorIfNoAPIClientSecret(t *testing.T) {
	os.Setenv(APIClientID, "EBAY_CLIENT_ID")

	_, err := NewClientCredentialsClient(oauth2.NoContext, SandBoxEndpoint.TokenURL, []string{ScopeBuyFeedAPI})
	assert.Error(t, err, "Environment variable EBAY_API_CLIENT_SECRET is not set")

	shutdown()
}

func TestIsNewSandboxClientCredentialsClientCreated(t *testing.T) {
	os.Setenv(APIClientID, "EBAY_CLIENT_ID")
	os.Setenv(APIClientSecret, "EBAY_CLIENT_SECRET")

	_, err := NewSandboxClientCredentialsClient(oauth2.NoContext, []string{ScopeBuyFeedAPI})
	assert.NilError(t, err)
}

func TestIsNewProdClientCredentialsClientCreated(t *testing.T) {
	os.Setenv(APIClientID, "EBAY_CLIENT_ID")
	os.Setenv(APIClientSecret, "EBAY_CLIENT_SECRET")

	_, err := NewProdClientCredentialsClient(oauth2.NoContext, []string{ScopeBuyFeedAPI})
	assert.NilError(t, err)
}
