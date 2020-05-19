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
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-querystring/query"
	"gotest.tools/v3/assert"
)

func newHTTPRequest(rangeLower, rangeUpper int64, params *feedParams, endpointURL string) *http.Request {

	u, _ := url.Parse(endpointURL)
	q, _ := query.Values(params)

	u.RawQuery = q.Encode()

	r, _ := http.NewRequest("GET", u.String(), nil)

	r.Header.Set(headerMarketplaceID, params.marketID)
	r.Header.Set(headerRange, fmt.Sprintf("bytes=%v-%v", rangeLower, rangeUpper))

	r.WithContext(context.Background())

	return r
}

func newHTTPResponse(statusCode int, rangeLower, rangeUpper, lenght int64, lastModified, body string) *http.Response {

	rs := &http.Response{Header: make(http.Header)}

	rs.Header.Set(headerContentRange, fmt.Sprintf("%v-%v/%v", rangeLower, rangeUpper, lenght))
	rs.Header.Set(headerLastModified, lastModified)

	u, _ := url.Parse(DefaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem)
	rq, _ := http.NewRequest("GET", u.String(), nil)

	rs.Request = rq
	rs.StatusCode = statusCode
	rs.Body = ioutil.NopCloser(strings.NewReader(body))

	return rs
}

func Test_NewSandboxFeedService(t *testing.T) {
	expHTTPClient := http.DefaultClient
	s := NewSandboxFeedService(expHTTPClient)

	assert.Assert(t, s.HTTPClient == expHTTPClient)
	assert.Assert(t, s.BaseURL == DefaultSandboxBaseURL)
	assert.Assert(t, s.Version == DefaultAPIVersion)
	assert.Assert(t, s.ChunkSize == DefaultSandboxMaxChunkSize)
}

func Test_NewProdFeedService(t *testing.T) {
	expHTTPClient := http.DefaultClient
	s := NewProdFeedService(expHTTPClient)

	assert.Assert(t, s.HTTPClient == expHTTPClient)
	assert.Assert(t, s.BaseURL == DefaultProdBaseURL)
	assert.Assert(t, s.Version == DefaultAPIVersion)
	assert.Assert(t, s.ChunkSize == DefaultProdMaxChunkSize)
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
		expEndpointURL string = DefaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
	)

	endpointURL, _ := url.Parse(expEndpointURL)

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
				endpointURL: endpointURL,
				params:      &feedParams{Scope: expScope, CategoryID: expCategoryID, marketID: expMarketID},
				rangeLower:  expRangeLower,
				rangeUpper:  expRangeUpper,
			},
			want:    newHTTPRequest(expRangeLower, expRangeUpper, &feedParams{Scope: expScope, CategoryID: expCategoryID, marketID: expMarketID}, expEndpointURL),
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

func Test_IsDownloadThreeChunks(t *testing.T) {

	var (
		expMarketID     string = "EBAY_US"
		expCategoryID   string = "1"
		expScope        string = scopeAllActive
		expLenght       int64  = 36
		expBodyChunk    string = "Hello World!"
		expLastModified string = "Wed, 21 Oct 2015 07:28:00 GMT"
		expEndpointURL  string = DefaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
		maxChunkSize    int64  = 12
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = maxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, &feedParams{Scope: expScope, CategoryID: expCategoryID, marketID: expMarketID}, expEndpointURL))).
		Return(newHTTPResponse(http.StatusPartialContent, rangeLower, rangeHigher, expLenght, expLastModified, expBodyChunk), nil)

	rangeLower = rangeHigher + 1
	rangeHigher = rangeHigher + maxChunkSize

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, &feedParams{Scope: expScope, CategoryID: expCategoryID, marketID: expMarketID}, expEndpointURL))).
		Return(newHTTPResponse(http.StatusPartialContent, rangeLower, rangeHigher, expLenght, expLastModified, expBodyChunk), nil)

	rangeLower = rangeHigher + 1
	rangeHigher = rangeHigher + maxChunkSize

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, &feedParams{Scope: expScope, CategoryID: expCategoryID, marketID: expMarketID}, expEndpointURL))).
		Return(newHTTPResponse(http.StatusPartialContent, rangeLower, rangeHigher, expLenght, expLastModified, expBodyChunk), nil)

	client := NewSandboxFeedService(m)
	client.ChunkSize = maxChunkSize

	buffer := new(bytes.Buffer)
	feedParams := &feedParams{Scope: scopeAllActive, marketID: expMarketID, CategoryID: expCategoryID, apiPath: pathGetItem}

	info, err := client.download(context.Background(), feedParams, buffer)
	assert.NilError(t, err)

	assert.Equal(t, buffer.String(), expBodyChunk+expBodyChunk+expBodyChunk)
	assert.Equal(t, info.CategoryID, expCategoryID)
	assert.Equal(t, info.MarketID, expMarketID)
	assert.Equal(t, info.Scope, expScope)
	assert.Equal(t, info.Size, expLenght)
	assert.Equal(t, info.LastModified.Format(time.RFC1123), expLastModified)
}

