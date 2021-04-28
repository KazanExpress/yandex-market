package models

// ExploreOffersResponse explore response structure.
type ExploreOffersResponse struct {
	Offers []Offer `json:"offers"`
	Pager  Pager   `json:"pager"`
}

// OfferExploreModel explore response offer model.
type OfferExploreModel struct {
	Bid              float64 `json:"bid"`
	Currency         string  `json:"currency"`
	CutPrice         bool    `json:"cutPrice"`
	Discount         int64   `json:"discount"`
	FeedID           int64   `json:"feedId"`
	ID               string  `json:"id"`
	MarketCategoryID int64   `json:"marketCategoryId"`
	ModelID          int64   `json:"modelId"`
	PreDiscountPrice float64 `json:"preDiscountPrice"`
	Price            float64 `json:"price"`
	ShopCategoryID   string  `json:"shopCategoryId"`
	Name             string  `json:"name"`
	URL              string  `json:"url"`
}

// Pager describes pagination status.
type Pager struct {
	CurrentPage int64 `json:"currentPage"`
	From        int64 `json:"from"`
	PagesCount  int64 `json:"pagesCount"`
	PageSize    int64 `json:"pageSize"`
	To          int64 `json:"to"`
	Total       int64 `json:"total"`
}
