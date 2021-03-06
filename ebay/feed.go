package ebay

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	// DefaultAPIVersion is the default Feed API version
	DefaultAPIVersion string = "v1_beta"

	// DefaultSandboxBaseURL is the default base url for the Feed API in the sandbox environment
	DefaultSandboxBaseURL string = "https://api.sandbox.ebay.com/buy/feed/"
	// DefaultSandboxMaxChunkSize is the default max chunk size supported in the sandbox environment
	DefaultSandboxMaxChunkSize int64 = 1048576

	// DefaultProdBaseURL is the default base url for the Feed API in the production environment
	DefaultProdBaseURL string = "https://api.ebay.com/buy/feed/"
	// DefaultProdMaxChunkSize is the default max chunk size supported in the production environment
	DefaultProdMaxChunkSize int64 = 10485760

	//  Feed API Path
	pathGetItem         string = "item"
	pathGetItemSnapshot string = "item_snapshot"
	pathGetItemGroup    string = "item_group"

	// Feed API Parameters
	scopeNewlyListed string = "NEWLY_LISTED"
	scopeAllActive   string = "ALL_ACTIVE"

	dateFormat         string = "20060102"
	snapshotDataFormat string = "2006-01-02T15:04:05.000Z"

	// Header
	headerRange         string = "Range"
	headerContentRange  string = "Content-range"
	headerLastModified  string = "Last-Modified"
	headerMarketplaceID string = "X-EBAY-C-MARKETPLACE-ID"
)

// FeedService handles the communication with eBay Feed API
// https://developer.ebay.com/api-docs/buy/feed/overview.html
type FeedService struct {
	// HTTPClient is the HTTP client instance
	HTTPClient HTTPClient
	// BaseURL is the Feed API base URL
	BaseURL string
	// Version is the API version to use
	Version string
	// ChunkSize is the size of the chunk used to download the file. Refers to DefaultProdMaxChunkSize and DefaultSandboxMaxChunkSize
	ChunkSize int64
}

// NewSandboxFeedService creates a new FeedService client pointing to eBay Sandbox environment.
func NewSandboxFeedService(httpClient HTTPClient) *FeedService {
	return &FeedService{
		HTTPClient: httpClient,
		BaseURL:    DefaultSandboxBaseURL,
		Version:    DefaultAPIVersion,
		ChunkSize:  DefaultSandboxMaxChunkSize,
	}
}

// NewProdFeedService creates a new FeedService client pointing to eBay Production environment.
func NewProdFeedService(httpClient HTTPClient) *FeedService {
	return &FeedService{
		HTTPClient: httpClient,
		BaseURL:    DefaultProdBaseURL,
		Version:    DefaultAPIVersion,
		ChunkSize:  DefaultProdMaxChunkSize,
	}
}

// FeedInfo containts information about the feed when the download is successful
// In case not content is found for the given feed criteria, the size will be zero
type FeedInfo struct {
	// Type is the type of the feed
	Type string
	// CategoryID is the eBay category ID associated to this feed
	CategoryID string
	// MarketID is the eBay marketplace ID associated to this feed
	MarketID string
	// Scope is the scope of the feed
	Scope string
	// LastModified is generated date of the feed
	LastModified time.Time
	// Size is the size of the feed file.
	Size int64
}

// WeeklyItemBoostrap downloads the latest weekly item boostrap feed for the given eBay market id and category id.
// The feed is written into the given destination which has to implement the io.Writer interface. The feed is encodedd in
// Tab Separated Value (TSV) format and gzip compressed: it is required to gunzip the feed before reading it.
// The function returns a FeedInfo object encoding the information abouth the downloaded file.
// https://developer.ebay.com/api-docs/buy/feed/resources/item/methods/getItemFeed
func (f *FeedService) WeeklyItemBoostrap(ctx context.Context, marketID, categoryID string, dst io.Writer) (*FeedInfo, error) {
	params := &feedParams{Scope: scopeAllActive, CategoryID: categoryID, marketID: marketID, apiPath: pathGetItem}
	return f.download(ctx, params, dst)
}

// DailyNewlyItems downloads the feed containing all the newly listed items for the given eBay market id, category id and date.
// The feed is written into the given destination which has to implement the io.Writer interface. The feed is encodedd in
// Tab Separated Value (TSV) format and gzip compressed: it is required to gunzip the feed before reading it.
// The function returns a FeedInfo object encoding the information abouth the downloaded file.
// https://developer.ebay.com/api-docs/buy/feed/resources/item/methods/getItemFeed
func (f *FeedService) DailyNewlyItems(ctx context.Context, marketID, categoryID string, date time.Time, dst io.Writer) (*FeedInfo, error) {
	params := &feedParams{Scope: scopeNewlyListed, CategoryID: categoryID, marketID: marketID, Date: date.Format(dateFormat), apiPath: pathGetItem}
	return f.download(ctx, params, dst)
}

