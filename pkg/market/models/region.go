package models

// RegionResponse describes region response body.
type RegionResponse struct {
	CommonResponse
	Regions []RegionObj `json:"regions"`
}

// RegionObj describes region.
type RegionObj struct {
	ID     int64      `json:"id"`
	Name   string     `json:"name"`
	Type   RegionType `json:"type"`
	Parent *RegionObj `json:"parent,omitempty"`
}

// RegionType is a type of region.
type RegionType string

const (

	// Area is a local region.
	Area RegionType = "AREA"
	// City is a major city.
	City RegionType = "CITY"
	// Continent is a continent.
	Continent RegionType = "CONTINENT"
	// Country is a country.
	Country RegionType = "COUNTRY"
	// District is a city district.
	District RegionType = "DISTRICT"
	// MonorailStation is a monorail station.
	MonorailStation RegionType = "MONORAIL_STATION"
	// OverseasTerritory is a territory of a state located on another continent.
	OverseasTerritory RegionType = "OVERSEAS_TERRITORY"
	// Region is a region.
	Region RegionType = "REGION"
	// Republic is a subject of the Russian Federation.
	Republic RegionType = "REPUBLIC"
	// RepublicArea is a district within a republic of the Russian federation.
	RepublicArea RegionType = "REPUBLIC_AREA"
	// SecondaryDistrict is a second-level city district.
	SecondaryDistrict RegionType = "SECONDARY_DISTRICT"
	// Settlement is a settlement.
	Settlement RegionType = "SETTLEMENT"
	// Sub is a suburb.
	Sub RegionType = "SUB"
	// SubwayStation is a metro (subway) station.
	SubwayStation RegionType = "SUBWAY_STATION"
	// Town is a town.
	Town RegionType = "TOWN"
	// Unknown is a unknown region.
	Unknown RegionType = "UNKNOWN"
)
