package models

// OfferHideRequest request body structure.
type OfferHideRequest struct {
	HiddenOffers []HiddenOffer `json:"hiddenOffers"`
}

type HiddenOffer struct {
	FeedID     int64  `json:"feedId"`
	OfferID    string `json:"offerId"`
	Comment    string `json:"comment"`
	TTLInHours int64  `json:"ttlInHours"`
}

type GetHiddenOfferResponse struct {
	Errors CommonErrors         `json:"errors"`
	Result GetHiddenOfferResult `json:"result"`
	Status Status               `json:"status"`
}

type GetHiddenOfferResult struct {
	HiddenOffers []HiddenOffer `json:"hiddenOffers"`
	Total        int64         `json:"total"`
	Paging       Paging        `json:"paging"`
}

type Paging struct {
	PrevPageToken string `json:"prevPageToken"`
	NextPageToken string `json:"nextPageToken"`
}

type OfferUnhideRequest struct {
	HiddenOffers []OfferToUnhide `json:"hiddenOffers"`
}

type OfferToUnhide struct {
	FeedID  int64  `json:"feedId"`
	OfferID string `json:"offerId"`
}
