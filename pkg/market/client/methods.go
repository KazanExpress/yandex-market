package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/KazanExpress/yandex-market/pkg/market/models"
)

// ListFeeds returns list of feeds placed in Yandex.Market for given campaign.
func (c *YandexMarketClient) ListFeeds(ctx context.Context, campaignID int64) ([]models.Feed, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/campaigns/%v/feeds", campaignID), "", nil)
	if err != nil {
		return nil, err
	}

	feedResponse := &models.FeedResponse{}

	err = c.executeRequest(req, feedResponse)

	return feedResponse.Feeds, err
}

// RefreshFeed tells Yandex.Market that feed was refreshed.
// After this, Yandex.Market starts updating feed data.
func (c *YandexMarketClient) RefreshFeed(ctx context.Context, campaignID, feedID int64) error {
	req, err := c.newRequest(ctx,
		http.MethodPost,
		fmt.Sprintf("/campaigns/%v/feeds/%v/refresh", campaignID, feedID), "", nil)
	if err != nil {
		return err
	}

	refreshResponse := &models.CommonResponse{}

	err = c.executeRequest(req, refreshResponse)
	if err != nil {
		return err
	}

	if refreshResponse.Status == models.StatusError {
		return fmt.Errorf("failed to refresh feed: %w", refreshResponse.Errors)
	}

	return nil
}

// SetOfferPrices overwrites prices from the feed.
// In single call allowed to set or delete no more than 2000 offers.
func (c *YandexMarketClient) SetOfferPrices(ctx context.Context, campaignID int64, offers []models.Offer) error {
	priceRequest := models.SetPriceRequest{Offers: offers}
	requestBody, err := json.Marshal(priceRequest)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost,
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
		return fmt.Errorf("failed to set prices: %w", setPriceResponse.Errors)
	}

	return nil
}

// GetOfferPrices returns prices set with SetOfferPrices.
func (c *YandexMarketClient) GetOfferPrices(ctx context.Context,
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

	req, err := c.newRequest(ctx, http.MethodGet,
		fmt.Sprintf("/v2/campaigns/%v/offer-prices", campaignID), query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	getPriceResponse := &models.GetPricesResponse{}

	err = c.executeRequest(req, getPriceResponse)
	if err != nil {
		return nil, err
	}

	if getPriceResponse.Status == models.StatusError {
		return nil, fmt.Errorf("failed to get prices: %w", getPriceResponse.Errors)
	}

	return getPriceResponse.Result.Offers, nil
}

// DeleteAllOffersPrices deletes all prices set with API.
// After deleting prices from the feed will be used.
func (c *YandexMarketClient) DeleteAllOffersPrices(ctx context.Context, campaignID int64) error {
	req, err := c.newRequest(ctx, http.MethodPost,
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
		return fmt.Errorf("failed to delete prices: %w", deletePricesResponse.Errors)
	}

	return nil
}

// HideOffers hides offers.
// Can hide up too 500 offers per call.
func (c *YandexMarketClient) HideOffers(
	ctx context.Context,
	campaignID int64,
	offersToHide []models.HiddenOffer,
) error {
	requestModel := models.OfferHideRequest{HiddenOffers: offersToHide}

	requestBody, err := json.Marshal(requestModel)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("v2/campaigns/%v/hidden-offers", campaignID),
		"",
		bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	hideOffersResponse := &models.CommonResponse{}

	err = c.executeRequest(req, hideOffersResponse)
	if err != nil {
		return err
	}

	if hideOffersResponse.Status == models.StatusError {
		return fmt.Errorf("failed to hide offers - %v", hideOffersResponse.Errors)
	}

	return nil
}

// GetHiddenOffers returns list of hidden offers for campaign.
func (c *YandexMarketClient) GetHiddenOffers(ctx context.Context,
	campaignID int64,
	opts ...models.GetHiddenOffersOption,
) (models.GetHiddenOfferResult, error) {
	o := models.GetHiddenOffersOptions{}
	for _, opt := range opts {
		opt(&o)
	}

	query := o.ToQueryArgs()

	req, err := c.newRequest(ctx, http.MethodGet,
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
		return models.GetHiddenOfferResult{}, fmt.Errorf("failed to hide offers: %w", getHiddenOffersResponse.Errors)
	}

	return getHiddenOffersResponse.Result, nil
}

// UnhideOffers unhides offers.
func (c *YandexMarketClient) UnhideOffers(
	ctx context.Context,
	campaignID int64,
	offersToUnhide []models.OfferToUnhide,
) error {
	requestModel := models.OfferUnhideRequest{HiddenOffers: offersToUnhide}
	requestBody, err := json.Marshal(requestModel)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodDelete,
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
		return fmt.Errorf("failed to unhide offers: %w", unhideResponse.Errors)
	}

	return nil
}

// ExploreOffers returns all offers that satisfy passed options.
func (c *YandexMarketClient) ExploreOffers(
	ctx context.Context,
	campaignID int64,
	options models.ExploreOptions,
) (models.ExploreOffersResponse, error) {
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

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/campaigns/%v/offers", campaignID), query.Encode(), nil)
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
