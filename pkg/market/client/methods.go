package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/KazanExpress/yandex-market/pkg/market/models"
)

// ListFeeds returns list of feeds placed in Yandex.Market for given campaign.
func (c *YandexMarketClient) ListFeeds(ctx context.Context, campaignID int64) ([]models.Feed, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/campaigns/%d/feeds", campaignID), url.Values{}, nil)
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
		fmt.Sprintf("/campaigns/%d/feeds/%d/refresh", campaignID, feedID), url.Values{}, nil)
	if err != nil {
		return err
	}

	refreshResponse := &models.CommonResponse{}

	err = c.executeRequest(req, refreshResponse)
	if err != nil {
		return err
	}

	if refreshResponse.Status.IsError() {
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
		fmt.Sprintf("/v2/campaigns/%d/offer-prices/updates", campaignID),
		url.Values{},
		bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	setPriceResponse := &models.CommonResponse{}

	err = c.executeRequest(req, setPriceResponse)
	if err != nil {
		return err
	}

	if setPriceResponse.Status.IsError() {
		return fmt.Errorf("failed to set prices: %w", setPriceResponse.Errors)
	}

	return nil
}

// GetOfferPrices returns prices set with SetOfferPrices.
func (c *YandexMarketClient) GetOfferPrices(ctx context.Context,
	campaignID int64,
	opts ...models.GetOfferPricesOption,
) ([]models.GetPriceOfferModel, error) {
	o := models.GetOfferPricesOptions{}

	for _, opt := range opts {
		opt(&o)
	}

	query := o.ToQueryArgs()

	req, err := c.newRequest(ctx, http.MethodGet,
		fmt.Sprintf("/v2/campaigns/%d/offer-prices", campaignID), query, nil)
	if err != nil {
		return nil, err
	}

	getPriceResponse := &models.GetPricesResponse{}

	err = c.executeRequest(req, getPriceResponse)
	if err != nil {
		return nil, err
	}

	if getPriceResponse.Status.IsError() {
		return nil, fmt.Errorf("failed to get prices: %w", getPriceResponse.Errors)
	}

	return getPriceResponse.Result.Offers, nil
}

