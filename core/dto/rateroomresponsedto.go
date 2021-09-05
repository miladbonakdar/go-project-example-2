package dto

import "hotel-engine/utils/indraframework"

type RateRoomResponseDto struct {
	Rooms          []RoomOptionDto                `json:"options"`
	HotelId        string                         `json:"hotelId"`
	SessionId      string                         `json:"sessionId"`
	Error          *indraframework.IndraException `json:"error"`
	RequestedRooms []RequestRoomDto               `json:"requestedRooms"`
}

func (a *RateRoomResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}

type RoomOptionDto struct {
	Id              string          `json:"id"`
	ProviderName    string          `json:"providerName"`
	Provider        string          `json:"provider"`
	NonRefundable   bool            `json:"nonRefundable"`
	MealPlan        MealPlanTypeDto `json:"mealPlan"`
	Currency        string          `json:"currency"`
	Rooms           []RoomDto       `json:"rooms"`
	Price           int64           `json:"price"`
	OldPrice        int64           `json:"oldPrice"`
	DiscountPercent int64           `json:"discountPercent"`
	DiscountPrice   int64           `json:"discountPrice"`
	Number          int             `json:"number"`
}

type RoomDto struct {
	ExtraCharge   []interface{} `json:"extraCharge"`
	PricePerNight int64         `json:"pricePerNight"`
	Name          string        `json:"name"`
	NameEn        string        `json:"nameEn"`
	Price         int64         `json:"price"`
	Number        int           `json:"number"`
}
