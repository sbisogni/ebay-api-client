package ebay

import (
	"bytes"
	"context"
	"ebay-api-client/ebay/mock_ebay"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-querystring/query"
	"gotest.tools/v3/assert"
)

func newHTTPRequest(rangeLower, rangeUpper int64, categoryID, marketID, endpointURL string) *http.Request {

	u, _ := url.Parse(endpointURL)
	q, _ := query.Values(feedParams{
		Scope: scopeAllActive, CategoryID: categoryID})

	u.RawQuery = q.Encode()

	r, _ := http.NewRequest("GET", u.String(), nil)

	r.Header.Set(headerMarketplaceID, marketID)
	r.Header.Set(headerRange, fmt.Sprintf("bytes=%v-%v", rangeLower, rangeUpper))

	r.WithContext(context.Background())

	return r
}

func newHTTPResponse(statusCode int, rangeLower, rangeUpper, lenght int64, lastModified, body string) *http.Response {

	rs := &http.Response{Header: make(http.Header)}

	rs.Header.Set(headerContentRange, fmt.Sprintf("%v-%v/%v", rangeLower, rangeUpper, lenght))
	rs.Header.Set(headerLastModified, lastModified)

	u, _ := url.Parse(defaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem)
	rq, _ := http.NewRequest("GET", u.String(), nil)

	rs.Request = rq
	rs.StatusCode = statusCode
	rs.Body = ioutil.NopCloser(strings.NewReader(body))

	return rs
}

func Test_NewSandboxFeedService(t *testing.T) {
	expHTTPClient := http.DefaultClient
	s := NewSandboxFeedService(expHTTPClient)

	assert.Assert(t, s.httpClient == expHTTPClient)
	assert.Assert(t, s.baseURL == defaultSandboxBaseURL)
	assert.Assert(t, s.version == DefaultAPIVersion)
	assert.Assert(t, s.maxChunkSize == defaultSandboxMaxChunkSize)
}

func Test_NewProdFeedService(t *testing.T) {
	expHTTPClient := http.DefaultClient
	s := NewProdFeedService(expHTTPClient)

	assert.Assert(t, s.httpClient == expHTTPClient)
	assert.Assert(t, s.baseURL == defaultProdBaseURL)
	assert.Assert(t, s.version == DefaultAPIVersion)
	assert.Assert(t, s.maxChunkSize == defaultProdMaxChunkSize)
}

func Test_processContentRange(t *testing.T) {

	var (
		expRangeLower int64 = 0
		expRangeUpper int64 = 1000
		expLenght     int64 = 20000
	)

	tests := []struct {
		name         string
		contentRange string
		want         int64
		want1        int64
		want2        int64
		wantErr      bool
	}{
		{
			name:         "Are content range limits and lenght extracted?",
			contentRange: fmt.Sprintf("%v-%v/%v", expRangeLower, expRangeUpper, expLenght),
			want:         expRangeLower,
			want1:        expRangeUpper,
			want2:        expLenght,
			wantErr:      false,
		},
		{
			name:         "Is error if content range missing?",
			contentRange: "",
			wantErr:      true,
		},
		{
			name:         "Is error if invalid content range - 1?",
			contentRange: fmt.Sprintf("%v-%v-%v", expRangeLower, expRangeUpper, expLenght),
			wantErr:      true,
		},
		{
			name:         "Is error if invalid content range - 2?",
			contentRange: fmt.Sprintf("%v%v/%v", expRangeLower, expRangeUpper, expLenght),
			wantErr:      true,
		},
		{
			name:         "Is error if invalid lenght?",
			contentRange: fmt.Sprintf("%v-%v/%v", expRangeLower, expRangeUpper, "expLenght"),
			wantErr:      true,
		},
		{
			name:         "Is error if invalid lower range?",
			contentRange: fmt.Sprintf("%v-%v/%v", "expRangeLower", expRangeUpper, expLenght),
			wantErr:      true,
		},
		{
			name:         "Is error if invalid higher range?",
			contentRange: fmt.Sprintf("%v-%v/%v", expRangeLower, "expRangeUpper", expLenght),
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := processContentRange(tt.contentRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("processHTTPResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
		})
	}
}

func Test_buildHTTPRequest(t *testing.T) {

	var (
		expMarketID    string = "EBAY_US"
		expCategoryID  string = "1"
		expScope       string = scopeAllActive
		expRangeLower  int64  = 0
		expRangeUpper  int64  = 1000
		expEndpointURL *url.URL
	)

	expEndpointURL, _ = url.Parse(defaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem)

	type args struct {
		endpointURL *url.URL
		params      *feedParams
		rangeLower  int64
		rangeUpper  int64
	}

	tests := []struct {
		name    string
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			name: "Is feed http  created?",
			args: args{
				endpointURL: expEndpointURL,
				params:      &feedParams{Scope: expScope, CategoryID: expCategoryID, marketID: expMarketID},
				rangeLower:  expRangeLower,
				rangeUpper:  expRangeUpper,
			},
			want:    newHTTPRequest(expRangeLower, expRangeUpper, expCategoryID, expMarketID, "https://api.sandbox.ebay.com/buy/feed/v1_beta/item?category_id=1&feed_scope=ALL_ACTIVE"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildHTTPRequest(tt.args.endpointURL, tt.args.params, tt.args.rangeLower, tt.args.rangeUpper)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildHTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildHTTPRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsWeeklyItemBoostrapPrecessingThreeChunks(t *testing.T) {

	var (
		expMarketID     string = "EBAY_US"
		expCategoryID   string = "1"
		expLenght       int64  = 2000
		expBodyChunk    string = "Hello World!"
		expLastModified string = " Wed, 21 Oct 2015 07:28:00 GMT"
		expEndpointURL  string = defaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = defaultSandboxMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndpointURL))).
		Return(newHTTPResponse(http.StatusPartialContent, rangeLower, rangeHigher, expLenght, expLastModified, expBodyChunk), nil)

	rangeLower = rangeHigher + 1
	rangeHigher = rangeHigher + defaultSandboxMaxChunkSize

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndpointURL))).
		Return(newHTTPResponse(http.StatusPartialContent, rangeLower, rangeHigher, expLenght, expLastModified, expBodyChunk), nil)

	rangeLower = rangeHigher + 1
	rangeHigher = rangeHigher + defaultSandboxMaxChunkSize

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndpointURL))).
		Return(newHTTPResponse(http.StatusOK, rangeLower, rangeHigher, expLenght, expLastModified, expBodyChunk), nil)

	client := NewSandboxFeedService(m)

	buffer := new(bytes.Buffer)
	info, err := client.WeeklyItemBoostrap(context.Background(), expMarketID, expCategoryID, buffer)
	assert.NilError(t, err)

	assert.Equal(t, buffer.String(), expBodyChunk+expBodyChunk+expBodyChunk)
	assert.Equal(t, info.CategoryID, expCategoryID)
	assert.Equal(t, info.MarketID, expMarketID)
	assert.Equal(t, info.Scope, scopeAllActive)
	assert.Equal(t, info.Size, expLenght)
	assert.Equal(t, info.LastModified, expLastModified)
}

