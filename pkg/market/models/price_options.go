package models

import (
	"net/url"
	"strconv"
)

// GetOfferPricesOptions describes pagination option for offers with prices set via API.
// Docs: https://yandex.ru/dev/market/partner/doc/dg/reference/get-campaigns-id-offer-prices.html
type GetOfferPricesOptions struct {
	CommonPagingOptions
}

// GetOfferPricesOption modifies GetOfferPricesOptions.
type GetOfferPricesOption func(*GetOfferPricesOptions)

// ToQueryArgs converts options to query args according to documentation of yandex market API.
func (o GetOfferPricesOptions) ToQueryArgs() url.Values {
	query := url.Values{}

	switch {
	case o.PageNumber != 0 && o.PageSize != 0:
		query.Add("page", strconv.Itoa(int(o.PageNumber)))
		query.Add("pageSize", strconv.Itoa(int(o.PageSize)))
	default:
		query.Add("limit", strconv.Itoa(int(o.Limit)))
		query.Add("offset", strconv.Itoa(int(o.Offset)))
	}

	return query
}

// WithLimitAndOffset sets limit and offset.
func WithLimitAndOffset(limit, offset int32) GetOfferPricesOption {
	return func(o *GetOfferPricesOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

// WithPageNumberAndSize sets page number and page size.
func WithPageNumberAndSize(number, size int32) GetOfferPricesOption {
	return func(o *GetOfferPricesOptions) {
		o.PageNumber = number
		o.PageSize = size
	}
}
