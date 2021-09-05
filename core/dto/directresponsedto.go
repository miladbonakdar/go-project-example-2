package dto

import "hotel-engine/utils/indraframework"

type DirectResponseDto struct {
	SessionId string                         `json:"sessionId"`
	HotelId   string                         `json:"hotelId"`
	Error     *indraframework.IndraException `json:"error"`
}

type DirectResult struct {
}

func (a *DirectResponseDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
