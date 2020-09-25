package market

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/KazanExpress/yandex-market/pkg/market/models"
)

type YandexMarketClient struct {
	options Options
}

type Options struct {
	OAuthToken    string
	OAuthClientID string
	APIEndpoint   string
	Client        *http.Client
}

func NewClient(opt Options) *YandexMarketClient {
	if opt.Client == nil {
		opt.Client = &http.Client{
			Timeout: 20 * time.Second,
		}
	}
	return &YandexMarketClient{
		options: opt,
	}
}

func (c *YandexMarketClient) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	var fullURL = c.options.APIEndpoint
	// safe concat of API endpoint and path
	if !strings.HasSuffix(fullURL, "/") {
		fullURL = fullURL + "/"
	}
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	fullURL = fullURL + path
	if !strings.HasSuffix(fullURL, ".json") {
		fullURL = fullURL + ".json"
	}
	var req, err = http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request - %w", err)
	}
	req.Header.Add("authorization", fmt.Sprintf("OAuth oauth_token=%s, oauth_client_id=%s", c.options.OAuthToken, c.options.OAuthClientID))
	req.Header.Add("user-agent", "KE/yandex-market client github.com/KazanExpress/yandex-market")
	req.Header.Add("accept", "*/*")

	return req, nil
}

func (c *YandexMarketClient) executeRequest(req *http.Request, jsonResponse interface{}) error {
	resp, err := c.options.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request - %w", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read from response body - %w", err)
	}

	err = json.Unmarshal(body, jsonResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshall json - %w", err)
	}
	return nil
}

// ListFeeds - Возвращает список прайс-листов, размещенных на Яндекс.Маркете для магазина. Также ресурс возвращает результаты автоматических проверок прайс-листов.
func (c *YandexMarketClient) ListFeeds(campaignID int64) ([]models.Feed, error) {
	var req, err = c.newRequest("GET", fmt.Sprintf("/v2/campaigns/%v/feeds", campaignID), nil)
	if err != nil {
		return nil, err
	}

	var feedResponse = &models.FeedResponse{}

	return feedResponse.Feeds, c.executeRequest(req, feedResponse)
}

// RefreshFeed - Позволяет сообщить, что магазин обновил прайс-лист. После этого Яндекс.Маркет начнет обновление данных на сервисе.
func (c *YandexMarketClient) RefreshFeed(campaignID, feedID int64) error {
	var req, err = c.newRequest("POST", fmt.Sprintf("/campaigns/%v/feeds/%v/refresh", campaignID, feedID), nil)
	if err != nil {
		return err
	}

	var refreshResponse = &models.RefreshResponse{}

	err = c.executeRequest(req, refreshResponse)
	if err != nil {
		return err
	}
	if refreshResponse.Status == models.StatusError {
		// todo: beautify list of errors
		return fmt.Errorf("failed to refresh feed - %v", refreshResponse.Errors)
	}
	return nil
}
