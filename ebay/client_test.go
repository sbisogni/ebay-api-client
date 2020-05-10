package ebay

import (
	"net/http"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_IsNewSandboxClientCreated(t *testing.T) {

	var expectedHTTPClient = &http.Client{}

	client, OK := NewSandboxClient(expectedHTTPClient)
	assert.Assert(t, OK)

	assert.Equal(t, client.baseURL.String(), defaultSandboxBaseURL)
	assert.Equal(t, client.httpClient, expectedHTTPClient)
}

func Test_IsNewSandboxClientReturningErrorIfNilHTTPClient(t *testing.T) {
	_, err := NewSandboxClient(nil)
	assert.ErrorContains(t, err, "httpClient is nil")
}

func Test_IsNewProdClient(t *testing.T) {

	var expectedHTTPClient = &http.Client{}

	client, OK := NewProdClient(expectedHTTPClient)
	assert.Assert(t, OK)

	assert.Equal(t, client.baseURL.String(), defaultProdBaseURL)
	assert.Equal(t, client.httpClient, expectedHTTPClient)
}

func Test_IsNewProdClientReturningErrorIfNilHTTPClient(t *testing.T) {
	_, err := NewProdClient(nil)
	assert.ErrorContains(t, err, "httpClient is nil")
}
