package dto

import "hotel-engine/utils/indraframework"

type HotelBadgeDto struct {
	ID              string                         `json:"id"`
	Text            string                         `json:"text"`
	Icon            string                         `json:"icon"`
	TextColor       string                         `json:"textColor"`
	BackgroundColor string                         `json:"backgroundColor"`
	Error           *indraframework.IndraException `json:"error"`
}

func (a *HotelBadgeDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
