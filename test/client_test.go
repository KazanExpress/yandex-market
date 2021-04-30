package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KazanExpress/yandex-market/pkg/market/client"
	"github.com/KazanExpress/yandex-market/pkg/market/models"
)

func TestMain(t *testing.M) {
	os.Exit(t.Run())
}

func getClient() *client.YandexMarketClient {
	return client.NewYandexMarketClient(
		// client.WithUserAgent("Yandex"),
		client.WithOAuth(os.Getenv("OAUTH_TOKEN"), os.Getenv("OAUTH_CLIENT_ID")),
		client.WithHTTPClient(&http.Client{
			Timeout: time.Second * 15,
		}),
	)
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

func getFeedID() int64 {
	feedID := os.Getenv("FEED_ID")
	if feedID == "" {
		return 1
	}

	res, err := strconv.ParseInt(feedID, 10, 64)
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
			got, err := c.ListFeeds(context.Background(), tt.args.campaignID)
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
	feeds, err := c.ListFeeds(context.Background(), campaignID)
	if err != nil {
		t.Fatal(err)
	}

	for _, feed := range feeds {
		tests = append(tests, test{
			name:    fmt.Sprintf("YandexMarketClient.RefreshFeed(%d, %d)", campaignID, feed.ID),
			args:    args{campaignID, feed.ID},
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.RefreshFeed(context.Background(), tt.args.campaignID, tt.args.feedID); (err != nil) != tt.wantErr {
				t.Errorf("YandexMarketClient.RefreshFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestYandexMarketClient_Prices(t *testing.T) {
	c := getClient()
	campaignID := getCampaign()
	offerID := os.Getenv("OFFER_ID")
	feedID := getFeedID()
	discountBase := 300.0
	price := 250.0

	err := c.SetOfferPrices(context.Background(), campaignID, []models.Offer{
		{
			Feed:   models.FeedObj{ID: feedID},
			Delete: false,
			ID:     offerID,
			Price: models.Price{
				CurrencyID:   models.CurrencyRUR,
				DiscountBase: discountBase,
				Value:        price,
			},
		},
	})

	assert.NoError(t, err)

	offerPrices, err := c.GetOfferPrices(context.Background(), campaignID)

	assert.NoError(t, err)
	assert.Len(t, offerPrices, 1, "there should be only 1 product set")

	offerPrice := offerPrices[0]

	assert.Equal(t, offerID, offerPrice.ID, "ids should match")
	assert.Equal(t, discountBase, offerPrice.Price.DiscountBase, "discountBase should match")
	assert.Equal(t, price, offerPrice.Price.Value, "price should match")
}

func TestYandexMarketClient_Hidden(t *testing.T) {
	c := getClient()
	campaignID := getCampaign()
	offerID := os.Getenv("OFFER_ID")
	feedID := getFeedID()
	comment := "Временно закончился на складе"

	initRes, err := c.GetHiddenOffers(context.Background(), campaignID)
	require.NoError(t, err)

	initalHidden := initRes.Total

	err = c.HideOffers(context.Background(), campaignID, []models.HiddenOffer{
		{
			FeedID:     feedID,
			OfferID:    offerID,
			TTLInHours: 12,
			Comment:    comment,
		},
	})
	assert.NoError(t, err)

	res, err := c.GetHiddenOffers(context.Background(), campaignID)

	assert.NoError(t, err)
	assert.NotZero(t, res.Total)

	err = c.UnhideOffers(context.Background(), campaignID, []models.OfferToUnhide{
		{
			FeedID:  feedID,
			OfferID: offerID,
		},
	})

	assert.NoError(t, err)
	assert.NotZero(t, res.Total)

	res, err = c.GetHiddenOffers(context.Background(), campaignID)

	assert.NoError(t, err)
	assert.Equal(t, res.Total, initalHidden)
}

func TestYandexMarketClient_Explore(t *testing.T) {
	c := getClient()
	campaignID := getCampaign()

	result, err := c.ExploreOffers(context.Background(), campaignID, models.WithPaginationExploreOption(1, 10))
	assert.NoError(t, err)

	assert.Greater(t, result.Pager.Total, int64(0))
}

func TestYandexMarketClient_FindRegions(t *testing.T) {
	c := getClient()

	regions, err := c.FindRegions(context.Background(), "Казань")
	require.NoError(t, err)

	t.Log(regions)

	require.NotEmpty(t, regions)
	assert.Equal(t, models.City, regions[0].Type)
	assert.Equal(t, "Казань", regions[0].Name)
}

func TestYandexMarketClient_PointOfSales(t *testing.T) {
	c := getClient()
	campaign := getCampaign()
	// moscowRegionID := int64(213)
	kazanRegionID := int64(43)

	initalOutlets, err := c.ListPointsOfSales(context.Background(), campaign)
	require.NoError(t, err, "failed to list initial outlets")

	initialOutletsCount := len(initalOutlets)

	newOutlet, err := c.CreatePointOfSale(context.Background(), campaign, models.PointOfSale{
		Name:       "KazanExpress",
		Visibility: models.Hidden,
		Type:       models.Depot,
		WorkingSchedule: models.WorkingSchedule{
			WorkInHoliday: false,
			ScheduleItems: []models.ScheduleItem{
				{
					StartDay:  models.Monday,
					EndDay:    models.Friday,
					StartTime: "09:00",
					EndTime:   "20:00",
				},
			},
		},
		Address: models.Address{
			RegionID: kazanRegionID,
			Street:   "ул. Петербургская",
			Number:   "1",
		},
		IsMain: false,
		Phones: []string{
			"+7 (401) 212-22-32",
		},
		DeliveryRules: []models.DeliveryRule{
			{
				Cost:            0,
				MinDeliveryDays: 1,
				MaxDeliveryDays: 2,
				PriceFreePickup: 100,
			},
		},
		Emails: []string{
			"test@mail.ru",
		},

		ShopOutletCode: "PVZ-test-1",
	})

	require.NoError(t, err, "failed to create outlet")

	t.Log(newOutlet)

	newOutlets, err := c.ListPointsOfSales(context.Background(), campaign)
	require.NoError(t, err)

	assert.Equal(t, initialOutletsCount+1, len(newOutlets), "number of outlets should increase after creating")

	err = c.DeletePointOfSale(context.Background(), campaign, newOutlet.OutletID)
	require.NoError(t, err, "failed to delete outlet")

	newOutlets, err = c.ListPointsOfSales(context.Background(), campaign)
	require.NoError(t, err)

	assert.Equal(t, initialOutletsCount, len(newOutlets), "number of outlets should be the same after deleting")
}
