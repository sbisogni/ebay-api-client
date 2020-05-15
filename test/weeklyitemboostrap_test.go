package test

import (
	"compress/gzip"
	"context"
	"ebay-api-client/ebay"
	"ebay-api-client/oauth2"
	"io/ioutil"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_WeekItemBoostrap(t *testing.T) {

	filename := "feed.tzv.gz"

	ctx := context.Background()
	httpClient, err := oauth2.NewSandboxClientCredentialsClient(ctx, []string{oauth2.ScopeBuyFeedAPI})
	assert.NilError(t, err)

	feedClient := ebay.NewSandboxFeedService(httpClient)
	//feedClient.ChunkSize = 1816

	feed, err := os.Create(filename)
	assert.NilError(t, err)

	defer feed.Close()
	defer os.Remove(filename)

	info, err := feedClient.WeeklyItemBoostrap(ctx, "EBAY_US", "1", feed)

	assert.NilError(t, err)
	assert.Assert(t, info.Size != 0)

	// Preparing for reading
	feed.Seek(0, 0)

	gunzip, err := gzip.NewReader(feed)
	assert.NilError(t, err)

	gunzip.Multistream(false)

	b, err := ioutil.ReadAll(gunzip)
	assert.NilError(t, err)

	assert.Assert(t, len(b) != 0)
	t.Logf(string(b))
}
