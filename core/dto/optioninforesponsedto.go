package dto

import (
	"hotel-engine/utils/indraframework"
)

type OptionInfoResponseDto struct {
	Detail OptionInfoDetailDto            `json:"detail"`
	Policy OptionCancellationDto          `json:"policy"`
	Error  *indraframework.IndraException `json:"error"`
}

func (a *OptionInfoResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}

type OptionInfoDetailDto struct {
	Price            int64                               `json:"price"`
	HotelId          string                              `json:"hotelId"`
	Provider         string                              `json:"provider"`
	RestrictedMarkup OptionInfoDetailRestrictedMarkupDto `json:"restrictedMarkup"`
	MealPlan         string                              `json:"mealPlan"`
	Currency         string                              `json:"currency"`
	Rooms            []OptionInfoDetailRoomDto           `json:"rooms"`
}

type OptionInfoDetailRoomDto struct {
	Name string `json:"name"`
}

type OptionInfoDetailRestrictedMarkupDto struct {
	Amount int64  `json:"amount"`
	Type   string `json:"type"`
}

type OptionCancellationDto struct {
	NonRefundable   bool     `json:"nonRefundable"`
	GeneralPolicies []string `json:"GeneralPolicies"`
}
