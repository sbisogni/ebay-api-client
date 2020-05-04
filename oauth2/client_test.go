package oauth2

import (
	"testing"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gotest.tools/v3/assert"
)

func TestGetToken(t *testing.T) {
	oAuthConf := clientcredentials.Config{
		ClientID:     "eBay-eBayAdsP-SBX-0196809ca-beaf40a7",
		ClientSecret: "SBX-196809cabc44-a6bb-4d0d-8ce6-3a43",
		Scopes: []string{"https://api.ebay.com/oauth/api_scope",
			"https://api.ebay.com/oauth/api_scope/buy.item.feed"},
		TokenURL: sandboxTokenURL,
	}

	_, OK := oAuthConf.TokenSource(oauth2.NoContext).Token()
	assert.Assert(t, OK)

}
