package ebay

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

const (

	// Item API

	itemPath        = "item"
	itemNewlyListed = "NEWLY_LISTED"
	itemAllActive   = "ALL_ACTIVE"

	// Header
	headerRange        = "Range"
	headerContentRange = "Content-range"
	headerLastModified = "Last-Modified"
)

// FeedService handles the communication with eBay Feed API
// https://developer.ebay.com/api-docs/buy/feed/overview.html
type FeedService service

// FeedInfo containts information about the feed when the download is successful
// In case not content is found for the given feed criteria, the size will be zero
type FeedInfo struct {
	// CategoryID is the eBay category ID associated to this feed
	CategoryID string
	// MarketID is the eBay marketplace ID associated to this feed
	MarketID string
	// Scope is the scope of the feed
	Scope string
	// LastModified is generated date of the feed
	LastModified string
	// Size is the size of the feed file.
	Size int64
}

// WeeklyItemBoostrap downloads the latest weekly item boostrap file
// https://developer.ebay.com/api-docs/buy/feed/resources/item/methods/getItemFeed
func (c *Client) WeeklyItemBoostrap(ctx context.Context, marketID, categoryID string, dst io.Writer) (*FeedInfo, error) {
	return c.download(&feedContext{
		ctx:      ctx,
		baseURL:  c.baseURL,
		pathURL:  itemPath,
		marketID: marketID,
		params: &queryParams{
			Scope:      itemAllActive,
			CategoryID: categoryID,
		},
	}, dst)
}

// queryParams is Feed API query parameters
type queryParams struct {
	Scope      string `url:"feed_scope"`
	CategoryID string `url:"category_id"`
	Date       string `url:"date,omitempty"`
}

// feedContext is a helper struct to reduce the number input parametes passed among the internal methods
type feedContext struct {
	ctx      context.Context
	baseURL  *url.URL
	pathURL  string
	marketID string
	params   *queryParams
}

// download is an helper function which implement the logic to download a multi-parts file feed
func (c *Client) download(feedCtx *feedContext, dst io.Writer) (*FeedInfo, error) {

	var (
		rangeLower int64 = 0
		rangeUpper int64 = c.maxChunkSize
		lenght     int64 = 0
	)

	info := &FeedInfo{
		CategoryID: feedCtx.params.CategoryID,
		Scope:      feedCtx.params.Scope,
		MarketID:   feedCtx.marketID,
	}

	// We want to make at least one request
	responseStatus := http.StatusPartialContent
	// Loop until all chunks are completed
	for responseStatus == http.StatusPartialContent {

		rq, err := buildHTTPRequest(feedCtx, rangeLower, rangeUpper)
		if err != nil {
			return nil, err
		}

		rs, err := c.httpClient.Do(rq)
		if err != nil {
			return nil, err
		}

		responseStatus = rs.StatusCode

		if responseStatus == http.StatusOK || responseStatus == http.StatusPartialContent {
			_, err := io.Copy(dst, rs.Body)
			if err != nil {
				return nil, fmt.Errorf("downloadFile(): impossible to copy response body")
			}

			rangeLower, rangeUpper, lenght, err = processContentRange(rs.Header.Get(headerContentRange))
			rangeLower = rangeUpper + 1
			rangeUpper = rangeUpper + c.maxChunkSize

			info.LastModified = rs.Header.Get(headerLastModified)
			info.Size = lenght
		} else if responseStatus == http.StatusNoContent {
			info.Size = 0
		} else {
			return nil, NewErrorResponse(rs)
		}

		rs.Body.Close()
	}

	return info, nil
}

// buildEndpointURL is an helper function to build the final endpoint URL including the query parameters
func buildEndpointURL(feedCtx *feedContext) (string, error) {
	u, err := feedCtx.baseURL.Parse(feedCtx.pathURL)
	if err != nil {
		return "", fmt.Errorf("buildEndpointURL(): cannot create endpoint URL: %v", err)
	}

	qs, err := query.Values(feedCtx.params)
	if err != nil {
		return "", fmt.Errorf("buildEndpointURL(): cannot parse query parameters: %v", err)
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// buildHTTPRequest is an helper function to build the Feed HTTP request
func buildHTTPRequest(feedCtx *feedContext, rangeLower, rangeUpper int64) (*http.Request, error) {

	endpointURL, err := buildEndpointURL(feedCtx)
	if err != nil {
		return nil, err
	}

	rq, err := http.NewRequest("GET", endpointURL, nil)
	if err != nil {
		return nil, err
	}

	rq.Header.Set(headerMarketplaceID, feedCtx.marketID)
	rq.Header.Set(headerRange, fmt.Sprintf("bytes=%v-%v", rangeLower, rangeUpper))

	rq.WithContext(feedCtx.ctx)

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