func Test_IsDownloadOneChunk(t *testing.T) {

	var (
		expMarketID     string = "EBAY_US"
		expCategoryID   string = "1"
		expScope        string = scopeAllActive
		expLenght       int64  = 2000
		expBodyChunk    string = "Hello World!"
		expLastModified string = "Wed, 21 Oct 2015 07:28:00 GMT"
		expEndpointURL  string = DefaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = DefaultSandboxMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, &feedParams{Scope: expScope, CategoryID: expCategoryID, marketID: expMarketID}, expEndpointURL))).
		Return(newHTTPResponse(http.StatusOK, rangeLower, rangeHigher, expLenght, expLastModified, expBodyChunk), nil)

	client := NewSandboxFeedService(m)

	buffer := new(bytes.Buffer)
	feedParams := &feedParams{Scope: scopeAllActive, marketID: expMarketID, CategoryID: expCategoryID, apiPath: pathGetItem}

	info, err := client.download(context.Background(), feedParams, buffer)
	assert.NilError(t, err)

	assert.Equal(t, buffer.String(), expBodyChunk)
	assert.Equal(t, info.CategoryID, expCategoryID)
	assert.Equal(t, info.MarketID, expMarketID)
	assert.Equal(t, info.Scope, expScope)
	assert.Equal(t, info.Size, expLenght)
	assert.Equal(t, info.LastModified.Format(time.RFC1123), expLastModified)
}

func Test_IsDownloadReturningErrorReponse(t *testing.T) {

	var (
		expMarketID      string = "EBAY_US"
		expCategoryID    string = "1"
		expScope         string = scopeAllActive
		expLenght        int64  = 2000
		expLastModified  string = "Wed, 21 Oct 2015 07:28:00 GMT"
		expEndpointURL   string = DefaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
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
		rangeHigher int64 = DefaultSandboxMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, &feedParams{Scope: expScope, CategoryID: expCategoryID, marketID: expMarketID}, expEndpointURL))).
		Return(newHTTPResponse(http.StatusBadRequest, rangeLower, rangeHigher, expLenght, expLastModified, expErrorResponse), nil)

	client := NewSandboxFeedService(m)

	buffer := new(bytes.Buffer)
	feedParams := &feedParams{Scope: scopeAllActive, marketID: expMarketID, CategoryID: expCategoryID, apiPath: pathGetItem}

	_, err := client.download(context.Background(), feedParams, buffer)

	assert.Error(t, err, "API Error\nGET https://api.sandbox.ebay.com/buy/feed/v1_beta/item HTTP/1.1\nHost: api.sandbox.ebay.com\nRespose Code: 400\n"+
		"Erros: [{ErrorID:13022 Domain:API_BROWSE Category:REQUEST Message:The 'category_id' 200 submitted is not supported. LongMessage:The 'category_id' 200 submitted is not supported. "+
		"Parameters:[{Name:categoryId Value:200}]}]\nWarnings: []")
}

func Test_IsDownloadReturningErrorIfNoContentFound(t *testing.T) {

	var (
		expMarketID    string = "EBAY_US"
		expCategoryID  string = "1"
		expScope       string = scopeAllActive
		expEndpointURL string = DefaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = DefaultSandboxMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, &feedParams{Scope: expScope, CategoryID: expCategoryID, marketID: expMarketID}, expEndpointURL))).
		Return(newHTTPResponse(http.StatusNoContent, 0, 0, 0, "", ""), nil)

	client := NewSandboxFeedService(m)

	buffer := new(bytes.Buffer)
	feedParams := &feedParams{Scope: scopeAllActive, marketID: expMarketID, CategoryID: expCategoryID, apiPath: pathGetItem}

	_, err := client.download(context.Background(), feedParams, buffer)

	assert.Error(t, err, "API No Content\nGET https://api.sandbox.ebay.com/buy/feed/v1_beta/item HTTP/1.1\nHost: api.sandbox.ebay.com\nRespose Code: 204\nErros: []\nWarnings: []")
}

func Test_IsDownloadReturningErrorIfHTTPError(t *testing.T) {
	var (
		expMarketID    string = "EBAY_US"
		expCategoryID  string = "1"
		expScope       string = scopeAllActive
		expEndpointURL string = DefaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = DefaultSandboxMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, &feedParams{Scope: expScope, CategoryID: expCategoryID, marketID: expMarketID}, expEndpointURL))).
		Return(nil, fmt.Errorf("HTTP Error"))

	client := NewSandboxFeedService(m)

	buffer := new(bytes.Buffer)
	feedParams := &feedParams{Scope: scopeAllActive, marketID: expMarketID, CategoryID: expCategoryID, apiPath: pathGetItem}

	_, err := client.download(context.Background(), feedParams, buffer)

	assert.Error(t, err, "HTTP Error")
}

