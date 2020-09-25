package market

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(t *testing.M) {
	_ = godotenv.Load()
	t.Run()
}

func getOptions() Options {
	var token = os.Getenv("OAUTH_TOKEN")
	var clientID = os.Getenv("OAUTH_CLIENT_ID")
	return Options{
		OAuthClientID: clientID,
		OAuthToken:    token,
	}
}

func getClient() *YandexMarketClient {
	return NewClient(getOptions())
}

func getCampaign() int64 {
	var camp = os.Getenv("CAMPAIGN")
	if camp == "" {
		return 1
	}
	var res, err = strconv.ParseInt(camp, 10, 64)
	if err != nil {
		return 2
	}
	return res
}

func TestYandexMarketClient_ListFeeds(t *testing.T) {
	type fields struct {
		options Options
	}
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
	var tests = []test{}
	var c = getClient()
	var campaignID = getCampaign()
	var feeds, err = c.ListFeeds(campaignID)
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