// DeleteAllOffersPrices deletes all prices set with API.
// After deleting prices from the feed will be used.
func (c *YandexMarketClient) DeleteAllOffersPrices(ctx context.Context, campaignID int64) error {
	req, err := c.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/v2/campaigns/%d/offer-prices/removals", campaignID),
		url.Values{},
		strings.NewReader(`{"removeAll": true}`))
	if err != nil {
		return err
	}

	deletePricesResponse := &models.CommonResponse{}

	err = c.executeRequest(req, deletePricesResponse)
	if err != nil {
		return err
	}

	if deletePricesResponse.Status.IsError() {
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
		fmt.Sprintf("v2/campaigns/%d/hidden-offers", campaignID),
		url.Values{},
		bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	hideOffersResponse := &models.CommonResponse{}

	err = c.executeRequest(req, hideOffersResponse)
	if err != nil {
		return err
	}

	if hideOffersResponse.Status.IsError() {
		return fmt.Errorf("failed to hide offers: %w", hideOffersResponse.Errors)
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
		fmt.Sprintf("/v2/campaigns/%d/hidden-offers", campaignID),
		query,
		nil)
	if err != nil {
		return models.GetHiddenOfferResult{}, err
	}

	getHiddenOffersResponse := &models.GetHiddenOfferResponse{}

	err = c.executeRequest(req, getHiddenOffersResponse)
	if err != nil {
		return models.GetHiddenOfferResult{}, err
	}

	if getHiddenOffersResponse.Status.IsError() {
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
		fmt.Sprintf("/v2/campaigns/%d/hidden-offers", campaignID),
		url.Values{},
		bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	unhideResponse := &models.CommonResponse{}

	err = c.executeRequest(req, unhideResponse)
	if err != nil {
		return err
	}

	if unhideResponse.Status.IsError() {
		return fmt.Errorf("failed to unhide offers: %w", unhideResponse.Errors)
	}

	return nil
}

// ExploreOffers returns all offers that satisfy passed options.
func (c *YandexMarketClient) ExploreOffers(
	ctx context.Context,
	campaignID int64,
	opts ...models.ExploreOption,
) (models.ExploreOffersResponse, error) {
	o := models.ExploreOptions{}

	for _, opt := range opts {
		opt(&o)
	}

	query := o.ToQueryArgs()

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/campaigns/%d/offers", campaignID), query, nil)
	if err != nil {
		return models.ExploreOffersResponse{}, err
	}

	response := models.ExploreOffersResponse{}

	err = c.executeRequest(req, &response)
	if err != nil {
		return models.ExploreOffersResponse{}, err
	}

	if response.Status.IsError() {
		return models.ExploreOffersResponse{}, fmt.Errorf("failed to explore offers: %w", response.Errors)
	}

	return response, nil
}

// FindRegions returns information about a region that meets the search conditions specified in the request.
// If multiple regions are found that meet the search conditions, it returns information on each found region
// (but no more than ten regions) in order to determine the desired region by parent regions.
func (c *YandexMarketClient) FindRegions(ctx context.Context, regionName string) ([]models.RegionObj, error) {
	query := url.Values{}
	query.Add("name", regionName)

	req, err := c.newRequest(ctx, http.MethodGet, "/v2/regions", query, nil)
	if err != nil {
		return nil, err
	}

	response := models.RegionResponse{}

	err = c.executeRequest(req, &response)
	if err != nil {
		return nil, err
	}

	if response.Status.IsError() {
		return nil, fmt.Errorf("failed to find region: %w", response.Errors)
	}

	return response.Regions, nil
}

// CreatePointOfSale creates point of sale with given params.
func (c *YandexMarketClient) CreatePointOfSale(
	ctx context.Context,
	campaignID int64,
	outlet models.PointOfSale,
) (models.CreateOutletResult, error) {
	requestBody, err := json.Marshal(outlet)
	if err != nil {
		return models.CreateOutletResult{}, fmt.Errorf("json marshal: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/v2/campaigns/%d/outlets", campaignID), url.Values{}, bytes.NewReader(requestBody))
	if err != nil {
		return models.CreateOutletResult{}, err
	}

	response := models.CreatePointOfSaleResponse{}

	err = c.executeRequest(req, &response)
	if err != nil {
		return models.CreateOutletResult{}, err
	}

	if response.Status.IsError() {
		return models.CreateOutletResult{}, fmt.Errorf("failed to create point of sale: %w", response.Errors)
	}

	return response.Result, nil
}

// UpdatePointOfSale updates point of sale.
func (c *YandexMarketClient) UpdatePointOfSale(
	ctx context.Context,
	campaignID int64,
	outlet models.PointOfSale,
) error {
	requestBody, err := json.Marshal(outlet)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPut,
		fmt.Sprintf("/v2/campaigns/%d/outlets", campaignID), url.Values{}, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	response := models.CommonResponse{}

	err = c.executeRequest(req, &response)
	if err != nil {
		return err
	}

	if response.Status.IsError() {
		return fmt.Errorf("failed to update point of sale: %w", response.Errors)
	}

	return nil
}

// GetPointOfSale returns point of sale by id.
func (c *YandexMarketClient) GetPointOfSale(
	ctx context.Context,
	campaignID, outletID int64,
) (models.PointOfSale, error) {
	req, err := c.newRequest(ctx, http.MethodGet,
		fmt.Sprintf("/v2/campaigns/%d/outlets/%d", campaignID, outletID), url.Values{}, nil)
	if err != nil {
		return models.PointOfSale{}, err
	}

	response := models.GetPointOfSaleResponse{}

	err = c.executeRequest(req, &response)
	if err != nil {
		return models.PointOfSale{}, err
	}

	if response.Status.IsError() {
		return models.PointOfSale{}, fmt.Errorf("failed to get point of sale: %w", response.Errors)
	}

	return response.PointOfSale, nil
}

// DeletePointOfSale deletes the store's point of sale in Yandex.Market.
func (c *YandexMarketClient) DeletePointOfSale(
	ctx context.Context,
	campaignID, outletID int64,
) error {
	req, err := c.newRequest(ctx, http.MethodDelete,
		fmt.Sprintf("/v2/campaigns/%d/outlets/%d", campaignID, outletID), url.Values{}, nil)
	if err != nil {
		return err
	}

	response := models.CommonResponse{}

	err = c.executeRequest(req, &response)
	if err != nil {
		return err
	}

	if response.Status.IsError() {
		return fmt.Errorf("failed to delete point of sale: %w", response.Errors)
	}

	return nil
}

// ListPointsOfSales returns a list of the store's points of sale.
func (c *YandexMarketClient) ListPointsOfSales(
	ctx context.Context,
	campaignID int64,
	opts ...models.GetPointsOfSaleOption,
) ([]models.PointOfSale, error) {
	o := models.GetPointsOfSaleOptions{}

	for _, opt := range opts {
		opt(&o)
	}

	req, err := c.newRequest(ctx, http.MethodGet,
		fmt.Sprintf("/v2/campaigns/%d/outlets", campaignID), o.ToQueryArgs(), nil)
	if err != nil {
		return nil, err
	}

	response := models.GetPointsOfSaleResponse{}

	err = c.executeRequest(req, &response)
	if err != nil {
		return nil, err
	}

	if response.Status.IsError() {
		return nil, fmt.Errorf("failed to get points of sale: %w", response.Errors)
	}

	return response.Outlets, nil
}
