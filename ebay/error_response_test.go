package ebay

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"gotest.tools/assert/cmp"
	"gotest.tools/v3/assert"
)

// Customer compare function as the default DeepEqual complains when
// comparing http.Response.Body value
func Compare(a, b *ErrorResponse) bool {
	return a == b || (cmp.DeepEqual(a.Errors, b.Errors)().Success() &&
		cmp.DeepEqual(a.Warnings, b.Warnings)().Success() &&
		a.Response == b.Response &&
		a.Message == b.Message)

}

func Test_IsNewErrorResponse(t *testing.T) {

	endpointURL, _ := url.Parse(sboxDefaultBaseURL)

	newHTTPResponse := func(status int, body string) *http.Response {
		return &http.Response{
			Request: &http.Request{
				Method: "GET",
				URL:    endpointURL,
			},
			StatusCode: status,
			Body:       ioutil.NopCloser(strings.NewReader(body)),
		}
	}

	httpReponseNotSupportedCategory := newHTTPResponse(http.StatusBadRequest, `{
		"errors": [{
			"errorId": 13022,
			"domain": "API_BROWSE",
			"category": "REQUEST",
			"message": "The 'category_id' 200 submitted is not supported.",
			"longMessage": "The 'category_id' 200 submitted is not supported.",
			"parameters": [{"name": "categoryId", "value": "200"}]
		}]
	}`)

	errorResponseNotSupportedCategory := ErrorResponse{
		Response: httpReponseNotSupportedCategory,
		Message:  "api error response",
		Errors: []ErrorData{
			{
				ErrorID:     13022,
				Domain:      "API_BROWSE",
				Category:    "REQUEST",
				Message:     "The 'category_id' 200 submitted is not supported.",
				LongMessage: "The 'category_id' 200 submitted is not supported.",
				Parameters:  []ErrorDataParam{{Name: "categoryId", Value: "200"}},
			},
		},
	}

	httpReponseInvalidJSONBody := newHTTPResponse(http.StatusBadRequest, "this is not a json body")

	errorResponseInvalidJSONBody := ErrorResponse{
		Response: httpReponseInvalidJSONBody,
		Message:  "this is not a json body",
	}

	tests := []struct {
		name string
		args *http.Response
		want *ErrorResponse
	}{

		{
			name: "is nil returned if error code is 200",
			args: newHTTPResponse(http.StatusOK, ""),
			want: nil,
		},
		{
			name: "is nil returned if error code is 299",
			args: newHTTPResponse(http.StatusOK, ""),
			want: nil,
		},
		{
			name: "is ErrorResponse inizialized from json body",
			args: httpReponseNotSupportedCategory,
			want: &errorResponseNotSupportedCategory,
		},
		{
			name: "is ErrorResponse.Message set if body is invalid json",
			args: httpReponseInvalidJSONBody,
			want: &errorResponseInvalidJSONBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewErrorResponse(tt.args)
			assert.Assert(t, Compare(got, tt.want), "%v != %v", got, tt.want)
		})
	}
}

func Test_IsErrorResponseToString(t *testing.T) {

	endpointURL, _ := url.Parse(sboxDefaultBaseURL)

	httpReponseNotSupportedCategory := &http.Response{
		Request: &http.Request{
			Method: "GET",
			URL:    endpointURL,
		},
		StatusCode: http.StatusBadRequest,
	}

	errorResponseNotSupportedCategory := &ErrorResponse{
		Response: httpReponseNotSupportedCategory,
		Message:  "api error response",
		Errors: []ErrorData{
			{
				ErrorID:     13022,
				Domain:      "API_BROWSE",
				Category:    "REQUEST",
				Message:     "The 'category_id' 200 submitted is not supported.",
				LongMessage: "The 'category_id' 200 submitted is not supported.",
				Parameters:  []ErrorDataParam{{Name: "categoryId", Value: "200"}},
			},
		},
	}

	assert.Error(t, errorResponseNotSupportedCategory,
		"api error response - GET https://api.sandbox.ebay.com/buy/feed/v1_beta/: 400 - errors: [{13022 API_BROWSE REQUEST The 'category_id' 200 submitted is not supported. The 'category_id' 200 submitted is not supported. [{categoryId 200}]}] - warning: []")
}
