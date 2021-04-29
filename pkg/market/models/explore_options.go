package models

import (
	"net/url"
	"strconv"
)

// ExploreOptions describes exploring option to get campaign offers.
// Doc: https://yandex.ru/dev/market/partner/doc/dg/reference/get-campaigns-id-offers.html.
type ExploreOptions struct {
	Currency       Currency
	FeedID         int64
	Matched        bool
	Query          string
	ShopCategoryID string
	PageNumber     int32
	PageSize       int32
}

// ToQueryArgs converts options to query args according to documentation of yandex market API.
func (o ExploreOptions) ToQueryArgs() url.Values {
	query := url.Values{}

	if o.Currency != "" {
		query.Add("currency", string(o.Currency))
	}

	if o.ShopCategoryID != "" {
		query.Add("shopCategoryId", o.ShopCategoryID)
	}

	if o.Query != "" {
		query.Add("query", o.Query)
	}

	if o.PageNumber > 0 {
		query.Add("page", strconv.Itoa(int(o.PageNumber)))
	}

	if o.PageSize > 0 {
		query.Add("pageSize", strconv.Itoa(int(o.PageSize)))
	}

	if o.FeedID > 0 {
		query.Add("feedId", strconv.Itoa(int(o.FeedID)))
	}

	query.Add("matched", strconv.FormatBool(o.Matched))

	return query
}

// ExploreOption modifies ExploreOptions.
type ExploreOption func(*ExploreOptions)

// WithCurrencyExploreOption sets currency.
func WithCurrencyExploreOption(currency Currency) ExploreOption {
	return func(o *ExploreOptions) {
		o.Currency = currency
	}
}

// WithFeedIDExploreOption sets feedID.
func WithFeedIDExploreOption(feedID int64) ExploreOption {
	return func(o *ExploreOptions) {
		o.FeedID = feedID
	}
}

// WithMatchedExploreOption sets Matched.
func WithMatchedExploreOption(matched bool) ExploreOption {
	return func(o *ExploreOptions) {
		o.Matched = matched
	}
}

// WithQueryExploreOption sets Query.
func WithQueryExploreOption(query string) ExploreOption {
	return func(o *ExploreOptions) {
		o.Query = query
	}
}

// WithShopCategoryIDExploreOption sets ShopCategoryID.
func WithShopCategoryIDExploreOption(shopCategoryID string) ExploreOption {
	return func(o *ExploreOptions) {
		o.ShopCategoryID = shopCategoryID
	}
}

// WithPaginationExploreOption sets pageNumber and pageSize.
func WithPaginationExploreOption(pageNumber, pageSize int32) ExploreOption {
	return func(o *ExploreOptions) {
		o.PageNumber = pageNumber
		o.PageSize = pageSize
	}
}
