package market

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/KazanExpress/yandex-market/pkg/market/models"
)

func TestMain(t *testing.M) {
	_ = godotenv.Load()

	os.Exit(t.Run())
}

func getOptions() Options {
	token := os.Getenv("OAUTH_TOKEN")
	clientID := os.Getenv("OAUTH_CLIENT_ID")

	return Options{
		OAuthClientID: clientID,
		OAuthToken:    token,
		APIEndpoint:   "https://api.partner.market.yandex.ru",
	}
}

func getClient() *YandexMarketClient {
	return NewClient(getOptions())
}

func getCampaign() int64 {
	camp := os.Getenv("CAMPAIGN")
	if camp == "" {
		return 1
	}

	res, err := strconv.ParseInt(camp, 10, 64)
	if err != nil {
		return 2
	}

	return res
}

func TestYandexMarketClient_ListFeeds(t *testing.T) {
	type args struct {
		campaignID int64
	}

	tests := []struct {
		name    string
		args    args
		wantLen int
		wantErr bool
	}{
		{
			name:    "simple test",
			args:    args{getCampaign()},
			wantErr: false,
			wantLen: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := getClient()
			got, err := c.ListFeeds(tt.args.campaignID)
			if (err != nil) != tt.wantErr {
				t.Errorf("YandexMarketClient.ListFeeds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("len(YandexMarketClient.ListFeeds()) = %v, want %v", got, tt.wantLen)
			}
		})
	}
}

func TestYandexMarketClient_RefreshFeed(t *testing.T) {
	type args struct {
		campaignID int64
		feedID     int64
	}

	type test struct {
		name    string
		args    args
		wantErr bool
	}

	tests := []test{}
	c := getClient()
	campaignID := getCampaign()
	feeds, err := c.ListFeeds(campaignID)
	if err != nil {
		t.Fatal(err)
	}

	for _, feed := range feeds {
		tests = append(tests, test{
			name:    fmt.Sprintf("YandexMarketClient.RefreshFeed(%v, %v)", campaignID, feed.ID),
			args:    args{campaignID, feed.ID},
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.RefreshFeed(tt.args.campaignID, tt.args.feedID); (err != nil) != tt.wantErr {
				t.Errorf("YandexMarketClient.RefreshFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestYandexMarketClient_Prices(t *testing.T) {
	c := getClient()
	campaignID := getCampaign()
	offerID := "169690W424150"
	feedID := int64(820450)
	discountBase := 300.0
	price := 250.0

	err := c.SetOfferPrices(campaignID, []models.Offer{
		{
			Feed:   models.FeedObj{ID: feedID},
			Delete: false,
			ID:     offerID,
			Price: models.Price{
				CurrencyID:   "RUR",
				DiscountBase: discountBase,
				Value:        price,
			},
		},
	})

	assert.NoError(t, err)

	offerPrices, err := c.GetOfferPrices(campaignID, nil, nil)

	assert.NoError(t, err)
	assert.Len(t, offerPrices, 1, "there should be only 1 product set")

	offerPrice := offerPrices[0]

	assert.Equal(t, offerID, offerPrice.ID, "ids should match")
	assert.Equal(t, discountBase, offerPrice.Price.DiscountBase, "discountBase should match")
	assert.Equal(t, price, offerPrice.Price.Value, "price should match")

	// WARN: uncomment accurately
	// err = c.DeleteAllOffersPrices(campaignID)
	// assert.NoError(t, err)

	// offerPrices, err = c.GetOfferPrices(campaignID, nil, nil)

	// assert.NoError(t, err)
	// assert.Len(t, offerPrices, 0, "there should be no product price set")
}

func TestYandexMarketClient_Hidden(t *testing.T) {
	c := getClient()
	campaignID := getCampaign()
	offerID := "169690W424150"
	feedID := int64(820450)
	comment := "Временно закончился на складе"

	err := c.HideOffers(campaignID, []models.HiddenOffer{
		{
			FeedID:     feedID,
			OfferID:    offerID,
			TTLInHours: 12,
			Comment:    comment,
		},
	})
	assert.NoError(t, err)

	res, err := c.GetHiddenOffers(campaignID, "", nil, nil)

	assert.NoError(t, err)
	assert.NotZero(t, res.Total)

	err = c.UnhideOffers(campaignID, []models.OfferToUnhide{
		{
			FeedID:  feedID,
			OfferID: offerID,
		},
	})

	assert.NoError(t, err)
	assert.NotZero(t, res.Total)

	res, err = c.GetHiddenOffers(campaignID, "", nil, nil)

	assert.NoError(t, err)
	assert.Zero(t, res.Total)
}

func TestYandexMarketClient_Explore(t *testing.T) {
	c := getClient()
	campaignID := getCampaign()

	result, err := c.ExploreOffers(campaignID, models.ExploreOptions{Page: 1})
	assert.NoError(t, err)

	assert.Greater(t, result.Pager.Total, int64(0))
}