func Test_IsDalyNewsItems(t *testing.T) {

	var (
		expMarketID     string = "EBAY_US"
		expCategoryID   string = "1"
		expScope        string = scopeNewlyListed
		expLenght       int64  = 36
		expBody         string = "Hello World!"
		expLastModified string = "Wed, 21 Oct 2015 07:28:00 GMT"
		expEndpointURL  string = DefaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
		expDate         string = "20200517"
		maxChunkSize    int64  = DefaultSandboxMaxChunkSize
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = maxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, &feedParams{Scope: expScope, CategoryID: expCategoryID, marketID: expMarketID, Date: expDate}, expEndpointURL))).
		Return(newHTTPResponse(http.StatusOK, rangeLower, rangeHigher, expLenght, expLastModified, expBody), nil)

	client := NewSandboxFeedService(m)

	date := time.Date(2020, time.May, 17, 0, 0, 0, 0, time.UTC)
	buffer := new(bytes.Buffer)

	info, err := client.DailyNewlyItems(context.Background(), expMarketID, expCategoryID, date, buffer)
	assert.NilError(t, err)

	assert.Equal(t, buffer.String(), expBody)
	assert.Equal(t, info.CategoryID, expCategoryID)
	assert.Equal(t, info.MarketID, expMarketID)
	assert.Equal(t, info.Scope, expScope)
	assert.Equal(t, info.Size, expLenght)
	assert.Equal(t, info.LastModified.Format(time.RFC1123), expLastModified)
}

func Test_IsWeeklyItemBoostrap(t *testing.T) {

	var (
		expMarketID     string = "EBAY_US"
		expCategoryID   string = "1"
		expScope        string = scopeAllActive
		expLenght       int64  = 36
		expBody         string = "Hello World!"
		expLastModified string = "Wed, 21 Oct 2015 07:28:00 GMT"
		expEndpointURL  string = DefaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItem
		maxChunkSize    int64  = DefaultSandboxMaxChunkSize
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = maxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, &feedParams{Scope: expScope, CategoryID: expCategoryID, marketID: expMarketID}, expEndpointURL))).
		Return(newHTTPResponse(http.StatusOK, rangeLower, rangeHigher, expLenght, expLastModified, expBody), nil)

	client := NewSandboxFeedService(m)

	buffer := new(bytes.Buffer)

	info, err := client.WeeklyItemBoostrap(context.Background(), expMarketID, expCategoryID, buffer)
	assert.NilError(t, err)

	assert.Equal(t, buffer.String(), expBody)
	assert.Equal(t, info.CategoryID, expCategoryID)
	assert.Equal(t, info.MarketID, expMarketID)
	assert.Equal(t, info.Scope, expScope)
	assert.Equal(t, info.Size, expLenght)
	assert.Equal(t, info.LastModified.Format(time.RFC1123), expLastModified)
}

func Test_IsItemSnapshot(t *testing.T) {

	var (
		expMarketID     string = "EBAY_US"
		expCategoryID   string = "1"
		expLenght       int64  = 36
		expBody         string = "Hello World!"
		expLastModified string = "Wed, 21 Oct 2015 07:28:00 GMT"
		expDate         string = "2020-05-17T16:00:00.000Z"
		expEndpointURL  string = DefaultSandboxBaseURL + DefaultAPIVersion + "/" + pathGetItemSnapshot
		maxChunkSize    int64  = DefaultSandboxMaxChunkSize
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = maxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, &feedParams{CategoryID: expCategoryID, marketID: expMarketID, SnapshotDate: expDate}, expEndpointURL))).
		Return(newHTTPResponse(http.StatusOK, rangeLower, rangeHigher, expLenght, expLastModified, expBody), nil)

	client := NewSandboxFeedService(m)

	date := time.Date(2020, time.May, 17, 16, 0, 0, 0, time.UTC)
	buffer := new(bytes.Buffer)

	info, err := client.ItemShapshot(context.Background(), expMarketID, expCategoryID, date, buffer)
	assert.NilError(t, err)

	assert.Equal(t, buffer.String(), expBody)
	assert.Equal(t, info.CategoryID, expCategoryID)
	assert.Equal(t, info.MarketID, expMarketID)
	assert.Equal(t, info.Size, expLenght)
	assert.Equal(t, info.LastModified.Format(time.RFC1123), expLastModified)
}
