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

	assert.Equal(t, client.baseURL.String(), sboxDefaultBaseURL)
	assert.Equal(t, client.maxChunkSize, sboxDefaultMaxChunkSize)
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

	assert.Equal(t, client.baseURL.String(), prodDefaultBaseURL)
	assert.Equal(t, client.maxChunkSize, prodDefaultMaxChunkSize)
	assert.Equal(t, client.httpClient, expectedHTTPClient)
}

func Test_IsNewProdClientReturningErrorIfNilHTTPClient(t *testing.T) {
	_, err := NewProdClient(nil)
	assert.ErrorContains(t, err, "httpClient is nil")
}
