package models

// FeedResponse feeds response structure.
type FeedResponse struct {
	Feeds []Feed `json:"feeds"`
}

// Feed feed structure.
type Feed struct {
	ID          int64       `json:"id"`
	URL         string      `json:"url"`
	Download    Download    `json:"download"`
	Content     Content     `json:"content"`
	Publication Publication `json:"publication"`
	Placement   Download    `json:"placement"`
}

// Content describes feed offer status.
type Content struct {
	Status              Status `json:"status"`
	TotalOffersCount    int64  `json:"totalOffersCount"`
	RejectedOffersCount int64  `json:"rejectedOffersCount"`
}

// Download describes download status.
type Download struct {
	Status Status `json:"status"`
}

// Publication describes publication status.
type Publication struct {
	Full                Time   `json:"full"`
	PriceAndStockUpdate Time   `json:"priceAndStockUpdate"`
	Status              Status `json:"status"`
}

// Time describes action time.
type Time struct {
	FileTime      string `json:"fileTime"`
	PublishedTime string `json:"publishedTime"`
}

// Status is a status.
type Status string

const (
	// StatusNa not available.
	StatusNa Status = "NA"
	// StatusOk everything is ok.
	StatusOk Status = "OK"
	// StatusError some error happened.
	StatusError Status = "ERROR"
)
