package models

// OfferHideRequest hide offers request body structure.
type OfferHideRequest struct {
	HiddenOffers []HiddenOffer `json:"hiddenOffers"`
}

// HiddenOffer is a structure of offer to hide.
type HiddenOffer struct {
	FeedID     int64  `json:"feedId"`
	OfferID    string `json:"offerId"`
	Comment    string `json:"comment"`
	TTLInHours int64  `json:"ttlInHours"`
}

// GetHiddenOfferResponse response structure.
type GetHiddenOfferResponse struct {
	Errors CommonErrors         `json:"errors"`
	Result GetHiddenOfferResult `json:"result"`
	Status Status               `json:"status"`
}

// GetHiddenOfferResult get hidden offers result structure.
type GetHiddenOfferResult struct {
	HiddenOffers []HiddenOffer `json:"hiddenOffers"`
	Total        int64         `json:"total"`
	Paging       Paging        `json:"paging"`
}

// Paging contains page tokens to use in further requests.
type Paging struct {
	PrevPageToken string `json:"prevPageToken"`
	NextPageToken string `json:"nextPageToken"`
}

// OfferUnhideRequest unhide request body.
type OfferUnhideRequest struct {
	HiddenOffers []OfferToUnhide `json:"hiddenOffers"`
}

// OfferToUnhide describes offer to unhide.
type OfferToUnhide struct {
	FeedID  int64  `json:"feedId"`
	OfferID string `json:"offerId"`
}
