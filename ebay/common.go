package ebay

const (
	// Sandobox environment
	sboxDefaultBaseURL            = "https://api.sandbox.ebay.com/buy/feed/v1_beta/"
	sboxDefaultMaxChunkSize int64 = 1048576

	// Production environment
	prodDefaultBaseURL            = "https://api.ebay.com/buy/feed/v1_beta/"
	prodDefaultMaxChunkSize int64 = 10485760

	// Header
	headerMarketplaceID = "X-EBAY-C-MARKETPLACE-ID"
)