// ItemShapshot downloads the hourly snapshot feed containing the details of all the items have changed within the specified GMT day and hour for the given category.
// The feed is written into the given destination which has to implement the io.Writer interface. The feed is encodedd in
// Tab Separated Value (TSV) format and gzip compressed: it is required to gunzip the feed before reading it.
// The function returns a FeedInfo object encoding the information abouth the downloaded file.
// https://developer.ebay.com/api-docs/buy/feed/resources/item_snapshot/methods/getItemSnapshotFeed
func (f *FeedService) ItemShapshot(ctx context.Context, marketID, categoryID string, date time.Time, dst io.Writer) (*FeedInfo, error) {
	params := &feedParams{CategoryID: categoryID, marketID: marketID, SnapshotDate: date.Format(snapshotDataFormat), apiPath: pathGetItemSnapshot}
	return f.download(ctx, params, dst)
}

// feedParams is Feed API query parameters
type feedParams struct {
	CategoryID   string `url:"category_id"`
	Scope        string `url:"feed_scope,omitempty"`
	Date         string `url:"date,omitempty"`
	SnapshotDate string `url:"snapshot_date,omitempty"`
	marketID     string
	apiPath      string
}

// download is an helper function which implement the logic to download a multi-parts file feed
func (f *FeedService) download(ctx context.Context, params *feedParams, dst io.Writer) (*FeedInfo, error) {

	var (
		rangeLower   int64  = 0
		rangeUpper   int64  = f.ChunkSize
		lenght       int64  = f.ChunkSize
		lastModified string = ""
	)

	info := &FeedInfo{
		CategoryID: params.CategoryID,
		Scope:      params.Scope,
		MarketID:   params.marketID,
	}

	endpointURL, err := url.Parse(f.BaseURL + f.Version + "/" + params.apiPath)
	if err != nil {
		return nil, fmt.Errorf("download(): cannot create endpoint URL: %v", err)
	}

	responseStatus := http.StatusPartialContent

	// Loop until response is partial and all chunks are completed
	for responseStatus == http.StatusPartialContent && rangeLower < lenght {

		rq, err := buildHTTPRequest(endpointURL, params, rangeLower, rangeUpper)
		if err != nil {
			return nil, err
		}

		rq.WithContext(ctx)

		rs, err := f.HTTPClient.Do(rq)
		if err != nil {
			return nil, err
		}

		responseStatus = rs.StatusCode

		if responseStatus == http.StatusOK || responseStatus == http.StatusPartialContent {
			_, err = io.Copy(dst, rs.Body)
			if err != nil {
				return nil, fmt.Errorf("download(): impossible to copy response body: %v", err)
			}

			rangeLower, rangeUpper, lenght, err = processContentRange(rs.Header.Get(headerContentRange))
			rangeLower = rangeUpper + 1
			rangeUpper = rangeUpper + f.ChunkSize

			lastModified = rs.Header.Get(headerLastModified)

		} else {
			return nil, NewErrorResponse(rs)
		}

		rs.Body.Close()
	}

	info.Size = lenght
	info.LastModified, err = time.Parse(time.RFC1123, lastModified)
	if err != nil {
		return nil, fmt.Errorf("download(): cannot read lastModified: %v", err)
	}

	return info, nil
}

// buildHTTPRequest is an helper function to build the Feed HTTP request
func buildHTTPRequest(endpointURL *url.URL, params *feedParams, rangeLower, rangeUpper int64) (*http.Request, error) {

	qs, err := query.Values(params)
	if err != nil {
		return nil, fmt.Errorf("buildHTTPRequest(): cannot parse query parameters: %v", err)
	}

	endpointURL.RawQuery = qs.Encode()

	rq, err := http.NewRequest("GET", endpointURL.String(), nil)
	if err != nil {
		return nil, err
	}

	rq.Header.Set(headerMarketplaceID, params.marketID)
	rq.Header.Set(headerRange, fmt.Sprintf("bytes=%v-%v", rangeLower, rangeUpper))

	return rq, nil
}

// processContentRange is an helper function to process the Content-Range paremeter from the HTTP response
// The fuciont returns the content range lower/upper limites and the total lenght of the file to download
func processContentRange(c string) (int64, int64, int64, error) {

	var (
		rangeLower int64
		rangeUpper int64
		lenght     int64
	)

	parts := strings.Split(c, "/")
	if len(parts) != 2 {
		return 0, 0, 0, fmt.Errorf("processContentRange(): %v has invalid format: %v", headerContentRange, c)
	}

	lenght, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	parts = strings.Split(parts[0], "-")
	if len(parts) != 2 {
		return 0, 0, 0, fmt.Errorf("processContentRange(): %v has invalid format: %v", headerContentRange, c)
	}

	rangeLower, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	rangeUpper, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	return rangeLower, rangeUpper, lenght, nil
}
