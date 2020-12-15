package market

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/prometheus/common/log"

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

const HTTPTimeout = 20 * time.Second

func NewClient(opt Options) *YandexMarketClient {
	if opt.Client == nil {
		opt.Client = &http.Client{
			Timeout: HTTPTimeout,
		}
	}

	if opt.APIEndpoint == "" {
		opt.APIEndpoint = "https://api.partner.market.yandex.ru/"
	}

	return &YandexMarketClient{
		options: opt,
	}
}

func (c *YandexMarketClient) newRequest(method, path, query string, body io.Reader) (*http.Request, error) {
	fullURL := c.options.APIEndpoint
	// safe concat of API endpoint and path
	if !strings.HasSuffix(fullURL, "/") {
		fullURL += "/"
	}

	path = strings.TrimPrefix(path, "/")
	fullURL += path

	if !strings.HasSuffix(fullURL, ".json") {
		fullURL += ".json"
	}

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request - %w", err)
	}

	req.Header.Add("authorization",
		fmt.Sprintf("OAuth oauth_token=%s, oauth_client_id=%s",
			c.options.OAuthToken, c.options.OAuthClientID))
	req.Header.Add("user-agent", "KE/yandex-market client github.com/KazanExpress/yandex-market")
	req.Header.Add("accept", "*/*")

	req.URL.RawQuery = query

	return req, nil
}

func (c *YandexMarketClient) executeRequest(req *http.Request, jsonResponse interface{}) error {
	resp, err := c.options.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request - %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return fmt.Errorf("failed to read from response body - %w", err)
	}
	log.Warnf("%s", string(body))
	err = json.Unmarshal(body, jsonResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshall json - %w", err)
	}

	return nil
}

// ListFeeds - Возвращает список прайс-листов, размещенных на Яндекс.Маркете для магазина.
// Также ресурс возвращает результаты автоматических проверок прайс-листов.
func (c *YandexMarketClient) ListFeeds(campaignID int64) ([]models.Feed, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/v2/campaigns/%v/feeds", campaignID), "", nil)
	if err != nil {
		return nil, err
	}

	feedResponse := &models.FeedResponse{}

	err = c.executeRequest(req, feedResponse)

	return feedResponse.Feeds, err
}

