package models

// SetPriceRequest is set price body.
type SetPriceRequest struct {
	Offers []Offer `json:"offers"`
}

// Offer describes offer structure.
type Offer struct {
	Feed   FeedObj `json:"feed"`
	ID     string  `json:"id"`
	Delete bool    `json:"delete"`
	Price  Price   `json:"price"`
}

// FeedObj describes feed.
type FeedObj struct {
	ID int64 `json:"id"`
}

// Price describes offer price.
type Price struct {
	CurrencyID   string  `json:"currencyId"`
	Value        float64 `json:"value"`
	DiscountBase float64 `json:"discountBase,omitempty"`
}

// SetPriceResponse set price response structure.
type SetPriceResponse struct {
	Status Status `json:"status"`
}

// GetPricesResponse get price response structure.
type GetPricesResponse struct {
	Errors CommonErrors `json:"errors"`
	Result Result       `json:"result"`
	Status Status       `json:"status"`
}

// Result is a get prices response result.
type Result struct {
	Offers []GetPriceOfferModel `json:"offers"`
	Total  int64                `json:"total"`
}

// GetPriceOfferModel offer model for get price response.
type GetPriceOfferModel struct {
	Feed      Feed   `json:"feed"`
	ID        string `json:"id"`
	Price     Price  `json:"price"`
	UpdatedAt string `json:"updatedAt"`
}
