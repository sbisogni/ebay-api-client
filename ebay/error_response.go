package ebay

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ErrorResponse reports errors or warning generated by the eBay API.Check for details.
// Check for details https://developer.ebay.com/api-docs/static/handling-error-messages.html
type ErrorResponse struct {
	Response *http.Response
	Message  string
	Errors   []ErrorData `json:"errors"`
	Warnings []ErrorData `json:"warnings"`
}

// ErrorData encodes API error or warning details.
// Check for details https://developer.ebay.com/api-docs/static/handling-error-messages.html
type ErrorData struct {
	ErrorID     int              `json:"errorId"`
	Domain      string           `json:"domain"`
	Category    string           `json:"category"`
	Message     string           `json:"message"`
	LongMessage string           `json:"longMessage"`
	Parameters  []ErrorDataParam `json:"parameters"`
}

// ErrorDataParam provides details about which parameter/value generated the error or the warning
type ErrorDataParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%v - %v %v: %d - errors: %v - warning: %v",
		e.Message, e.Response.Request.Method, e.Response.Request.URL,
		e.Response.StatusCode, e.Errors, e.Warnings)
}

func (e *ErrorData) Error() string {
	return fmt.Sprintf("errorId: %v domain: %v category: %v message: %v",
		e.ErrorID, e.Domain, e.Category, e.Message)
}

// NewErrorResponse creates a new ErrorResponse from the http response
func NewErrorResponse(rs *http.Response) *ErrorResponse {

	status := rs.StatusCode
	// if response is successful (200-299), do nothing
	if http.StatusOK <= status && status < http.StatusMultiStatus {
		return nil
	}

	errorResponse := &ErrorResponse{
		Response: rs,
		Message:  "api error response",
	}
	data, err := ioutil.ReadAll(rs.Body)
	if err == nil && data != nil {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}
