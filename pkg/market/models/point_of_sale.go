package models

import (
	"net/url"
	"strconv"
)

// PointOfSale is structure needed to create outlet.
type PointOfSale struct {
	Name string     `json:"name"`
	Type OutletType `json:"type"`
	// Coords is a point of sale coordinates.
	// Format: longitude, latitude. Separators: comma and / or space. For example, 20.4522144, 54.7104264.
	Coords          string           `json:"coords,omitempty"`
	IsMain          bool             `json:"isMain"`
	ShopOutletCode  string           `json:"shopOutletCode,omitempty"`
	Visibility      OutletVisibility `json:"visibility"`
	Address         Address          `json:"address"`
	Phones          []string         `json:"phones"`
	WorkingSchedule WorkingSchedule  `json:"workingSchedule"`
	DeliveryRules   []DeliveryRule   `json:"deliveryRules"`
	Emails          []string         `json:"emails"`
}

// Address describes outlet address.
type Address struct {
	RegionID int64  `json:"regionId"`
	Street   string `json:"street"`
	Number   string `json:"number"`
}

// DeliveryRule describes delivery rules.
type DeliveryRule struct {
	Cost              int64 `json:"cost"`
	MinDeliveryDays   int64 `json:"minDeliveryDays"`
	MaxDeliveryDays   int64 `json:"maxDeliveryDays"`
	DeliveryServiceID int64 `json:"deliveryServiceId,omitempty"`
	OrderBefore       int64 `json:"orderBefore,omitempty"`
	PriceFreePickup   int64 `json:"priceFreePickup,omitempty"`
}

// WorkingSchedule contains working schedule details.
type WorkingSchedule struct {
	WorkInHoliday bool           `json:"workInHoliday"`
	ScheduleItems []ScheduleItem `json:"scheduleItems"`
}

// ScheduleItem describes schedule rule.
type ScheduleItem struct {
	StartDay  Day    `json:"startDay"`
	EndDay    Day    `json:"endDay"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// OutletVisibility is a enum for outlet state.
type OutletVisibility string

const (

	// Hidden the point of sale is turned off.
	Hidden OutletVisibility = "HIDDEN"
	// Visible the point of sale turned on.
	Visible OutletVisibility = "VISIBLE"
)

// OutletType is a enum for outlet types.
type OutletType string

const (
	// Depot is pickup point.
	Depot OutletType = "DEPOT"
	// Mixed is point of sale of a mixed type (a retail space and a pickup point).
	Mixed OutletType = "MIXED"
	// Retail is a retail space.
	Retail OutletType = "RETAIL"
)

// Day describes day of week.
type Day string

const (
	// Monday day of week.
	Monday Day = "MONDAY"
	// Tuesday day of week.
	Tuesday Day = "TUESDAY"
	// Wednesday day of week.
	Wednesday Day = "WEDNESDAY"
	// Thursday day of week.
	Thursday Day = "THURSDAY"
	// Friday day of week.
	Friday Day = "FRIDAY"
	// Saturday day of week.
	Saturday Day = "SATURDAY"
	// Sunday day of week.
	Sunday Day = "SUNDAY"
)

// CreatePointOfSaleResponse describes response of create outlet response.
type CreatePointOfSaleResponse struct {
	CommonResponse
	Result CreateOutletResult `json:"result"`
}

// CreateOutletResult contains newly created outlet information.
type CreateOutletResult struct {
	OutletID int64 `json:"id"`
}

// GetPointOfSaleResponse describes response for get point of sale request.
type GetPointOfSaleResponse struct {
	CommonResponse
	PointOfSale
}

// GetPointsOfSaleResponse describes response body for get point of sales request.
type GetPointsOfSaleResponse struct {
	CommonResponse
	Outlets []PointOfSale `json:"outlets"`
	Paging  Paging        `json:"paging"`
	Pager   Pager         `json:"pager"`
}

// GetPointsOfSaleOptions describes options needed to query point of sales.
type GetPointsOfSaleOptions struct {
	PageNumber int32
	PageSize   int32
	Limit      int32
	RegionID   int64

	SellerOutletCode string
	PageToken        string
}

// ToQueryArgs converts options to query args.
func (o GetPointsOfSaleOptions) ToQueryArgs() url.Values {
	query := url.Values{}

	if o.PageToken != "" {
		query.Add("page_token", o.PageToken)
		query.Add("limit", strconv.Itoa(int(o.Limit)))
	} else if o.PageNumber > 0 && o.PageSize > 0 {
		query.Add("page", strconv.Itoa(int(o.PageNumber)))
		query.Add("pageSize", strconv.Itoa(int(o.PageSize)))
	}

	if o.RegionID > 0 {
		query.Add("region_id", strconv.FormatInt(o.RegionID, 10))
	}

	if o.SellerOutletCode != "" {
		query.Add("shop_outlet_code", o.SellerOutletCode)
	}

	return query
}

// GetPointsOfSaleOption modifies GetPointsOfSaleOptions.
type GetPointsOfSaleOption func(*GetPointsOfSaleOptions)

// WithPageNumbersAndSize sets page number and page size.
func WithPageNumbersAndSize(pageNumber, pageSize int32) GetPointsOfSaleOption {
	return func(o *GetPointsOfSaleOptions) {
		o.PageNumber = pageNumber
		o.PageSize = pageSize
	}
}

// WithRegionID sets regionID.
func WithRegionID(regionID int64) GetPointsOfSaleOption {
	return func(o *GetPointsOfSaleOptions) {
		o.RegionID = regionID
	}
}

// WithPageTokenAndLimit sets page token and limit.
func WithPageTokenAndLimit(token string, limit int32) GetPointsOfSaleOption {
	return func(o *GetPointsOfSaleOptions) {
		o.Limit = limit
		o.PageToken = token
	}
}

// WithSellerOutletCode sets seller outlet code.
func WithSellerOutletCode(code string) GetPointsOfSaleOption {
	return func(o *GetPointsOfSaleOptions) {
		o.SellerOutletCode = code
	}
}