// RefreshFeed - Позволяет сообщить, что магазин обновил прайс-лист.
// После этого Яндекс.Маркет начнет обновление данных на сервисе.
func (c *YandexMarketClient) RefreshFeed(campaignID, feedID int64) error {
	req, err := c.newRequest("POST", fmt.Sprintf("/campaigns/%v/feeds/%v/refresh", campaignID, feedID), "", nil)
	if err != nil {
		return err
	}

	refreshResponse := &models.CommonResponse{}

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

// SetOfferPrices - перезаписывает цены из фидов.
// В одном запросе можно установить или удалить цены не более чем для 2000 предложений.
func (c *YandexMarketClient) SetOfferPrices(campaignID int64, offers []models.Offer) error {
	priceRequest := models.SetPriceRequest{Offers: offers}
	requestBody, err := json.Marshal(priceRequest)
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST",
		fmt.Sprintf("/v2/campaigns/%v/offer-prices/updates", campaignID),
		"",
		bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	setPriceResponse := &models.CommonResponse{}

	err = c.executeRequest(req, setPriceResponse)
	if err != nil {
		return err
	}

	if setPriceResponse.Status == models.StatusError {
		// todo: beautify list of errors
		return fmt.Errorf("failed to set prices - %v", setPriceResponse.Errors)
	}

	return nil
}

// GetOfferPrices - возвращает цены установленные через SetOfferPrices.
func (c *YandexMarketClient) GetOfferPrices(
	campaignID int64,
	limit *models.LimitOptions,
	paging *models.PagingOptions) ([]models.GetPriceOfferModel, error) {
	query := url.Values{}

	if limit != nil {
		query.Add("limit", fmt.Sprintf("%v", limit.Limit))
		query.Add("offset", fmt.Sprintf("%v", limit.Offset))
	} else if paging != nil {
		query.Add("page", fmt.Sprintf("%v", paging.Page))
		query.Add("pageSize", fmt.Sprintf("%v", paging.Size))
	}

	req, err := c.newRequest("GET", fmt.Sprintf("/v2/campaigns/%v/offer-prices", campaignID), query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	getPriceResponse := &models.GetPricesResponse{}

	err = c.executeRequest(req, getPriceResponse)
	if err != nil {
		return nil, err
	}

	if getPriceResponse.Status == models.StatusError {
		// todo: beautify list of errors
		return nil, fmt.Errorf("failed to get prices - %v", getPriceResponse.Errors)
	}

	return getPriceResponse.Result.Offers, nil
}

// DeleteAllOffersPrices - удаляет все цены на предложения, установленные через API.
// После удаления начнут действовать цены из прайс-листов.
func (c *YandexMarketClient) DeleteAllOffersPrices(campaignID int64) error {
	req, err := c.newRequest("POST",
		fmt.Sprintf("/v2/campaigns/%v/offer-prices/removals", campaignID),
		"",
		bytes.NewReader([]byte(`{"removeAll": true}`)))
	if err != nil {
		return err
	}

	deletePricesResponse := &models.CommonResponse{}

	err = c.executeRequest(req, deletePricesResponse)
	if err != nil {
		return err
	}

	if deletePricesResponse.Status == models.StatusError {
		// todo: beautify list of errors
		return fmt.Errorf("failed to delete prices - %v", deletePricesResponse.Errors)
	}

	return nil
}

// HideOffers - скрывает предложения магазина.
// Можно передать от одного до 500 предложений.
func (c *YandexMarketClient) HideOffers(campaignID int64, offersToHide []models.HiddenOffer) error {
	requestModel := models.OfferHideRequest{HiddenOffers: offersToHide}

	requestBody, err := json.Marshal(requestModel)
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST",
		fmt.Sprintf("v2/campaigns/%v/hidden-offers", campaignID),
		"",
		bytes.NewReader(requestBody))

	hideOffersResponse := &models.CommonResponse{}

	err = c.executeRequest(req, hideOffersResponse)
	if err != nil {
		return err
	}

	if hideOffersResponse.Status == models.StatusError {
		// todo: beautify list of errors
		return fmt.Errorf("failed to hide offers - %v", hideOffersResponse.Errors)
	}

	return nil
}

// GetHiddenOffers - Возвращает список скрытых предложений магазина.
func (c *YandexMarketClient) GetHiddenOffers(
	campaignID int64,
	pageToken string,
	limits *models.LimitOptions,
	paging *models.PagingOptions) (models.GetHiddenOfferResult, error) {
	query := url.Values{}

	if pageToken != "" {
		query.Set("page_token", pageToken)
	} else if limits != nil {
		query.Add("limit", fmt.Sprintf("%v", limits.Limit))
		query.Add("offset", fmt.Sprintf("%v", limits.Offset))
	} else if paging != nil {
		query.Add("page_number", fmt.Sprintf("%v", paging.Page))
		query.Add("page_size", fmt.Sprintf("%v", paging.Size))
	}

	req, err := c.newRequest("GET",
		fmt.Sprintf("/v2/campaigns/%v/hidden-offers", campaignID),
		query.Encode(),
		nil)
	if err != nil {
		return models.GetHiddenOfferResult{}, err
	}

	getHiddenOffersResponse := &models.GetHiddenOfferResponse{}

	err = c.executeRequest(req, getHiddenOffersResponse)
	if err != nil {
		return models.GetHiddenOfferResult{}, err
	}

	if getHiddenOffersResponse.Status == models.StatusError {
		// todo: beautify list of errors
		return models.GetHiddenOfferResult{}, fmt.Errorf("failed to hide offers - %v", getHiddenOffersResponse.Errors)
	}

	return getHiddenOffersResponse.Result, nil
}

func (c *YandexMarketClient) UnhideOffers(campaignID int64, offersToUnhide []models.OfferToUnhide) error {
	requestModel := models.OfferUnhideRequest{HiddenOffers: offersToUnhide}
	requestBody, err := json.Marshal(requestModel)
	if err != nil {
		return err
	}

	req, err := c.newRequest("DELETE",
		fmt.Sprintf("/v2/campaigns/%v/hidden-offers", campaignID),
		"",
		bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	unhideResponse := &models.CommonResponse{}

	err = c.executeRequest(req, unhideResponse)
	if err != nil {
		return err
	}

	if unhideResponse.Status == models.StatusError {
		// todo: beautify list of errors
		return fmt.Errorf("failed to unhide offers - %v", unhideResponse.Errors)
	}

	return nil
}

func (c *YandexMarketClient) ExploreOffers(campaignID int64, options models.ExploreOptions) (models.ExploreOffersResponse, error) {
	query := url.Values{}

	if options.Currency != "" {
		query.Add("currency", options.Currency)
	}

	if options.ShopCategoryID != "" {
		query.Add("shopCategoryId", options.ShopCategoryID)
	}

	if options.Query != "" {
		query.Add("query", options.Query)
	}

	if options.Page > 0 {
		query.Add("page", fmt.Sprintf("%v", options.Page))
	}

	if options.PageSize > 0 {
		query.Add("pageSize", fmt.Sprintf("%v", options.PageSize))
	}

	if options.FeedID > 0 {
		query.Add("feedId", fmt.Sprintf("%v", options.FeedID))
	}

	query.Add("matched", fmt.Sprintf("%v", options.Matched))

	req, err := c.newRequest("GET", fmt.Sprintf("/v2/campaigns/%v/offers", campaignID), query.Encode(), nil)
	if err != nil {
		return models.ExploreOffersResponse{}, err
	}

	response := models.ExploreOffersResponse{}

	err = c.executeRequest(req, &response)
	if err != nil {
		return models.ExploreOffersResponse{}, err
	}

	return response, nil
}