func Test_IsWeeklyItemBoostrapPrecessingOneChunk(t *testing.T) {

	var (
		expMarketID     string = "EBAY_US"
		expCategoryID   string = "1"
		expLenght       int64  = 2000
		expBodyChunk    string = "Hello World!"
		expLastModified string = "Wed, 21 Oct 2015 07:28:00 GMT"
		expEndpointURL  string = defaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = defaultSandboxMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndpointURL))).
		Return(newHTTPResponse(http.StatusOK, rangeLower, rangeHigher, expLenght, expLastModified, expBodyChunk), nil)

	client := NewSandboxFeedService(m)

	buffer := new(bytes.Buffer)
	info, err := client.WeeklyItemBoostrap(context.Background(), expMarketID, expCategoryID, buffer)
	assert.NilError(t, err)

	assert.Equal(t, buffer.String(), expBodyChunk)
	assert.Equal(t, info.CategoryID, expCategoryID)
	assert.Equal(t, info.MarketID, expMarketID)
	assert.Equal(t, info.Scope, scopeAllActive)
	assert.Equal(t, info.Size, expLenght)
	assert.Equal(t, info.LastModified, expLastModified)
}

func Test_IsWeeklyItemBoostrapReturningErrorReponse(t *testing.T) {

	var (
		expMarketID      string = "EBAY_US"
		expCategoryID    string = "1"
		expLenght        int64  = 2000
		expLastModified  string = "Wed, 21 Oct 2015 07:28:00 GMT"
		expEndpointURL   string = defaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
		expErrorResponse string = `{
		"errors": [{
			"errorId": 13022,
			"domain": "API_BROWSE",
			"category": "REQUEST",
			"message": "The 'category_id' 200 submitted is not supported.",
			"longMessage": "The 'category_id' 200 submitted is not supported.",
			"parameters": [{"name": "categoryId", "value": "200"}]
		}]}`
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = defaultSandboxMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndpointURL))).
		Return(newHTTPResponse(http.StatusBadRequest, rangeLower, rangeHigher, expLenght, expLastModified, expErrorResponse), nil)

	client := NewSandboxFeedService(m)

	buffer := new(bytes.Buffer)
	_, err := client.WeeklyItemBoostrap(context.Background(), expMarketID, expCategoryID, buffer)

	assert.Error(t, err, "API Error\nGET https://api.sandbox.ebay.com/buy/feed/v1_beta/item HTTP/1.1\nHost: api.sandbox.ebay.com\nRespose Code: 400\n"+
		"Erros: [{ErrorID:13022 Domain:API_BROWSE Category:REQUEST Message:The 'category_id' 200 submitted is not supported. LongMessage:The 'category_id' 200 submitted is not supported. "+
		"Parameters:[{Name:categoryId Value:200}]}]\nWarnings: []")
}

func Test_IsWeeklyItemBoostrapSizeZeroIfNoContentFound(t *testing.T) {

	var (
		expMarketID    string = "EBAY_US"
		expCategoryID  string = "1"
		expLenght      int64  = 0
		expEndpointURL string = defaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = defaultSandboxMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndpointURL))).
		Return(newHTTPResponse(http.StatusNoContent, 0, 0, 0, "", ""), nil)

	client := NewSandboxFeedService(m)

	buffer := new(bytes.Buffer)
	info, err := client.WeeklyItemBoostrap(context.Background(), expMarketID, expCategoryID, buffer)
	assert.NilError(t, err)

	assert.Equal(t, buffer.String(), "")
	assert.Equal(t, info.CategoryID, expCategoryID)
	assert.Equal(t, info.MarketID, expMarketID)
	assert.Equal(t, info.Scope, scopeAllActive)
	assert.Equal(t, info.Size, expLenght)
	assert.Equal(t, info.LastModified, "")
}

func Test_IsWeeklyItemBoostrapReturningErrorIfHTTPError(t *testing.T) {
	var (
		expMarketID    string = "EBAY_US"
		expCategoryID  string = "1"
		expEndpointURL string = defaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = defaultSandboxMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndpointURL))).
		Return(nil, fmt.Errorf("HTTP Error"))

	client := NewSandboxFeedService(m)

	buffer := new(bytes.Buffer)
	_, err := client.WeeklyItemBoostrap(context.Background(), expMarketID, expCategoryID, buffer)
	assert.Error(t, err, "HTTP Error")
}
