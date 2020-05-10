package oauth2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
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

func TestIsTokenSettingTokenTypeToBearer(t *testing.T) {

	tkRs := `{
		"access_token": "v^1.1#i^1#p^1#f^0#r^0#I^3#t^H4sIAAAAAAAAAOVYa2wUVRTu9qWI5SEPCRCyjC9onZk7s52d2bG7YUtpuzzawpZSq4DzuFOG7s4Mc2do1xAtjaKJvzAgf0xsY9TEKOgPjRKMMaJGCAEMD0HEREKUgKItESXGeGd3KdtKoNDVNHF/7OSee+653/nOOXPuXNBdOq58S/2Wy2W+Owp7u0F3oc/HjAfjSksqJhQVziwpADkKvt7u+7uLe4p+rEJSMmGJKyCyTANBf1cyYSAxLQwTrm2IpoR0JBpSEiLRUcR4dNlSkaWAaNmmYypmgvDHasKEwvKyFgjIwRCUgxKPhcZVk81mmGAEQVNlVuYCLCsAAeJ5hFwYM5AjGU6YYAELSMCRDGhmeBEIIsNRXDDURvhboI1008AqFCAiabRieq2dA/XGSCWEoO1gI0QkFq2NN0ZjNYsamqvoHFuRLA1xR3JcNHS00FShv0VKuPDG26C0thh3FQUiRNCRzA5DjYrRq2BuA36aacjwlWyID3GSBhWFD+SFylrTTkrOjXF4El0ltbSqCA1Hd1I3YxSzIa+HipMdNWATsRq/91juSgld06EdJhZVRx+NNjUREVgtpUjvL6qiJjJe3UoCJhQUQEiRSBlKWiWQ+OwmGUtZioftstA0VN0jDPkbTKcaYsRwKC9BkcvhBSs1Go12VHM8NLl6oUH+uDYvoJkIus46w4spTGIS/OnhzdkfXO04ti67Dhy0MHwiTU+YkCxLV4nhk+k8zKZOFwoT6xzHEmm6s7OT6gxQpt1OswAwdOuypXFlHUxKxFVdr9aRfvMFpJ52RcEl2oV00UlZGEsXzlMMwGgnIoFKFvCVWd6HwooMl/5DkOMzPbQa8lUdjCJ4aSLwrKbBgKLlozoi2QSlPRxQxjmalOwO6FgJSYGkgvPMTUJbV8UAp7EBQYOkGgxpZGVI00iZU4Mko0EIIJRlJST8X4pkpGkeh4oNnfzleT5yvE6hU2Z9clVbq7T+yQ5jZR3TGWhgO1prNqwSjM5FVh1Hxw26tSZRFwuPtBKu77xiWrDJTOhKKp8MeLWeh0q31SbJdlLVbgqP4zCRwI9RuYs8d8dWqL31CBuQLJ3ySptSzCRtSvid7onWphH7R6JEy26K0nGqU7jY1VHRFLWsWDLpOpKcgLE8toD//vV/Xfd0fDgaUz7haGbCqquZUw2Vji2FNiqUDZHp2vhARzV6jb7Z7IAGfnU6tplIQLuFGXWgRxtfr9bzzMctdJjb8zvPR5sxktdKQsfps3asefavR1OXxlj7ZjghxARAZYgblV8L0/FsTo2dllVsZhysN5Ez8h5zC2dweuhlQKQg/WN6fB+BHt+HhT4f4AHJVID5pUUri4vuJhBueBSSDFU2uyhd0iiktxv4Y9eGVAdMWZJuF5b69JNHlN9zriF6V4MZgxcR44qY8Tm3EmD2tZkSZuK9ZZgajgEMDwSGawP3XZstZqYXTz1Y8fWkbc9UDMwJ6pPbBpzaK8dfuQTKBpV8vpKC4h5fweKoypef2nFC+GFu8r2TF0ovnqnfv7lhl/D6+S/X/rW//dLW8/buN3a6R4+/VnXn5s8X3LPk2SnbfvY3GPxkff3DPWsOuW9u1WbIFz55R27/vi/CPXvs+XOXu4/t6D9Y17dx6lNnS42Wb5Yc7p8/afVz5Jq2yEPvd8z9bkvBhLPE0SNXOqrYvVb5hv5pB4q01J5v5yzf1CbPW/Hq9r59s87t3rEhvukR2Dtrz7S9Fw83H6rrW/wiXT4QOch/0c0veHvfMrl/dq1wVnx5bsO+vrIzp+3HPlu566Wntz/xLv8peYqqe2vm4yc2HTd+/eMu9oFLO79q/GDSyQdf+I0dgH9OPDd9wpRfPj5dUP5TS/2BeZkw/g22lOQZIBIAAA==",
		"expires_in": 7200,
		"token_type": "Application Access Token"
	}`

	apiHandler := http.NewServeMux()
	apiHandler.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(tkRs))
	})

	server := httptest.NewServer(apiHandler)
	defer server.Close()

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, server.Client())
	ts := newTokenSource(ctx, &clientcredentials.Config{
		ClientID:     "EBAY_API_CLIENT_ID",
		ClientSecret: "EBAY_API_CLIENT_SECRET",
		Scopes:       []string{ScopeBuyFeedAPI},
		TokenURL:     server.URL + "/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
	})

	tk, err := ts.Token()
	assert.NilError(t, err)
	assert.Assert(t, tk.TokenType == "bearer")
}
