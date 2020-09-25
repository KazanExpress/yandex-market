package models

// Generated by https://quicktype.io

type FeedResponse struct {
	Feeds []Feed `json:"feeds"`
}

type Feed struct {
	ID          int64       `json:"id"`
	URL         string      `json:"url"`
	Download    Download    `json:"download"`
	Content     Content     `json:"content"`
	Publication Publication `json:"publication"`
	Placement   Download    `json:"placement"`
}

type Content struct {
	Status              Status `json:"status"`
	TotalOffersCount    int64  `json:"totalOffersCount"`
	RejectedOffersCount int64  `json:"rejectedOffersCount"`
}

type Download struct {
	Status Status `json:"status"`
}

type Publication struct {
	Full                Full   `json:"full"`
	PriceAndStockUpdate Full   `json:"priceAndStockUpdate"`
	Status              Status `json:"status"`
}

type Full struct {
	FileTime      string `json:"fileTime"`
	PublishedTime string `json:"publishedTime"`
}

type Status string

const (
	StatusNa    Status = "NA"
	StatusOk    Status = "OK"
	StatusError Status = "ERROR"
)