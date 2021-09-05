package dto

import (
	"hotel-engine/utils/indraframework"
)

type SearchResponseDto struct {
	Result    []SearchResponseHotelDto       `json:"result"`
	Filters   []FilterOutput                 `json:"filters"`
	Sorts     []SortDto                      `json:"sorts"`
	TotalHits int                            `json:"totalHits"`
	Error     *indraframework.IndraException `json:"error"`
}

type SearchResponseHotelDto struct {
	PlaceID         string                      `json:"place_id"`
	Id              string                      `json:"id"`
	Type            string                      `json:"type"`
	Kind            string                      `json:"kind"`
	MinNight        int                         `json:"min_night"`
	ReservationType string                      `json:"reservation_type"`
	PaymentType     string                      `json:"paymentType"`
	Name            string                      `json:"name"`
	NameEn          string                      `json:"nameEn"`
	Description     string                      `json:"description"`
	Region          string                      `json:"region"`
	Images          []string                    `json:"images"`
	Image           string                      `json:"image"`
	MinPrice        float64                     `json:"min_price"`
	PricePerNight   float64                     `json:"pricePerNight"`
	Star            int                         `json:"star"`
	Location        SearchResponseLocationDto   `json:"location"`
	Tags            []string                    `json:"tags"`
	Verified        bool                        `json:"verified"`
	RateReview      SearchResponseRateReviewDto `json:"rateReview"`
	Amenities       []string                    `json:"amenities"`
	OldPrice        *int64                      `json:"oldPrice,omitempty"`
	Discount        int                         `json:"discount,omitempty"`
	DiscountPrice   int                         `json:"discountPrice,omitempty"`
	Badges          []SearchResponseBadgeDto    `json:"badges"`
}

type SearchResponseRateReviewDto struct {
	Count int     `json:"count"`
	Score float64 `json:"score"`
}

type SearchResponseBadgeDto struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type SearchResponseLocationDto struct {
	City     string                       `json:"city"`
	CityEn   string                       `json:"cityEn"`
	Geo      SearchResponseLocationGeoDto `json:"geo"`
	Province string                       `json:"province"`
}
type SearchResponseLocationGeoDto struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (a *SearchResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
