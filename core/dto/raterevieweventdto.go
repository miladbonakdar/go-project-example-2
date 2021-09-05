package dto

type RateReviewEventDto struct {
	PlaceId      string  `json:"PlaceId"`
	Rating       float64 `json:"Rating"`
	ReviewsCount int     `json:"ReviewsCount"`
	ProductType  int     `json:"ProductType"`
}
