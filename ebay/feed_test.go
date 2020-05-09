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

func newHTTPRequest(rangeLower, rangeUpper int64, categoryID, marketID, endPoint string) *http.Request {

	u, _ := url.Parse(endPoint)
	q, _ := query.Values(queryParams{
		Scope: itemAllActive, CategoryID: categoryID})

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

	u, _ := url.Parse(sboxDefaultBaseURL + itemPath)
	rq, _ := http.NewRequest("GET", u.String(), nil)

	rs.Request = rq
	rs.StatusCode = statusCode
	rs.Body = ioutil.NopCloser(strings.NewReader(body))

	return rs
}

func Test_buildEndpointURL(t *testing.T) {

	var expBaseURL, _ = url.Parse(sboxDefaultBaseURL)

	type args struct {
		feedCtx *feedContext
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "is endpoint URL created without date param",
			args: args{
				feedCtx: &feedContext{
					baseURL: expBaseURL,
					pathURL: itemPath,
					params: &queryParams{
						Scope:      itemAllActive,
						CategoryID: "1",
					},
				},
			},
			want:    "https://api.sandbox.ebay.com/buy/feed/v1_beta/item?category_id=1&feed_scope=ALL_ACTIVE",
			wantErr: false,
		},
		{
			name: "is endpoint URL created with date param",
			args: args{
				feedCtx: &feedContext{
					baseURL: expBaseURL,
					pathURL: itemPath,
					params:  &queryParams{Scope: itemAllActive, CategoryID: "1", Date: "20200419"},
				},
			},
			want:    "https://api.sandbox.ebay.com/buy/feed/v1_beta/item?category_id=1&date=20200419&feed_scope=ALL_ACTIVE",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildEndpointURL(tt.args.feedCtx)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildEndpointURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("buildEndpointURL() = %v, want %v", got, tt.want)
			}
		})
	}
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
			name:         "are content range limits and lenght extracted?",
			contentRange: fmt.Sprintf("%v-%v/%v", expRangeLower, expRangeUpper, expLenght),
			want:         expRangeLower,
			want1:        expRangeUpper,
			want2:        expLenght,
			wantErr:      false,
		},
		{
			name:         "is error if content range missing?",
			contentRange: "",
			wantErr:      true,
		},
		{
			name:         "is error if invalid content range - 1?",
			contentRange: fmt.Sprintf("%v-%v-%v", expRangeLower, expRangeUpper, expLenght),
			wantErr:      true,
		},
		{
			name:         "is error if invalid content range - 2?",
			contentRange: fmt.Sprintf("%v%v/%v", expRangeLower, expRangeUpper, expLenght),
			wantErr:      true,
		},
		{
			name:         "is error if invalid lenght?",
			contentRange: fmt.Sprintf("%v-%v/%v", expRangeLower, expRangeUpper, "expLenght"),
			wantErr:      true,
		},
		{
			name:         "is error if invalid lower range?",
			contentRange: fmt.Sprintf("%v-%v/%v", "expRangeLower", expRangeUpper, expLenght),
			wantErr:      true,
		},
		{
			name:         "is error if invalid higher range?",
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
		expCtx        context.Context = context.Background()
		expBaseURL    string          = sboxDefaultBaseURL
		expPathURL    string          = itemPath
		expMarketID   string          = "EBAY_US"
		expCategoryID string          = "1"
		expScope      string          = itemAllActive
		expRangeLower int64           = 0
		expRangeUpper int64           = 1000
	)

	type args struct {
		feedCtx    *feedContext
		rangeLower int64
		rangeUpper int64
	}

	arguments := func() args {

		u, _ := url.Parse(expBaseURL)

		return args{
			feedCtx: &feedContext{
				ctx:      expCtx,
				baseURL:  u,
				pathURL:  expPathURL,
				params:   &queryParams{Scope: expScope, CategoryID: expCategoryID},
				marketID: expMarketID,
			},
			rangeLower: expRangeLower,
			rangeUpper: expRangeUpper,
		}
	}

	tests := []struct {
		name    string
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			name:    "is feed http response created?",
			args:    arguments(),
			want:    newHTTPRequest(expRangeLower, expRangeUpper, expCategoryID, expMarketID, "https://api.sandbox.ebay.com/buy/feed/v1_beta/item?category_id=1&feed_scope=ALL_ACTIVE"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildHTTPRequest(tt.args.feedCtx, tt.args.rangeLower, tt.args.rangeUpper)
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
		expEndPoint     string = sboxDefaultBaseURL + itemPath
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = sboxDefaultMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndPoint))).
		Return(newHTTPResponse(http.StatusPartialContent, rangeLower, rangeHigher, expLenght, expLastModified, expBodyChunk), nil)

	rangeLower = rangeHigher + 1
	rangeHigher = rangeHigher + sboxDefaultMaxChunkSize

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndPoint))).
		Return(newHTTPResponse(http.StatusPartialContent, rangeLower, rangeHigher, expLenght, expLastModified, expBodyChunk), nil)

	rangeLower = rangeHigher + 1
	rangeHigher = rangeHigher + sboxDefaultMaxChunkSize

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndPoint))).
		Return(newHTTPResponse(http.StatusOK, rangeLower, rangeHigher, expLenght, expLastModified, expBodyChunk), nil)

	client, _ := NewSandboxClient(m)

	buffer := new(bytes.Buffer)
	info, err := client.WeeklyItemBoostrap(context.Background(), expMarketID, expCategoryID, buffer)
	assert.NilError(t, err)

	assert.Equal(t, buffer.String(), expBodyChunk+expBodyChunk+expBodyChunk)
	assert.Equal(t, info.CategoryID, expCategoryID)
	assert.Equal(t, info.MarketID, expMarketID)
	assert.Equal(t, info.Scope, itemAllActive)
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
		expEndPoint     string = sboxDefaultBaseURL + itemPath
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = sboxDefaultMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndPoint))).
		Return(newHTTPResponse(http.StatusOK, rangeLower, rangeHigher, expLenght, expLastModified, expBodyChunk), nil)

	client, _ := NewSandboxClient(m)

	buffer := new(bytes.Buffer)
	info, err := client.WeeklyItemBoostrap(context.Background(), expMarketID, expCategoryID, buffer)
	assert.NilError(t, err)

	assert.Equal(t, buffer.String(), expBodyChunk)
	assert.Equal(t, info.CategoryID, expCategoryID)
	assert.Equal(t, info.MarketID, expMarketID)
	assert.Equal(t, info.Scope, itemAllActive)
	assert.Equal(t, info.Size, expLenght)
	assert.Equal(t, info.LastModified, expLastModified)
}

