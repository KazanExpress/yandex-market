package models

import (
	"net/url"
	"strconv"
)

// GetHiddenOffersOptions describes pagination options for get hidden offers request.
// Docs: https://yandex.ru/dev/market/partner/doc/dg/reference/get-campaigns-id-hidden-offers.html .
type GetHiddenOffersOptions struct {
	PageNumber int32
	PageSize   int32

	Limit  int32
	Offset int32

	PageToken string

	OfferID string
	FeedID  int64
}

// ToQueryArgs converts options to query args according to documentation of yandex market API.
func (o *GetHiddenOffersOptions) ToQueryArgs() url.Values {
	query := url.Values{}

	switch {
	case o.PageToken != "":
		query.Set("page_token", o.PageToken)
		query.Add("offset", strconv.Itoa(int(o.Offset)))
	case o.PageNumber != 0 && o.PageSize != 0:
		query.Add("page_number", strconv.Itoa(int(o.PageNumber)))
		query.Add("page_size", strconv.Itoa(int(o.PageSize)))
	default:
		query.Add("limit", strconv.Itoa(int(o.Limit)))
		query.Add("offset", strconv.Itoa(int(o.Offset)))
	}

	if o.FeedID > 0 {
		query.Add("feed_id", strconv.FormatInt(o.FeedID, 10))
	}

	if o.OfferID != "" {
		query.Add("offer_id", o.OfferID)
	}

	return query
}

// GetHiddenOffersOption modifies PaginationOptions.
type GetHiddenOffersOption func(*GetHiddenOffersOptions)

// WithOffset sets offset.
func WithOffset(offset int32) GetHiddenOffersOption {
	return func(o *GetHiddenOffersOptions) {
		o.Offset = offset
	}
}

// WithLimit sets limit.
func WithLimit(limit int32) GetHiddenOffersOption {
	return func(o *GetHiddenOffersOptions) {
		o.Limit = limit
	}
}

// WithPageToken sets page token.
func WithPageToken(token string) GetHiddenOffersOption {
	return func(o *GetHiddenOffersOptions) {
		o.PageToken = token
	}
}

// WithPageSizeAndNumber sets page size and number.
func WithPageSizeAndNumber(size, number int32) GetHiddenOffersOption {
	return func(o *GetHiddenOffersOptions) {
		o.PageNumber = number
		o.PageSize = size
	}
}

// WithFeedID sets feed id.
func WithFeedID(feedID int64) GetHiddenOffersOption {
	return func(o *GetHiddenOffersOptions) {
		o.FeedID = feedID
	}
}

// WithOfferID sets offer id.
func WithOfferID(offerID string) GetHiddenOffersOption {
	return func(o *GetHiddenOffersOptions) {
		o.OfferID = offerID
	}
}