func Test_IsWeeklyItemBoostrapReturningErrorReponse(t *testing.T) {

	var (
		expMarketID      string = "EBAY_US"
		expCategoryID    string = "1"
		expLenght        int64  = 2000
		expLastModified  string = "Wed, 21 Oct 2015 07:28:00 GMT"
		expEndPoint      string = sboxDefaultBaseURL + itemPath
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
		rangeHigher int64 = sboxDefaultMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndPoint))).
		Return(newHTTPResponse(http.StatusBadRequest, rangeLower, rangeHigher, expLenght, expLastModified, expErrorResponse), nil)

	client, _ := NewSandboxClient(m)

	buffer := new(bytes.Buffer)
	_, err := client.WeeklyItemBoostrap(context.Background(), expMarketID, expCategoryID, buffer)

	assert.Error(t, err, "api error response - GET https://api.sandbox.ebay.com/buy/feed/v1_beta/item: 400 - errors: [{13022 API_BROWSE REQUEST The 'category_id' 200 submitted is not supported. The 'category_id' 200 submitted is not supported. [{categoryId 200}]}] - warning: []")
}

func Test_IsWeeklyItemBoostrapSizeZeroIfNoContentFound(t *testing.T) {

	var (
		expMarketID   string = "EBAY_US"
		expCategoryID string = "1"
		expLenght     int64  = 0
		expEndPoint   string = sboxDefaultBaseURL + itemPath
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_ebay.NewMockHTTPClient(ctrl)

	var (
		rangeLower  int64 = 0
		rangeHigher int64 = sboxDefaultMaxChunkSize
	)

	m.EXPECT().
		Do(gomock.Eq(newHTTPRequest(rangeLower, rangeHigher, expCategoryID, expMarketID, expEndPoint))).
		Return(newHTTPResponse(http.StatusNoContent, 0, 0, 0, "", ""), nil)

	client, _ := NewSandboxClient(m)

	buffer := new(bytes.Buffer)
	info, err := client.WeeklyItemBoostrap(context.Background(), expMarketID, expCategoryID, buffer)
	assert.NilError(t, err)

	assert.Equal(t, buffer.String(), "")
	assert.Equal(t, info.CategoryID, expCategoryID)
	assert.Equal(t, info.MarketID, expMarketID)
	assert.Equal(t, info.Scope, itemAllActive)
	assert.Equal(t, info.Size, expLenght)
	assert.Equal(t, info.LastModified, "")
}
